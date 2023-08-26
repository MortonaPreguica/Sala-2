package utils

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"golang.org/x/image/draw"
)

func CropImage(img image.Image, block types.Block) (image.Image, error) {
	var signatureFieldBoundingBox image.Rectangle

	signatureFieldBoundingBox = image.Rect(
		int(block.Geometry.BoundingBox.Left*float32(img.Bounds().Dx())),
		int(block.Geometry.BoundingBox.Top*float32(img.Bounds().Dy())),
		int((block.Geometry.BoundingBox.Left+block.Geometry.BoundingBox.Width)*float32(img.Bounds().Dx())),
		int((block.Geometry.BoundingBox.Top+block.Geometry.BoundingBox.Height)*float32(img.Bounds().Dy())),
	)

	if signatureFieldBoundingBox.Empty() {
		return nil, errors.New("bounding box is empty")
	}

	croppedImg := image.NewRGBA(image.Rect(0, 0, signatureFieldBoundingBox.Dx(), signatureFieldBoundingBox.Dy()))
	draw.Copy(croppedImg, image.Point{}, img, signatureFieldBoundingBox, draw.Src, nil)

	return croppedImg, nil
}

func ImageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImageToBytesJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MakeWhiteTransparent(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()

			whiteThreshold := uint32(0xffff) - 0x100
			if r > whiteThreshold && g > whiteThreshold && b > whiteThreshold {
				a = 0 // Define a transparência para 0 (totalmente transparente)
			}

			newImg.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	return newImg
}
