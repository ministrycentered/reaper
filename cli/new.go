package main

import (
	"bufio"
	"errors"
	"html/template"
	"os"
	"path/filepath"

	"github.com/maddiesch/reaper"
)

func newAppHandler(c *reaper.Context) error {
	name := c.Argument("name")
	if name == "" {
		return errors.New("missing required name argument")
	}

	path, err := filepath.Abs(c.FlagString("path"))
	if err != nil {
		return err
	}

	c.Outputf("Creating `%s` at %s", name, path)

	return createNewApp(name, path, c.FlagString("github"))
}

var newProjectStructure = map[string]string{
	"Makefile":   fileNewMakefile,
	".gitignore": fileNewGitignore,
	"go.mod":     fileNewModule,
	"main.go":    fileNewMain,
}

func createNewApp(name, path, ghname string) error {
	root := filepath.Join(path, name)
	if exists(root) {
		return reaper.NewErrorf(84, "can't create directory at %s, it already exists", root)
	}

	err := os.Mkdir(root, os.ModePerm)
	if err != nil {
		return err
	}

	info := templateInfo{
		Name:          name,
		Github:        ghname,
		ReaperVersion: reaper.Version,
	}

	for name, content := range newProjectStructure {
		err := write(info, filepath.Join(root, name), content)
		if err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

type templateInfo struct {
	Name          string
	Github        string
	ReaperVersion string
}

func write(info templateInfo, path, content string) error {
	temp, err := template.New(filepath.Base(path)).Parse(content)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buff := bufio.NewWriter(f)

	err = temp.Execute(buff, info)

	buff.Flush()

	return nil
}
