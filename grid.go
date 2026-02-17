package imagewatermark

import (
	"fmt"
	"image"
	"runtime"
	"sync"
)

// ApplyGrid applies a grid pattern of watermarks to an input image based on the provided configuration.
//
// The function performs the following steps:
//  1. Validates the GridConfig.
//  2. Preprocesses the watermark (resize, rotate, and apply opacity).
//  3. Generates grid positions based on spacing and offset settings.
//  4. Applies the watermark at each grid position.
//  5. Returns the final image with the grid watermarks applied.
//
// The parallel image loading uses goroutines to improve performance for I/O operations.
// If the offset or spacing values are negative, the grid pattern will be shifted or compressed accordingly.
//
// Parameters:
//   - inputImg: The input image to which the grid watermark will be applied.
//   - watermarkImg: The watermark image to overlay on the input image.
//   - config: GridConfig struct containing spacing, offset, and watermark appearance settings.
//
// Returns:
//   - An image.Image containing the final image with the grid watermark pattern applied.
//   - An error if configuration is invalid or image loading fails.
//
// Example:
//
//	config := GridConfig{
//		GeneralConfig: GeneralConfig{
//			WatermarkWidthPercent: 20,
//			OpacityAlpha:          0.5,
//			RotationDegrees:       45,
//		},
//		GridSpacingX: 75,
//		GridSpacingY: 75,
//		OffsetX:      -50,
//		OffsetY:      -50,
//	}
//	result, err := ApplyGrid(config)
//	if err != nil {
//		log.Fatal(err)
//	}
func ApplyGrid(
	inputImg image.Image,
	watermarkImg image.Image,
	config GridConfig,
) (image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid grid watermark configuration: %w", err)
	}

	preparedWM := watermarkImg

	if config.OpacityAlpha < 1 {
		preparedWM = applyOpacity(preparedWM, config.OpacityAlpha)
	}

	if config.RotationDegrees != 0 {
		preparedWM = rotateImage(preparedWM, config.RotationDegrees)
	}

	preparedWM = resizeWatermark(preparedWM, inputImg, config.GeneralConfig)
	positions := generateGridPositions(inputImg, preparedWM, config)

	result := applyGridWatermarks(inputImg, preparedWM, positions)

	return result, nil
}

func BatchApplyGrid(
	inputImgs []image.Image,
	watermarkImg image.Image,
	config GridConfig,
) ([]image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid grid watermark configuration: %w", err)
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
			positions := generateGridPositions(currImage, currentWM, config)

			canvas := generateBaseCanvas(currImage)

			result := applyGridWatermarks(canvas, currentWM, positions)

			results[index] = result
		}(i)
	}

	wg.Wait()

	return results, nil
}

// generateGridPositions calculates all positions where watermarks should be placed in a grid pattern.
//
// This function creates a grid of GridPosition objects by iterating through the input image dimensions
// using the GridSpacingX and GridSpacingY values. The OffsetX and OffsetY values determine the starting
// position of the grid.
//
// The grid continues until it exceeds the input image boundaries. If spacing or offset values are negative,
// the grid pattern will be shifted or compressed accordingly, allowing for overlapping or inverted patterns.
//
// Parameters:
//   - inputImg: The input image to determine grid boundaries.
//   - watermarkImg: The watermark image to get its dimensions.
//   - config: GridConfig containing spacing and offset settings.
//
// Returns:
//   - A slice of GridPosition objects representing all positions where watermarks should be placed.
func generateGridPositions(inputImg, watermarkImg image.Image, config GridConfig) []image.Point {
	inputW, inputH := inputImg.Bounds().Dx(), inputImg.Bounds().Dy()
	wmW, wmH := watermarkImg.Bounds().Dx(), watermarkImg.Bounds().Dy()

	if config.OffsetX >= inputW || config.OffsetY >= inputH {
		return nil
	}

	stepX := wmW + config.GridSpacingX
	stepY := wmH + config.GridSpacingY

	countX := (inputW - config.OffsetX + stepX - 1) / stepX
	countY := (inputH - config.OffsetY + stepY - 1) / stepY

	totalPoints := countX * countY
	if totalPoints <= 0 {
		return nil
	}

	positions := make([]image.Point, 0, totalPoints)

	for y := config.OffsetY; y < inputH; y += stepY {
		for x := config.OffsetX; x < inputW; x += stepX {
			positions = append(positions, image.Point{X: x, Y: y})
		}
	}

	return positions
}

// applyGridWatermarks applies watermarks at each position in the provided grid.
//
// This function creates a copy of the input image and then overlays the watermark image
// at each image.Point. The watermarks are composited using the "Over" operator, which means
// that watermarks in front will cover watermarks behind them if they overlap.
//
// The function does not use goroutines as the bottleneck is typically the drawing operations
// rather than I/O. Sequential application ensures consistent ordering and simpler logic.
//
// Parameters:
//   - inputImg: The input image to which watermarks will be applied.
//   - watermarkImg: The preprocessed watermark image to be applied.
//   - positions: A slice of image.Point objects indicating where to place each watermark.
//
// Returns:
//   - An image.Image containing the input image with watermarks applied at all grid positions.
func applyGridWatermarks(inputImg, watermarkImg image.Image, positions []image.Point) image.Image {
	canvas := generateBaseCanvas(inputImg)

	for _, pos := range positions {
		drawWatermarkAtPosition(canvas, watermarkImg, pos)
	}

	return canvas
}
