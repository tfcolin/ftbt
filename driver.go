package ftbt

import (
	"io"
      "os"
)

const (
	TST_NORMAL = -1
	TST_WIN = 0
	TST_LOSS = 1
	TST_QUIT = 2
)

var ui UI 

var intro string

func LoadIntro (fin_name string) {
      buf, _ := os.ReadFile (fin_name)
      intro = string(buf)
}

func Init (ui_note UINote, ui_select [2]UISelect, ftech, fmap io.Reader, fintro_name string) bool {

	ui.note = ui_note
	ui.sel = ui_select

	InitMiniPath ()

	if LoadDX (ftech) {
		return true
	}
	if LoadTech (ftech) {
		return true
	}
      LoadIntro (fintro_name)
	if LoadMap (fmap, ui.note) {
		return true
	}
	if LoadPeople (fmap, ui.note) {
		return true
	}
	return false
}

func DoOneTurn () (status int) {
	StartTurnPeople ()
	MapUpdataAux ()
	for opt := 0; opt < 2; opt ++ {
		MapUpdateWeight (opt)
		ui.sel[opt].TurnStart (opt)
		for {
			MapClearPath(true)
			obj_loc := ui.sel[opt].SelObj (SOS_PEOPLE, -1, nil)
			if obj_loc == -1 {
				is_conf := ui.sel[opt].Confirm("Quit?")
				if is_conf {
					return TST_QUIT 
				} else {
					continue
				}
			}
			if obj_loc == -2 {
				is_conf := ui.sel[opt].Confirm("End turn?")
				if is_conf {
					break
				} else {
					continue
				}
			}
			ip := gmap[obj_loc].people
			if ip != -1 && people[ip].opt == opt {
				DoPeople (ip, ui)
				if np_in[1 - opt] == 0 {
					return opt
				}
				if np_in[opt] == 0 {
					return 1 - opt
				}
				if np_finish[opt] == np_in[opt] {
					break
				}
			}
		}
	}
	return TST_NORMAL
}

func Do () (st int) {
	ui.note.GameStart ()
	for {
		st = DoOneTurn ()
		if st != TST_NORMAL {
			break
		}
	}

	if st == TST_QUIT {
		ui.note.GameOver (-1)
	} else {
		ui.note.GameOver (st)
	}
	return
}
