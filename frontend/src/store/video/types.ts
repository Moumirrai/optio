export type Stats = {
  byteCount: number;
  imageCount: number;
  timeCount: number;
};

export type Session = {
  count: number;
  savings: number;
  time: number;
};

export type State = {
  stats: Stats;
  session: Session;
  files: VideoFile[];
  totalSize: number;
  loading: boolean;
  current: {
    file: VideoFile | null;
    eta: number;
  }
};

export type VideoFile = {
  name: string;
  id: string;
  size: number;
  dateCreated: string;
  path: string;
  converted: boolean;
  convertedPath: string;
  convertedSize: number;
  error: string;
  width: number;
  height: number;
  duration: number;
  bitrate: number;
  framerate: string;
  progress?: number;
};

export type VideoSessionFiles = {
  videoFiles: Record<string, VideoFile>;
  totalVideoSize: number;
};
