package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Strings []string

func (s Strings) Empty() bool {
	return len(s) == 0
}

func (s Strings) Contains(target string) bool {
	for _, val := range s {
		if target == val {
			return true
		}
	}
	return false
}

func findWordByFile(path, findWord string) (bool, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return bytes.Contains(input, []byte(findWord)), nil
}

func replaceFile(path, findWord, replaceWord string) (bool, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	isFind := bytes.Contains(input, []byte(findWord))
	if !isFind {
		return false, nil
	}
	output := bytes.Replace(input, []byte(findWord), []byte(replaceWord), -1)

	if err = ioutil.WriteFile(path, output, 0666); err != nil {
		fmt.Println(err)
		return isFind, err
	}
	return isFind, nil
}

func main() {
	var depth int = 1
	var files, ignores Strings
	var targetStr, replaceStr string
	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		if v == "--help" || v == "-h" {
			fmt.Println("usage: replacer [--depth=xxxx] [--target=xxxx] [--replace=xxxx] [--files=aaaa,bbbb] [--ignore=cccc,dddd]")
			return
		} else if strings.HasPrefix(v, "--depth=") {
			tmpVal, err := strconv.ParseUint(strings.Split(v, "=")[1], 10, 64)
			if err != nil {
				panic(err)
			}
			depth = int(tmpVal)
		} else if strings.HasPrefix(v, "--target=") {
			targetStr = strings.Split(v, "=")[1]
		} else if strings.HasPrefix(v, "--replace=") {
			replaceStr = strings.Split(v, "=")[1]
		} else if strings.HasPrefix(v, "--files=") {
			files = strings.Split(strings.Split(v, "=")[1], ",")
		} else if strings.HasPrefix(v, "--ignore=") {
			ignores = strings.Split(strings.Split(v, "=")[1], ",")
		} else {
			fmt.Printf("invalid argument, %s\n", v)
		}
	}
	isFind := targetStr != ""
	isReplace := replaceStr != ""

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	maxDepth := strings.Count(baseDir, string(os.PathSeparator)) + depth

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath := "." + strings.Split(path, baseDir)[1]

		if info.IsDir() {
			if strings.Count(path, string(os.PathSeparator)) > maxDepth {
				return filepath.SkipDir
			} else if ignores.Contains(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !files.Empty() && !files.Contains(info.Name()) {
			return nil
		}

		if isFind {
			isFindTarget, err := findWordByFile(path, targetStr)
			if err != nil {
				return err
			} else if isFindTarget {
				fmt.Printf("path: %#v\n", relativePath)
			}

			if isFindTarget && isReplace {
				_, err := replaceFile(path, targetStr, replaceStr)
				if err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("path: %#v\n", relativePath)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
