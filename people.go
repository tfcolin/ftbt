package ftbt

import (
	"fmt"
	"io"
)

type People struct {
	name string
	opt int

	abi int
	def int

	hp int

	tech []int

	speed int

	loc int
	is_finish bool
}

var people []People
var np_finish [2]int
var np_in [2]int 

func LoadPeople (fin io.Reader, ui UINote) bool {
	var n_people, n, ntech, x, y int
	n, _ = fmt.Fscan (fin, &n_people)
	if n != 1 {
		return true
	}
	people = make ([]People, n_people)
	for i := 0; i < n_people; i ++ {
		n, _ = fmt.Fscan (fin, 
		&people[i].name, &people[i].opt, &people[i].abi, &people[i].def,
		&people[i].hp, &people[i].speed, &x, &y, &ntech)
		if n != 9 {
			return true
		}
		people[i].tech = make ([]int, ntech)
		for j := 0; j < ntech; j ++ {
			n, _ = fmt.Fscan (fin, &people[i].tech[j])
			if n != 1 {
				return true
			}
		}
		people[i].is_finish = false
		people[i].loc = -1
		SetPeople (i, XY2Loc (x, y), ui)
	}

	return false
}

func StartTurnPeople () {
	for i, _ := range people {
		people[i].is_finish = false
	}
	for i := 0; i < 2; i ++ {
		np_finish[i] = 0
	}
}

func MovePeople (ip int, iloc int, ui UINote) {
	if people[ip].loc == -1 {
		return
	}
	ui.PeopleMove (ip, people[ip].loc, iloc)
	gmap[people[ip].loc].people = -1
	people[ip].loc = iloc
	gmap[iloc].people = ip
}

func SetPeople (ip int, iloc int, ui UINote) {
	if people[ip].loc != -1 {
		return
	}
	ui.PeopleIn (ip, iloc)
	people[ip].loc = iloc
	gmap[iloc].people = ip
	np_in[people[ip].opt] ++
}

func RemovePeople (ip int, ui UINote) {
	p := &people[ip]
	if p.loc == -1 {
		return
	}
	ui.PeopleOut (ip, p.loc)
	gmap[p.loc].people = -1
	p.loc = -1
	np_in[p.opt] --
	if p.is_finish {
		np_finish[p.opt] --
	}
}

func DoPeople (ip int, ui UI) bool {
	var obj int

	p := &people[ip]
	if p.is_finish {
		return true
	}
	opt := p.opt

	valid_loc := SetMovingRange (ip)
	path := make ([]int, 1)
	for {
		obj = ui.sel[opt].SelObj (SOS_MOVE_OBJ, -1, valid_loc)
		if obj == -1 {
			return true
		}
		if (obj == p.loc || gmap[obj].people == -1) && gmap[obj].path_last_loc != -2 {
			break
		}
	}

	/* triger people step event */
	path[0] = obj
	for i := 0; path[i] >= 0; i ++ {
		path = append (path, gmap[path[i]].path_last_loc)
	}
	path = path[:len(path) - 1]
	for i := len (path) - 1; i > 0; i -- {
		var dir int
		if IAbs (path[i - 1] - path[i]) == 1 {
			dir = 0
		} else {
			dir = 2
		}
		if path[i - 1] < path[i] {
			dir ++ 
		}
		ui.note.PeopleStep (ip, path[i], dir)
	}

	if len(path) > 1 {
		MovePeople (ip, obj, ui.note)
	}

	var tech_obj int = -1
	for ; tech_obj == -1 ; {
		MapClearPath (true)
		itech := ui.sel[opt].SelTech (ip)
		if itech == -1 {
			break
		}
		SetTechRange (ip, itech)
		for {  
			tech_obj = ui.sel[opt].SelObj (SOS_TECH_OBJ, people[ip].tech[itech], nil)
			if tech_obj == -1 {
				break
			}
			is_err := DoTech (ip, itech, tech_obj, ui.note)
			if !is_err {
				break
			}
		}
	}

	if p.loc != -1 {
		p.is_finish = true
		np_finish[opt] ++
	}

	return false
}

