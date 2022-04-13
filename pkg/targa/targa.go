package targa

import (
	// "fmt"
	"errors"
	"fmt"
	"image/color"
	"os"
)

type ImageType byte

const (
	Empty ImageType = iota
	UncompressedColorMapped
	UncompressedRGB
	UncompressedGrayscale

	RunLengthColorMapped    = 9
	RunLengthColorRGB       = 10
	RunLengthColorGrayscale = 11
)

type Header struct {
	Type       ImageType
	Width      int
	Height     int
	PixelSize  byte
	Descriptor byte
}

func (h Header) GetBytes() []byte {
	// TODO: account for other options, dont use magic number 18
	header := make([]byte, 18)
	header[2] = byte(h.Type)

	header[12] = byte(h.Width & 0x00FF)         // Width lo
	header[13] = byte((h.Width & 0xFF00) / 256) // Width hi

	header[14] = byte(h.Height & 0x00FF)         // Height lo
	header[15] = byte((h.Height & 0xFF00) / 256) // Height  hi

	header[16] = h.PixelSize
	header[17] = h.Descriptor
	return header
}

func (h Header) ByteSizePerPixel() int {
	pixelByteSize := 0 // Empty
	// TODO: account for other enums
	switch h.Type {
	case UncompressedGrayscale, RunLengthColorGrayscale:
		pixelByteSize = 2
	case UncompressedRGB, RunLengthColorRGB:
		pixelByteSize = 3

	}
	return pixelByteSize
}

type TGAImage struct {
	Header
	Pixels []byte
}

func NewImage(width, height int, imgType ImageType, depth byte) *TGAImage {
	// TODO: account for other descriptor
	header := Header{imgType, width, height, depth, 0}
	pixels := make([]byte, header.ByteSizePerPixel()*width*height)
	return &TGAImage{
		header,
		pixels,
	}
}

func (t TGAImage) FillRGB(col color.RGBA) {
	// TODO: account for grayscale and mapped color
	counter := 0
	for i := 0; i < t.Width; i++ {
		for j := 0; j < t.Height; j++ {
			t.Pixels[counter] = col.B
			counter++

			t.Pixels[counter] = col.G
			counter++

			t.Pixels[counter] = col.R
			counter++
		}
	}
}

func (t TGAImage) SetPixelRGB(x, y int, col color.RGBA) error {
	if x >= t.Width || x < 0 || y >= t.Height || y < 0 {
		return errors.New(fmt.Sprintf("Cell (%d, %d) is out of bound", x, y))
	}

	cell := (y*t.Width + x) * 3
	t.Pixels[cell] = col.B
	t.Pixels[cell+1] = col.G
	t.Pixels[cell+2] = col.R
	return nil
}

func (t TGAImage) GetBytes() []byte {
	return append(t.Header.GetBytes(), t.Pixels...)
}

func (t TGAImage) WriteToFile(name string) error {
	return os.WriteFile(name, t.GetBytes(), os.FileMode(0777))
}
