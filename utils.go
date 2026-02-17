package imagewatermark

import (
	"image"
	"image/color"
	"image/draw"
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
func getWatermarkPosition(watermarkImage, inputImage image.Image, verticalAlign VerticalAlign, horizontalAlign HorizontalAlign, spacing int) image.Point {
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

// rotateImage rotates the input image by the specified degrees in the configuration.
// The function uses the "imaging" library to perform the rotation, which handles the necessary calculations
// to rotate the image around its center and fills any empty areas with transparency.
//
// Parameters:
//   - img: The input image to be rotated.
//   - rotationDegrees: The angle in degrees to rotate the image. Positive values rotate clockwise.
//
// Returns:
//   - An image.Image containing the rotated image with transparent background for any new areas created by the rotation.
func rotateImage(img image.Image, rotationDegrees float64) image.Image {
	return imaging.Rotate(img, rotationDegrees, image.Transparent)
}

// resizeWatermark resizes the watermark image based on the specified percentage of the base image width.
//
// This function calculates the new width for the watermark using the getNewWatermarkWidth function
// and then resizes the watermark image using the "imaging" library. The height is automatically
// adjusted to maintain the aspect ratio. The resampling filter can be specified in the configuration,
// and if not provided, it defaults to CatmullRom for high-quality resizing.
//
// Parameters:
//   - watermarkImg: The original watermark image to be resized.
//   - baseImage: The input image used as reference for width calculation.
//   - config: GeneralConfig containing the WatermarkWidthPercent and optional ResampleFilter.
//
// Returns:
//   - An image.Image containing the resized watermark, ready to be applied to the input image.
func resizeWatermark(watermarkImg image.Image, baseImage image.Image, config GeneralConfig) image.Image {
	watermarkWidth := getNewWatermarkWidth(baseImage, config.WatermarkWidthPercent)

	resampleFilter := config.ResampleFilter
	if resampleFilter.Support <= 0 {
		resampleFilter = imaging.CatmullRom
	}

	return imaging.Resize(watermarkImg, watermarkWidth, 0, resampleFilter)
}

// generateBaseCanvas creates a new RGBA canvas based on the input image dimensions and draws the input image onto it.
//
// This function is used to create a mutable canvas that can be modified with watermarks.
// It initializes a new *image.RGBA with the same dimensions as the input image and uses the "draw" package
// to copy the input image onto the canvas. The resulting RGBA image allows for proper alpha compositing when applying watermarks.
//
// Parameters:
//   - inputImg: The original input image that serves as the base for the canvas.
//
// Returns:
//   - A pointer to an image.RGBA containing the input image drawn onto a new canvas, ready for watermark application.
func generateBaseCanvas(inputImg image.Image) *image.RGBA {
	bounds := inputImg.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, inputImg, bounds.Min, draw.Src)
	return canvas
}

// drawWatermarkAtPosition is a helper function that draws the watermark image onto the canvas at a specific position.
//
// This function calculates the destination rectangle based on the provided position and the dimensions of the watermark image.
// It then uses the "draw" package to composite the watermark onto the canvas using the "Over" operator, which ensures proper blending.
//
// Parameters:
//   - canvas: The RGBA image onto which the watermark will be drawn.
//   - watermarkImg: The watermark image to be drawn.
//   - pos: The image.Point representing the top-left corner where the watermark should be placed.
func drawWatermarkAtPosition(canvas *image.RGBA, watermarkImg image.Image, pos image.Point) {
	wmBounds := watermarkImg.Bounds()
	dr := image.Rectangle{
		Min: pos,
		Max: pos.Add(wmBounds.Size()),
	}
	draw.Draw(canvas, dr, watermarkImg, image.Point{0, 0}, draw.Over)
}
