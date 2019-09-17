package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")
	var uri = flag.String("uri", "", "a vaild iiif uri, e.g. colorize.jpg/full/full/90/default.webp")
	var prefix = flag.String("prefix", "", "a path to an image folder")
	var showProfile = flag.Bool("profile", false, "print profile")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	endpoint := "http://localhost"

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	p, err := iiifparser.NewIIIFQueryParser(*uri, nil)
	check(err)

	identifier, _ := p.GetIIIFParameter("identifier")
	format, _ := p.GetIIIFParameter("format")
	//fmt.Println("%v", iiifparams)

	body, err := ioutil.ReadFile(*prefix + identifier)
	check(err)

	image, err := iiifimage.NewNativeImage(identifier, body)
	check(err)

	level, err := iiiflevel.NewLevelFromConfig(config, endpoint)
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

		_, err := image.Transform(transformation)
		check(err)

	}

	opts := iiifimage.EncodingOptions{
		Format:  format,
		Quality: 70,
	}

	data, err := image.Encode(&opts)
	check(err)

	err = ioutil.WriteFile("./test."+opts.Format, data, 0644)
	check(err)

}
