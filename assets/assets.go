package assets

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func GetImage(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file %s: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", path, err)
	}

	fmt.Println("Loaded image:", path)
	return ebiten.NewImageFromImage(img), nil
}
