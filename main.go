package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	iiifconfig "github.com/jdobber/go-iiif-mod/config"
	lib "github.com/jdobber/go-iiif-mod/lib"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")
	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	uri := "colorize.jpg/full/full/90/default.webp"
	paramNames := [6]string{"identifier", "region", "size", "rotation", "quality", "format"}
	params := strings.Split(uri, "/")

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
	/*
		p := lib.IIIFQueryParser{
			Opts: sanitize.DefaultOptions(),
			Vars: vars,
		}

		//fmt.Println("%v", p)

		iiifparams, _ := p.GetIIIFParameters()
		//fmt.Println("%v", iiifparams)
	*/
	body, err := ioutil.ReadFile("/home/jens/Bilder/colorize.jpg")
	check(err)

	image, err := lib.NewVIPSImage(body)
	check(err)

	level, err := lib.NewLevelFromConfig(config, "endpoint")
	check(err)

	profile, err := lib.NewProfile("endpoint", image, level)
	check(err)

	payload, _ := json.Marshal(profile)
	fmt.Println(string(payload))
}
