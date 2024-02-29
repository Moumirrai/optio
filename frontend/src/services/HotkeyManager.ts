import {
  useMainStore,
  useVideoStore,
  useImageStore,
  editorMode,
} from "@/store";
import { useGlobalToast } from "@/plugins/toast";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { formatTime, formatSize, formatDate } from "@/utils/format";

class HotkeyManager {
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
    this.altListener();
  }
  private altListener() {
    document.addEventListener("keydown", (e) => {
      if (!e.altKey) return;

      const isProcessing = this.mainStore.isProcessing;
      const isImageMode = this.isImageMode();
      const store = isImageMode ? this.imageStore : this.videoStore;

      if (isProcessing && e.key === "s") {
        store.stopProcess();
        return;
      }

      if (!isProcessing) {
        switch (e.key) {
          case "a":
            if (store.loading) return;
            store.addFiles();
            break;
          case "i":
            this.mainStore.mode = editorMode.IMAGE;
            break;
          case "v":
            this.mainStore.mode = editorMode.VIDEO;
            break;
          case "o":
            this.mainStore.setOutDirImages();
            break;
          case "e":
            this.mainStore.openOutDirImages();
            break;
          case "c":
            store.clear();
            break;
          //on enter press
          case "Enter":
            if (store.files.length == 0) return;
            store.process();
            break;
        }
      }
    });
  }

  private isImageMode(): boolean {
    return this.mainStore.mode === editorMode.IMAGE;
  }
}

export { HotkeyManager };

type ConversionInfo = {
  id: string;
  newSize: number;
  ratio: number;
  time: number;
};

type VideoConversionProgress = {
  id: string;
  progress: number;
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
