package ftbt

import (
	"gitee.com/tfcolin/dsg"
)

type MPEdge struct {
	start int
	end int
	length int
}

var mpq * dsg.PQueue

func mp_cmp (v1_, v2_ dsg.Value) int {
	v1 := v1_.(MPEdge)
	v2 := v2_.(MPEdge)
	if v1.length < v2.length {
		return -1
	} else {
		return 1
	}
}

func next_loc (loc int, dir int) (nloc int) {
	nloc = 0
	x, y := Loc2XY (loc)
	pn := 1 - dir % 2 * 2
	xy := dir / 2
	if xy == 0 {
		x += pn
		if (x < 0 || x >= MAP_SIZE[0]) {
			nloc = -1
		}
	} else {
		y += pn
		if (y < 0 || y >= MAP_SIZE[1]) {
			nloc = -1
		}
	}
	if nloc == -1 {
		return
	}
	nloc = XY2Loc (x, y)
	return
}

func loc_plus (cx, cy int, dp [2]int) (loc int) {
	x, y := cx + dp[0], cy + dp[1]
	if x < 0 || x >= MAP_SIZE[0] || y < 0 || y >= MAP_SIZE[1] {
		loc = -1
	} else {
		loc = XY2Loc (x, y)
	}
	return loc
}

func InitMiniPath () {
	mpq = dsg.InitPQueue (mp_cmp)
}

func CalDist (center int, weight []int, dist []int) {
	var edge MPEdge 
	mpq.Flush()
	for i := 0; i < MAP_SIZE[0] * MAP_SIZE[1]; i ++ {
		dist[i] = -1
	}

	edge.start = -1
	edge.end = center
	edge.length = 0
	mpq.Add (edge)

	for {
		p_edge := mpq.Pop()
		if p_edge == nil {
			break
		}
		edge = p_edge.(MPEdge)
		loc := edge.end
		if dist[loc] != -1 {
			continue
		}
		dist[loc] = edge.length

		for dir := 0; dir < 4; dir ++ {
			nloc := next_loc(loc, dir)
			if nloc == -1 || dist[nloc] != -1 {
				continue
			}

			new_edge := MPEdge {
				start : loc,
				end : nloc,
				length : edge.length + weight[nloc],
			}

			mpq.Add (new_edge)
		}
	}
}

func SetMovingRange (ip int) (valid_loc []int) {
	var edge MPEdge 
	MapClearPath (false)
	p := people[ip]
	if p.loc == -1 {
		return
	}
	valid_loc = make ([]int, 0)

	mpq.Flush()

	edge.start = -1
	edge.end = p.loc
	edge.length = 0
	mpq.Add (edge)

	for {
		p_edge := mpq.Pop()
		if p_edge == nil {
			break
		}
		edge = p_edge.(MPEdge)
		loc := edge.end
		last_loc := edge.start
		if gmap[loc].path_last_loc != -2 {
			continue
		}
		gmap[loc].path_last_loc = last_loc
		valid_loc = append (valid_loc, loc)

		var dir_has_neighbour [4]bool 
		for dir := 0; dir < 4; dir ++ {
			nloc := next_loc(loc, dir)
			if nloc == -1 || gmap[nloc].path_last_loc != -2 {
				continue
			}
			if gmap[nloc].weight == -1 {
				dir_has_neighbour[dir] = true
			}
		}

		for dir := 0; dir < 4; dir ++ {
			nloc := next_loc(loc, dir)
			if nloc == -1 || gmap[nloc].path_last_loc != -2 || gmap[nloc].weight == -1 {
				continue
			}

			weight := gmap[nloc].weight
			/* 侧向经过对方人员, 移动消耗加倍 */
			dir1 := (1 - dir / 2) * 2
			for i := 0 ; i < 2; i ++ {
				if dir_has_neighbour[dir1] {
					weight *= 2
					break
				}
				dir1 ++
			}

			new_edge := MPEdge {
				start : loc,
				end : nloc,
				length : edge.length + weight,
			}
			if new_edge.length <= p.speed {
				mpq.Add (new_edge)
			}
		}
	}
	return
}

func SetTechRange (ip int, itech int) {
	MapClearPath (false)
	center := people[ip].loc
	tech := gtech[people[ip].tech[itech]]
	for _, dp := range tech.rel_range {
		x, y := Loc2XY (center)
		loc := loc_plus (x, y, dp)
		if loc == -1 {
			continue
		}
		if tech.dx[gmap[loc].dx] {
			gmap[loc].path_last_loc = -1
		}
	}
}

func LocDist (l1, l2 int) int {
	x1, y1 := Loc2XY (l1)
	x2, y2 := Loc2XY (l2)
	return IAbs (x1 - x2) + IAbs (y1 - y2)
}

func GetAoeRange (cloc int, d int) (aloc []int) {
	cx, cy := Loc2XY (cloc)
	for dy := -d ; dy <= d; dy ++ {
		i := d - IAbs(dy)
		ay := cy + dy
		if ay < 0 || ay >= MAP_SIZE[1] {
			continue
		}
		for dx := -i ; dx <= i; dx ++ {
			ax := cx + dx
			if ax < 0 || ax >= MAP_SIZE[0] {
				continue
			}
			loc := XY2Loc (ax, ay)
			aloc = append (aloc, loc)
		}
	}
	return
}

