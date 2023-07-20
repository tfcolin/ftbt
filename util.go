package ftbt

func IAbs (a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}

func Loc2XY (loc int) (x, y int) {
	if loc == -1 {
		return -1, -1
	}
	x = loc % MAP_SIZE[0]
	y = loc / MAP_SIZE[0]
	return
}

func XY2Loc (x, y int) (loc int) {
	loc = y * MAP_SIZE[0] + x
	return 
}
