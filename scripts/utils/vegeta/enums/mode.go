package enums

import (
	"fmt"
	"strings"
)

type Mode string
type ModeTemplate map[Mode][]string

const (
	PPTRead     Mode = "ppt-read"
	PPTWrite    Mode = "ppt"
	NestedRead  Mode = "nested-read"
	NestedWrite Mode = "nested"
	Recursive   Mode = "recursive-read"
	All         Mode = "all"

	recursiveTemplate string = "%s/ppt/%s"
	pptTemplate       string = "%s/ppt/%s?recursive=false"
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
	if strings.Contains(string(m), "-read") {
		new := strings.ReplaceAll(string(m), "-read", "")
		return new
	}
	return string(m)
}

func (m Mode) IsRead() bool {
	return m == PPTRead || m == NestedRead || m == Recursive
}

func (m Mode) Template() string {
	return modeTemplates()[m]
}

func modes() []Mode {
	return []Mode{PPTRead, NestedRead, Recursive, All}
}

func (m Mode) Expand() []Mode {
	return modesMap()[m]
}

func modesMap() map[Mode][]Mode {
	return map[Mode][]Mode{
		"ppt":    []Mode{PPTRead, PPTWrite, Recursive},
		"nested": []Mode{NestedRead, NestedWrite},
		"all":    []Mode{PPTRead, PPTWrite, NestedRead, NestedWrite, Recursive},
	}
}

func modeTemplates() map[Mode]string {
	return map[Mode]string{
		PPTRead:     pptTemplate,
		PPTWrite:    pptTemplate,
		NestedWrite: nestedTemplate,
		NestedRead:  nestedTemplate,
		Recursive:   recursiveTemplate,
	}
}
