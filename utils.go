package imagewatermark

import (
	"image"
	"math/rand"
	"os"
	"time"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

func loadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func getNewWatermarkWidth(originalImage image.Image, watermarkWidthPercentage float64) int {
	originalImageWidth := float64(originalImage.Bounds().Dx())

	return int(originalImageWidth * (watermarkWidthPercentage / 100))
}

func getWatermarkPosition(inputImage, watermarkImage image.Image, verticalAlign VerticalAlign, horizontalAlign HorizontalAlign, spacing int) image.Point {
	inputImageBounds := inputImage.Bounds()
	watermarkImageBounds := watermarkImage.Bounds()

	var position image.Point

	switch horizontalAlign {
	case HorizontalLeft:
		position.X = spacing
	case HorizontalMiddle:
		position.X = int(inputImageBounds.Dx()/2) - (watermarkImageBounds.Dx() / 2)
	case HorizontalRight:
		position.X = inputImageBounds.Dx() - watermarkImageBounds.Dx() - spacing
	case HorizontalRandom:
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		position.X = int(r.Float64()*(float64(inputImageBounds.Dx())-float64(watermarkImageBounds.Dx()))) + spacing
	default:
		return image.Point{}
	}

	switch verticalAlign {
	case VerticalTop:
		position.Y = spacing
	case VerticalMiddle:
		position.Y = int((inputImageBounds.Dy() / 2) - (watermarkImageBounds.Dy() / 2))
	case VerticalBottom:
		position.Y = inputImageBounds.Dy() - watermarkImageBounds.Dy() - spacing
	case VerticalRandom:
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		position.Y = int(r.Float64()*(float64(inputImageBounds.Dy())-float64(watermarkImageBounds.Dy()))) + spacing
	default:
		return image.Point{}
	}

	return position
}
