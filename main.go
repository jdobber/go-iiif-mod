package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"
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

	identifier, _ := p.GetIIIFParameter("identifier")
	//fmt.Println("%v", iiifparams)

	body, err := ioutil.ReadFile("/home/jens/Bilder/" + identifier)
	check(err)

	image, err := iiifimage.NewVIPSImage(body)
	check(err)

	level, err := iiiflevel.NewLevelFromConfig(config, "endpoint")
	check(err)

	profile, err := iiifprofile.NewProfile("endpoint", image, level)
	check(err)

	payload, _ := json.Marshal(profile)
	fmt.Println(string(payload))
}
