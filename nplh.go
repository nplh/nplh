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

func warning(msg string, a ...interface{}) {
	red := color.New(color.FgRed).SprintfFunc()
	fmt.Println(red("✘ "+msg, a...))
}

func done(msg string) {
	gray := color.New(color.FgHiBlack).SprintfFunc()
	green := color.New(color.FgGreen).SprintfFunc()
	fmt.Println(green("✔") + " " + gray(msg))
}

func resolvePath(path string) (newpath string) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func fileExists(path string) (exists bool) {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

type Symlink struct {
	Source  string
	Targets []string
}

func readConfig(path string) (config []Symlink, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	yaml.Unmarshal(b, &m)
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
			targetCurrentLink, err := filepath.EvalSymlinks(resolvePath(target))
			absoluteSource := filepath.Join(dotfileDirectory, line.Source)
			if err == nil && targetCurrentLink != absoluteSource {
				warning(target + " already exists, not overriding")
			} else if !fileExists(resolvePath(target)) {
				os.MkdirAll(filepath.Dir(resolvePath(target)), 0777)
				done(absoluteSource + " → " + target)
				os.Symlink(absoluteSource, resolvePath(target))
			}
		}
	}
	done("Done linking files")
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
