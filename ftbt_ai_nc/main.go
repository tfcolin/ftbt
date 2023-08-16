package main

import (
	"fmt"
	"os"
	"gitee.com/tfcolin/ftbt"
)

func main() {

	if len (os.Args) < 4 {
		fmt.Printf ("Usage: ftbt_ai_nc rule_file map_file intro_file\n")
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
	ui_sel[0] = ftbt.InitNCursesUISelect (0, ui_note)
	ui_sel[1] = ftbt.InitNCursesUIAISelect (1, ui_note)

	is_err := ftbt.Init (ui_note, ui_sel, ftech, fmap, os.Args[3])
	if is_err {
		fmt.Printf ("初始化错误\n")
		ui_note.End ()
		return
	}
	ftbt.InitAI ()
	is_err = ftbt.LoadAI (ftech, fmap)
	if is_err {
		fmt.Printf ("AI初始化错误\n")
		ui_note.End ()
		return
	}

	/* start trigger UI event */
	ftbt.Do()

	ui_note.End ()
}
