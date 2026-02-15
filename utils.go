package imagewatermark

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

// OpenImage loads an image from the specified path with auto-orientation support.
//
// The function automatically corrects the image orientation based on EXIF metadata,
// ensuring that images taken with different device orientations display correctly.
//
// Parameters:
//   - path: The file path to the image to be loaded.
//
// Returns:
//   - An image.Image containing the loaded image data.
//   - An error if the file cannot be read or the format is not supported.
func OpenImage(path string) (image.Image, error) {
	return imaging.Open(path, imaging.AutoOrientation(true))
}

// getNewWatermarkWidth calculates the new width of the watermark based on a percentage of the original image width.
//
// This function is used to scale the watermark proportionally to the input image dimensions.
// For example, if the input image is 1000px wide and watermarkWidthPercentage is 20,
// the function will return 200.
//
// Parameters:
//   - originalImage: The original input image used as reference for width calculation.
//   - watermarkWidthPercentage: The desired watermark width as a percentage of the original image width (0-100).
//
// Returns:
//   - An int representing the new watermark width in pixels.
func getNewWatermarkWidth(originalImage image.Image, watermarkWidthPercentage float64) int {
	originalImageWidth := float64(originalImage.Bounds().Dx())

	return int(originalImageWidth * (watermarkWidthPercentage / 100))
}

// preprocessWatermark resizes and rotates the watermark image based on the provided configuration.
//
// The function performs the following operations in order:
//  1. Calculates the new watermark width based on the GeneralConfig percentage.
//  2. Resizes the watermark image using the CatmullRom interpolation filter.
//  3. Rotates the watermark if RotationDegrees is not zero.
//  4. Applies opacity/transparency to the watermark if OpacityAlpha is less than 1.
//
// Parameters:
//   - originalImage: The input image used as reference for width calculation.
//   - watermarkImage: The watermark image to be processed.
//   - config: GeneralConfig containing size, rotation, and opacity settings.
//
// Returns:
//   - An image.Image containing the processed watermark, ready to be applied to the input image.
func preprocessWatermark(originalImage image.Image, watermarkImage image.Image, config GeneralConfig) image.Image {

	watermarkWidth := getNewWatermarkWidth(originalImage, config.WatermarkWidthPercent)

	resampleFilter := config.ResampleFilter
	if resampleFilter.Support <= 0 {
		resampleFilter = imaging.CatmullRom
	}

	resizedWatermarkImage := imaging.Resize(watermarkImage, watermarkWidth, 0, resampleFilter)

	if config.RotationDegrees != 0 {
		resizedWatermarkImage = imaging.Rotate(resizedWatermarkImage, config.RotationDegrees, image.Transparent)
	}

	if config.OpacityAlpha < 1 {
		resizedWatermarkImage = applyOpacity(resizedWatermarkImage, config.OpacityAlpha)
	}

	return resizedWatermarkImage
}

// getWatermarkPosition calculates the position of the single watermark on the input image based on alignment and spacing settings.
//
// This function computes the X and Y coordinates where the watermark should be placed.
// It supports eight different alignment options (four cardinal directions and four random options)
// and applies spacing/padding from the calculated edge.
//
// For random alignment options (VerticalRandom/HorizontalRandom), a random position is selected
// within valid bounds. If there isn't enough space, the function defaults to the edge position.
//
// Parameters:
//   - inputImage: The input image where the watermark will be placed.
//   - watermarkImage: The watermark image whose dimensions are used for boundary calculations.
//   - verticalAlign: Vertical alignment option (top, middle, bottom, or random).
//   - horizontalAlign: Horizontal alignment option (left, middle, right, or random).
//   - spacing: Padding in pixels from the aligned edge.
//
// Returns:
//   - An image.Point containing the calculated X and Y coordinates for the watermark placement.
func getWatermarkPosition(inputImage, watermarkImage image.Image, verticalAlign VerticalAlign, horizontalAlign HorizontalAlign, spacing int) image.Point {
	inBounds := inputImage.Bounds()
	wmBounds := watermarkImage.Bounds()

	inW, inH := inBounds.Dx(), inBounds.Dy()
	wmW, wmH := wmBounds.Dx(), wmBounds.Dy()

	var position image.Point

	switch horizontalAlign {
	case HorizontalLeft:
		position.X = spacing
	case HorizontalMiddle:
		position.X = (inW - wmW) / 2
	case HorizontalRight:
		position.X = inW - wmW - spacing
	case HorizontalRandom:
		minX := spacing
		maxX := inW - wmW - spacing
		if maxX > minX {
			position.X = rand.Intn(maxX-minX+1) + minX
		} else {
			position.X = minX
		}
	default:
		position.X = (inW - wmW) / 2
	}

	switch verticalAlign {
	case VerticalTop:
		position.Y = spacing
	case VerticalMiddle:
		position.Y = (inH - wmH) / 2
	case VerticalBottom:
		position.Y = inH - wmH - spacing
	case VerticalRandom:
		minY := spacing
		maxY := inH - wmH - spacing
		if maxY > minY {
			position.Y = rand.Intn(maxY-minY+1) + minY
		} else {
			position.Y = minY
		}
	default:
		position.Y = (inH - wmH) / 2
	}

	return position
}

// applyOpacity modifies the transparency/opacity of an image by multiplying the alpha channel of each pixel.
//
// This function iterates through every pixel in the input image and adjusts its alpha channel
// based on the provided opacity value. An opacity of 1.0 means fully opaque (no change),
// while lower values (e.g., 0.5) make the image more transparent.
//
// The function returns a new *image.NRGBA with the modified opacity, without altering the original image.
//
// Parameters:
//   - img: The input image whose opacity should be modified.
//   - opacityAlpha: The opacity multiplier (0.0 to 1.0). For example, 0.5 reduces opacity to 50%.
//
// Returns:
//   - A pointer to an image.NRGBA containing the image with the modified opacity applied to all pixels.
func applyOpacity(img image.Image, opacityAlpha float64) *image.NRGBA {
	bounds := img.Bounds()
	result := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newA := uint8(float64(a>>8) * opacityAlpha)
			result.Set(x, y, color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: newA,
			})
		}
	}

	return result
}
