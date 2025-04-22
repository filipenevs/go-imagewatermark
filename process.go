package imagewatermark

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

type VerticalAlign string
type HorizontalAlign string

const (
	VerticalTop    VerticalAlign = "top"
	VerticalMiddle VerticalAlign = "mid"
	VerticalBottom VerticalAlign = "bottom"
	VerticalRandom VerticalAlign = "rand"
)

const (
	HorizontalLeft   HorizontalAlign = "left"
	HorizontalMiddle HorizontalAlign = "mid"
	HorizontalRight  HorizontalAlign = "right"
	HorizontalRandom HorizontalAlign = "rand"
)

// WatermarkConfig holds configuration for watermarking a image.
// It includes input paths and watermark appearance options.
type WatermarkConfig struct {
	InputPath             string
	WatermarkPath         string
	OpacityAlpha          float64
	WatermarkWidthPercent float64
	VerticalAlign         VerticalAlign
	HorizontalAlign       HorizontalAlign
	Spacing               int
	RotationDegrees       float64
}

// validate checks if the WatermarkConfig has valid values for all fields.
// It returns an error describing the first invalid value found,
// or nil if all fields are valid.
func (c WatermarkConfig) validate() error {
	if c.InputPath == "" || c.WatermarkPath == "" {
		return fmt.Errorf("input image and watermark paths must be provided")
	}

	if c.OpacityAlpha <= 0 || c.OpacityAlpha > 1 {
		return fmt.Errorf("opacity must be greater than 0 and less than or equal to 1: %f", c.OpacityAlpha)
	}

	if c.WatermarkWidthPercent <= 0 || c.WatermarkWidthPercent > 100 {
		return fmt.Errorf("watermark width percent must be greater than 0 and at most 100: %f", c.WatermarkWidthPercent)
	}

	if c.Spacing < 0 {
		return fmt.Errorf("spacing must be a non-negative integer: %d", c.Spacing)
	}

	if c.RotationDegrees < 0 || c.RotationDegrees >= 360 {
		return fmt.Errorf("rotation degrees must be between 0 and less than 360 (360 result in no rotation): %f", c.RotationDegrees)
	}

	return nil
}

// ProcessImageWithWatermark overlays a watermark onto an input image based on the provided configuration.
//
// Parameters:
//   - config: WatermarkConfig struct containing paths, opacity, size, alignment, and rotation settings.
//
// Returns:
//   - An image.Image containing the final image with the watermark applied.
//   - An error if any step fails, such as invalid configuration or failure to load images.
func ProcessImageWithWatermark(
	config WatermarkConfig,
) (image.Image, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid watermark configuration: %w", err)
	}

	inputImage, err := loadImageFromFile(config.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load input image: %v", err)
	}

	watermarkImage, err := loadImageFromFile(config.WatermarkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load watermark image: %v", err)
	}

	watermarkWidth := getNewWatermarkWidth(inputImage, config.WatermarkWidthPercent)
	watermarkImage = imaging.Resize(watermarkImage, watermarkWidth, 0, imaging.Lanczos)
	if config.RotationDegrees != 0 {
		watermarkImage = imaging.Rotate(watermarkImage, config.RotationDegrees, image.Transparent)
	}

	watermarkPosition := getWatermarkPosition(inputImage, watermarkImage, config.VerticalAlign, config.HorizontalAlign, config.Spacing)

	outputImage := imaging.Overlay(inputImage, watermarkImage, watermarkPosition, config.OpacityAlpha)

	return outputImage, nil

}
