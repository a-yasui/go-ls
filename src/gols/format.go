package gols

import (
	"fmt"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"log"
	"os"
	"os/user"
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

// 最も長いファイル名を取得する
func _GetlongestFileName(file []FileStruct) string {
	longest_name := ""

	for _, _file_ := range file {
		if len(longest_name) < len(_file_.info.Name()) {
			longest_name = _file_.info.Name()
		}
	}
	return longest_name
}

// ファイルパーミッションを文字列にして返す
func _GetPermissionString(file os.FileInfo) string {
	return fmt.Sprintf("%s", file.Mode().String())
}

// 一行詳細表示をする
func FormatPrintOneLine(file FileStruct) () {
	var file_size int64
	permission := _GetPermissionString(file.info)
	file_size = file.info.Size()
	owner := ""
	group := ""

	// https://qiita.com/gorilla0513/items/ce0657e2e7de4f46ab2d
	// こんなの知らないとかけないやーん
	if stat, ok := file.info.Sys().(*syscall.Stat_t); ok {
		// syscall.Stat_t.UidがユーザIDになります
		uid := strconv.Itoa(int(stat.Uid))
		// os/userパッケージのLookupIdにUidを渡すとユーザ名を取得できます
		u, err := user.LookupId(uid)
		if err != nil {
			owner = uid
		} else {
			owner = u.Username
		}

		// syscall.Stat_t.GidがグループIDになります
		gid := strconv.Itoa(int(stat.Gid))
		// os/userパッケージのLookupGroupIdにGidを渡すとグループ名を取得できます
		g, err := user.LookupGroupId(gid)
		if err != nil {
			group = gid
		} else {
			group = g.Name
		}
	}

	// ファイルの時は 1
	link_count := 1

	file_time := file.info.ModTime()

	fmt.Printf(
		"%s %d %s %s %4d %s %s\n",
		permission,
		link_count,
		owner,
		group,
		file_size,
		fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			file_time.Year(),
			file_time.Month(),
			file_time.Day(),
			file_time.Hour(),
			file_time.Minute(),
			file_time.Second(),
		),
		file.info.Name(),
	)
}

// ファイル名のみ表示させる
func FormatPrintOnlyNames(file []FileStruct) () {
	// 最長ファイル名から右寄せで表示させる。
	var root string = ""
	longest_name := _GetlongestFileName(file)

	// 無いものは表示できぬ
	if longest_name == "" {
		return
	}

	var term_size = getWindowSize(syscall.Stdout)
	var oneline_word = int((term_size.Column / (len(longest_name) + 1)))
	var word_width = strconv.Itoa(len(longest_name) + 1)
	var need_line = len(file) % oneline_word

	if len(file) < oneline_word {
		need_line = 1
	}

	log.Printf("[debug] 最長文字 %d", len(longest_name))
	log.Printf("[debug] window size %v", term_size)
	log.Printf("[debug] ファイル名の表示に使う横幅 %s", word_width)
	log.Printf("[debug] 一行に表示させるワード数 %d", oneline_word)
	log.Printf("[debug] 使用する行数 %d", need_line)
	log.Print("[debug] Format: " + " %-" + (word_width) + "s ")

	buff := fmt.Sprintf("%s\n", root)
	for i := 0; i < len(file); i++ {
		var _file_ = file[i]
		var format = "%-" + word_width + "s "
		if (i+1)%oneline_word == 0 {
			format = "%s"
		}

		buff += fmt.Sprintf(format, _file_.info.Name())
		if (i+1)%oneline_word == 0 {
			buff += "\n"
			log.Printf("[debug] oneline %d at i{%d}", (i+1)%oneline_word, i)
		}
	}

	fmt.Print(buff)
	fmt.Print("\n")
}
