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
    files: [],
    totalSize: 0,
    loading: false,
  }),
  actions: {
    async addFiles(): Promise<void> {
      const rawData = await imageManager.AddFiles();
      this.loading = false;
      if (rawData === "") {
        return;
      }
      const data = JSON.parse(rawData) as ImageSessionFiles;
      this.files = Object.values(data.imageFiles);
      this.totalSize = data.totalImageSize;
    },
    async clear(): Promise<void> {
      await imageManager.Clear();
      this.files = [];
      this.totalSize = 0;
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
  },
});
