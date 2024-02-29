package video

import (
	"fmt"

	"github.com/xfrr/goffmpeg/transcoder"
)

/* type VideoFileInfo struct {
	Name          string    `json:"name"`
	ID            string    `json:"id"`
	Size          int64     `json:"size"`
	DateCreated   time.Time `json:"dateCreated"`
	Path          string    `json:"path"`
	ConvertedPath string    `json:"convertedPath"`
	Converted     bool      `json:"converted"`
	Error         string    `json:"error"`
	ConvertedSize int64     `json:"convertedSize"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	Duration      float64   `json:"duration"`
	Bitrate       int       `json:"bitrate"`
}
*/

type Options struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Bitrate int    `json:"bitrate"`
	Codec   string `json:"codec"`
}

func Reencode(input string, output string, options Options) (*transcoder.Transcoder, <-chan transcoder.Progress, <-chan error) {
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(input, output)
	if err != nil {
		panic(err)
	}
	trans.MediaFile().SetVideoCodec("h264_nvenc")
	trans.MediaFile().SetVideoBitRate(fmt.Sprintf("%dk", options.Bitrate))

	done := trans.Run(true)
	progress := trans.Output()

	return trans, progress, done
}
