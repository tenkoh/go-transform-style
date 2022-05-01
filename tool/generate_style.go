package main

import (
	"bytes"
	"go/format"
	"html/template"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const tmpl = `// Code generated by tool/generate_style.go; DO NOT EDIT.
package textstyle
const ({{range .}}
    {{.Label}}_LOWER_OFFSET = {{printf "0x%x" .LowerOffset}}
    {{.Label}}_UPPER_OFFSET = {{printf "0x%x" .UpperOffset}}
    {{.Label}}_DIGIT_OFFSET = {{printf "0x%x" .DigitOffset}}{{end}}
)
{{range .}}
var {{.TitleLabel}} = NewTransformer(
    NewSimpleReplacer(
        {{.Label}}_LOWER_OFFSET,
        {{.Label}}_UPPER_OFFSET,
        {{.Label}}_DIGIT_OFFSET,
    ),
)
{{end}}`

type Style struct {
	Label       string
	TitleLabel  string
	LowerOffset uint32
	UpperOffset uint32
	DigitOffset uint32
}

func byteDiff(a, b string) uint32 {
	a32 := stringToUint(a)
	b32 := stringToUint(b)
	return a32 - b32
}

func stringToUint(s string) uint32 {
	b := []byte(s)
	var sum uint32
	for _, bb := range b {
		sum = sum << 8
		sum += uint32(bb)
	}
	return sum
}

// replace must be aA0 with a specific text style.
func NewStyle(label string, replace string) Style {
	ss := make([]string, 3)
	n := 0
	for _, r := range replace {
		if n > 2 {
			panic("argument `replace` must be each text style's aA0")
		}
		ss[n] = string(r)
		n++
	}
	if n != 3 {
		panic("argument `replace` must be each text style's aA0")
	}
	// for generate title case
	c := cases.Title(language.English)
	s := Style{
		strings.ToUpper(label),
		c.String(label),
		byteDiff(ss[0], "a"),
		byteDiff(ss[1], "A"),
		byteDiff(ss[2], "0"),
	}
	return s
}

func main() {
	styles := []Style{
		NewStyle("bold", "𝐚𝐀𝟎"),
		NewStyle("italic", "𝑎𝐴0"),
	}
	t := template.Must(template.New("style").Parse(tmpl))
	f, err := os.Create("../styles.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, styles); err != nil {
		panic(err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = f.Write(formatted)
	if err != nil {
		panic(err)
	}
}
