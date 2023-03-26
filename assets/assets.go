package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
)

var (
	//go:embed Boden.png
	bodenData []byte
	BodenImg  image.Image

	//go:embed Tisch.png
	tischData []byte
	TischImg  image.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(bodenData))
	if err != nil {
		panic(err)
	}
	BodenImg = img

	img, _, err = image.Decode(bytes.NewReader(tischData))
	if err != nil {
		panic(err)
	}
	TischImg = img
}
