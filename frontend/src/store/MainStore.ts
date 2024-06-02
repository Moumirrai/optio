// Utilities
import { defineStore } from "pinia";
import {codec, Config, editorMode, Size, State, target} from "./types";
import {
  GetAppConfig,
  OpenOutputDir,
  SetConfig,
  SetOutDir,
} from "../../wailsjs/go/config/Config";
import { useGlobalToast } from "@/plugins/toast";

const globalToast = useGlobalToast();

let saveTimeout: NodeJS.Timeout | null = null;

export const useMainStore = defineStore("main", {
  state: (): State => ({
    config: {
      outDir: "",
      target: target.JPG,
      prefix: "",
      suffix: "",
      sizes: [],
      imageOpt: {
        jpegOpt: {
          quality: 0,
          preserveMetadata: false,
        },
        pngOpt: {
          quality: 0,
        },
        webpOpt: {
          quality: 0,
          lossless: false,
        },
      },
      videoOpt: {
        bitrate: 50000,
        height: 0,
        width: 0,
        codec: codec.X264,
        percentageMode: true,
        percentage: 50,
      },
      activeSize: undefined,
      preserveCreationTime: false,
    },
    saveings: 0,
    sizeModal: false,
    configModal: false,
    configLoaded: false,
    ongoingProcess: false,
    ongoingCancelation: false,
    progress: 0,
    mode: editorMode.IMAGE,
  }),
  getters: {
    //getConfig: (state): Config => state.config,
    /* getStats: (state): Stats => state.stats,
        getSession: (state): Session => state.session, */
    isProcessing: (state): boolean =>
      state.ongoingProcess || state.ongoingCancelation,
    canSwitchMode: (state): boolean =>
      !state.ongoingProcess && !state.ongoingCancelation,
  },
  actions: {
    async getConfig(): Promise<Config> {
      try {
        let config = (await GetAppConfig()) as Config;
        if (config.sizes === null) config.sizes = [];
        if (config.activeSize == "") config.activeSize = undefined;
        console.log(config);
        this.config = config;
        console.log(this.config);
        this.configLoaded = true;
        return config;
      } catch (err) {
        console.error(err);
        return this.config;
      }
    },
    async setConfig(): Promise<void> {
      if (saveTimeout) {
        console.log("Clearing timeout");
        clearTimeout(saveTimeout);
      }

      saveTimeout = setTimeout(async () => {
        console.log("Saving config");
        console.log(JSON.stringify(this.config));
        await SetConfig(JSON.stringify(this.config));
        saveTimeout = null;
      }, 200);
    },
    async setOutDirImages(): Promise<void> {
      const outDir = await SetOutDir();
      this.config.outDir = outDir;
    },
    async addSize(size: Size): Promise<void> {
      if (typeof size.width !== "number") {
        size.width = Number(size.width);
      }
      if (typeof size.height !== "number") {
        size.height = Number(size.height);
      }
      console.log(this.config.sizes);
      this.config.sizes.push(size);
      await this.setConfig();
    },
    async removeSize(index: number): Promise<void> {
      this.config.sizes.splice(index, 1);
      await this.setConfig();
    },
    openOutDirImages(): void {
      OpenOutputDir();
    },
  },
});
