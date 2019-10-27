package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/rakyll/statik/fs"
	"github.com/savaki/kafka-protocol-gen/protocol"

	//_ "github.com/savaki/kafka-protocol-gen/render/statik"
	"github.com/urfave/cli"
)

const suffix = ".go"

var opts struct {
	dir       string
	module    string
	src       string // src dir of protocol json files
	templates string // templates contains optional directory of templates
	last      int    // only include the last N versions; 0 means include all versions
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Value:       ".",
			Usage:       "output directory",
			Destination: &opts.dir,
		},
		cli.IntFlag{
			Name:        "last",
			Usage:       "last N versions",
			Destination: &opts.last,
		},
		cli.StringFlag{
			Name:        "module",
			Usage:       "module name",
			Destination: &opts.module,
		},
		cli.StringFlag{
			Name:        "src",
			Value:       ".",
			Usage:       "directory containing json kafka protocol definition",
			Destination: &opts.src,
		},
		cli.StringFlag{
			Name:        "templates",
			Usage:       "optional directory of templates to render",
			Destination: &opts.templates,
		},
	}
	app.EnableBashCompletion = true
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func action(_ *cli.Context) error {
	dir, err := writeTemplates()
	if err != nil {
		return err
	}
	if opts.templates == "" {
		defer os.RemoveAll(dir)
	}

	var filenames []string
	fn := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filenames = append(filenames, path)
		}
		return nil
	}
	if err := filepath.Walk(dir, fn); err != nil {
		return fmt.Errorf("unable to read directory: %w", err)
	}

	all, err := template.New("templates").Funcs(funcMap).ParseFiles(filenames...)
	if err != nil {
		return fmt.Errorf("unable to load templates: %w", err)
	}
	all = all.Funcs(funcMap)

	fmt.Println(all.DefinedTemplates())

	var messages []protocol.Message
	callback := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("unable to open file, %v: %w", path, err)
		}
		defer f.Close()

		message, err := protocol.Parse(f)
		if err != nil {
			return fmt.Errorf("unable to parse file, %v: %w", path, err)
		}

		messages = append(messages, message)
		return nil
	}

	if err := filepath.Walk(opts.src, callback); err != nil {
		return err
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "_") {
			return nil
		}

		t := all.Lookup(filepath.Base(path))
		if t == nil {
			return fmt.Errorf("unable to lookup template, %v: %w", filepath.Base(path), err)
		}

		switch {
		case strings.Contains(path, "{{.MessageName}}"):
			for _, message := range messages {
				versions := protocol.ValidVersions{To: message.ValidVersions.To}
				if opts.last > 0 {
					if from := message.ValidVersions.To - opts.last + 1; from > 0 {
						versions.From = from
					}
				}

				fn := func() error {
					rel := path
					if strings.HasPrefix(rel, opts.templates) {
						rel = rel[len(opts.templates):]
					}
					filename, err := interpolate(filepath.Join(opts.dir, rel), message, message.ApiKey)
					if err != nil {
						return err
					}

					if ext := filepath.Ext(filename); strings.HasPrefix(ext, suffix) && len(ext) > len(suffix) {
						filename = filename[0:len(filename)-len(ext)] + "." + ext[len(suffix):]
					}

					if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
						return err
					}

					f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						return err
					}
					defer f.Close()

					data := map[string]interface{}{
						"Module":   opts.module,
						"Package":  "v0",
						"Message":  message,
						"Versions": versions,
					}

					if err := t.Execute(f, data); err != nil {
						return err
					}
					fmt.Println("wrote", filename)

					return nil
				}

				if err := fn(); err != nil {
					return err
				}
			}
		default:
			fn := func() error {
				rel := path
				if strings.HasPrefix(rel, opts.templates) {
					rel = rel[len(opts.templates):]
				}
				filename := filepath.Join(opts.dir, rel)

				if ext := filepath.Ext(filename); strings.HasPrefix(ext, suffix) && len(ext) > len(suffix) {
					filename = filename[0:len(filename)-len(ext)] + "." + ext[len(suffix):]
				}

				if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
					return err
				}

				f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				defer f.Close()

				data := map[string]interface{}{
					"Module": opts.module,
				}

				if err := t.Execute(f, data); err != nil {
					return err
				}
				fmt.Println("wrote", filename)

				return nil
			}

			if err := fn(); err != nil {
				return err
			}
		}

		return nil
	}

	return filepath.Walk(dir, walkFunc)
}

type VersionFields struct {
	ApiKey   int
	Fields   []protocol.Field
	Name     string
	Versions protocol.ValidVersions
}

var funcMap = template.FuncMap{
	"baseType":         baseType,
	"capitalize":       capitalize,
	"findStructs":      findStructs,
	"findStructFields": findStructFields,
	"forVersion": func(versions protocol.ValidVersions, fields []protocol.Field) []protocol.Field {
		var valid []protocol.Field

	loop:
		for _, f := range fields {
			field := f
			for version := versions.From; version <= versions.To; version++ {
				if f.Versions.IsValid(version) {
					valid = append(valid, field)
					continue loop
				}
			}
		}
		return valid
	},
	"goType":           goType,
	"isArray":          isArray,
	"isBytes":          isBytes,
	"isPartialOverlap": isPartialOverlap,
	"isPrimitiveArray": isPrimitiveArray,
	"isString":         isString,
	"isStructArray":    isStructArray,
	"structName": func(a string) string {
		return strings.ReplaceAll(a, "[]", "")
	},
	"toVersionFields": func(versions protocol.ValidVersions, message protocol.Message) VersionFields {
		return VersionFields{
			ApiKey:   message.ApiKey,
			Fields:   message.Fields,
			Name:     message.Name,
			Versions: versions,
		}
	},
	"type": func(v string) string { return strings.ReplaceAll(v, "[]", "") },
}

func baseType(v string) string {
	return strings.ReplaceAll(v, "[]", "")
}

func capitalize(v string) string {
	if len(v) == 0 {
		return ""
	}

	return strings.ToUpper(v[0:1]) + v[1:]
}

func goType(t string) string {
	switch t {
	case "bytes":
		return "[]byte"
	default:
		return t
	}
}

func isArray(t string) bool {
	return strings.Contains(t, "[]")
}

func isBytes(t string) bool {
	return t == "bytes"
}

func isPartialOverlap(valid protocol.ValidVersions, versions protocol.Versions) bool {
	var matches int
	for version := valid.From; version <= valid.To; version++ {
		if versions.IsValid(version) {
			matches++
		}
	}

	want := (valid.To - valid.From) + 1
	return matches != want
}

func isPrimitiveArray(t string) bool {
	return t == "[]string" || t == "[]int32" || t == "[]int64"
}

func isString(t string) bool {
	return t == "string"
}

func isStructArray(t string) bool {
	return isArray(t) && !isPrimitiveArray(t)
}

var re = regexp.MustCompile(`^[^A-Za-z0-9]*([A-Z0-9]*)([a-z0-9]*)`)

func kebabCase(v string) string {
	remain := v
	updated := make([]byte, 0, 2*len(v))
	for remain != "" {
		var (
			match        = re.FindStringSubmatch(remain)
			upper, lower = match[1], match[2]
		)
		remain = remain[len(match[0]):]

		if upper == "" && lower == "" {
			continue
		}
		if len(updated) > 0 {
			updated = append(updated, '-')
		}
		updated = append(updated, strings.ToLower(upper)...)
		updated = append(updated, lower...)
	}
	return string(updated)
}

func snakeCase(v string) string {
	remain := v
	updated := make([]byte, 0, 2*len(v))
	for remain != "" {
		var (
			match        = re.FindStringSubmatch(remain)
			upper, lower = match[1], match[2]
		)
		remain = remain[len(match[0]):]

		if upper == "" && lower == "" {
			continue
		}
		if len(updated) > 0 {
			updated = append(updated, '_')
		}
		updated = append(updated, strings.ToLower(upper)...)
		updated = append(updated, lower...)
	}
	return string(updated)
}

func interpolate(path string, message protocol.Message, apiKey int) (string, error) {
	t, err := template.New("path").Parse(path)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	data := map[string]interface{}{
		"MessageName": snakeCase(message.Name),
		"ApiKey":      apiKey,
	}
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func findStructFields(apiKey int, versions protocol.ValidVersions, fields []protocol.Field) []VersionFields {
	var structFields []VersionFields
	for _, f := range fields {
		if len(f.Fields) == 0 {
			continue
		}
		if !f.Versions.IsValidVersions(versions) {
			continue
		}

		item := VersionFields{
			ApiKey:   apiKey,
			Fields:   f.Fields,
			Name:     baseType(f.Type) + strconv.Itoa(apiKey),
			Versions: versions,
		}
		structFields = append(structFields, item)
		structFields = append(structFields, findStructFields(apiKey, versions, f.Fields)...)
	}
	return structFields
}

func findStructs(apiKey int, versions protocol.ValidVersions, message protocol.Message) []VersionFields {
	structFields := findStructFields(apiKey, versions, message.Fields)
	for _, f := range message.CommonStructs {
		item := VersionFields{
			ApiKey:   apiKey,
			Fields:   f.Fields,
			Name:     f.Name + strconv.Itoa(apiKey),
			Versions: versions,
		}
		structFields = append(structFields, item)
	}
	return structFields
}

func writeTemplates() (string, error) {
	if opts.templates != "" {
		return opts.templates, nil
	}

	dir, err := ioutil.TempDir(os.TempDir(), "templates-")
	if err != nil {
		return "", fmt.Errorf("unable to create temporary directory: %w", err)
	}

	hfs, err := fs.New()
	if err != nil {
		return "", fmt.Errorf("unable to load static assets: %w", err)
	}

	callback := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		filename := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Base(filename), 0755); err != nil {
			return fmt.Errorf("unable to write directoy, %v: %w", filepath.Base(path), err)
		}

		out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("unable to create file, %v: %w", filename, err)
		}
		defer out.Close()

		in, err := hfs.Open(path)
		if err != nil {
			return fmt.Errorf("unable to open template, %v: %w", path, err)
		}
		defer in.Close()

		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf("unable to copy content: %w", err)
		}

		return nil
	}
	if err := fs.Walk(hfs, "/", callback); err != nil {
		return "", err
	}

	return dir, nil
}
