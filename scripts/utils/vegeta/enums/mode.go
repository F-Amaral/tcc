package enums

import (
	"fmt"
	"strings"
)

type Mode string
type ModeTemplate map[Mode][]string

const (
	PPT       Mode = "ppt"
	Nested    Mode = "nested"
	Recursive Mode = "recursive"
	All       Mode = "all"

	recursiveTemplate string = "%s/ppt/%s"
	pptTemplate       string = "%s/ppt/%s=?recursive=false"
	nestedTemplate    string = "%s/nested/%s"
)

func NameOf(value string) (*Mode, error) {
	for _, format := range modes() {
		if strings.EqualFold(value, format.String()) {
			return &format, nil
		}
	}
	return nil, fmt.Errorf("Invalid mode: %s", value)

}

func (m Mode) String() string {
	if m.Is(Recursive) {
		return fmt.Sprintf("%s-%s", string(PPT), string(Recursive))
	}
	return string(m)
}

func (m Mode) Is(mode Mode) bool {
	return m == mode
}

func (m Mode) Template() string {
	return modeTemplates()[m]
}

func modes() []Mode {
	return []Mode{PPT, Nested, Recursive, All}
}

func (m Mode) Expand() []Mode {
	return modesMap()[m]
}

func modesMap() map[Mode][]Mode {
	return map[Mode][]Mode{
		"ppt":       []Mode{PPT},
		"nested":    []Mode{Nested},
		"recursive": []Mode{Recursive},
		"all":       []Mode{PPT, Nested, Recursive},
	}
}

func modeTemplates() map[Mode]string {
	return map[Mode]string{
		PPT:       pptTemplate,
		Nested:    nestedTemplate,
		Recursive: recursiveTemplate,
	}
}
