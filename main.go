package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"
	"github.com/whosonfirst/go-sanitize"
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
	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	paramNames := [6]string{"identifier", "region", "size", "rotation", "quality", "format"}
	params := strings.Split(*uri, "/")

	vars := make(map[string]string)

	for idx, name := range paramNames {
		if name == "quality" {
			a := strings.Split(params[idx], ".")
			vars["quality"] = a[0]
			vars["format"] = a[1]
			break
		} else {
			vars[name] = params[idx]
		}

	}

	p := iiifparser.IIIFQueryParser{
		Opts: sanitize.DefaultOptions(),
		Vars: vars,
	}

	//fmt.Println("%v", p)

	endpoint := "Endpoint"

	identifier, _ := p.GetIIIFParameter("identifier")
	format, _ := p.GetIIIFParameter("format")
	//fmt.Println("%v", iiifparams)

	body, err := ioutil.ReadFile(*prefix + identifier)
	check(err)

	//image, err := iiifimage.NewVIPSImage(body)
	//check(err)

	image, err := iiifimage.NewNativeImage(body)
	check(err)

	level, err := iiiflevel.NewLevelFromConfig(config, endpoint)
	check(err)
	/*
		profile, err := iiifprofile.NewProfile(endpoint, image, level)
		check(err)

		payload, _ := json.Marshal(profile)
		fmt.Println(string(payload))
	*/

	iiifparams, err := p.GetIIIFParameters()
	check(err)

	transformation, err := iiifimage.NewTransformation(level,
		iiifparams.Region,
		iiifparams.Size,
		iiifparams.Rotation,
		iiifparams.Quality,
		iiifparams.Format)
	check(err)

	//uri, err := transformation.ToURI(params.Identifier)
	//body, err := derivatives_cache.Get(uri)

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
