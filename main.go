package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/savaki/kafka-protocol-gen/protocol"
	"github.com/savaki/kafka-protocol-gen/render"
	"github.com/urfave/cli"
)

var opts struct {
	dir string
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Value:       ".",
			Usage:       "directory containing json kafka protocol definition",
			Destination: &opts.dir,
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

		for version := message.ValidVersions.From; version <= message.ValidVersions.To; version++ {
			fmt.Println("------------------------")
			err := render.Message(os.Stdout, message, version)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return filepath.Walk(opts.dir, callback)
}
