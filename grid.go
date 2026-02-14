package imagewatermark

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
)

// GridPosition represents a single position in the grid where a watermark will be applied.
//
// Fields:
//   - X: The horizontal coordinate (in pixels) where the watermark should be placed.
//   - Y: The vertical coordinate (in pixels) where the watermark should be placed.
type GridPosition struct {
	X, Y int
}

// ApplyGrid applies a grid pattern of watermarks to an input image based on the provided configuration.
//
// The function performs the following steps:
//  1. Validates the GridConfig.
//  2. Loads the input image and watermark image in parallel using goroutines.
//  3. Preprocesses the watermark (resize, rotate, and apply opacity).
//  4. Generates grid positions based on spacing and offset settings.
//  5. Applies the watermark at each grid position.
//  6. Returns the final image with the grid watermarks applied.
//
// The parallel image loading uses goroutines to improve performance for I/O operations.
// If the offset or spacing values are negative, the grid pattern will be shifted or compressed accordingly.
//
// Parameters:
//   - config: GridConfig struct containing paths, spacing, offset, and watermark appearance settings.
//
// Returns:
//   - An image.Image containing the final image with the grid watermark pattern applied.
//   - An error if configuration is invalid or image loading fails.
//
// Example:
//
//	config := GridConfig{
//		GeneralConfig: GeneralConfig{
//			InputPath:             "input.jpg",
//			WatermarkPath:         "logo.png",
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
func ApplyGrid(config GridConfig) (image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid grid watermark configuration: %w", err)
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

	positions := generateGridPositions(inputImg, processedWatermark, config)

	resultImg := applyGridWatermarks(inputImg, processedWatermark, positions)

	return resultImg, nil
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
func generateGridPositions(inputImg, watermarkImg image.Image, config GridConfig) []GridPosition {
	inputBounds := inputImg.Bounds()
	watermarkBounds := watermarkImg.Bounds()

	var positions []GridPosition

	for y := config.OffsetY; y < inputBounds.Dy(); y += config.GridSpacingY + watermarkBounds.Dy() {
		for x := config.OffsetX; x < inputBounds.Dx(); x += config.GridSpacingX + watermarkBounds.Dx() {
			positions = append(positions, GridPosition{X: x, Y: y})
		}
	}

	return positions
}

// applyGridWatermarks applies watermarks at each position in the provided grid.
//
// This function creates a copy of the input image and then overlays the watermark image
// at each GridPosition. The watermarks are composited using the "Over" operator, which means
// that watermarks in front will cover watermarks behind them if they overlap.
//
// The function does not use goroutines as the bottleneck is typically the drawing operations
// rather than I/O. Sequential application ensures consistent ordering and simpler logic.
//
// Parameters:
//   - inputImg: The input image to which watermarks will be applied.
//   - watermarkImg: The preprocessed watermark image to be applied.
//   - positions: A slice of GridPosition objects indicating where to place each watermark.
//
// Returns:
//   - An image.Image containing the input image with watermarks applied at all grid positions.
func applyGridWatermarks(inputImg, watermarkImg image.Image, positions []GridPosition) image.Image {

	bounds := inputImg.Bounds()
	canvas := image.NewRGBA(bounds)

	draw.Draw(canvas, bounds, inputImg, bounds.Min, draw.Src)

	wmBounds := watermarkImg.Bounds()

	for _, pos := range positions {
		dr := image.Rectangle{
			Min: image.Point{X: pos.X, Y: pos.Y},
			Max: image.Point{X: pos.X + wmBounds.Dx(), Y: pos.Y + wmBounds.Dy()},
		}
		draw.Draw(canvas, dr, watermarkImg, image.Point{0, 0}, draw.Over)
	}

	return canvas
}
