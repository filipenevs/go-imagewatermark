package imagewatermark

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
)

// ApplySingle overlays a single watermark onto an input image at a specific position based on the provided configuration.
//
// The function performs the following steps:
//  1. Validates the SingleConfig to ensure all settings are valid.
//  2. Loads the input image and watermark image from disk.
//  3. Preprocesses the watermark (resize, rotate, and apply opacity).
//  4. Calculates the watermark position based on vertical/horizontal alignment and spacing.
//  5. Creates an RGBA copy of the input image.
//  6. Overlays the watermark onto the input image using the "Over" compositing operator.
//
// The "Over" operator ensures correct alpha blending, making the watermark appear
// with the correct transparency and blending with any existing pixels.
//
// Parameters:
//   - config: SingleConfig struct containing paths, opacity, size, alignment, and rotation settings.
//
// Returns:
//   - An image.Image containing the final image with the watermark applied.
//   - An error if any step fails, such as invalid configuration or failure to load images.
//
// Example:
//
//	config := SingleConfig{
//		GeneralConfig: GeneralConfig{
//			InputPath:             "input.jpg",
//			WatermarkPath:         "logo.png",
//			WatermarkWidthPercent: 20,
//			OpacityAlpha:          0.6,
//			RotationDegrees:       0,
//		},
//		VerticalAlign:   VerticalBottom,
//		HorizontalAlign: HorizontalRight,
//		Spacing:         10,
//	}
//	result, err := ApplySingle(config)
//	if err != nil {
//		log.Fatal(err)
//	}
func ApplySingle(
	config SingleConfig,
) (image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid single watermark configuration: %w", err)
	}

	var inputImg, watermarkImg image.Image
	var inputErr, watermarkErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		inputImg, inputErr = openImage(config.InputPath)
	}()

	go func() {
		defer wg.Done()
		watermarkImg, watermarkErr = openImage(config.WatermarkPath)
	}()

	wg.Wait()

	if inputErr != nil {
		return nil, fmt.Errorf("failed to load input image: %w", inputErr)
	}
	if watermarkErr != nil {
		return nil, fmt.Errorf("failed to load watermark image: %w", watermarkErr)
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
