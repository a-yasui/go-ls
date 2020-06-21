package gols

import (
	"fmt"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"log"
	"os"
	"strconv"
	"syscall"
)

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
}

/**
 * 取得した生データ
 * see : https://flaviocopes.com/go-list-files/
 */
type FileStruct struct {
	path string
	info os.FileInfo
	err  error
}

// 表示用に使うデータ構造
func NewFile(path string, info os.FileInfo) *FileStruct {
	result := FileStruct{path: path, info: info}
	return &result
}

/**
 * see : https://qiita.com/t0w4/items/29798299155713f15a83
 */
func getWindowSize(fd int) *Size {
	ws := &Window{}
	Height, _ := terminal.Height()
	Width, _ := terminal.Width()

	ws.Column = int(Width)
	ws.Row = int(Height)
	return &ws.Size
}

// 一行詳細表示をする
func FormatPrintOneLine(display bool, file FileStruct) () {
	fmt.Println(file.info.Name())
}

// ファイル名のみ表示させる
func FormatPrintOnlyNames(hiddenFileDisplay bool, file []FileStruct) () {
	// 最長ファイル名から右寄せで表示させる。
	var longest_name string = ""
	var root string = ""
	for _, _file_ := range file {
		root = _file_.path
		if len(longest_name) < len(_file_.info.Name()) {
			longest_name = _file_.info.Name()
		}
	}

	// 無いものは表示できぬ
	if longest_name == "" {
		return
	}

	var term_size = getWindowSize(syscall.Stdout)
	var oneline_word = int((term_size.Column / (len(longest_name) + 1)))
	var column_width = strconv.Itoa(len(longest_name) + 1)

	log.Printf("[debug] 最長文字 %d", len(longest_name))
	log.Printf("[debug] window size %v", term_size)
	log.Printf("[debug] ファイル名の表示に使う横幅 %s", column_width)
	log.Printf("[debug] 一行に表示させるワード数 %d", oneline_word)
	log.Print("[debug] Format: " + " %-" + (column_width) + "s ")

	fmt.Println(root)
	for i, _file_ := range file {
		fmt.Printf("%-"+(column_width)+"s ", _file_.info.Name())
		if (i + 1)%oneline_word == 0 && i != 0 {
			fmt.Print("\n")
		}
	}
	fmt.Print("\n")
}
