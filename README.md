# Optio

Optio is an image and video transcoding utility with a user-friendly GUI interface. It allows you to select images in JPG, JPEG, WEBP, or PNG formats and customize their output format, size, aspect ratio, compression level, and more. Optio also supports video transcoding with customizable bitrate and codec options, including support for NVENC for fast transcoding on capable GPUs.

## Features

- **Image Conversion**: Supports JPG, JPEG, WEBP, and PNG formats.
- **Customization Options**: Set desired format, size, aspect ratio, compression level, and more.
- **Metadata Preservation**: Option to preserve original file creation time, with future plans to preserve additional metadata.
- **Efficient Loading**: Loads images only when needed to avoid freezes and wait times.
- **Video Transcoding**: Supports video transcoding with customizable bitrate and codec options. Requires FFmpeg to be installed.
- **NVENC Support**: Achieve fast transcoding on supported GPUs with NVENC.

## Inspiration

Optio was inspired by the [Optimus](https://github.com/Splode/optimus) project but has been built from the ground up with improvements. Unlike Optimus, Optio loads images only when needed, ensuring a smoother user experience without freezes and wait times.

## Requirements

- FFmpeg (for video transcoding)

## Live Development

To run in live development mode, run `wails dev` in the project directory.

## Building

To build a redistributable, production mode package, use `wails build`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your improvements.
