package imagewatermark

import (
	"fmt"

	"github.com/disintegration/imaging"
)

// VerticalAlign defines the vertical alignment options for watermark placement.
//
// Supported values:
//   - VerticalTop: Aligns the watermark to the top of the image.
//   - VerticalMiddle: Centers the watermark vertically.
//   - VerticalBottom: Aligns the watermark to the bottom of the image.
//   - VerticalRandom: Places the watermark at a random vertical position.
type VerticalAlign int

// HorizontalAlign defines the horizontal alignment options for watermark placement.
//
// Supported values:
//   - HorizontalLeft: Aligns the watermark to the left of the image.
//   - HorizontalMiddle: Centers the watermark horizontally.
//   - HorizontalRight: Aligns the watermark to the right of the image.
//   - HorizontalRandom: Places the watermark at a random horizontal position.
type HorizontalAlign int

const (
	VerticalTop VerticalAlign = iota
	VerticalMiddle
	VerticalBottom
	VerticalRandom
)

const (
	HorizontalLeft HorizontalAlign = iota
	HorizontalMiddle
	HorizontalRight
	HorizontalRandom
)

// GeneralConfig holds common configuration settings for all watermarking operations.
//
// This struct contains the paths to the input image and watermark, along with general
// appearance settings that apply to any watermarking method (single or grid).
//
// Fields:
//   - OpacityAlpha: Transparency level of the watermark (0.0 to 1.0, where 1.0 is fully opaque).
//   - WatermarkWidthPercent: Desired watermark width as a percentage of the input image width (0-100).
//   - RotationDegrees: Rotation angle for the watermark in degrees (0-360).
//   - ResampleFilter: Resampling filter to use when resizing the watermark. (Default is CatmullRom)
//   - MaxWorkers: Maximum number of concurrent workers for batch processing (Default is number of CPU cores).
type GeneralConfig struct {
	OpacityAlpha          float64
	WatermarkWidthPercent float64
	RotationDegrees       float64
	ResampleFilter        imaging.ResampleFilter
	MaxWorkers            int
}

// validate checks if the GeneralConfig has valid values for all fields.
//
// It performs the following validations:
//   - OpacityAlpha must be greater than 0 and less than or equal to 1.
//   - WatermarkWidthPercent must be greater than 0 and at most 100.
//   - RotationDegrees must be between 0 and less than 360.
//   - MaxWorkers must be a non-negative integer.
//
// Returns:
//   - An error describing the first invalid value found, or nil if all fields are valid.
func (c GeneralConfig) validate() error {
	if c.OpacityAlpha <= 0 || c.OpacityAlpha > 1 {
		return fmt.Errorf("opacity must be greater than 0 and less than or equal to 1: %f", c.OpacityAlpha)
	}

	if c.WatermarkWidthPercent <= 0 || c.WatermarkWidthPercent > 100 {
		return fmt.Errorf("watermark width percent must be greater than 0 and at most 100: %f", c.WatermarkWidthPercent)
	}

	if c.RotationDegrees < 0 || c.RotationDegrees > 360 {
		return fmt.Errorf("rotation degrees must be between 0 and 360 (360 result in no rotation): %f", c.RotationDegrees)
	}

	if c.MaxWorkers < 0 {
		return fmt.Errorf("max workers must be a non-negative integer: %d", c.MaxWorkers)
	}

	return nil
}

// SingleConfig holds all configuration settings for applying a single watermark to an image.
//
// This struct extends GeneralConfig with alignment and spacing options specific to
// single watermark placement.
//
// Fields:
//   - GeneralConfig: Embedded struct containing common watermarking settings.
//   - VerticalAlign: Vertical alignment of the watermark (top, middle, bottom, or random).
//   - HorizontalAlign: Horizontal alignment of the watermark (left, middle, right, or random).
//   - Spacing: Distance in pixels between the watermark and the aligned edge.
type SingleConfig struct {
	GeneralConfig
	VerticalAlign   VerticalAlign
	HorizontalAlign HorizontalAlign
	Spacing         int
}

// validate checks if the SingleConfig has valid values for all fields.
//
// It first validates the embedded GeneralConfig, then checks:
//   - Spacing must be a non-negative integer.
//
// Returns:
//   - An error describing the first invalid value found, or nil if all fields are valid.
func (c SingleConfig) validate() error {
	if err := c.GeneralConfig.validate(); err != nil {
		return err
	}

	if c.Spacing < 0 {
		return fmt.Errorf("spacing must be a non-negative integer: %d", c.Spacing)
	}

	return nil
}

// GridConfig holds all configuration settings for applying a grid pattern of watermarks to an image.
//
// This struct extends GeneralConfig with grid-specific options for spacing and positioning.
// The grid pattern can be controlled using positive values for regular grids, or negative values
// for overlapping or inverted patterns.
//
// Fields:
//   - GeneralConfig: Embedded struct containing common watermarking settings.
//   - GridSpacingX: Horizontal spacing between watermarks in the grid (in pixels). Can be negative for overlapping.
//   - GridSpacingY: Vertical spacing between watermarks in the grid (in pixels). Can be negative for overlapping.
//   - OffsetX: Initial horizontal offset for the grid starting position (in pixels).
//   - OffsetY: Initial vertical offset for the grid starting position (in pixels).
type GridConfig struct {
	GeneralConfig
	GridSpacingX int
	GridSpacingY int
	OffsetX      int
	OffsetY      int
}

// validate checks if the GridConfig has valid values for all fields.
//
// It validates the embedded GeneralConfig. GridSpacingX, GridSpacingY, OffsetX, and OffsetY
// can be any integer value (including negative), so no specific validation is performed on them.
//
// Returns:
//   - An error if the GeneralConfig validation fails, or nil if all fields are valid.
func (c GridConfig) validate() error {
	if err := c.GeneralConfig.validate(); err != nil {
		return err
	}

	return nil
}
