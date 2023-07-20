package ftbt

import (
	"fmt"
	"io"
)

const (
	NEIGHBOUR_WEIGHT_SCALE = 2
)

var (
	MAP_SIZE [2]int
)

type Location struct {
	people int
	dx int
	weight int

	/* 用于指示路径和是否有效: -1: 起点, 有效, -2: 无效 
	* 移动时指示路径中的上一个点 */
	path_last_loc int 

	/* [操作方][攻防] */
	buf_value [2][2] int
	buf_dur [2][2] int
	fly_weight [2]int
	fly_weight_dur [2]int
	disable_long [2]bool
	disable_long_dur [2]int
}

var gmap []Location 
var dx_weight [] int
var dx_name []string

func LoadDX (fin io.Reader) bool {
	var ndx, n int
	n, _ = fmt.Fscan (fin, &ndx)
	if n != 1 || ndx <= 0 {
		return true
	}
	dx_name = make ([]string, ndx)
	dx_weight = make ([]int, ndx)
	for i := 0; i < ndx; i ++ {
		n, _ = fmt.Fscan (fin, &dx_name[i], &dx_weight[i])
		if n != 2 {
			return true
		}
	}
	return false
}

func LoadMap (fin io.Reader, ui UINote) bool {
	var nx, ny int
	n, _ := fmt.Fscan (fin, &nx, &ny)
	if n != 2 {
		return true
	}
	MAP_SIZE[0] = nx
	MAP_SIZE[1] = ny
	gmap = make ([]Location, nx * ny)
	for i := 0; i < nx * ny; i ++ {
		n, _ := fmt.Fscan (fin, &gmap[i].dx)
		if n != 1 {
			return true
		}
		gmap[i].people = -1
		gmap[i].weight = -1
	}
	MapClearPath (true)
	ui.LoadMapDone ()
	return false
}

func MapClearPath (is_valid bool) {
	for i, _ := range gmap {
		if is_valid {
			gmap[i].path_last_loc = -1
		} else {
			gmap[i].path_last_loc = -2
		}
	}
}

func MapUpdateWeight (opt int) {
	for i, loc := range gmap {
		if gmap[i].people != -1 && people[gmap[i].people].opt != opt {
			gmap[i].weight = -1
		} else {
			if loc.fly_weight[opt] != 0 {
				gmap[i].weight = loc.fly_weight[opt]
			} else {
				gmap[i].weight = dx_weight[loc.dx]
			} 
		}
	}
}

func MapUpdataAux () {
	for i, _ := range gmap {
		loc := &gmap[i]
		for opt := 0; opt < 2; opt ++ {
			for j := 0; j < 2; j ++ {
				if loc.buf_dur[opt][j] > 0 {
					loc.buf_dur[opt][j] --
					if loc.buf_dur[opt][j] == 0 {
						loc.buf_value[opt][j] = 0
					}
				}
			}
			if loc.fly_weight_dur[opt] > 0 {
				loc.fly_weight_dur[opt] --
				if loc.fly_weight_dur[opt] == 0 {
					loc.fly_weight[opt] = 0
				}
			}
			if loc.disable_long_dur[opt] > 0 {
				loc.disable_long_dur[opt] --
				if loc.disable_long_dur[opt] == 0 {
					loc.disable_long[opt] = false
				}
			}
		}
	}
}

