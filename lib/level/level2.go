package level

import (
	_ "fmt"
	_ "log"

	"sort"

	iiifcompliance "github.com/jdobber/go-iiif-mod/lib/compliance"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
)

type Level2 struct {
	Level      `json:"-"`
	Formats    []string                  `json:"formats"`
	Qualities  []string                  `json:"qualities"`
	Supports   []string                  `json:"supports"`
	compliance iiifcompliance.Compliance `json:"-"`
}

func NewLevel2(config *iiifconfig.Config, endpoint string) (*Level2, error) {

	compliance, err := iiifcompliance.NewLevel2Compliance(config)

	if err != nil {
		return nil, err
	}

	formats := compliance.Formats()
	sort.Strings(formats)

	qualities := compliance.Qualities()
	sort.Strings(qualities)

	supports := compliance.Supports()
	sort.Strings(supports)

	l := Level2{
		Formats:    formats,
		Qualities:  qualities,
		Supports:   supports,
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level2) Compliance() iiifcompliance.Compliance {
	return l.compliance
}
