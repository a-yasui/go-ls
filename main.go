package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/yookoala/realpath"
)

/**
 *
 */
func __realpath(path string)(new_path string ) {
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
		return []string{ __realpath(dir) }
	}

	// 引数の中身を正規化
	var argc = len(lsdir)
	var result []string
	for i := 0; i < argc; i++ {
		result = append(result, __realpath(lsdir[i]))
	}
	return lsdir
}

func init_runner () {
	flag.Parse()
}

func main() {
	init_runner()

	fmt.Println(getLookDirectory())
}
