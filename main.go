package main

import (
	"./src/gols"
	"flag"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var ALL_SHOW bool    // -a option
var LISTING_OPT bool // -l option
var DebugOpt bool    // -D option

// TODO: -C
// TODO: -F
// TODO: -R
// TODO: -c
// TODO: -d
// TODO: -i
// TODO: -q
// TODO: -r
// TODO: -t
// TODO: -u
// TODO: -1




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

func getFilesAtDirecotry(root string, walkFn filepath.WalkFunc) error {
	f, err := os.Open(root)
	if err != nil {
		return err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	// Sortable
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})

	for f := 0; f < len(list); f++ {
		if list[f].Name()[0:1] == "." && ALL_SHOW == false {
			continue
		}

		err = walkFn(root, list[f], nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func displayTheFiles(load_path string) error {
	var files []gols.FileStruct

	// 一覧取得
	err := getFilesAtDirecotry(load_path, func(path string, info os.FileInfo, err error) error {
		lll := gols.NewFile(path, info)
		files = append(files, *lll)
		return nil
	})

	if err != nil {
		return err
	}

	// 表示をする
	if LISTING_OPT {
		for i := 0; i < len(files); i++ {
			gols.FormatPrintOneLine(ALL_SHOW, files[i])
		}
	} else {
		gols.FormatPrintOnlyNames(ALL_SHOW, files)
	}

	return nil
}

/**
 *
 */
func init_runner() {
	// -a オプション
	flag.BoolVar(&ALL_SHOW, "a", false, "list all files")

	// -l オプション
	flag.BoolVar(&LISTING_OPT, "l", false, "list display")

	// -D オプション :: Debug
	flag.BoolVar(&DebugOpt, "D", false, "Debug Option")
	flag.Parse()


	if DebugOpt == false {
		// ログを黙らせる
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	init_runner()

	look_directory := getLookDirectory()
	for l := 0; l < len(look_directory); l++ {
		var err = displayTheFiles(look_directory[l])
		if err != nil {
			log.Fatal(err)
		}
	}
}
