package models

import (
	"optio/backend/image/jpeg"
	"optio/backend/image/png"
	"optio/backend/image/webp"
)

type VideoOptions struct {
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Bitrate        int    `json:"bitrate"`
	Codec          string `json:"codec"`
	PercentageMode bool   `json:"percentageMode"`
	Percentage     int    `json:"percentage"`
}

type ImageOptions struct {
	JpegOpt *jpeg.Options `json:"jpegOpt"`
	PngOpt  *png.Options  `json:"pngOpt"`
	WebpOpt *webp.Options `json:"webpOpt"`
}
