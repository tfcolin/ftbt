package ftbt

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

const (
	TC_PEOPLE = 0
	/* 时间, 范围均不可变 */
	TC_ABI_BUF = 1
	TC_DEF_BUF = 2
	/* 时间固定, 范围可变 */
	TC_FLY = 3
	/* 时间, 范围均可变 */
	TC_DISABLE_LONG = 4
)

var (
	MAX_FLY_WEIGHT int
	MAX_ABI int
)

type Tech struct {
	name string
	class int
	/* 是否为远距离 */
	is_long bool
	/* 可选择的作用范围 */
	rel_range [][2]int
	/* aoe 范围, 可变时，表示最大值 */
	aoe_range int
	/* aoe 比例 (%)*/
	aoe_ratio int
	/* aoe 是否影响己方单位 */
	aoe_self bool
	/* 持续时间, 可变时, 表示最大值 */
	duration int
	/* 支持地形 */
	dx []bool

	/* 效果值:
	* 对于 TC_PEOPLE, 为基础伤害值, 
	* 对于 其他类型, 为使用者 abi 的百分比值
	*/
	value int

}

var gtech []Tech

func LoadTech (fin io.Reader) bool {
	var n, ntech, n_range, ndx int
	ndx = len (dx_weight)
	if ndx == 0 {
		return true
	}
	n, _ = fmt.Fscan (fin, &MAX_FLY_WEIGHT, &MAX_ABI, &ntech)
	if n != 3 {
		return true
	}
	gtech = make ([]Tech, ntech)
	for i := 0; i < ntech; i ++ {
		n, _ = fmt.Fscan (fin, &gtech[i].name, &gtech[i].class, &gtech[i].is_long, &n_range)
		if n != 4 {
			return true
		}
		gtech[i].rel_range = make ([][2]int, n_range)
		for j := 0; j < n_range; j ++ {
			n, _ = fmt.Fscan (fin, &gtech[i].rel_range[j][0], &gtech[i].rel_range[j][1]) 
			if n != 2 {
				return true
			}
		}
		n, _ = fmt.Fscan (fin, &gtech[i].aoe_range,
		&gtech[i].aoe_ratio, &gtech[i].aoe_self, &gtech[i].duration)
		if n != 4 {
			return true
		}
		gtech[i].dx = make ([]bool, ndx)
		for j := 0; j < ndx; j ++ {
			n, _ = fmt.Fscan (fin, &gtech[i].dx[j])
			if n != 1 {
				return true
			}
		}
		n, _ = fmt.Fscan (fin, &gtech[i].value)
		if n != 1 {
			return true
		}
	}
	return false
}

/* 如果目标类型错误, 返回 true */
func DoTech (ip int, itech int, obj_loc int, ui UINote) (is_err bool) {
	var is_update_weight bool

	if gmap[obj_loc].path_last_loc != -1 {
		return true
	}
	p := &people[ip]
	tech := gtech[p.tech[itech]]
	switch {
	case tech.class == TC_PEOPLE:
		ip_obj := gmap[obj_loc].people
		if ip_obj == -1 {
			return true
		}
		p_obj := &people[ip_obj]
		if p_obj.opt == p.opt {
			return true
		}
		ui.DoTech (p.tech[itech], ip, ip_obj, obj_loc)
		dhp_c := p.abi - p_obj.def + tech.value 
		dhp_c += gmap[p.loc].buf_value[p.opt][0]
		if dhp_c < 1 {
			dhp_c = 1
		}
		dhp_aoe := dhp_c * tech.aoe_ratio / 100
		if dhp_aoe < 1 {
			dhp_aoe = 1
		}
		aoe_range := GetAoeRange (obj_loc, tech.aoe_range)
		for _, loc := range aoe_range {
			var dhp int
			if tech.dx[gmap[loc].dx] && gmap[loc].people >= 0 {
				ip_aoe := gmap[loc].people
				p_aoe := &people[ip_aoe]
				if !tech.aoe_self && p_aoe.opt == p.opt {
					continue
				}
				if tech.is_long && gmap[loc].disable_long[p_aoe.opt] {
					continue
				}
				if loc == obj_loc {
					dhp = dhp_c 
				} else {
					dhp = dhp_aoe
				}
				dhp += int(math.Floor (0.15 * rand.Float64() * float64(dhp)))
				dhp -= gmap[loc].buf_value[1 - p.opt][1]
				if dhp < 1 {
					dhp = 1
				}
				ui.HPDown (ip_aoe, loc, p_aoe.hp, dhp)
				p_aoe.hp -= dhp
				if (p_aoe.hp <= 0) {
					p_aoe.hp = 0
					RemovePeople (ip_aoe, ui)
					if p_aoe.opt != p.opt {
						is_update_weight = true
					}
				}
			}
		}
	case tech.class >= TC_ABI_BUF && tech.class < TC_FLY:
		ui.DoTech (p.tech[itech], ip, -1, obj_loc)
		iad := tech.class - TC_ABI_BUF
		aoe_range := GetAoeRange (obj_loc, tech.aoe_range)
		buf_c := p.abi * tech.value / 100
		for _, loc := range aoe_range {
			var buf int
			if !tech.dx[gmap[loc].dx] {
				continue
			}
			if loc == obj_loc {
				buf = buf_c
			} else {
				buf = buf_c * tech.aoe_ratio / 100
			}
			ui.LocBuf (loc, iad, buf, tech.duration)
			gmap[loc].buf_value[p.opt][iad] = buf
			gmap[loc].buf_dur[p.opt][iad] = tech.duration
		}
	case tech.class == TC_FLY:
		ui.DoTech (p.tech[itech], ip, -1, obj_loc)
		aoe_dist := tech.aoe_range * p.abi * tech.value / (100 * MAX_ABI)
		aoe_range := GetAoeRange (obj_loc, aoe_dist)
		abi_c := p.abi * tech.value / 100
		for _, loc := range aoe_range {
			var abi, weight int
			if !tech.dx[gmap[loc].dx] {
				continue
			}
			if loc == obj_loc {
				abi = abi_c
			} else {
				abi = abi_c * tech.aoe_ratio / 100
			}
			weight = (MAX_FLY_WEIGHT - 1) * (MAX_ABI - abi) / MAX_ABI + 1
			ui.LocFly (loc, weight, tech.duration)
			gmap[loc].fly_weight[p.opt] = weight
			gmap[loc].fly_weight_dur[p.opt] = tech.duration
		}
		is_update_weight = true
	case tech.class == TC_DISABLE_LONG:
		ui.DoTech (p.tech[itech], ip, -1, obj_loc)
		aoe_dist := tech.aoe_range * p.abi / MAX_ABI
		aoe_range := GetAoeRange (obj_loc, aoe_dist)
		dur := tech.duration * p.abi / MAX_ABI
		if dur < 1 {
			dur = 1
		}
		for _, loc := range aoe_range {
			if !tech.dx[gmap[loc].dx] {
				continue
			}
			ui.LocDisableLong (loc, dur)
			gmap[loc].disable_long[p.opt] = true
			gmap[loc].disable_long_dur[p.opt] = dur
		}
	}
	if is_update_weight {
		MapUpdateWeight (p.opt)
	}
	return false
}

