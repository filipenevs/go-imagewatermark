
# 🧊 Go Image Watermarking Package

This Go package provides a simple and configurable way to overlay watermarks on images. It lets you control watermark opacity, size, alignment and rotation.

---

## 📄 Usage

```go
import "github.com/filipenevs/go-imagewatermark"
```

### Basic Example

```go
result, err := imagewatermark.ProcessImageWithWatermark(imagewatermark.WatermarkConfig{
	InputPath:             "input.png",
	WatermarkPath:         "logo.png",
	OpacityAlpha:          0.6,
	WatermarkWidthPercent: 20,
	VerticalAlign:         imagewatermark.VerticalBottom,
	HorizontalAlign:       imagewatermark.HorizontalRight,
	Spacing:               20,
	RotationDegrees:       45,
})

if err != nil {
	log.Fatalf("Failed to add watermark: %v", err)
}
```

---

## 📚 API Reference

### Types

#### `WatermarkConfig`

| Field                  | Type             | Description                                                                 |
|------------------------|------------------|-----------------------------------------------------------------------------|
| `InputPath`            | `string`         | Path to the input image file                                                |
| `WatermarkPath`        | `string`         | Path to the watermark image (supports transparency)                         |
| `OpacityAlpha`         | `float64`        | Opacity level (e.g., `0.5` = 50% transparency)                              |
| `WatermarkWidthPercent`| `float64`        | Width of the watermark as % of image width                                  |
| `VerticalAlign`        | `VerticalAlign`  | Vertical alignment: `top`, `mid`, `bottom`, `rand`                          |
| `HorizontalAlign`      | `HorizontalAlign`| Horizontal alignment: `left`, `mid`, `right`, `rand`                        |
| `Spacing`              | `int`            | Padding in pixels from aligned edge                                         |
| `RotationDegrees`      | `float64`        | Rotation angle in degrees                                                   |

---

## 📜 License

MIT License. See `LICENSE` file for details.

---

## 🤝 Contributing

Pull requests and issues are welcome!