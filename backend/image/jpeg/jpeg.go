package jpeg

import (
	"bufio"
	"image"
	"image/jpeg"
	"io"
)

// Options represent JPEG encoding options.
type Options struct {
	Quality          int  `json:"quality"`
	PreserveMetadata bool `json:"preserveMetadata"`
}

// DecodeJPEG decodes a JPEG file and return an image.
func DecodeJPEG(r io.Reader) (image.Image, error) {
	i, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// EncodeJPEG encodes an image into JPEG and returns a buffer.
func EncodeJPEG(i image.Image, writer *bufio.Writer, o *Options) (err error) {
	err = jpeg.Encode(writer, i, &jpeg.Options{Quality: o.Quality})
	return err
}
