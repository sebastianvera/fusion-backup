package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GoProFile represents multiple files that have the same name but different
// extensions.
//
// i.e:
// foobar.ext1 foobar.ext2 foobar.ext3 foobar.info
//
// Is represented by:
// GoProFile{Name: "foobar", Extensions: ["ext1", "ext2", "ext3", "info"]}
type GoProFile struct {
	Name       string
	Extensions []string
}

func main() {
	// TODO: parse srcDir and outDir from args
	srcDir, err := homedir.Expand("~/Downloads/gopro")
	if err != nil {
		panic(err)
	}
	outDir, err := homedir.Expand("~/Downloads/gopro/out")
	if err != nil {
		panic(err)
	}

	if err := createFolder(outDir); err != nil {
		panic(err)
	}

	goProFiles := readGoProFiles(srcDir)
	for _, goProFile := range goProFiles {
		filename := goProFile.Name
		srcPath := path.Join(srcDir, filename)
		dstPath := path.Join(outDir, filename)

		if err := createFolder(dstPath); err != nil {
			panic(err)
		}

		for _, extension := range goProFile.Extensions {
			from := fmt.Sprintf("%s.%s", srcPath, extension)
			to := fmt.Sprintf("%s/%s.%s", dstPath, filename, extension)

			if err := os.Rename(from, to); err != nil {
				panic(err)
			}
		}
	}
}

func createFolder(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func readGoProFiles(srcPath string) []*GoProFile {
	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		panic(err)
	}

	uniqueNames := make(map[string][]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		separateName := strings.Split(file.Name(), ".")
		if len(separateName) != 2 || len(separateName[0]) == 0 {
			continue
		}
		name := separateName[0]
		ext := separateName[1]

		uniqueNames[name] = append(uniqueNames[name], ext)
	}

	i := 0
	goProFiles := make([]*GoProFile, len(uniqueNames))
	for name, extensions := range uniqueNames {
		goProFiles[i] = &GoProFile{
			Name:       name,
			Extensions: extensions,
		}

		i++
	}

	return goProFiles
}
