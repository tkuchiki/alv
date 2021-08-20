package visualiser

import (
	"io"
	"os"
	"regexp"

	"github.com/tkuchiki/alv/cmd/alv/option"

	"github.com/tkuchiki/alp/helpers"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alv/diagram"
	"github.com/tkuchiki/alv/parser"
)

type Visualizer struct {
	parser            parser.Parser
	inReader          *os.File
	outWriter         *os.File
	uriMatchingGroups []*regexp.Regexp
	sankey            *diagram.Sankey
	logFormat         string
}

func NewVisualizer(logFormat string) *Visualizer {
	return &Visualizer{
		sankey:    diagram.NewSankey(),
		logFormat: logFormat,
	}
}

func (v *Visualizer) Render() error {
	data := map[string][]parser.ParsedLog{}

Loop:
	for {
		s, err := v.parser.Parse()
		if err != nil {
			if err == io.EOF {
				break
			} else if err == errors.SkipReadLineErr {
				continue Loop
			}

			return err
		}

		var uri string
		if len(v.uriMatchingGroups) > 0 {
			uri = s.Uri
			for _, re := range v.uriMatchingGroups {
				if ok := re.Match([]byte(uri)); ok {
					pattern := re.String()
					uri = pattern
				}
			}
			s.Uri = uri
		}

		data[s.User] = append(data[s.User], *s)
	}
	return v.sankey.Render(v.outWriter, data)
}

func (v *Visualizer) SetParser(opts *option.Options) {
	var p parser.Parser
	switch opts.LogFormat {
	case "json":
		keys := parser.NewJSONKeys(opts.KeyOption.UriKey, opts.KeyOption.MethodKey, opts.KeyOption.UserKey)
		p = parser.NewJSONParser(v.GetInReader(), keys, opts.QueryString, opts.QueryStringIgnoreValues)
	}

	v.parser = p
}

func (v *Visualizer) SetInReader(filename string) error {
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		v.inReader = f
	} else {
		v.inReader = os.Stdin
	}

	return nil
}

func (v *Visualizer) GetInReader() *os.File {
	return v.inReader
}

func (v *Visualizer) CloseInReader() {
	v.inReader.Close()
}

func (v *Visualizer) SetOutWriter(filename string) error {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}

		v.outWriter = f
	} else {
		v.outWriter = os.Stdout
	}

	return nil
}

func (v *Visualizer) CloseOutWriter() {
	v.outWriter.Close()
}

func (v *Visualizer) SetUriMatchingGroups(groups []string) error {
	uriGroups, err := helpers.CompileUriMatchingGroups(groups)
	if err != nil {
		return err
	}

	v.uriMatchingGroups = uriGroups

	return nil
}
