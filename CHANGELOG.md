# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.0.0] - 2026-02-13

### Added
- Grid pattern support with customizable spacing and offsets.
- High-performance blending engine using native 'image/draw'.
- Custom byte-level alpha blending (faster than 'imaging').
- Concurrent image loading via Goroutines.

### Changed
- Default interpolator changed to Catmull-Rom for better speed/quality balance.
- Rendering pipeline now uses concrete types (*image.RGBA) for CPU optimization.

### Refactored
- Alignment constants migrated to 'iota' (int) for faster comparisons.
- Optimized positioning logic.
- Centralized image loading for improved reusability.
- Eliminated redundant image cloning to reduce memory overhead.

---

## [1.0.0] - 2025-04-22

### Added
- Initial stable release 🎉

### Changed
- Refactored image loading logic to use `imaging.Open`, ensuring automatic correction of image orientation based on EXIF metadata.

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
