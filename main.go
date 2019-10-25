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
	_ "github.com/savaki/kafka-protocol-gen/render/statik"
	"github.com/urfave/cli"
)

const suffix = ".go"

var opts struct {
	dir       string
	module    string
	src       string // src dir of protocol json files
	templates string // templates contains optional directory of templates
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
	//fileSystem, err := fs.New()
	//if err != nil {
	//	return fmt.Errorf("unable to load static assets: %w", err)
	//}

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

	var openFunc func(string) (io.ReadCloser, error)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		in, err := openFunc(path)
		if err != nil {
			return err
		}
		defer in.Close()

		data, err := ioutil.ReadAll(in)
		if err != nil {
			return err
		}

		t, err := template.New("code").Funcs(funcMap).Parse(string(data))
		if err != nil {
			return err
		}

		switch {
		case strings.Contains(path, "{{.MessageName}}"):
			for _, message := range messages {
				for version := message.ValidVersions.From; version <= message.ValidVersions.To; version++ {
					fn := func() error {
						rel := path
						if strings.HasPrefix(rel, opts.templates) {
							rel = rel[len(opts.templates):]
						}
						filename, err := interpolate(filepath.Join(opts.dir, rel), message, version)
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
							"Module":  opts.module,
							"Package": "v" + strconv.Itoa(version),
							"Message": message,
							"Structs": findStructFields(message.Fields),
							"Version": version,
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

	if opts.templates != "" {
		openFunc = func(path string) (io.ReadCloser, error) { return os.Open(path) }
		return filepath.Walk(opts.templates, walkFunc)
	}

	hfs, err := fs.New()
	if err != nil {
		return fmt.Errorf("unable to load static assets: %w", err)
	}

	openFunc = func(path string) (io.ReadCloser, error) { return hfs.Open(path) }
	return fs.Walk(hfs, "/", walkFunc)
}

var funcMap = template.FuncMap{
	"capitalize": capitalize,
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
	"isArray": isArray,
	"structName": func(a string) string {
		return strings.ReplaceAll(a, "[]", "")
	},
	"type": func(v string) string { return strings.ReplaceAll(v, "[]", "") },
}

func capitalize(v string) string {
	if len(v) == 0 {
		return ""
	}

	return strings.ToUpper(v[0:1]) + v[1:]
}

func isArray(t string) bool {
	return strings.Contains(t, "[]")
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

func interpolate(path string, message protocol.Message, version int) (string, error) {
	t, err := template.New("path").Parse(path)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	data := map[string]interface{}{
		"MessageName": snakeCase(message.Name),
		"Version":     version,
	}
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
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
