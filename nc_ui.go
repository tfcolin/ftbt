package ftbt

import (
	"log"
	"fmt"
	"gitee.com/tfcolin/goncurses"
)

const (
      NUM_WTECH int = 10
)

type NCursesUI struct {

	is_first_msg bool 
      opt_name [2]string 

	w_all  * goncurses.Window
	w_map  * goncurses.Window
	p_map  * goncurses.Pad
	w_base * goncurses.Window
	w_aux  * goncurses.Window
	w_atech * goncurses.Window
	w_scroll * goncurses.Window
      w_msg * goncurses.Window

	w_dx 		* goncurses.Window
	w_weight 	* goncurses.Window
	w_people 	* goncurses.Window
	w_opt 	* goncurses.Window
	w_hp 		* goncurses.Window
	w_is_finish * goncurses.Window
	w_abi		* goncurses.Window
	w_def		* goncurses.Window
	w_speed	* goncurses.Window
	w_buf_v [2][2]* goncurses.Window
	w_buf_d [2][2]* goncurses.Window
	w_fly_v [2]	* goncurses.Window
	w_fly_d [2]	* goncurses.Window
	w_dl_v [2]	* goncurses.Window
	w_dl_d [2]	* goncurses.Window
	w_tech [NUM_WTECH] * goncurses.Window
	w_sb   [2]  * goncurses.Window
	w_msg_in    * goncurses.Window

	wl_dx             * goncurses.Window
	wl_weight         * goncurses.Window
	wl_people         * goncurses.Window
	wl_opt 		* goncurses.Window
	wl_hp 		* goncurses.Window
	wl_is_finish 	* goncurses.Window
	wl_abi		* goncurses.Window
	wl_def		* goncurses.Window
	wl_speed		* goncurses.Window
	wl_aux_opt [2]    * goncurses.Window
	wl_aux_val [2]    * goncurses.Window
	wl_aux_dur [2]    * goncurses.Window
	wl_buf [2]   	* goncurses.Window
	wl_fly 		* goncurses.Window
	wl_dl 		* goncurses.Window
	wl_tech           * goncurses.Window
	wl_sb   [2]       * goncurses.Window

      sc_base[2]           int
      sc_win [2]           int
      sc_len [2]           int

      cx, cy         int
}

type NCursesUISelect struct {
      * NCursesUI
      opt int
}

type NCursesUIAISelect struct {
	* NCursesUISelect
	/* current people id */
	cp int
}

func InitNCursesUISelect (opt int, nui * NCursesUI) * NCursesUISelect {
      var nuis NCursesUISelect
      nuis.NCursesUI = nui
      nuis.opt = opt
      return &nuis
}

func InitNCursesUIAISelect (opt int, nui * NCursesUI) * NCursesUIAISelect {
      var nuis NCursesUIAISelect
	nuis.NCursesUISelect = InitNCursesUISelect (opt, nui)
      return &nuis
}

func (nui * NCursesUI) End () {
	goncurses.End ()
}

func InitNCursesUI () * NCursesUI {
	var nui NCursesUI

	nui.opt_name = [2]string{"OptA", "OptB"}
	nui.is_first_msg = true

      var err error
	/* init ncurses */
	nui.w_all, err = goncurses.Init("zh_CN.UTF-8")
	// nui.w_all, err = goncurses.Init()
	if err != nil {
		goncurses.End()
		log.Fatal ("NCUI Init error.\n")
		return nil
	}

	goncurses.Cursor (0)
	goncurses.Echo (false)
	goncurses.MouseMask (goncurses.M_ALL, nil)
	has_color := goncurses.HasColors()
	cc_color := goncurses.CanChangeColor ()
	goncurses.StartColor ()
	n_color := goncurses.Colors ()
	if !has_color || !cc_color || n_color < 16 {
		goncurses.End()
		log.Fatal ("This terminal do not have enough color support.")
		return nil
	}

	/* init colors */
	goncurses.InitColor (0, 0, 0, 0)
	goncurses.InitColor (1, 1000, 0, 0)
	goncurses.InitColor (2, 0, 1000, 0)
	goncurses.InitColor (3, 0, 0, 1000)
	goncurses.InitColor (4, 1000, 1000, 0)
	goncurses.InitColor (5, 1000, 0, 1000)
	goncurses.InitColor (6, 0, 1000, 1000)
	goncurses.InitColor (7, 1000, 1000, 1000)
	goncurses.InitColor (8, 1000, 500, 0)
	goncurses.InitColor (9, 400, 400, 400)
	goncurses.InitColor (10, 0, 500, 1000)
	goncurses.InitColor (11, 200, 300, 400)
	goncurses.InitColor (12, 1000, 600, 600)
	goncurses.InitColor (13, 600, 600, 1000)

	/* white/black */
	goncurses.InitPair (0, 7, 0)

	/* show people and aux */
	goncurses.InitPair (1, 12, 0)
	goncurses.InitPair (2, 13, 0)
	goncurses.InitPair (3, 12, 9)
	goncurses.InitPair (4, 13, 9)
	goncurses.InitPair (5, 12, 11)
	goncurses.InitPair (6, 13, 11)

	/* show scroll bar */
	goncurses.InitPair (7, 0, 8)
	goncurses.InitPair (8, 0, 10)

	/* show dx */
	goncurses.InitPair (9, 3, 4) // PL
	goncurses.InitPair (10, 1, 6) // GR
	goncurses.InitPair (11, 5, 2) // FO
	goncurses.InitPair (12, 2, 5) // MO
	goncurses.InitPair (13, 4, 3) // CI
	goncurses.InitPair (14, 6, 1) // WA

	/* set background */
	nui.w_all.Color (0)
	/* enable special keys */
	nui.w_all.Keypad (true)
	/* Refresh */
	nui.w_all.Refresh ()

	ny, nx := nui.w_all.MaxYX ()
	if ny < 32 || nx < 164 {
		goncurses.End()
		log.Fatal ("Not enough space. Please use a larger terminal.")
		return nil
	}

      /* create all sub windows */
	nui.w_map, _ = goncurses.NewWindow (32, 122, 0, 0)
	nui.w_base, _ = goncurses.NewWindow (7, 42, 0, 122)
	nui.w_aux, _ = goncurses.NewWindow (8, 42, 7, 122)
	nui.w_atech, _ = goncurses.NewWindow (7, 42, 15, 122)
	nui.w_scroll, _ = goncurses.NewWindow (4, 42, 22, 122)
      nui.w_msg, _ = goncurses.NewWindow (6, 42, 26, 122)
	nui.w_msg_in, _ = goncurses.NewWindow (4, 40, 27, 123)
      nui.w_msg_in.ScrollOk (true)

	nui.wl_dx, _ = goncurses.NewWindow       (1, 12, 1, 123)
	nui.w_dx, _ = goncurses.NewWindow 	 (1, 8, 1, 135)
	nui.wl_weight, _ = goncurses.NewWindow   (1, 12, 1, 143)
	nui.w_weight, _ = goncurses.NewWindow    (1, 8, 1, 155)

	nui.wl_people, _ = goncurses.NewWindow   (1, 12, 2, 123)
	nui.w_people, _ = goncurses.NewWindow    (1, 8, 2, 135)
	nui.wl_opt, _ = goncurses.NewWindow      (1, 12, 2, 143)
	nui.w_opt, _ = goncurses.NewWindow       (1, 8, 2, 155)

	nui.wl_hp, _ = goncurses.NewWindow 		(1, 12, 3, 123)
	nui.w_hp, _ = goncurses.NewWindow             (1, 8, 3, 135)
	nui.wl_is_finish, _ = goncurses.NewWindow     (1, 12, 3, 143)
	nui.w_is_finish, _ = goncurses.NewWindow      (1, 8, 3, 155)

	nui.wl_abi, _ = goncurses.NewWindow 	(1, 12, 4, 123)
	nui.w_abi, _ = goncurses.NewWindow      (1, 8, 4, 135)
	nui.wl_def, _ = goncurses.NewWindow     (1, 12, 4, 143)
	nui.w_def, _ = goncurses.NewWindow      (1, 8, 4, 155)

	nui.wl_speed, _ = goncurses.NewWindow  (1, 12, 5, 123)
	nui.w_speed, _ = goncurses.NewWindow   (1, 8, 5, 135)
                                         
      nui.wl_aux_opt[0], _ = goncurses.NewWindow   (1, 4, 8, 123 + 8 + 6)
      nui.wl_aux_opt[1], _ = goncurses.NewWindow   (1, 4, 8, 123 + 8 * 3 + 6)
      nui.wl_aux_val[0], _ = goncurses.NewWindow   (1, 8, 9, 123 + 8)
      nui.wl_aux_dur[0], _ = goncurses.NewWindow   (1, 8, 9, 123 + 8 * 2)
      nui.wl_aux_val[1], _ = goncurses.NewWindow   (1, 8, 9, 123 + 8 * 3)
      nui.wl_aux_dur[1], _ = goncurses.NewWindow   (1, 8, 9, 123 + 8 * 4)

      nui.wl_buf[0], _ = goncurses.NewWindow   (1, 8, 10, 123)
      nui.wl_buf[1], _ = goncurses.NewWindow   (1, 8, 11, 123)
      nui.wl_fly, _    = goncurses.NewWindow   (1, 8, 12, 123)
      nui.wl_dl, _     = goncurses.NewWindow   (1, 8, 13, 123)

      for i := 0; i < 2; i ++ {
            for j := 0; j < 2; j ++ {
                  nui.w_buf_v[i][j], _ = goncurses.NewWindow   (1, 8, 10 + j, 123 + 8 + 16 * i)
                  nui.w_buf_d[i][j], _ = goncurses.NewWindow   (1, 8, 10 + j, 123 + 16 + 16 * i)
            }
            nui.w_fly_v[i], _ = goncurses.NewWindow   (1, 8, 12, 123 + 8 + 16 * i)
            nui.w_fly_d[i], _ = goncurses.NewWindow   (1, 8, 12, 123 + 16 + 16 * i)
            nui.w_dl_v[i], _ = goncurses.NewWindow    (1, 8, 13, 123 + 8 + 16 * i)
            nui.w_dl_d[i], _ = goncurses.NewWindow    (1, 8, 13, 123 + 16 + 16 * i)
      }

      nui.wl_tech, _ = goncurses.NewWindow   (1, 10, 16, 123)
	nsh := (NUM_WTECH - 1) / 3 + 1
      for i := 0; i < nsh; i ++ {
            for j := 0; j < 3; j ++ {
                  ind := i * 3 + j
			if ind < NUM_WTECH {
				nui.w_tech[ind], _ = goncurses.NewWindow   (1, 10, 16 + i, 123 + 10 + 10 * j)
			}
            }
      }

      for i := 0; i < 2; i ++ {
            nui.wl_sb[i], _ = goncurses.NewWindow  (1, 10, 23 + i, 123)
            nui.w_sb[i], _  = goncurses.NewWindow  (1, 30, 23 + i, 123 + 10)
      }

      /* draw window borders */
      nui.w_map.Border ('|', '|', '-', '-', '*', '*', '*', '*')
      nui.w_base.Border ('|', '|', '-', '-', '*', '*', '*', '*')
      nui.w_aux.Border ('|', '|', '-', '-', '*', '*', '*', '*')
      nui.w_atech.Border ('|', '|', '-', '-', '*', '*', '*', '*')
      nui.w_scroll.Border ('|', '|', '-', '-', '*', '*', '*', '*')
      nui.w_msg.Border ('|', '|', '-', '-', '*', '*', '*', '*')

      nui.w_map.Refresh()
      nui.w_base.Refresh()
      nui.w_aux.Refresh()
      nui.w_atech.Refresh()
      nui.w_scroll.Refresh()
      nui.w_msg.Refresh()

      /* write labels */
      nui.wl_dx.Printf ("%-12s", "Ground")
      nui.wl_weight.Printf ("%-12s", "Weight")
      nui.wl_people.Printf ("%-12s", "Role")
      nui.wl_opt.Printf ("%-12s", "Control")
      nui.wl_hp.Printf ("%-12s", "HP")
      nui.wl_is_finish.Printf ("%-12s", "Finished")
      nui.wl_abi.Printf ("%-12s", "Ability")
      nui.wl_def.Printf ("%-12s", "Defence")
      nui.wl_speed.Printf ("%-12s", "Speed")
      nui.wl_aux_opt[0].Printf ("%4s", "OptA")
      nui.wl_aux_opt[1].Printf ("%4s", "OptB")
      for i := 0; i < 2; i ++ {
            nui.wl_aux_val[i].Printf ("%-8s", " Value")
            nui.wl_aux_dur[i].Printf ("%-8s", "  Dur.")
      }
      nui.wl_buf[0].Printf ("%-8s", "Abi.")
      nui.wl_buf[1].Printf ("%-8s", "Def.")
      nui.wl_fly.Printf ("%-8s", "Fly")
      nui.wl_dl.Printf ("%-8s", "Dis.L")
      nui.wl_tech.Printf ("%-10s", "  Tech")

      nui.wl_sb[0].Printf ("%-10s", "XScroll")
      nui.wl_sb[1].Printf ("%-10s", "YScroll")

      nui.wl_dx.Refresh()           
      nui.wl_weight.Refresh()       
      nui.wl_people.Refresh()       
      nui.wl_opt.Refresh() 	    
      nui.wl_hp.Refresh() 	    
      nui.wl_is_finish.Refresh()    
      nui.wl_abi.Refresh()	    
      nui.wl_def.Refresh()	    
      nui.wl_speed.Refresh()	    
      for i := 0; i < 2; i ++ {
            nui.wl_aux_opt[i].Refresh()  
            nui.wl_aux_val[i].Refresh()  
            nui.wl_aux_dur[i].Refresh()  
            nui.wl_buf[i].Refresh()      
            nui.wl_sb[i].Refresh()     
      }
      nui.wl_fly.Refresh() 	    
      nui.wl_dl.Refresh() 	    
      nui.wl_tech.Refresh()         

	return &nui                        
}

func (nui * NCursesUI) ShowMsg (msg string) {
	if nui.is_first_msg {
		nui.is_first_msg = false
	} else {
		nui.w_msg_in.Printf ("\n")
	}
      nui.w_msg_in.Printf ("%s", msg)
	nui.w_msg_in.Refresh ()
}

func (nui * NCursesUI) UpdateScroll () {

      /* adjust sc_base */
      for i := 0; i < 2; i ++ {
            if nui.sc_base[i] < 0 {
                  nui.sc_base[i] = 0
            }
            if nui.sc_base[i] + nui.sc_win[i] > nui.sc_len[i] {
                  nui.sc_base[i] = nui.sc_len[i] - nui.sc_win[i]
            }
      }

      /* draw scroll bar */
      var sb_len int = 30
      var sb_s, sb_e int
      for i := 0; i < 2; i ++ {
            sb_s = sb_len * nui.sc_base[i] / nui.sc_len[i]
            sb_e = sb_len * (nui.sc_base[i] + nui.sc_win[i]) / nui.sc_len[i] 
            nui.w_sb[i].ColorOn (7)
            for j := 0; j < sb_s; j ++ {
                  nui.w_sb[i].MoveAddChar (0, j, ' ')
            }
            nui.w_sb[i].ColorOff (7)
            nui.w_sb[i].ColorOn (8)
            for j := sb_s; j < sb_e; j ++ {
                  nui.w_sb[i].MoveAddChar (0, j, ' ')
            }
            nui.w_sb[i].ColorOff (8)
            nui.w_sb[i].ColorOn (7)
            for j := sb_e; j < sb_len; j ++ {
                  nui.w_sb[i].MoveAddChar (0, j, ' ')
            }
            nui.w_sb[i].ColorOff (7)
            nui.w_sb[i].Refresh ()
      }
                  
      nui.p_map.Refresh (nui.sc_base[1] * 2, nui.sc_base[0] * 6, 1, 1, nui.sc_win[1] * 2, nui.sc_win[0] * 6)
}

func (nui * NCursesUI) DrawInfo () {
      iloc := nui.cy * MAP_SIZE[0] + nui.cx
      loc := gmap[iloc]

      /* show dx and weight */
      nui.w_dx.MovePrintf (0, 0, "%-8s", dx_name[loc.dx])
      if loc.weight >= 0 {
            nui.w_weight.MovePrintf (0, 0, "%-8d", loc.weight)
      } else {
            nui.w_weight.Erase ()
      }

      /* show aux */
      for i := 0; i < 2; i ++ {
            for j := 0; j < 2; j ++ {
                  nui.w_buf_v[i][j].MovePrintf (0, 0, "%-8d", loc.buf_value[i][j])
                  nui.w_buf_d[i][j].MovePrintf (0, 0, "%-8d", loc.buf_dur[i][j])
                  nui.w_buf_v[i][j].Refresh ()
                  nui.w_buf_d[i][j].Refresh ()
            }
            nui.w_fly_v[i].MovePrintf (0, 0, "%-8d", loc.fly_weight[i])
            nui.w_fly_d[i].MovePrintf (0, 0, "%-8d", loc.fly_weight_dur[i])
            nui.w_dl_v[i].MovePrintf (0, 0, "%-8v", loc.disable_long[i])
            nui.w_dl_d[i].MovePrintf (0, 0, "%-8d", loc.disable_long_dur[i])
            nui.w_fly_v[i].Refresh ()
            nui.w_fly_d[i].Refresh ()
            nui.w_dl_v[i].Refresh ()
            nui.w_dl_d[i].Refresh ()
      }

      /* show people */
      if loc.people >= 0 {
            p := people[loc.people]
            nui.w_people.MovePrintf (0, 0, "%-8s", p.name)
            nui.w_opt.MovePrintf (0, 0, "%-8s", nui.opt_name[p.opt])
            nui.w_hp.MovePrintf (0, 0, "%-8d", p.hp)
            nui.w_is_finish.MovePrintf (0, 0, "%-8t", p.is_finish)
            nui.w_abi.MovePrintf (0, 0, "%-8d", p.abi)
            nui.w_def.MovePrintf (0, 0, "%-8d", p.def)
            nui.w_speed.MovePrintf (0, 0, "%-8d", p.speed)
            for i, it := range p.tech {
                  if i < NUM_WTECH {
                        nui.w_tech[i].MovePrintf (0, 0, "%-10s", gtech[it].name)
                  }
            }
            for i := len(p.tech) ; i < NUM_WTECH; i ++ {
                  nui.w_tech[i].Erase ()
            }
      } else {
            nui.w_people.Erase ()
            nui.w_opt.Erase ()
            nui.w_hp.Erase ()
            nui.w_is_finish.Erase ()
            nui.w_abi.Erase ()
            nui.w_def.Erase ()
            nui.w_speed.Erase ()
            for i := 0; i < NUM_WTECH; i ++ {
                  nui.w_tech[i].Erase ()
            }
      }

      nui.w_dx.Refresh ()
      nui.w_weight.Refresh ()
      nui.w_people.Refresh ()
      nui.w_opt.Refresh ()
      nui.w_hp.Refresh ()
      nui.w_is_finish.Refresh ()
      nui.w_abi.Refresh ()
      nui.w_def.Refresh ()
      nui.w_speed.Refresh ()
      for i := 0; i < NUM_WTECH; i ++ {
            nui.w_tech[i].Refresh ()
      }

}

func (nui * NCursesUI) DrawMapStatic () {
      for i := 0; i < MAP_SIZE[1]; i ++ {
            by := i * 2
            for j := 0; j < MAP_SIZE[0]; j ++ {
                  iloc := i * MAP_SIZE[0] + j
                  loc := gmap[iloc]
                  bx := j * 6
                  nui.p_map.ColorOn (int16(9 + loc.dx))
                  nui.p_map.MovePrint (by, bx, dx_name[loc.dx])
			nui.p_map.ColorOff (int16(9 + loc.dx))
		}
	}
}

func (nui * NCursesUI) DrawGrid (x, y int, hp_down int) {
      loc := gmap[y * MAP_SIZE[0] + x]
      by := y * 2
      bx := x * 6
      var ic int
      if x == nui.cx && y == nui.cy {
            ic = 5
      } else {
            if loc.path_last_loc == -2 {
                  ic = 3
            } else {
                  ic = 1
            }
      }

      for opt := 0; opt < 2; opt ++ {
            nui.p_map.ColorOn (int16(ic + opt))
		if hp_down >= 0 {
			hp_down = hp_down % 100
			nui.p_map.MovePrintf (by, bx +2, "%2d", hp_down)
		} else {
			if loc.people != -1 && people[loc.people].opt == opt {
				nui.p_map.MoveAddChar (by, bx + 2 + opt, 'P')
			} else {
				nui.p_map.MoveAddChar (by, bx + 2 + opt, ' ')
			}
		}
            if loc.buf_dur[opt][0] != 0 {
                  nui.p_map.MoveAddChar (by + 1, bx + opt, 'A')
            } else {
                  nui.p_map.MoveAddChar (by + 1, bx + opt, ' ')
            }
            if loc.buf_dur[opt][1] != 0 {
                  nui.p_map.MoveAddChar (by + 1, bx + 2 + opt, 'D')
            } else {
                  nui.p_map.MoveAddChar (by + 1, bx + 2 + opt, ' ')
            }
            if loc.fly_weight_dur[opt] != 0 {
                  nui.p_map.MoveAddChar (by, bx + 4 + opt, 'F')
            } else {
                  nui.p_map.MoveAddChar (by, bx + 4 + opt, ' ')
            }
            if loc.disable_long_dur[opt] != 0 {
                  nui.p_map.MoveAddChar (by + 1, bx + 4 + opt, 'S')
            } else {
                  nui.p_map.MoveAddChar (by + 1, bx + 4 + opt, ' ')
            }
            nui.p_map.ColorOff (int16(ic + opt))
      }
}

func (nui * NCursesUI) Update () {
      for i := 0; i < MAP_SIZE[1]; i ++ {
            for j := 0; j < MAP_SIZE[0]; j ++ {
                  nui.DrawGrid (j, i, -1)
            }
      }
      nui.DrawInfo ()
      nui.UpdateScroll ()
}

func (nui * NCursesUI) ShowGrid (x, y int) {
      if x >= nui.sc_base[0] && x < nui.sc_base[0] + nui.sc_win[0] && 
	y >= nui.sc_base[1] && y < nui.sc_base[1] + nui.sc_win[1] {
            bx := (x - nui.sc_base[0]) * 6
            by := (y - nui.sc_base[1]) * 2
            nui.p_map.Refresh (y * 2, x * 6, 1 + by, 1 + bx, by + 2, bx + 6)
      }
}

func (nui * NCursesUI) UpdateCursor (cx, cy int, hp_down int) {
	ocx, ocy := nui.cx, nui.cy
	nui.cx, nui.cy = cx, cy
      nui.DrawGrid (cx, cy, hp_down)
	nui.ShowGrid (cx, cy)
	if ocx != cx || ocy != cy {
		nui.DrawGrid (ocx, ocy, -1)
		nui.ShowGrid (ocx, ocy)
	}
      nui.DrawInfo ()
}

func (nui * NCursesUI) LoadMapDone () {
      nui.p_map, _ = goncurses.NewPad (MAP_SIZE[1] * 2, MAP_SIZE[0] * 6)

      nui.sc_win[0] = 20
      nui.sc_win[1] = 15
      for i := 0; i < 2; i ++ {
            nui.sc_len[i] = MAP_SIZE[i]
            nui.sc_base[i] = 0
            nui.cx, nui.cy = 0, 0
		if nui.sc_len[i] < nui.sc_win[i] {
			nui.sc_win[i] = nui.sc_len[i]
		}
      }

      nui.DrawMapStatic ()
      nui.Update ()
}

func (nui * NCursesUI) GameStart () {
      nui.ShowMsg ("游戏开始")
}

func (nui * NCursesUI) TurnStart (opt int) {
      msg := fmt.Sprintf ("%s turn start. There are total %d active roles in battle field.", 
      nui.opt_name[opt], np_in[opt])
      nui.ShowMsg (msg)
}

func (nuis * NCursesUISelect) TurnStart (opt int) {
	nuis.NCursesUI.TurnStart (opt)
}

func (nui * NCursesUI) GameOver (win_opt int) {
	var msg string = "Game over."
	if win_opt >= 0 {
		msg += fmt.Sprintf (" %s win.", nui.opt_name[win_opt])
	}
	nui.ShowMsg (msg)
	nui.ShowMsg ("Press any key to exit")
	nui.w_all.GetChar ()
}

func (nui * NCursesUI) DoTech (itech, p_src, p_obj, l_obj int) {
      msg := fmt.Sprintf ("%s use tech %s", people[p_src].name, gtech[itech].name) 
      if p_obj != -1 {
            msg += fmt.Sprintf (" on %s", people[p_obj].name) 
      }
	msg += ". Press any key to continue"
      nui.ShowMsg (msg)
	x, y := Loc2XY (people[p_obj].loc)
	nui.UpdateCursor (x, y, -1)
	nui.w_all.GetChar ()
}

func (nui * NCursesUI) HPDown (ip, loc, hp, dhp int) {
      x, y := Loc2XY (loc)
      msg := fmt.Sprintf ("%s at location (%d, %d) drop hp %d from %d, press any key to continue", 
      people[ip].name, x, y, dhp, hp)
      nui.ShowMsg (msg)
	nui.UpdateCursor (x, y, dhp)
	nui.w_all.GetChar ()
}

func (nui * NCursesUI) PeopleOut (ip, loc int) {
      x, y := Loc2XY (loc)
      msg := fmt.Sprintf ("%s evacuate from location (%d, %d)", people[ip].name, x, y)
      nui.ShowMsg (msg)

}

func (nui * NCursesUI) PeopleIn (ip, loc int) {
}
func (nui * NCursesUI) PeopleMove (ip, loc_start, loc_end int) {
}
func (nui * NCursesUI) PeopleStep (ip, loc, dir int) {
}
func (nui * NCursesUI) LocBuf (loc, att_def, val, dur int) {
}
func (nui * NCursesUI) LocFly (loc, weight, dur int) {
}
func (nui * NCursesUI) LocDisableLong (loc, dur int) {
}

func next_p (ip int) int {
	ip = ip + 1
	if ip >= len(people) {
		ip = 0
	}
	return ip
}

func (nuis * NCursesUISelect) BrowseOperation (c goncurses.Key, cur_p * int) {
	switch {
	case c == 'l' || c == goncurses.KEY_RIGHT:
		if nuis.cx < MAP_SIZE[0] - 1 {
			nuis.UpdateCursor (nuis.cx + 1, nuis.cy, -1)
		}
	case c == 'h' || c == goncurses.KEY_LEFT:
		if nuis.cx > 0 {
			nuis.UpdateCursor (nuis.cx - 1, nuis.cy, -1)
		}
	case c == 'k' || c == goncurses.KEY_UP:
		if nuis.cy > 0 {
			nuis.UpdateCursor (nuis.cx, nuis.cy - 1, -1)
		}
	case c == 'j' || c == goncurses.KEY_DOWN:
		if nuis.cy < MAP_SIZE[1] - 1 {
			nuis.UpdateCursor (nuis.cx, nuis.cy + 1, -1)
		}
	case c == 'w':
		nuis.sc_base[0] += nuis.sc_win[0] / 2
		nuis.UpdateScroll ()
		if nuis.cx < nuis.sc_base[0] {
			nuis.UpdateCursor(nuis.sc_base[0], nuis.cy, -1)
		}
	case c == 'b':
		nuis.sc_base[0] -= nuis.sc_win[0] / 2
		nuis.UpdateScroll ()
		if nuis.cx >= nuis.sc_base[0] + nuis.sc_win[0] {
			nuis.UpdateCursor(nuis.sc_base[0] + nuis.sc_win[0] - 1, nuis.cy, -1)
		}
	case c == 'd' || c == 4: /* 4 == Ctrl-D */
		nuis.sc_base[1] += nuis.sc_win[1] / 2
		nuis.UpdateScroll ()
		if nuis.cy < nuis.sc_base[1] {
			nuis.UpdateCursor(nuis.cx, nuis.sc_base[1], -1)
		}
	case c == 'u' || c == 21: /* 21 == Ctrl-U */
		nuis.sc_base[1] -= nuis.sc_win[1] / 2
		nuis.UpdateScroll ()
		if nuis.cy >= nuis.sc_base[1] + nuis.sc_win[1] {
			nuis.UpdateCursor(nuis.cx, nuis.sc_base[1] + nuis.sc_win[1] - 1, -1)
		}
	case c == '\t':
		if np_in[nuis.opt] == 0 {
			return 
		}

		var p * People
		for {
			*cur_p = next_p (*cur_p)
			p = &people[*cur_p]
			if p.loc != -1 && p.opt == nuis.opt && !p.is_finish {
				break
			}
		}

		x, y := Loc2XY (p.loc)
		nuis.UpdateCursor (x, y, -1)
		nuis.CenterCursor ()
		nuis.UpdateScroll()
	}

	if nuis.cx < nuis.sc_base[0] {
		nuis.sc_base[0] -= nuis.sc_win[0] / 2
		nuis.UpdateScroll ()
	}
	if nuis.cx >= nuis.sc_base[0] + nuis.sc_win[0] {
		nuis.sc_base[0] += nuis.sc_win[0] / 2
		nuis.UpdateScroll ()
	}
	if nuis.cy < nuis.sc_base[1] {
		nuis.sc_base[1] -= nuis.sc_win[1] / 2
		nuis.UpdateScroll ()
	}
	if nuis.cy >= nuis.sc_base[1] + nuis.sc_win[1] {
		nuis.sc_base[1] += nuis.sc_win[1] / 2
		nuis.UpdateScroll ()
	}
}

func (nuis * NCursesUISelect) SelObj (status SelObjStatus, itech int, valid_loc []int) int {
      nuis.Update ()

	var msg string
	switch status {
	case SOS_PEOPLE:
		msg = "select people to operate"
	case SOS_MOVE_OBJ:
		msg = "select moving object location"
	case SOS_TECH_OBJ:
		msg = fmt.Sprintf ("select object for tech %s", gtech[itech].name)
	}
	nuis.ShowMsg (msg)

      var cur_p int = len(people) - 1

      gc: for {
            c := nuis.w_all.GetChar ()
            switch c {
            case 'q':
                  return -1
		case 't':
			return -2
            case '\n':
                  break gc
		}
		nuis.BrowseOperation (c, &cur_p)
      }

      return XY2Loc (nuis.cx, nuis.cy)
}

func (nui * NCursesUI) CenterCursor () {
	nui.sc_base[0] = nui.cx - nui.sc_win[0] / 2
      nui.sc_base[1] = nui.cy - nui.sc_win[1] / 2
}

func (nuis * NCursesUISelect) SelTech (ip int) int {
      p := people[ip]
      if p.loc == -1 {
            return -1
      }

	msg := fmt.Sprintf ("Select tech (0 - %d)", len(p.tech) - 1)
	nuis.ShowMsg (msg) 
	nuis.CenterCursor ()
      nuis.Update ()

      var cur_p int = len(people) - 1

	var it int
      for {
            c := nuis.w_all.GetChar ()
		if c == 'q' {
			return -1
		}
            it = int (c - '0')
		if it >= 0 && it < len(p.tech) {
			break
		}
		nuis.BrowseOperation (c, &cur_p)

      }

	return it
}

func (nuis * NCursesUISelect) Confirm (msg string) bool {
	nuis.ShowMsg (msg)
	for {
		c := nuis.w_all.GetChar ()
		if c == 'y' || c == 'Y' {
			return true
		}
		if c == 'n' || c == 'N' {
			return false
		}
	}
}

func (nuis * NCursesUIAISelect) TurnStart (opt int) {
	nuis.NCursesUI.TurnStart (opt)
	nuis.cp = len(people) - 1
	AICalAllDist ()
}

/* return -1 to quit, -2 to end turn (only for SOS_PEOPLE) */
func (nuis * NCursesUIAISelect)	SelObj (status SelObjStatus, itech int, valid_loc []int) int {
	var msg string
	var res int = -1
      var cur_p int

	cur_p = nuis.cp
	switch status {
	case SOS_PEOPLE:
		var p * People
		for {
			nuis.cp = next_p (nuis.cp)
			p = &people[nuis.cp]
			if p.opt == nuis.opt && p.loc != -1 {
				break
			}
		}
		cur_p = nuis.cp
		res = p.loc
		msg = fmt.Sprintf ("people %s action by AI. press c to continue.", p.name)
		nuis.cx, nuis.cy = Loc2XY (p.loc)
		nuis.CenterCursor ()
	case SOS_MOVE_OBJ:
		AIMakeDecision (nuis.cp, valid_loc)
		msg = fmt.Sprintf ("people %s move by AI. press c to continue.", people[nuis.cp].name)
		res = ai_people[nuis.cp].opt_move_loc
		if res != -1 {
			nuis.cx, nuis.cy = Loc2XY (res)
		} 
	case SOS_TECH_OBJ:
		p := people[nuis.cp]
		msg = fmt.Sprintf ("people %s do tech %s by AI. press c to continue.", p.name, gtech[itech].name)
		res = ai_people[nuis.cp].opt_tech_obj
		nuis.cx, nuis.cy = Loc2XY (p.loc)
	}

	nuis.Update ()
	nuis.ShowMsg (msg)

      for {
            c := nuis.w_all.GetChar ()
		if c == 'c' {
			break
		}
		nuis.BrowseOperation (c, &cur_p)
      }

	return res
}

/* return -1 to give up tech use */
func (nuis * NCursesUIAISelect)	SelTech (ip int) int {
	return ai_people[ip].opt_itech
}

/* return true to confirm quit */
func (nuis * NCursesUIAISelect)	Confirm (msg string) bool {
	return true
}
