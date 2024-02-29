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
    jpegOpt: JPEGOpt;
    pngOpt: PNGOpt;
    webpOpt: WebPOpt;
    videoOpt: VideoOpt;
    preserveCreationTime: boolean;
};

export enum resizeStrategy {
    Fill = 0,
    Fit = 1,
    Smart = 2,
}

export type VideoOpt = {
    bitrate: number;
    height: number;
    width: number;
    codec: string;
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
    mode: editorMode;
};

export enum editorMode {
    IMAGE = "image",
    VIDEO = "video",
}

