export enum target {
    JPG = "jpg",
    PNG = "png",
    WEBP = "webp",
}

export type Config = {
    outDir: string;
    target: target;
    prefix: string;
    suffix: string;
    sizes: Size[];
    activeSize?: string;
    imageOpt: {
      jpegOpt: JPEGOpt;
      pngOpt: PNGOpt;
      webpOpt: WebPOpt;
    };
    videoOpt: VideoOpt;
    preserveCreationTime: boolean;
};

export enum resizeStrategy {
    Fill = 0,
    Fit = 1,
    Smart = 2,
}

export enum codec {
    H264 = "h264_nvenc",
    H265 = "hevc_nvenc",
    VP9 = "vp9",
    AV1 = "av1",
    X264 = "libx264",
    X265 = "libx265",
    VP8 = "libvpx",
}

export type VideoOpt = {
    bitrate: number;
    height: number;
    width: number;
    codec: codec;
    percentageMode: boolean;
    percentage: number;
};

export type Size = {
    height: number | null;
    width: number | null;
    strategy: resizeStrategy;
    name: string;
};

export type JPEGOpt = {
    quality: number;
    preserveMetadata: boolean;
};

export type PNGOpt = {
    quality: number;
};

export type WebPOpt = {
    quality: number;
    lossless: boolean;
};

export type State = {
    config: Config;
    sizeModal: boolean;
    configModal: boolean;
    configLoaded: boolean;
    ongoingProcess: boolean;
    ongoingCancelation: boolean;
    progress: number;
    saveings: number;
    mode: editorMode;
};

export enum editorMode {
    IMAGE = "image",
    VIDEO = "video",
}

