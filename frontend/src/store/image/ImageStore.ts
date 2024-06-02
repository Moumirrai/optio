import {defineStore} from "pinia";
import {ImageSessionFiles, State} from "./types";
import * as imageManager from "../../../wailsjs/go/imageManager/ImageManager";
import {useMainStore} from "../MainStore";

export const useImageStore = defineStore("image", {
  state: (): State => ({
    stats: {
      byteCount: 0,
      imageCount: 0,
      timeCount: 0,
    },
    session: {
      count: 0,
      savings: 0,
      time: 0,
    },
    hash: "",
    files: [],
    totalSize: 0,
    loading: false,
    progress: {
      percentage: 0,
      processed: 0,
      left: 0,
      eta: {
        seconds: 0,
        minutes: 0
      }
    }
  }),
  actions: {
    async addFiles(): Promise<void> {
      const rawData = await imageManager.AddFiles();
      if (rawData === null) {
        return;
      }
      //const data = await JSON.parse(rawData) as ImageSessionFiles;

      //this.files = Object.values(rawData.imageFiles);
      this.files = rawData.imageFiles;
      this.totalSize = rawData.totalImageSize;
      this.hash = Math.random().toString(36).substring(7);
      this.loading = false;
    },
    async clear(): Promise<void> {
      await imageManager.Clear();
      this.files = [];
      this.totalSize = 0;
      this.hash = Math.random().toString(36).substring(7);
    },
    async test(): Promise<void> {
      const pepe = await imageManager.Debug();
      console.log(pepe);
    },
    process(): void {
      const mainStore = useMainStore();
      if (mainStore.ongoingProcess || this.files.length === 0) {
        return;
      }
      imageManager.StartConversion();
      mainStore.ongoingProcess = true;
      this.startTime = Date.now()
      this.resetProgress();
      console.log("Processing");
    },
    stopProcess(): void {
      const mainStore = useMainStore();
      if (!mainStore.ongoingProcess) {
        return;
      }
      imageManager.StopConversion();
      mainStore.ongoingCancelation = true;
      console.log("Stopped");
    },
    resetProgress(): void {
      this.progress = {
        percentage: 0,
        processed: 0,
        left: 0,
        eta: {
          seconds: 0,
          minutes: 0
        }
      }
    }
  },
});

