package image

import (
	"bufio"
	"bytes"
	stdimage "image"
	"image/color"
	jpeg "image/jpeg"
	"io"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// NativeImage ...
type NativeImage struct {
	Image
	id     string
	bimg   stdimage.Image
	config stdimage.Config
	format string
}

type EncodingOptions struct {
	Quality int
	Format  string
}

// Height ...
func (im *NativeImage) Height() int {
	return im.config.Height
}

// Width ....
func (im *NativeImage) Width() int {
	return im.config.Width
}

// Identifier ...
func (im *NativeImage) Identifier() string {
	return im.id
}

// Identifier ...
func (im *NativeImage) Encode(opts *EncodingOptions) ([]byte, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	var err error

	switch opts.Format {
	case "jpg":
		err = jpeg.Encode(w, im.bimg,
			&jpeg.Options{
				Quality: opts.Quality,
			})
	case "webp":
		err = webp.Encode(&b, im.bimg,
			&webp.Options{
				Lossless: false,
				Quality:  float32(opts.Quality),
			})
	}

	return b.Bytes(), err
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
		id:     id,
		bimg:   bimg,
		config: config,
		format: format,
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

		im.bimg = imaging.Crop(im.bimg, stdimage.Rect(rgi.X, rgi.Y, rgi.X+rgi.Width, rgi.Y+rgi.Height))

	}

	if t.Size != "max" && t.Size != "full" {

		si, err := t.SizeInstructions(im)
		//fmt.Println("%v", si)

		if err != nil {
			return nil, err
		}

		im.bimg = imaging.Resize(im.bimg, si.Width, si.Height, imaging.Box)
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
		im.bimg = imaging.FlipH(im.bimg)
	}

	if ri.Angle != 0 {
		im.bimg = imaging.Rotate(im.bimg, -ri.Angle, color.White)
	}

	/*
		QUALITY
	*/

	switch t.Quality {
	case "grey":
		fallthrough
	case "gray":
		im.bimg = imaging.Grayscale(im.bimg)
	case "bitonal":
		im.bimg = imaging.Grayscale(im.bimg)
		im.bimg = imaging.AdjustContrast(im.bimg, 100)
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
