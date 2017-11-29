package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := os.Chdir("test")
	if err != nil {
		fmt.Printf("error: could not chdir into 'test': %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

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
	err := os.RemoveAll("other-vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".zshrc")
	assert.Nil(t, err)

	err = run([]string{"nplh.test", "--directory", "dotfiles"})
	assert.Nil(t, err)

	checkFile(".vimrc", "this is vimrc", t)
	checkFile("other-vimrc", "this is vimrc", t)
	checkFile(".zshrc", "this is zshrc", t)

	err = os.RemoveAll("other-vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".zshrc")
	assert.Nil(t, err)
}

func TestMissingConfig(t *testing.T) {
	err := run([]string{"nplh.test", "--directory", "adslkj"})

	if err == nil {
		t.Error("run was supposed to throw an error for missing config")
	}
}

func TestNoOverride(t *testing.T) {
	err := os.RemoveAll("other-vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".zshrc")
	assert.Nil(t, err)

	err = ioutil.WriteFile(".vimrc", []byte("this is original vimrc"), 0644)
	assert.Nil(t, err)

	err = run([]string{"nplh.test", "--directory", "dotfiles"})
	assert.Nil(t, err)

	checkFile(".vimrc", "this is original vimrc", t)
	checkFile("other-vimrc", "this is vimrc", t)
	checkFile(".zshrc", "this is zshrc", t)

	err = os.RemoveAll("other-vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".vimrc")
	assert.Nil(t, err)
	err = os.RemoveAll(".zshrc")
	assert.Nil(t, err)
}

func TestResolvePath(t *testing.T) {
	path, err := resolvePath("~/foobar")
	if err != nil {
		t.Errorf("error: resolving path: %v\n", err)
	}

	usr, err := user.Current()
	if err != nil {
		t.Errorf("error: getting current user: %v\n", err)
	}

	dir := usr.HomeDir
	expected := dir + "/foobar"

	if path != expected {
		t.Error(path + " does not match " + dir + expected)
	}
}

func TestFileExists(t *testing.T) {
	err := ioutil.WriteFile("foobar", []byte(""), 0644)
	assert.Nil(t, err)

	if !fileExists("foobar") {
		t.Error("fileExists should be true")
	}

	err = os.RemoveAll("foobar")
	assert.Nil(t, err)
}
