package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"
	"github.com/whosonfirst/go-sanitize"
)

var (
	totalTime time.Duration
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	totalTime += elapsed
	log.Printf("%s took %s [total: %s]", name, elapsed, totalTime)
}

func main() {
	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")
	var uri = flag.String("uri", "", "a vaild iiif uri, e.g. colorize.jpg/full/full/90/default.webp")
	var prefix = flag.String("prefix", "", "a path to an image folder")
	var showProfile = flag.Bool("profile", false, "print profile")
	var err error

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	endpoint := "http://localhost"

	var config *iiifconfig.Config
	config, err = func(s string) (*iiifconfig.Config, error) {
		defer timeTrack(time.Now(), "NewConfigFromFlag")
		return iiifconfig.NewConfigFromFlag(s)
	}(*cfg)

	var p *iiifparser.IIIFQueryParser
	p, err = func(uri string, opts *sanitize.Options) (*iiifparser.IIIFQueryParser, error) {
		defer timeTrack(time.Now(), "NewIIIFQueryParser")
		return iiifparser.NewIIIFQueryParser(uri, opts)
	}(*uri, nil)
	check(err)

	identifier, _ := p.GetIIIFParameter("identifier")
	format, _ := p.GetIIIFParameter("format")
	//fmt.Println("%v", iiifparams)

	var body []byte
	body, err = func(filename string) ([]byte, error) {
		defer timeTrack(time.Now(), "ReadFile")
		return ioutil.ReadFile(filename)
	}(*prefix + identifier)
	check(err)

	var image *iiifimage.NativeImage
	image, err = func(id string, body []byte) (*iiifimage.NativeImage, error) {
		defer timeTrack(time.Now(), "NewNativeImage")
		return iiifimage.NewNativeImage(id, body)
	}(identifier, body)
	check(err)

	var level iiiflevel.Level
	level, err = func(config *iiifconfig.Config, endpoint string) (iiiflevel.Level, error) {
		defer timeTrack(time.Now(), "NewLevelFromConfig")
		return iiiflevel.NewLevelFromConfig(config, endpoint)
	}(config, endpoint)
	check(err)

	if *showProfile {
		profile, err := iiifprofile.NewProfile(endpoint, image, level)
		check(err)

		payload, _ := json.Marshal(profile)
		fmt.Println(string(payload))
	}

	iiifparams, err := p.GetIIIFParameters()
	check(err)

	transformation, err := iiifimage.NewTransformation(level,
		iiifparams.Region,
		iiifparams.Size,
		iiifparams.Rotation,
		iiifparams.Quality,
		iiifparams.Format)
	check(err)

	if transformation.HasTransformation() {

		_, err = func(t *iiifimage.Transformation) (*iiifimage.NativeImage, error) {
			defer timeTrack(time.Now(), "Transform")
			return image.Transform(t)
		}(transformation)
		check(err)

	}

	opts := iiifimage.EncodingOptions{
		Format:  format,
		Quality: 70,
	}

	var data []byte
	data, err = func(o *iiifimage.EncodingOptions) ([]byte, error) {
		defer timeTrack(time.Now(), "Encode")
		return image.Encode(o)
	}(&opts)
	check(err)

	err = ioutil.WriteFile("./test."+opts.Format, data, 0644)
	log.Printf("Wrote %d bytes to %s", len(data), "./test."+opts.Format)
	check(err)	

}
