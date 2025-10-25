package site

import (
	"bytes"
	"html/template"
	"os"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

type Page struct {
	Title   string `yaml:"title"`
	Type    string `yaml:"type"`
	Studio  string `yaml:"studio"`
	Summary string `yaml:"summary"`
	Date    string `yaml:"date"`

	Body        template.HTML  `yaml:"-"`
	FrontMatter map[string]any `yaml:"-"`
}

func LoadPage(path string) (*Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fm map[string]any
	body, err := frontmatter.Parse(bytes.NewReader(data), &fm)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := decodeFrontMatter(fm, &page); err != nil {
		return nil, err
	}

	md := markdownRenderer()
	var buf bytes.Buffer
	if err := md.Convert(body, &buf); err != nil {
		return nil, err
	}

	page.Body = template.HTML(buf.String())
	page.FrontMatter = fm
	return &page, nil
}

func decodeFrontMatter(src map[string]any, dst any) error {
	if src == nil {
		return nil
	}
	raw, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(raw, dst)
}
