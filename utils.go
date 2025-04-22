package imagewatermark

import (
	"image"
	"math/rand"
	"time"
)

// getNewWatermarkWidth calculates the new watermark width based on a percentage.
func getNewWatermarkWidth(originalImage image.Image, watermarkWidthPercentage float64) int {
	originalImageWidth := float64(originalImage.Bounds().Dx())

	return int(originalImageWidth * (watermarkWidthPercentage / 100))
}

// getWatermarkPosition calculates the position of the watermark based on alignment and spacing.
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
		minX := spacing
		maxX := max(inputImageBounds.Dx()-watermarkImageBounds.Dx()-spacing, minX)
		position.X = r.Intn(maxX-minX+1) + minX
	default:
		position.X = int(inputImageBounds.Dx()/2) - (watermarkImageBounds.Dx() / 2)
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
		minY := spacing
		maxY := max(inputImageBounds.Dy()-watermarkImageBounds.Dy()-spacing, minY)
		position.Y = r.Intn(maxY-minY+1) + minY
	default:
		position.Y = int((inputImageBounds.Dy() / 2) - (watermarkImageBounds.Dy() / 2))
	}

	return position
}
