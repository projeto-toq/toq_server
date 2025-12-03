package imageprocessing

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// readOrientation extracts the EXIF orientation flag from raw bytes.
func readOrientation(data []byte) int {
	if len(data) == 0 {
		return 1
	}
	reader := bytes.NewReader(data)
	tags, err := exif.Decode(reader)
	if err != nil {
		return 1
	}
	orientationTag, err := tags.Get(exif.Orientation)
	if err != nil {
		return 1
	}
	value, err := orientationTag.Int(0)
	if err != nil {
		return 1
	}
	return value
}

// normalizeOrientation rotates or flips the image according to EXIF orientation spec.
func normalizeOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img) // 90° clockwise
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img) // 90° counter-clockwise
	default:
		return img
	}
}
