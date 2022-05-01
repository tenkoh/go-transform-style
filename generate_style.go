//go:build ignore
// +build ignore

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

const tmpl = `// Code generated by generate_style.go; DO NOT EDIT.
package textstyle
const ({{range .}}
    {{.Label}}_LOWER_OFFSET = {{printf "%d" .LowerOffset}}
    {{.Label}}_UPPER_OFFSET = {{printf "%d" .UpperOffset}}
    {{.Label}}_DIGIT_OFFSET = {{printf "%d" .DigitOffset}}{{end}}
)
{{range .}}
func {{.TitleLabel}}() *Transformer{
    return NewTransformer(
        NewSimpleReplacer(
            "{{.TitleLabel}}",
            {{.Label}}_LOWER_OFFSET,
            {{.Label}}_UPPER_OFFSET,
            {{.Label}}_DIGIT_OFFSET,
        ),
    )
}
{{end}}`

type Style struct {
	Label       string
	TitleLabel  string
	LowerOffset rune
	UpperOffset rune
	DigitOffset rune
}

// edit here to add some styles with simpleReplacer.
// Label could contain space, but it is transformed as below.
// Style.Label => All upper case with underscore
// Style.TitleLabel => Title case with no space
var Styles = []Style{
	NewStyle("bold", "𝐚𝐀𝟎"),
	NewStyle("italic", "𝑎𝐴0"),
	NewStyle("bold italic", "𝒂𝑨𝟎"),
	NewStyle("script", "𝒶𝒜0"),
	NewStyle("bold script", "𝓪𝓐𝟎"),
	NewStyle("fraktur", "𝔞𝔄0"),
	NewStyle("bold fraktur", "𝖆𝕬𝟎"),
	NewStyle("double struck", "𝕒𝔸𝟘"),
	NewStyle("sans serif", "𝖺𝖠𝟢"),
	NewStyle("sans serif bold", "𝗮𝗔𝟬"),
	NewStyle("sans serif italic", "𝘢𝘈𝟢"),
	NewStyle("sans serif bold italic", "𝙖𝘼𝟬"),
	NewStyle("monospace", "𝚊𝙰𝟶"),
}

// replacePattern must be aA0 with a specific text style.
func NewStyle(label string, replacePattern string) Style {
	ss := make([]rune, 3)
	var n int
	for _, r := range replacePattern {
		if n > 2 {
			panic("argument `replace` must be each text style's aA0")
		}
		ss[n] = r
		n++
	}
	if n != 3 {
		panic("argument `replace` must be each text style's aA0")
	}
	s := Style{
		toLabel(label),
		toTitleLabel(label),
		ss[0] - 'a',
		ss[1] - 'A',
		ss[2] - '0',
	}
	return s
}

func toLabel(s string) string {
	// replace space by underscore
	ss := strings.Replace(s, " ", "_", -1)
	return strings.ToUpper(ss)
}

func toTitleLabel(s string) string {
	// transform into title case, then remove space
	c := cases.Title(language.English)
	ss := c.String(s)
	return strings.Replace(ss, " ", "", -1)
}

func main() {
	t := template.Must(template.New("style").Parse(tmpl))
	f, err := os.Create("styles.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, Styles); err != nil {
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
