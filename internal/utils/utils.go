package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"golang.org/x/image/draw"
)

func CropImage(img image.Image, block types.Block) error {
	var signatureFieldBoundingBox image.Rectangle

	signatureFieldBoundingBox = image.Rect(
		int(block.Geometry.BoundingBox.Left*float32(img.Bounds().Dx())),
		int(block.Geometry.BoundingBox.Top*float32(img.Bounds().Dy())),
		int((block.Geometry.BoundingBox.Left+block.Geometry.BoundingBox.Width)*float32(img.Bounds().Dx())),
		int((block.Geometry.BoundingBox.Top+block.Geometry.BoundingBox.Height)*float32(img.Bounds().Dy())),
	)

	if signatureFieldBoundingBox.Empty() {
		return errors.New("bounding box is empty")
	}

	croppedImg := image.NewRGBA(image.Rect(0, 0, signatureFieldBoundingBox.Dx(), signatureFieldBoundingBox.Dy()))
	draw.Copy(croppedImg, image.Point{}, img, signatureFieldBoundingBox, draw.Src, nil)

	croppedFile, err := os.Create("cropped_signature_field.jpg")
	if err != nil {
		return err
	}
	defer croppedFile.Close()
	jpeg.Encode(croppedFile, croppedImg, nil)

	return nil
}
