
# 🧊 Go Image Watermarking Package

This Go package provides a simple and configurable way to overlay watermarks on images. It lets you control watermark opacity, size, alignment and rotation.

Supports both single watermark placement and grid patterns with full customization options.

---

## ⬇️ Installation

```bash
go get github.com/filipenevs/go-imagewatermark
```

## 📄 Quick Start

```go
import "github.com/filipenevs/go-imagewatermark"
```

### Single Watermark Example


```go
package main

import (
    "log"
    "github.com/filipenevs/go-imagewatermark"
)

func main() {
    config := imagewatermark.SingleConfig{
        GeneralConfig: imagewatermark.GeneralConfig{
            InputPath:             "input.jpg",
            WatermarkPath:         "logo.png",
            WatermarkWidthPercent: 20,
            OpacityAlpha:          0.5,
            RotationDegrees:       0,
        },
        VerticalAlign:   imagewatermark.VerticalBottom,
        HorizontalAlign: imagewatermark.HorizontalRight,
        Spacing:         30,
    }

    result, err := imagewatermark.ApplySingle(config)
    if err != nil {
        log.Fatal(err)
    }

		err = SaveImageToFile(result, "png", "output_image.png")
		if err != nil {
				log.Fatal(err)
		}
}
```
<img width="600" height="400" alt="output_image_example_1" src="https://github.com/user-attachments/assets/78cc0e97-8ca4-4e4a-b2cf-e2acb911ef76" />

### Grid Watermark Example

```go
package main

import (
    "log"
    "github.com/filipenevs/go-imagewatermark"
)

func main() {
    config := imagewatermark.GridConfig{
        GeneralConfig: imagewatermark.GeneralConfig{
            InputPath:             "input.jpg",
            WatermarkPath:         "logo.png",
            WatermarkWidthPercent: 10,
            OpacityAlpha:          0.3,
            RotationDegrees:       45,
        },
        GridSpacingX: 40,
        GridSpacingY: 40,
        OffsetX:      -15,
        OffsetY:      -15,
    }

    result, err := imagewatermark.ApplyGrid(config)
    if err != nil {
        log.Fatal(err)
    }

    err = SaveImageToFile(result, "png", "output_image.png")
		if err != nil {
				log.Fatal(err)
		}
}
```
<img width="600" height="400" alt="output_image_example_2" src="https://github.com/user-attachments/assets/7f06ce27-0335-4205-a1f9-8f897159b12f" />

---

## ⚙️ Config

### General Configuration

The `GeneralConfig` struct contains common settings for all watermarking operations:

| Field | Type | Description | Range |
|-------|------|-------------|-------|
| `InputPath` | string | Path to the input image file | Any valid file path |
| `WatermarkPath` | string | Path to the watermark image file | Any valid file path |
| `OpacityAlpha` | float64 | Transparency level of the watermark | (0.0 - 1.0] |
| `WatermarkWidthPercent` | float64 | Watermark width as percentage of input image | (0 - 100] |
| `RotationDegrees` | float64 | Rotation angle for the watermark | [0 - 360] |

### Single Watermark Configuration

The `SingleConfig` extends `GeneralConfig` with alignment options:

| Field | Type | Description | Options |
|-------|------|-------------|---------|
| `VerticalAlign` | VerticalAlign | Vertical alignment position | `VerticalTop`, `VerticalMiddle`, `VerticalBottom`, `VerticalRandom` |
| `HorizontalAlign` | HorizontalAlign | Horizontal alignment position | `HorizontalLeft`, `HorizontalMiddle`, `HorizontalRight`, `HorizontalRandom` |
| `Spacing` | int | Distance from aligned edge (pixels) | Any non-negative integer |

**Example:**

```go
config := imagewatermark.SingleConfig{
    GeneralConfig: imagewatermark.GeneralConfig{
        InputPath:             "input.jpg",
        WatermarkPath:         "logo.png",
        WatermarkWidthPercent: 20,
        OpacityAlpha:          0.7,
        RotationDegrees:       0,
    },
    VerticalAlign:   imagewatermark.VerticalBottom,
    HorizontalAlign: imagewatermark.HorizontalRight,
    Spacing:         15,
}

result, err := imagewatermark.ApplySingle(config)
```

### Grid Watermark Configuration

The `GridConfig` extends `GeneralConfig` with grid-specific options:

| Field | Type | Description | Notes |
|-------|------|-------------|-------|
| `GridSpacingX` | int | Horizontal spacing between watermarks (pixels) | Can be negative for overlapping |
| `GridSpacingY` | int | Vertical spacing between watermarks (pixels) | Can be negative for overlapping |
| `OffsetX` | int | Initial horizontal offset for grid start (pixels) | Can be negative to shift grid left |
| `OffsetY` | int | Initial vertical offset for grid start (pixels) | Can be negative to shift grid up |

**Example:**

```go
config := imagewatermark.GridConfig{
    GeneralConfig: imagewatermark.GeneralConfig{
        InputPath:             "input.jpg",
        WatermarkPath:         "logo.png",
        WatermarkWidthPercent: 10,
        OpacityAlpha:          0.3,
        RotationDegrees:       0,
    },
    GridSpacingX: 200,
    GridSpacingY: 200,
    OffsetX:      -50,
    OffsetY:      -50,
}

result, err := imagewatermark.ApplyGrid(config)
```

## 🧠 Advanced Examples

### Random Placement

```go
config := imagewatermark.SingleConfig{
    GeneralConfig: imagewatermark.GeneralConfig{
        InputPath:             "input.jpg",
        WatermarkPath:         "logo.png",
        WatermarkWidthPercent: 15,
        OpacityAlpha:          0.5,
        RotationDegrees:       0,
    },
    VerticalAlign:   imagewatermark.VerticalRandom,   // Random vertically
    HorizontalAlign: imagewatermark.HorizontalRandom, // Random horizontally
}

result, err := imagewatermark.ApplySingle(config)
```

### Rotated Grid Pattern

```go
config := imagewatermark.GridConfig{
    GeneralConfig: imagewatermark.GeneralConfig{
        InputPath:             "input.jpg",
        WatermarkPath:         "logo.png",
        WatermarkWidthPercent: 8,
        OpacityAlpha:          0.25,
        RotationDegrees:       45, // Rotate watermarks 45 degrees
    },
    GridSpacingX: 250,
    GridSpacingY: 250,
    OffsetX:      -75,
    OffsetY:      -75,
}

result, err := imagewatermark.ApplyGrid(config)
```

## Error Handling

The library provides detailed error messages for common issues:

```go
result, err := imagewatermark.ApplySingle(config)
if err != nil {
    log.Printf("Error: %v", err)
    // Possible errors:
    // - "invalid watermark configuration: ..."
    // - "failed to load input image: ..."
    // - "failed to load watermark image: ..."
}
```

---

## 📜 License

MIT License. See `LICENSE` file for details.

---

## 🤝 Contributing

Pull requests and issues are welcome!
