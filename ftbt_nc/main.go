package main

import (
	"fmt"
	"os"
	"tfgolib/ftbt"
)

func main() {

	if len (os.Args) < 3 {
		fmt.Printf ("Usage: ftbt_nc rule_file map_file\n")
		return
	}

	ftech, err := os.Open (os.Args[1])
	if err != nil {
		panic ("cannot open tech file")
	}
	fmap, err := os.Open (os.Args[2])
	if err != nil {
		panic ("cannot open map file")
	}

	ui_note := ftbt.InitNCursesUI ()
	var ui_sel [2]ftbt.UISelect
	for i := 0; i < 2; i ++ {
		ui_sel[i] = ftbt.InitNCursesUISelect (i, ui_note)
	}

	is_err := ftbt.Init (ui_note, ui_sel, ftech, fmap)
	if is_err {
		fmt.Printf ("初始化错误\n")
		ui_note.End ()
		return
	}
	ftbt.Do()

	ui_note.End ()
}
