package golog

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ColorMatch struct {
	Text  string
	Color string

	c Color
}
type ColorFile struct {
	Rule []*ColorMatch
}

func (self *ColorFile) ColorFromText(text string) Color {

	for _, rule := range self.Rule {
		if strings.Contains(text, rule.Text) {
			return rule.c
		}
	}

	return NoColor
}

func (self *ColorFile) Load(data string) error {

	err := json.Unmarshal([]byte(data), self)
	if err != nil {
		return err
	}

	for _, rule := range self.Rule {

		rule.c = matchColor(rule.Color)

		if rule.c == NoColor {
			return fmt.Errorf("color name not exists: %s", rule.Text)
		}

	}

	return nil
}

func NewColorFile() *ColorFile {
	return &ColorFile{}
}
