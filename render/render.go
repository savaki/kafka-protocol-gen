package render

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/rakyll/statik/fs"
	"github.com/savaki/kafka-protocol-gen/protocol"
	_ "github.com/savaki/kafka-protocol-gen/render/statik"
)

var _t *template.Template

var funcMap = template.FuncMap{
	"forVersion": func(version int, fields []protocol.Field) []protocol.Field {
		var valid []protocol.Field
		for _, f := range fields {
			field := f
			if version >= field.Versions.From && (version <= field.Versions.To || field.Versions.UpToCurrent) {
				valid = append(valid, field)
			}
		}
		return valid
	},
	"structName": func(a string) string {
		return strings.ReplaceAll(a, "[]", "")
	},
}

func getTemplate() (*template.Template, error) {
	if _t == nil {
		fileSystem, err := fs.New()
		if err != nil {
			return nil, fmt.Errorf("unable to load static assets: %w", err)
		}

		f, err := fileSystem.Open("/content.gogo")
		if err != nil {
			return nil, err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		_t = template.Must(template.New("code").Funcs(funcMap).Parse(string(data)))
	}

	return _t, nil
}

// Render specific version of message to writer
func Message(w io.Writer, message protocol.Message, version int) error {
	data := map[string]interface{}{
		"Package": "v" + strconv.Itoa(version),
		"Message": message,
		"Structs": findStructFields(message.Fields),
		"Version": version,
	}

	t, err := getTemplate()
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}

func findStructFields(fields []protocol.Field) []protocol.Field {
	var structFields []protocol.Field
	for _, f := range fields {
		field := f
		if len(field.Fields) == 0 {
			continue
		}

		structFields = append(structFields, field)
		structFields = append(structFields, findStructFields(field.Fields)...)
	}
	return structFields
}
