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
  files: ImageFile[];
  totalSize: number;
  loading: boolean;
};

export type ImageFile = {
  name: string;
  id: string;
  size: number;
  type: string;
  dateCreated: string;
  path: string;
  converted: boolean;
  convertedPath: string;
  convertedSize: number;
  ratio: number;
  error: string;
};

export type ImageSessionFiles = {
  imageFiles: Record<string, ImageFile>;
  totalImageSize: number;
};
