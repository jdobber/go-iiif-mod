package level

/*

Things I am not sure about include the relationship of level/*.go and compliance/*.go which are
very much related but somehow seem like they should be in separate namespaces. I'm not sure...
(20160912/thisisaaronland)

*/

import (
	"errors"
	"fmt"
	_ "log"

	iiifcompliance "github.com/jdobber/go-iiif-mod/lib/compliance"
	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
)

type Level interface {
	Compliance() iiifcompliance.Compliance
}

func NewLevelFromConfig(config *iiifconfig.Config, endpoint string) (Level, error) {

	compliance := "2" //config.Level.Compliance

	if compliance == "0" {

		message := fmt.Sprintf("Unsupported compliance level '%s'", compliance)
		return nil, errors.New(message)

	} else if compliance == "1" {

		message := fmt.Sprintf("Unsupported compliance level '%s'", compliance)
		return nil, errors.New(message)
	} else if compliance == "2" {

		return NewLevel2(config, endpoint)

	} else {

		message := fmt.Sprintf("Invalid compliance level '%s'", compliance)
		return nil, errors.New(message)

	}
}
