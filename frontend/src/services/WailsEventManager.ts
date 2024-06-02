import { useMainStore, useVideoStore, useImageStore } from "@/store";
import { useGlobalToast } from "@/plugins/toast";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { formatTime, formatSize, formatDate } from "@/utils/format";

enum NotificationType {
  INFO = "info",
  SUCCESS = "success",
  ERROR = "error",
}

type Notification = {
  msg: string;
  type: NotificationType;
};

class WailsEventManager {
  mainStore: ReturnType<typeof useMainStore>;
  videoStore: ReturnType<typeof useVideoStore>;
  imageStore: ReturnType<typeof useImageStore>;
  globalToast = useGlobalToast();
  constructor(
    mainStore: ReturnType<typeof useMainStore>,
    imageStore: ReturnType<typeof useImageStore>,
    videoStore: ReturnType<typeof useVideoStore>
  ) {
    this.mainStore = mainStore;
    this.videoStore = videoStore;
    this.imageStore = imageStore;
    this.notificationListener();
    this.completeListener();
  }
  private notificationListener() {
    EventsOn("notify", (data: Notification) => {
      console.log(data);
      console.log(data.msg);
      switch (data.type) {
        case NotificationType.INFO:
          this.globalToast(data.msg);
          break;
        case NotificationType.SUCCESS:
          this.globalToast.success(data.msg);
          break;
        case NotificationType.ERROR:
          this.globalToast.error(data.msg);
          break;
      }
    });
  }

  private completeListener() {
    let test = 0;
    let startTime: number | null = null;
    EventsOn("conversion:image:progress", (data: ConversionInfo) => {
      const file = this.imageStore.files.find((f) => f.id === data.id);
      if (file) {
        if (!startTime) {
          startTime = Date.now();
        }
        file.converted = true;
        file.convertedSize = data.newSize;
        file.ratio = data.ratio;
        const percentage = (file.size / this.imageStore.totalSize) * 100;
        this.mainStore.progress += percentage;
        const saved = file.size - data.newSize;
        //if saves is a number
        if (!isNaN(saved)) {
          console.log(saved);
          this.mainStore.saveings += saved;
        }

        test += file.size;

        if (startTime) {
          const elapsedTime = (Date.now() - startTime) / 1000;
          const speed = test / elapsedTime;
          // Calculate the remaining size in bytes
          const remainingSize = this.imageStore.totalSize - test;

          // Calculate the ETA in seconds
          const eta = remainingSize / speed;

          // Convert the ETA to minutes and seconds
          const etaMinutes = Math.floor(eta / 60);
          const etaSeconds = Math.floor(eta % 60);

          this.imageStore.progress.eta.minutes = etaMinutes;
          this.imageStore.progress.eta.seconds = etaSeconds;

          console.log(`ETA: ${etaMinutes} minutes, ${etaSeconds} seconds`);
        }
      }
    });

    EventsOn("conversion:image:complete", (data: ConversionComplete) => {
      console.log(data);
      this.globalToast.success(
        `Converted ${data.length} images in ${formatTime(data.time)}`
      );
      this.mainStore.ongoingProcess = false;
      this.mainStore.ongoingCancelation = false;
      this.mainStore.progress = 0;
      test = 0;
      startTime = null;
    });

    EventsOn("conversion:video:progress", (data: VideoConversionProgress) => {
      const file = this.videoStore.files.find((f) => f.id === data.id);
      if (!file) return
      file!.progress = data.progress;
      this.videoStore.current = {
        eta: data.eta,
        file,
      }
      console.log(data);
    });

    EventsOn("conversion:video:file", (data: VideoConversionInfo) => {
      const file = this.videoStore.files.find((f) => f.id === data.id);
      file!.converted = true;
      file!.progress = 0;
      console.log(data);
    });

    EventsOn("conversion:video:complete", () => {
      this.globalToast.success("Video conversion complete");
      this.mainStore.ongoingProcess = false;
      this.mainStore.ongoingCancelation = false;
      this.mainStore.progress = 0;
      this.videoStore.current = {
        eta: 0,
        file: null,
      }
      test = 0;
      startTime = null;
      console.log("Video conversion complete");
    });

    EventsOn("video:addingFiles", () => {
      this.videoStore.loading = true;
    });

    EventsOn("image:addingFiles", () => {
      this.imageStore.loading = true;
    });
  }
}

export { WailsEventManager };

type ConversionInfo = {
  id: string;
  newSize: number;
  ratio: number;
  time: number;
};

type VideoConversionProgress = {
  id: string;
  progress: number;
  eta: number;
};

type VideoConversionInfo = {
  id: string;
  time: number;
};

type ConversionComplete = {
  length: number;
  time: number;
};

//TODO: DELETE THIS
