package image

import (
	"bufio"
	"bytes"
	stdimage "image"
	"image/color"
	jpeg "image/jpeg"
	"io"
	"fmt"
	"errors"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// NativeImage ...
type NativeImage struct {
	Image
	ID     string
	Bimg   stdimage.Image
	Config stdimage.Config
	Format string
}

type EncodingOptions struct {
	Quality int
	Format  string
}

// Height ...
func (im *NativeImage) Height() int {
	return im.Config.Height
}

// Width ....
func (im *NativeImage) Width() int {
	return im.Config.Width
}

// Identifier ...
func (im *NativeImage) Identifier() string {
	return im.ID
}

// Identifier ...
func (im *NativeImage) Encode(opts *EncodingOptions) ([]byte, error) {
	var buf bytes.Buffer	
	defer buf.Reset()

	var err error

	switch opts.Format {
	case "jpg":
		w := bufio.NewWriter(&buf)
		err = jpeg.Encode(w, im.Bimg,
			&jpeg.Options{
				Quality: opts.Quality,
			})
	case "webp":
		err = webp.Encode(&buf, im.Bimg,
			&webp.Options{
				Lossless: false,
				Quality:  float32(opts.Quality),
			})
	default:
		return nil, errors.New(fmt.Sprintf("%s is unknown format for encoding", opts.Format))		
	}

	return buf.Bytes(), err
}

// NewNativeImage ...
func NewNativeImage(id string, body []byte) (*NativeImage, error) {

	reader := bytes.NewReader(body)
	bimg, format, err := stdimage.Decode(reader)

	reader.Seek(0, io.SeekStart)
	config, _, err := stdimage.DecodeConfig(reader)

	if err != nil {
		return nil, err
	}

	im := NativeImage{
		ID:     id,
		Bimg:   bimg,
		Config: config,
		Format: format,
	}

	return &im, nil
}

func (im *NativeImage) Transform(t *Transformation) (*NativeImage, error) {

	if t.Region != "full" {

		rgi, err := t.RegionInstructions(im)
		//fmt.Println("%v", rgi)

		if err != nil {
			return nil, err
		}

		im.Bimg = imaging.Crop(im.Bimg, stdimage.Rect(rgi.X, rgi.Y, rgi.X+rgi.Width, rgi.Y+rgi.Height))

	}

	if t.Size != "max" && t.Size != "full" {

		si, err := t.SizeInstructions(im)
		//fmt.Println("%v", si)

		if err != nil {
			return nil, err
		}

		im.Bimg = imaging.Resize(im.Bimg, si.Width, si.Height, imaging.Box)
	}

	/*
		ROTATION
	*/

	ri, err := t.RotationInstructions(im)
	//fmt.Println("%v", ri)

	if err != nil {
		return nil, err
	}

	if ri.Flip {
		im.Bimg = imaging.FlipH(im.Bimg)
	}

	if ri.Angle != 0 {
		im.Bimg = imaging.Rotate(im.Bimg, -ri.Angle, color.White)
	}

	/*
		QUALITY
	*/

	switch t.Quality {
	case "grey":
		fallthrough
	case "gray":
		im.Bimg = imaging.Grayscale(im.Bimg)
	case "bitonal":
		im.Bimg = imaging.Grayscale(im.Bimg)
		im.Bimg = imaging.AdjustContrast(im.Bimg, 100)
	}

	/*
		FORMAT
	*/

	_, err = t.FormatInstructions(im)
	//fmt.Println("%v", fi)

	if err != nil {
		return nil, err
	}

	return im, nil
}
