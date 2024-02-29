package videoManager

import (
	"context"
	"optio/backend/metadata"
	"os"
	"strconv"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

type VideoFileInfo struct {
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
	Framerate     string    `json:"framerate"`
}

func GetFileInfo(path string, ctx context.Context) (VideoFileInfo, error) {
	stat, err := os.Stat(path)
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		println(err.Error())
		return VideoFileInfo{}, err
	}
	stream := data.Streams[0]
	duration, err := strconv.ParseFloat(data.Streams[0].Duration, 64)
	framerate := data.Streams[0].RFrameRate
	bitrate, err := strconv.Atoi(data.Streams[0].BitRate)
	if err != nil {
		println(err.Error())
		return VideoFileInfo{}, err
	}
	return VideoFileInfo{
		Name:          stat.Name(),
		ID:            stat.Name(),
		Size:          stat.Size(),
		DateCreated:   metadata.GetCTime(stat),
		Path:          data.Format.Filename,
		ConvertedPath: "",
		Converted:     false,
		Error:         "",
		ConvertedSize: 0,
		Width:         stream.Width,
		Height:        stream.Height,
		Duration:      duration,
		Bitrate:       bitrate,
		Framerate:     framerate,
	}, nil
}
