package imagewatermark

import (
	"fmt"
	"image"
	"image/draw"
)

// ApplySingle overlays a single watermark onto an input image at a specific position based on the provided configuration.
//
// The function performs the following steps:
//  1. Validates the SingleConfig to ensure all settings are valid.
//  2. Preprocesses the watermark (resize, rotate, and apply opacity).
//  3. Calculates the watermark position based on vertical/horizontal alignment and spacing.
//  4. Creates an RGBA copy of the input image.
//  5. Overlays the watermark onto the input image using the "Over" compositing operator.
//
// The "Over" operator ensures correct alpha blending, making the watermark appear
// with the correct transparency and blending with any existing pixels.
//
// Parameters:
//   - inputImg: The input image to which the watermark will be applied.
//   - watermarkImg: The watermark image to overlay on the input image.
//   - config: SingleConfig struct containing opacity, size, alignment, and rotation settings.
//
// Returns:
//   - An image.Image containing the final image with the watermark applied.
//   - An error if any step fails, such as invalid configuration or failure to load images.
//
// Example:
//
//	config := SingleConfig{
//		GeneralConfig: GeneralConfig{
//			WatermarkWidthPercent: 20,
//			OpacityAlpha:          0.6,
//			RotationDegrees:       0,
//		},
//		VerticalAlign:   VerticalBottom,
//		HorizontalAlign: HorizontalRight,
//		Spacing:         10,
//	}
//	result, err := ApplySingle(inputImg, watermarkImg, config)
//	if err != nil {
//		log.Fatal(err)
//	}
func ApplySingle(
	inputImg image.Image,
	watermarkImg image.Image,
	config SingleConfig,
) (image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid single watermark configuration: %w", err)
	}

	processedWatermark := preprocessWatermark(inputImg, watermarkImg, config.GeneralConfig)

	watermarkPosition := getWatermarkPosition(inputImg, processedWatermark, config.VerticalAlign, config.HorizontalAlign, config.Spacing)

	bounds := inputImg.Bounds()
	outputRGBA := image.NewRGBA(bounds)

	draw.Draw(outputRGBA, bounds, inputImg, bounds.Min, draw.Src)
	wmBounds := processedWatermark.Bounds()
	sr := wmBounds
	dp := watermarkPosition

	dr := image.Rectangle{Min: dp, Max: dp.Add(wmBounds.Size())}

	draw.Draw(outputRGBA, dr, processedWatermark, sr.Min, draw.Over)

	return outputRGBA, nil
}
