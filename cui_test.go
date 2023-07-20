package ftbt

import (
	"testing"
	"fmt"
	"os"
)

type TestUI struct {
	opt int
}

var test_ui  [2]TestUI
var test_uis [2]UISelect

/* DoTech triger before HPDown, Loc.. */
func (ui TestUI) DoTech (itech, p_src, p_obj, l_obj int) {
	x, y := Loc2XY (l_obj)
	fmt.Printf ("%s 发动 %s 技能, 目标位置 (%d, %d)", people[p_src].name, gtech[itech].name, x, y) 
	if p_obj != -1 {
		fmt.Printf (", 目标人物 %s", people[p_obj].name)
	}
	fmt.Printf ("\n")
}

func (ui TestUI) HPDown (ip, loc, hp, dhp int) {
	x, y := Loc2XY (loc)
	fmt.Printf ("位于 (%d, %d) 的 %s 体力从 %d 下降 %d\n", x, y, people[ip].name, hp, dhp)
}

func (ui TestUI) PeopleIn (ip, loc int) {
	x, y := Loc2XY (loc)
	fmt.Printf ("%s 加入战斗位置 (%d, %d)\n", people[ip].name, x, y)
}

func (ui TestUI) PeopleOut (ip, loc int) {
	x, y := Loc2XY (loc)
	fmt.Printf ("%s 从位置 (%d, %d) 撤退\n", people[ip].name, x, y)
}

/* PeopleMove trigger after PeopleStep */
func (ui TestUI) PeopleMove (ip, loc_start, loc_end int) {
	x1, y1 := Loc2XY (loc_start)
	x2, y2 := Loc2XY (loc_end)
	name := people[ip].name

	fmt.Printf ("%s 从位置 (%d, %d) 移动到位置 (%d, %d)\n", name, x1, y1, x2, y2)
}

func (ui TestUI) PeopleStep (ip, loc, dir int) {
	dir_str := []string {"右", "左", "下", "上"}
	name := people[ip].name
	x, y := Loc2XY (loc)
	fmt.Printf ("%s 从位置 (%d, %d) 向 %s 移动一步\n", name, x, y, dir_str[dir])
}

/* att_def: 0: att; 1: def */
func (ui TestUI) LocBuf (loc, att_def, val, dur int) {
	ad_str := []string{"攻击", "防御"}
	x, y := Loc2XY (loc)
	fmt.Printf ("位置 (%d, %d) 获得 %s 辅助值 %d, 辅助时长 %d\n", x, y, ad_str[att_def], val, dur)
}

func (ui TestUI) LocFly (loc, weight, dur int) {
	x, y := Loc2XY (loc)
	fmt.Printf ("位置 (%d, %d) 获得飞行辅助, 移动消耗 %d, 辅助时长 %d\n", x, y, weight, dur)
}

func (ui TestUI) LocDisableLong (loc, dur int) {
	x, y := Loc2XY (loc)
	fmt.Printf ("位置 (%d, %d) 获得远程无效辅助, 辅助时长 %d\n", x, y, dur)
}

func (ui TestUI) DrawMap () {
	var vname, dname, pname string
	fmt.Printf ("%6s", " Map")
	for x := 0; x < MAP_SIZE[0]; x ++ {
		fmt.Printf ("%6d", x)
	}
	fmt.Println ()
	for y := 0; y < MAP_SIZE[1]; y ++ {
		fmt.Printf ("%6d", y)
		for x := 0; x < MAP_SIZE[0]; x ++ {
			i := XY2Loc (x, y)
			dname = " " + dx_name[gmap[i].dx]
			if gmap[i].path_last_loc == -2 {
				vname = "*"
			} else {
				vname = ""
			}
			fmt.Printf ("%5s%1s", dname, vname)
		}
		fmt.Printf ("\n")
		fmt.Printf ("%6s", "")
		for x := 0; x < MAP_SIZE[0]; x ++ {
			i := XY2Loc (x, y)
			fmt.Printf ("%6d", gmap[i].weight)
		}
		fmt.Printf ("\n")
		fmt.Printf ("%6s", "")
		for x := 0; x < MAP_SIZE[0]; x ++ {
			i := XY2Loc (x, y)
			ip := gmap[i].people
			if ip == -1 {
				pname = " "
			} else {
				pname = " " + people[ip].name
			}
			fmt.Printf ("%6s", pname)
		}
		fmt.Printf ("\n")
	}
}

func (ui TestUI) ShowLoc (x, y int) {
	if x < 0 || x >= MAP_SIZE[0] || y < 0 || y >= MAP_SIZE[1] {
		fmt.Println ("输入错误")
		return
	}
	iloc := XY2Loc (x, y)
	loc := gmap[iloc]
	opt := ui.opt
	fmt.Printf ("位置 (%d, %d) 的辅助信息:\n", x, y)
	fmt.Printf ("    攻击辅助, 值 %d, 时长 %d\n", loc.buf_value[opt][0], loc.buf_dur[opt][0])
	fmt.Printf ("    防御辅助, 值 %d, 时长 %d\n", loc.buf_value[opt][1], loc.buf_dur[opt][1])
	fmt.Printf ("    飞行辅助, 消耗 %d, 时长 %d\n", loc.fly_weight[opt], loc.fly_weight_dur[opt])
	fmt.Printf ("    远程无效辅助: %v, 时长 %d\n", loc.disable_long[opt], loc.disable_long_dur[opt])
}

func (ui TestUI) ShowPeople (ip int) {
	if ip < 0 || ip >= len (people) {
		fmt.Println ("输入错误")
		return
	}
	p := people[ip]
	fmt.Printf ("人员 %s 的信息:\n", p.name)
	fmt.Printf ("    操作方: %d\n", p.opt)
	fmt.Printf ("    攻击: %d\n", p.abi)
	fmt.Printf ("    防御: %d\n", p.def)
	fmt.Printf ("    体力: %d\n", p.hp)
	fmt.Printf ("    技术:")
	for _, itech := range p.tech {
		fmt.Printf ("%s ", gtech[itech].name)
	}
	fmt.Printf("\n")
	fmt.Printf ("    速度: %d\n", p.speed)
	x, y := Loc2XY (p.loc)
	fmt.Printf ("    位置: (%d, %d)\n", x, y)
	fmt.Printf ("    是否完成行动: %v\n", p.is_finish)
}

func (ui TestUI) ShowTech (itech int) {

	tech_class_name := []string {
		"攻击",
		"攻击辅助",
		"防御辅助",
		"飞行辅助",
		"远程无效辅助",
	}

	if itech < 0 || itech >= len (gtech) {
		fmt.Println ("输入错误")
		return
	}

	tech := gtech[itech]
	fmt.Printf ("技能 %s 的信息:\n", tech.name)
	fmt.Printf ("    类型: %s\n", tech_class_name[tech.class])
	fmt.Printf ("    是否远程: %v\n", tech.is_long)
	fmt.Printf ("    AOE 范围: %d\n", tech.aoe_range)
	fmt.Printf ("    AOE 比例: %d%%\n", tech.aoe_ratio)
	fmt.Printf ("    AOE 误伤: %v\n", tech.aoe_self)
	fmt.Printf ("    持续时间: %d 回合\n", tech.duration)
	fmt.Printf ("    支持地形:")
	for i, dname := range dx_name {
		if tech.dx[i] {
			fmt.Printf (" %s ", dname)
		}
	}
	fmt.Printf ("\n")
	fmt.Printf ("    效果值: %d\n", tech.value)
}

func (ui TestUI) SelObj (status SelObjStatus, itech int, valid_loc []int) int {
	ui.DrawMap ()

	var msg string
	switch status {
	case SOS_PEOPLE:
		msg = "请选择操作角色 (X Y): (-1 1: 回合结束) (-1 0: 退出)"
	case SOS_MOVE_OBJ:
		msg = "请选择移动目的地 (X Y): (-1 0: 取消)"
	case SOS_TECH_OBJ:
		msg = fmt.Sprintf ("请为技术 %s 选择目标 (X Y): (-1 0: 放弃)", gtech[itech].name)
	}

	for {
		fmt.Printf ("%s \n  (-2 0 查看地形)\n  (-3 ip 查看人物)\n  (-4 itech 查看技术)\n", msg)
		var x, y, cx, cy int
		fmt.Scan (&x, &y)
		if y < 0 || y >= MAP_SIZE[1] {
			fmt.Println ("输入错误")
			continue
		}
		switch {
		case x == -1 :
			if status == SOS_PEOPLE && y == 1 {
				return -2 
			} else {
				return -1
			}
		case x == -2 :
			fmt.Println ("请输入查询位置 (X Y):")
			fmt.Scan (&cx, &cy)
			ui.ShowLoc (cx, cy)
		case x == -3 :
			ui.ShowPeople (y)
		case x == -4 :
			ui.ShowTech (y)
		case x >= 0:
			return XY2Loc (x, y)
		}
	}
}

func (ui TestUI) SelTech (ip int) int {
	ui.DrawMap ()
	fmt.Printf ("请选择使用技能:\n")
	fmt.Printf ("    %d : %s\n", -1, "取消")
	for i, itech := range people[ip].tech {
		fmt.Printf ("    %d : %s\n", i, gtech[itech].name)
	}
	ntech := len (people[ip].tech)

	var itech int
	fmt.Scan (&itech)
	if itech < 0 || itech >= ntech {
		return -1
	} else {
		return itech
	}
}

func (ui TestUI) Confirm (msg string) bool {
	yes_str := []string {"Y", "y", "Yes", "yes", "YES"}
	no_str := []string {"N", "n", "No", "no", "NO"}
	for {
		fmt.Printf ("确认 %s 吗 (Y/N)?", msg)
		var q_str string
		n, _ := fmt.Scan ("%s", &q_str)
		if n != 1 {
			continue
		}
		for _, str := range yes_str {
			if q_str == str {
				return true
			}
		}
		for _, str := range no_str {
			if q_str == str {
				return false
			}
		}
	}
}

func (ui TestUI) LoadMapDone () {
	fmt.Println ("地图加载完成")
	ui.DrawMap ()
}

func (ui TestUI)	GameStart () {
	fmt.Println ("游戏开始")
}

func (ui TestUI) TurnStart (opt int) {
	fmt.Printf ("操作方 %d 回合开始\n", opt)
}

/* win_opt == -1: 退出*/ 
func (ui TestUI) GameOver (win_opt int) {
	if win_opt == -1 {
		fmt.Printf ("退出游戏\n")
	} else {
		ui.DrawMap ()
		fmt.Printf ("操作方 %d 获胜\n", win_opt)
	}
}

func InitTestUI () {
	for i := 0; i < 2; i ++ {
		test_ui[i].opt = i
		test_uis[i] = test_ui[i]
	}
}

func TestFtbtCui (t * testing.T) {
	ftech, err := os.Open ("test.tech")
	if err != nil {
		fmt.Println ("无法打开测试技术文件");
		return
	}
	fmap, err := os.Open ("test.map")
	if err != nil {
		fmt.Println ("无法打开测试关卡文件");
		return
	}
	InitTestUI ()

	is_err := Init (test_ui[0], test_uis, ftech, fmap)
	if is_err {
		fmt.Printf ("初始化错误\n")
		return
	}
	Do()
}

func TestNCCui (t * testing.T) {
	ftech, err := os.Open ("star.tech")
	if err != nil {
		panic ("cannot open tech file")
	}
	fmap, err := os.Open ("m1.map")
	if err != nil {
		panic ("cannot open map file")
	}

	ui_note := InitNCursesUI ()
	var ui_sel [2]UISelect
	for i := 0; i < 1; i ++ {
		ui_sel[i] = InitNCursesUISelect (i, ui_note)
	}

	is_err := Init (ui_note, ui_sel, ftech, fmap)
	if is_err {
		fmt.Printf ("初始化错误\n")
		return
	}
	Do()
}
