package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func exit(err string) {
	fmt.Println("Error: " + err)
	os.Exit(1)
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

func repoUrl(repository string) (repo string) {
	_, err := url.ParseRequestURI(repository)

	if err == nil {
		return repository
	}
	return "https://github.com/" + repository
}

type Symlink struct {
	Source  string
	Targets []string
}

func readConfig(path string) (config []Symlink) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	lines := strings.Split(string(b), "\n")
	symlinks := []Symlink{}
	for _, line := range lines {
		if line != "" {
			words := strings.Split(line, " ")
			symlinks = append(symlinks, Symlink{
				Source:  words[0],
				Targets: words[1:],
			})
		}
	}
	return symlinks
}

func main() {
	usr, _ := user.Current()
	dotfileDirectory := filepath.Join(usr.HomeDir, "dotfiles")
	configPath := filepath.Join(dotfileDirectory, "nplh")

	app := cli.NewApp()
	app.Name = "No Place Like Home"
	app.Usage = "A quick dotfile linker"
	app.Version = "0.0.9"

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "setup a dotfiles directory",
			Action: func(c *cli.Context) {
				fmt.Println("init")
			},
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install a dotfiles repo",
			Action: func(c *cli.Context) {
				if c.Args().Get(0) == "" {
					exit("Please specify a repository")
				}
				repository := repoUrl(c.Args().Get(0))
				if fileExists(dotfileDirectory) {
					exit("directory " + dotfileDirectory + " exists")
				}
				fmt.Printf("Installing %s\n", repository)
				cmd := exec.Command("git", "clone", repository, dotfileDirectory)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "link out the files to their corresponding homes",
			Action: func(c *cli.Context) {
				for _, line := range readConfig(configPath) {
					for _, target := range line.Targets {
						targetCurrentLink, err := os.Readlink(resolvePath(target))
						absoluteSource := filepath.Join(dotfileDirectory, line.Source)
						if err == nil && targetCurrentLink != absoluteSource {
							fmt.Println(target + " already exists, not overriding")
						} else if !fileExists(resolvePath(target)) {
							fmt.Println(absoluteSource + " -> " + target)
							os.Symlink(absoluteSource, resolvePath(target))
						}
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
