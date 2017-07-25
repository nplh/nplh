package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"testing"
)

func checkFile(file string, expected string, t *testing.T) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Error(err)
	}

	if strings.TrimSpace(string(b)) != strings.TrimSpace(expected) {
		t.Errorf("in " + file + ", \"" + string(b) + "\" did not match \"" + expected + "\"")
	}
}

func TestLink(t *testing.T) {
	os.Chdir("test")
	os.RemoveAll("other-vimrc")
	os.RemoveAll(".vimrc")
	os.RemoveAll(".zshrc")
	run([]string{"nplh.test", "--directory", "dotfiles"})

	checkFile(".vimrc", "this is vimrc", t)
	checkFile("other-vimrc", "this is vimrc", t)
	checkFile(".zshrc", "this is zshrc", t)

	os.RemoveAll("other-vimrc")
	os.RemoveAll(".vimrc")
	os.RemoveAll(".zshrc")
}

func TestMissingConfig(t *testing.T) {
	os.Chdir("test")
	err := run([]string{"nplh.test", "--directory", "adslkj"})

	if err == nil {
		t.Error("Was supposed to throw an error for missing config")
	}
}

func TestNoOverride(t *testing.T) {
	os.Chdir("test")
	os.RemoveAll("other-vimrc")
	os.RemoveAll(".vimrc")
	os.RemoveAll(".zshrc")

	ioutil.WriteFile(".vimrc", []byte("this is original vimrc"), 0644)

	run([]string{"nplh.test", "--directory", "dotfiles"})

	checkFile(".vimrc", "this is original vimrc", t)
	checkFile("other-vimrc", "this is vimrc", t)
	checkFile(".zshrc", "this is zshrc", t)

	os.RemoveAll("other-vimrc")
	os.RemoveAll(".vimrc")
	os.RemoveAll(".zshrc")
}

func TestResolvePath(t *testing.T) {
	path := resolvePath("~/foobar")
	usr, _ := user.Current()
	dir := usr.HomeDir
	expected := dir + "/foobar"

	if path != expected {
		t.Error(path + " does not match " + dir + expected)
	}
}

func TestFileExists(t *testing.T) {
	ioutil.WriteFile("foobar", []byte(""), 0644)
	if !fileExists("foobar") {
		t.Error("fileExists should be true")
	}
	os.RemoveAll("foobar")
}
