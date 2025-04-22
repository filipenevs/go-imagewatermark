# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.1.1] - 2025-04-22

### Fixed
- Correctly apply spacing when using random alignment.

### Refactored
- Fallback to center alignment when an unknown alignment is provided.
- Extracted watermark alignment validation into a dedicated method for better code clarity.
- Simplified `loadImageFromFile` function and improved internal documentation.
- Skip rotation logic when rotation angle is zero to avoid unnecessary processing.

---

## [0.1.0] - 2025-04-21

### Added
- Initial release of `go-imagewatermark`.
- Function `ProcessImageWithWatermark` to overlay an image on another image.
- Full configuration options for watermark:
  - Opacity control (`OpacityAlpha`)
  - Dynamic width scaling (`WatermarkWidthPercent`)
  - Alignment settings (`HorizontalAlign`, `VerticalAlign`)
  - Edge spacing in pixels (`Spacing`)
  - Rotation (`RotationDegrees`)

### Notes
- This is an early release (`v0.x.x`) and the API may change before reaching `v1.0.0`.
- Use with caution in production environments. Feedback is welcome!
