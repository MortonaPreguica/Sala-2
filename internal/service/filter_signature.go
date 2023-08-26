package service

import (
	"image"
	"image/color"
)

func toBlackAndWhite(originalImage image.Image, whiteThreshold uint8) image.Image {
	size := originalImage.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)

	modifiedImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := originalImage.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			modifiedColorValue := originalColor.R

			if modifiedColorValue >= whiteThreshold {
				modifiedColorValue = 255
			} else {
				modifiedColorValue = 0
			}

			modifiedColor := color.RGBA{
				R: modifiedColorValue,
				G: modifiedColorValue,
				B: modifiedColorValue,
				A: originalColor.A,
			}

			modifiedImg.Set(x, y, modifiedColor)
		}
	}

	return modifiedImg
}

func FilterSignature(img image.Image) image.Image {
	const whiteThreshold = uint8(80)
	blackAndWhiteImg := toBlackAndWhite(img, whiteThreshold)
	return blackAndWhiteImg
}
