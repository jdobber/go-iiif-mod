package lib

import (
	"errors"
	"fmt"
	gohttp "net/http"
	gourl "net/url"
	"strings"

	"github.com/gorilla/mux"
	echo "github.com/labstack/echo/v4"
	"github.com/whosonfirst/go-sanitize"
)

// IIIFParameters ...
type IIIFParameters struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

// IIIFQueryParser ...
type IIIFQueryParser struct {
	Opts *sanitize.Options
	Vars map[string]string
}

func NewIIIFQueryParser(r *gohttp.Request) (*IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := mux.Vars(r)

	p := IIIFQueryParser{
		Opts: opts,
		Vars: vars,
	}

	return &p, nil
}

func NewIIIFQueryParser2(c echo.Context) (*IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := make(map[string]string)

	for _, name := range c.ParamNames() {
		if name == "quality" {
			a := strings.Split(c.Param(name), ".")
			vars["quality"] = a[0]
			vars["format"] = a[1]
		} else {
			vars[name] = c.Param(name)
		}

	}

	p := IIIFQueryParser{
		Opts: opts,
		Vars: vars,
	}

	return &p, nil
}

func (p *IIIFQueryParser) GetIIIFParameter(key string) (string, error) {

	var err error

	value := p.Vars[key]

	value, err = sanitize.SanitizeString(value, p.Opts)

	if err != nil {
		return "", err
	}

	value, err = gourl.QueryUnescape(value)

	if err != nil {
		return "", err
	}

	// This should be already be stripped out by the time we get here but just
	// in case... (20160926/thisisaaronland)

	if strings.Contains(value, "../") {
		msg := fmt.Sprintf("Invalid key %s", key)
		err := errors.New(msg)
		return "", err
	}

	return value, nil
}

func (p *IIIFQueryParser) GetIIIFParameters() (*IIIFParameters, error) {

	id, err := p.GetIIIFParameter("identifier")

	if err != nil {
		return nil, err
	}

	region, err := p.GetIIIFParameter("region")

	if err != nil {
		return nil, err
	}

	size, err := p.GetIIIFParameter("size")

	if err != nil {
		return nil, err
	}

	rotation, err := p.GetIIIFParameter("rotation")

	if err != nil {
		return nil, err
	}

	quality, err := p.GetIIIFParameter("quality")

	if err != nil {
		return nil, err
	}

	format, err := p.GetIIIFParameter("format")

	if err != nil {
		return nil, err
	}

	params := IIIFParameters{
		Identifier: id,
		Region:     region,
		Size:       size,
		Rotation:   rotation,
		Quality:    quality,
		Format:     format,
	}

	return &params, nil
}
