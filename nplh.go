package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type messageType int

const (
	warning messageType = iota
	done
)

var (
	redSprintf   = color.New(color.FgRed).SprintfFunc()
	graySprintf  = color.New(color.FgHiBlack).SprintfFunc()
	greenSprintf = color.New(color.FgGreen).SprintfFunc()
)

func termPrintf(mt messageType, msg string, a ...interface{}) (n int, err error) {
	if mt == warning {
		return fmt.Printf(redSprintf("✘ "+msg, a...))
	} else if mt == done {
		return fmt.Printf(greenSprintf("✔") + " " + graySprintf(msg, a...))
	}
	return fmt.Printf(msg, a...)
}

func resolvePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	return path, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

type Symlink struct {
	Source  string
	Targets []string
}

func readConfig(path string) ([]Symlink, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})

	if err := yaml.Unmarshal(b, &m); err != nil {
		return []Symlink{}, err
	}

	symlinks := []Symlink{}
	for k, v := range m {
		targets := []string{}
		value, ok := v.(string)
		if ok {
			targets = []string{value}
		} else {
			for _, val := range v.([]interface{}) {
				targets = append(targets, val.(string))
			}
		}

		symlinks = append(symlinks, Symlink{
			Source:  k,
			Targets: targets,
		})
	}
	return symlinks, nil
}

func link(dotfileDirectory string) (err error) {
	configPath := filepath.Join(dotfileDirectory, "nplh.yml")
	config, err := readConfig(configPath)
	if err != nil {
		return err
	}
	for _, line := range config {
		for _, target := range line.Targets {
			path, err := resolvePath(target)
			if err != nil {
				return err
			}

			targetCurrentLink, err := filepath.EvalSymlinks(path)
			absoluteSource := filepath.Join(dotfileDirectory, line.Source)
			if err == nil && targetCurrentLink != absoluteSource {
				termPrintf(warning, target+" already exists, not overriding")
			} else if !fileExists(path) {
				if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
					return err
				}
				termPrintf(done, absoluteSource+" → "+target)
				if err := os.Symlink(absoluteSource, path); err != nil {
					return err
				}
			}
		}
	}
	termPrintf(done, "Done linking files\n")
	return nil
}

func run(args []string) error {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("error: getting current user: %v\n", err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "No Place Like Home"
	app.Usage = "A quick dotfile linker"
	app.Version = "1.0.0"

	cli.AppHelpTemplate = `
	NAME:
		 {{.Name}} - {{.Usage}}

	USAGE:
		 nplh [options]

	OPTIONS:
		 {{range .VisibleFlags}}{{.}}
		 {{end}}
	VERSION:
		 {{.Version}}
`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "directory, d",
			Value: filepath.Join(usr.HomeDir, "dotfiles"),
			Usage: "your dotfiles directory",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		dotfileDirectory := c.String("directory")
		return link(dotfileDirectory)
	}

	return app.Run(args)
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Printf("error: running app: %v\n", err)
	}
}
