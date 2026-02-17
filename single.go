package imagewatermark

import (
	"fmt"
	"image"
	"runtime"
	"sync"
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

	preparedWM := watermarkImg

	if config.OpacityAlpha < 1 {
		preparedWM = applyOpacity(preparedWM, config.OpacityAlpha)
	}

	if config.RotationDegrees != 0 {
		preparedWM = rotateImage(preparedWM, config.RotationDegrees)
	}

	preparedWM = resizeWatermark(preparedWM, inputImg, config.GeneralConfig)
	watermarkPosition := getWatermarkPosition(preparedWM, inputImg, config.VerticalAlign, config.HorizontalAlign, config.Spacing)

	canvas := generateBaseCanvas(inputImg)

	drawWatermarkAtPosition(canvas, preparedWM, watermarkPosition)

	return canvas, nil
}

// BatchApplySingle applies a single watermark to a batch of input images concurrently based on the provided configuration.
//
// This function performs the following steps:
//  1. Validates the SingleConfig to ensure all settings are valid.
//  2. Preprocesses the watermark (rotate and apply opacity) once for efficiency.
//  3. Uses a worker pool to process multiple images concurrently, applying the watermark to each image.
//  4. Calculates the watermark position for each image based on alignment and spacing settings.
//  5. Returns a slice of images with the watermark applied.
//
// Parameters:
//   - inputImg: A slice of input images to which the watermark will be applied.
//   - watermarkImg: The watermark image to overlay on each input image.
//   - config: SingleConfig struct containing opacity, size, alignment, rotation, and concurrency settings.
//
// Returns:
//   - A slice of image.Image objects containing the final images with the watermark applied.
//   - An error if any step fails, such as invalid configuration.
//
// Example:
//
//	cfg := SingleConfig{
//		GeneralConfig: GeneralConfig{
//			WatermarkWidthPercent: 20,
//			OpacityAlpha:          0.6,
//			RotationDegrees:       0,
//		},
//		VerticalAlign:   VerticalBottom,
//		HorizontalAlign: HorizontalRight,
//		Spacing:         10,
//		MaxWorkers:      4,
//	}
//	results, err := BatchApplySingle(inputImages, watermarkImg, cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
func BatchApplySingle(
	inputImgs []image.Image,
	watermarkImg image.Image,
	config SingleConfig,
) ([]image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid single watermark configuration: %w", err)
	}

	maxWorkers := config.MaxWorkers
	if maxWorkers <= 0 {
		maxWorkers = runtime.NumCPU()
	}

	preparedWM := watermarkImg
	if config.OpacityAlpha < 1 {
		preparedWM = applyOpacity(preparedWM, config.OpacityAlpha)
	}
	if config.RotationDegrees != 0 {
		preparedWM = rotateImage(preparedWM, config.RotationDegrees)
	}

	numImages := len(inputImgs)
	results := make([]image.Image, numImages)

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for i := 0; i < numImages; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			currImage := inputImgs[index]
			currentWM := resizeWatermark(preparedWM, currImage, config.GeneralConfig)

			watermarkPosition := getWatermarkPosition(currentWM, currImage, config.VerticalAlign, config.HorizontalAlign, config.Spacing)

			canvas := generateBaseCanvas(currImage)
			drawWatermarkAtPosition(canvas, currentWM, watermarkPosition)

			results[index] = canvas
		}(i)
	}

	wg.Wait()

	return results, nil
}
