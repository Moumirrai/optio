package webp

import (
	"bufio"
	"image"
	"io"

	"github.com/chai2010/webp"
)

// Options represent WebP encoding options.
type Options struct {
	Lossless bool `json:"lossless"`
	Quality  int  `json:"quality"`
}

// DecodeWebp a webp file and return an image.
func DecodeWebp(r io.Reader) (image.Image, error) {
	i, err := webp.Decode(r)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// EncodeWebp encodes an image into webp and returns a buffer.
func EncodeWebp(i image.Image, writer *bufio.Writer, o *Options) (err error) {
	err = webp.Encode(writer, i, &webp.Options{Lossless: o.Lossless, Quality: float32(o.Quality)})
	return err
}
