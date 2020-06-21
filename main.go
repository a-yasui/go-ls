package main

import (
	"flag"
	"fmt"
	"github.com/yookoala/realpath"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"path/filepath"
	"sort"
)

/**
 * see : https://flaviocopes.com/go-list-files/
 */

type lsFile struct {
	path string
	info os.FileInfo
	err  error
}

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
}

var ALL_SHOW bool    // -a option
var LISTING_OPT bool // -l option

/**
 *
 */
func __realpath(path string) (new_path string) {
	var _new_path, err = realpath.Realpath(path)
	if err != nil {
		log.Fatal(err)
	}
	return _new_path
}

/**
 *
 */
func getLookDirectory() (lookDirectory []string) {
	var lsdir = flag.Args()
	if len(lsdir) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return []string{__realpath(dir)}
	}

	// 引数の中身を正規化
	var argc = len(lsdir)
	var result []string
	for i := 0; i < argc; i++ {
		result = append(result, __realpath(lsdir[i]))
	}
	return lsdir
}


func getFilesAtDirecotry (root string,walkFn filepath.WalkFunc ) error{
	f, err := os.Open(root)
	if err != nil {
		return err
	}

	list, err:= f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	// Sortable
	sort.Slice(list, func(i,j int) bool {
		return list[i].Name() < list[j].Name()
	})

	for f:=0; f < len(list);f++ {
		if list[f].Name()[0:1] == "." && ALL_SHOW == false {
			continue
		}

		err = walkFn(root, list[f], nil)
		if err != nil {
			return err
		}
	}

	return nil;
}

func displayTheFiles(load_path string) error {
	var files []lsFile

	// 一覧取得
	err := getFilesAtDirecotry(load_path, func(path string, info os.FileInfo, err error) error {
		var lll = lsFile{path, info, err}
		files = append(files, lll)
		return nil
	})

	if err != nil {
		return err
	}

	// 表示をする
	for _, file := range files {
		fmt.Print(file.info.Name() + "\n")
	}

	return nil
}

/**
 *
 */
func getWindowSize(fd int) *Size {
	ws := &Window{}
	var err error
	ws.Row, ws.Column, err = terminal.GetSize(fd)
	if err != nil {
		log.Fatal(err)
	}
	return &ws.Size
}

/**
 *
 */
func init_runner() {
	// -a オプション
	flag.BoolVar(&ALL_SHOW, "a", false, "list all files")

	// -l オプション
	flag.BoolVar(&LISTING_OPT, "l", false, "list display")

	flag.Parse()
}

func main() {
	init_runner()

	var look_directory = getLookDirectory()
	for l := 0; l < len(look_directory); l++ {
		var err = displayTheFiles(look_directory[l])
		if err != nil {
			log.Fatal(err)
		}
	}
}
