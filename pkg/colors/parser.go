package colors

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"regexp"
	"strings"

	"github.com/muesli/gamut"
)

type Parser struct {
	Colors []color.Color
	re     *regexp.Regexp
}

// inits a color parser
func NewParser() *Parser {
	// re := regexp.MustCompile(`^#?[a-fA-F0-9]{3,6}$`)

	// pattern := "A#([A-Fa-f0-9]{6})$"
	pattern := "(#|0x)?[0-9a-fA-F]{6}"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regular expression:", err)
	}

	return &Parser{
		re: regex,
	}
}

// delete all colors in parser
func (p *Parser) ClearColors() {
	p.Colors = p.Colors[:0] // clear slice
}

// parse hex colors from a file or generic reader (from a request, stdin, a string reader etc...)
func (p *Parser) ParseString(s string) error {
	matches := p.re.FindAllString(s, -1)
	if matches != nil {
		for _, c := range matches {
			clr := gamut.Hex(c)
			p.Colors = append(p.Colors, clr)
		}
	}

	return nil
}

// parse hex colors from a file or generic reader (from a request, stdin, a string reader etc...)
func (p *Parser) ParseFile(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		matches := p.re.FindAllString(line, -1)

		if matches != nil {
			for _, c := range matches {
				// println(c)
				clr := gamut.Hex(strings.TrimSpace(c))
				p.Colors = append(p.Colors, clr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
