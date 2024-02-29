// Utilities
import { defineStore } from "pinia";
import { State, VideoSessionFiles } from "./types";
import * as videoManager from "../../../wailsjs/go/videoManager/VideoManager";
import { useMainStore } from "../MainStore";

export const useVideoStore = defineStore("video", {
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
  getters: {
    //getConfig: (state): Config => state.config,
    /* getStats: (state): Stats => state.stats,
    getSession: (state): Session => state.session, */
  },
  actions: {
    async addFiles(): Promise<void> {
      const rawData = await videoManager.AddFiles();
      this.loading = false;
      if (rawData === "") {
        return;
      }
      const data = JSON.parse(rawData) as VideoSessionFiles;
      this.files = Object.values(data.videoFiles);
      this.totalSize = data.totalVideoSize;
    },
    async clear(): Promise<void> {
      await videoManager.Clear();
      this.files = [];
      this.totalSize = 0;
    },
    process(): void {
      const mainStore = useMainStore();
      if (mainStore.ongoingProcess || this.files.length === 0) {
        return;
      }
      videoManager.StartReencoding();
      mainStore.ongoingProcess = true;
      console.log("Processing");
    },
    stopProcess(): void {
      const mainStore = useMainStore();
      if (!mainStore.ongoingProcess) {
        return;
      }
      videoManager.StopReencoding();
      mainStore.ongoingCancelation = true;
      console.log("Stopped");
    },
  },
});
