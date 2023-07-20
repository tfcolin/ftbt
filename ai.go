package ftbt

import (
	"fmt"
	"io"
	"gitee.com/tfcolin/dsg"
)

var (
	AI_MAX_MOVE_DIST int
	AI_MAX_MOVE_SCORE int
	AI_MAX_BUF_VALUE int
	AI_MAX_BUF_SCORE int
	AI_MAX_DEF_DIST int
	AI_MAX_DEF_SCORE int
	AI_MAX_TECH_DIST int
	AI_MAX_TECH_SCORE int
	AI_MAX_TECH_BUF_VMD int
	AI_MAX_TECH_FLY_VMD int
	AI_MAX_TECH_DL_VMD int
)

/* 用于保存距离的工作空间 */
var (
	w_map []int
	ai_people []AIPeople
)


type AIPeople struct {
	goal int
	move_dist []int
	goal_dist []int

	/* 可移动目标 */
	move_loc []int
	/* 各项技术的允许目标 [it (in people)][] 
	* TC_PEOPLE: 为 people 的索引
	* 其他类型: 为 gmap 的索引
	*/
	tech_loc [][]int

	/* 最优决策目标 */
	opt_move_loc int
	opt_itech int
	opt_tech_obj int

	tech_people_score_  []int   
	tech_loc_score_ 	  []int   
	tech_people_score   [][]int 
	tech_loc_score      [][]int 

	/* 决策因子 */
	fact_att, fact_def, fact_move int
}

func LoadAIBase (fin io.Reader) bool {
	var n int
	n, _ = fmt.Fscan (fin, &	AI_MAX_MOVE_DIST  )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_MOVE_SCORE )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_BUF_VALUE  )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_BUF_SCORE  )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_DEF_DIST   )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_DEF_SCORE  )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_TECH_DIST )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_TECH_SCORE )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_TECH_BUF_VMD )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_TECH_FLY_VMD )
	if n != 1 { return true }
	n, _ = fmt.Fscan (fin, &	AI_MAX_TECH_DL_VMD )
	if n != 1 { return true }
	return false
}

func InitAI () {
	msize := MAP_SIZE[0] * MAP_SIZE[1]
	w_map = make ([]int, msize)
	ai_people = make ([]AIPeople, len(people))

	for i, p := range people {
		aip := &ai_people[i]
		aip.move_loc = make ([]int ,0)
		aip.tech_loc = make ([][]int, len(p.tech))
		aip.goal = -1
		aip.move_dist = make ([]int, msize)
		aip.goal_dist = make ([]int, msize)
		aip.opt_move_loc = -1
		aip.opt_tech_obj = -1
		aip.tech_people_score_ = make ([]int, len(p.tech) * len(people))
		aip.tech_loc_score_ = make ([]int, len(p.tech) * msize)
		aip.tech_people_score = dsg.IndexIArray2D (aip.tech_people_score_, []int{len(p.tech), len(people)})
		aip.tech_loc_score = dsg.IndexIArray2D (aip.tech_loc_score_, []int{len(p.tech), msize})
		aip.fact_att = 30
		aip.fact_def = 30
		aip.fact_move = 40
	}
	for i := 0; i < msize; i ++ {
		w_map[i] = dx_weight[gmap[i].dx]
	}

}

func LoadAI (ftech, fmap io.Reader) bool {
	if (LoadAIBase (ftech)) {
		return true
	}
	if (LoadAIPeople (fmap)) {
		return true
	}
	return false
}

func LoadAIPeople (fin io.Reader) bool {
	var n int
	var goal, goal_x, goal_y, att, def, mov int
	for i, p := range people {
		n, _ = fmt.Fscan (fin, &att, &def, &mov)
		if n != 3 {return true}
		n, _ = fmt.Fscan (fin, &goal_x, &goal_y)
		if n != 2 {return true}
		if goal_x >= 0 && goal_y >= 0 {
			goal = XY2Loc (goal_x, goal_y)
		} else {
			goal = -1
		}
		if p.opt == 0 {
			goal, att, def, mov = -1, -1, -1, -1
		}
		if p.opt == 1 && att == -1 {
			att, def, mov = 30, 30, 40
		}
		AISetArg (i, goal, att, def, mov)
	}
	return false
}

func AISetArg (ip int, goal int, fact_att, fact_def, fact_move int) {
	aip := &ai_people[ip]
	aip.goal = goal
	aip.fact_att = fact_att
	if aip.fact_att == -1 {
		aip.fact_def = -1
		aip.fact_move = -1
	} else {
		aip.fact_def = fact_def
		aip.fact_move = fact_move
	}
}

/* called for all people */
func AICalAllDist () {
	for i, p := range people {
		aip := &ai_people[i]
		if p.loc == -1 {
			continue 
		}
		CalDist (p.loc, w_map, aip.move_dist)
		if aip.fact_att == -1 {
			continue
		}
		if aip.goal != -1 {
			CalDist (aip.goal, w_map, aip.goal_dist)
		}
	}
}

func AICalTechObjs (ip int, move_loc []int) {
	aip := &ai_people[ip]
	if aip.fact_att == -1 {
		return 
	}
	aip.move_loc = move_loc
	p := people[ip]

	pset := dsg.InitSet (len(people))
	mset := dsg.InitSet (MAP_SIZE[0] * MAP_SIZE[1])

	for i, itech := range p.tech {
		pset.Empty()
		mset.Empty()
		aip.tech_loc[i] = make ([]int, 0)
		for _, cloc := range aip.move_loc {
			x, y := Loc2XY (cloc)
			tech := gtech[itech]
			for _, dp := range tech.rel_range {
				loc := loc_plus (x, y, dp)
				if loc == -1 { continue }
				if !tech.dx[gmap[loc].dx] { continue }
				if tech.class == TC_PEOPLE {
					ip_obj := gmap[loc].people 
					if ip_obj != -1 && people[ip_obj].opt != p.opt {
						if !tech.is_long || !gmap[loc].disable_long[people[ip_obj].opt] {
							if !pset.GetLabel (ip_obj) {
								aip.tech_loc[i] = append (aip.tech_loc[i], ip_obj)
								pset.SetLabel (ip_obj, true)
							}
						}
					}
				} else {
					if !mset.GetLabel (loc) {
						aip.tech_loc[i] = append (aip.tech_loc[i], loc)
						mset.SetLabel (loc, true)
					}
				}
			}
		}
	}
}

/* called for only AI People */
func AIMakeDecision (ip int, move_loc []int) {
	aip := &ai_people[ip]
	p := people[ip]
	if p.loc == -1 || aip.fact_att == -1 { return }

	AICalTechObjs (ip, move_loc)

	loc_score := make ([]int, len(aip.move_loc))
	def_score := make ([]int, len(aip.move_loc))
	att_score := make ([]int, len(aip.move_loc))
	att_tech := make ([]int, len(aip.move_loc))
	att_tech_loc := make ([]int, len(aip.move_loc))

	if aip.goal == -1 {
		ip_min_d := -1
		min_d := 0
		for irp, rp := range people {
			if p.opt != rp.opt && rp.loc != -1 {
				if ip_min_d == -1 || aip.move_dist[rp.loc] < min_d {
					min_d = aip.move_dist[rp.loc] 
					ip_min_d = irp
				}
			}
		}
		if ip_min_d == -1 {
			panic ("no enermy for AI")
		}
		CalDist (people[ip_min_d].loc, w_map, aip.goal_dist)
	}

	nrp := 0
	for irp, rp := range people {
		airp := ai_people[irp]
		if p.opt != rp.opt && rp.loc != -1 {
			for i, loc := range aip.move_loc {
				if airp.move_dist[loc] < AI_MAX_DEF_DIST {
					def_score[i] += airp.move_dist[loc] * AI_MAX_DEF_SCORE / AI_MAX_DEF_DIST
				} else {
					def_score[i] += AI_MAX_DEF_SCORE
				}
			}
			nrp ++
		}
	}

	for i, loc := range aip.move_loc {
		if aip.goal_dist[loc] > AI_MAX_MOVE_DIST {
			loc_score[i] = 0
		} else {
			loc_score[i] = (AI_MAX_MOVE_DIST - aip.goal_dist[loc]) * AI_MAX_MOVE_SCORE / AI_MAX_MOVE_DIST
		}
	}

	for i, loc := range aip.move_loc {
		def_score[i] /= nrp
		for j := 0; j < 2; j ++ {
			loc_score[i] += gmap[loc].buf_value[p.opt][j] * AI_MAX_BUF_SCORE / AI_MAX_BUF_VALUE 
		}
		att_score[i] = 0
	}

	for i, objs := range aip.tech_loc {
		tech := gtech[p.tech[i]]
		switch {
			case tech.class == TC_PEOPLE: 
			for _, ip := range objs {
				p_obj := people[ip]
				dhp_c := p.abi - p_obj.def + tech.value 
				if dhp_c < 1 {
					dhp_c = 1
				}
				dhp_aoe := dhp_c * tech.aoe_ratio / 100
				if dhp_aoe < 1 {
					dhp_aoe = 1
				}
				aoe_range := GetAoeRange (p_obj.loc, tech.aoe_range)
				score := 0
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
						if loc == p_obj.loc {
							dhp = dhp_c 
						} else {
							dhp = dhp_aoe
						}
						if dhp < 1 {
							dhp = 1
						}
						if dhp > p_aoe.hp {
							dhp = p_aoe.hp
						}
						if p_aoe.opt != p.opt {
							score += dhp * AI_MAX_TECH_SCORE / p_aoe.hp
						} else {
							score -= dhp * AI_MAX_TECH_SCORE / p_aoe.hp
						}
					}
				}
				if score < 0 { score = 0 } 
				aip.tech_people_score[i][ip] = score
			}
		case tech.class >= TC_ABI_BUF && tech.class < TC_FLY:
			for _, obj_loc := range objs {
				aoe_range := GetAoeRange (obj_loc, tech.aoe_range)
				iad := tech.class - TC_ABI_BUF
				buf_c := p.abi * tech.value / 100
				vmd := 0
				for _, loc := range aoe_range {
					var buf int
					if !tech.dx[gmap[loc].dx] {
						continue
					}
					dur := tech.duration - gmap[loc].buf_dur[p.opt][iad] 
					if dur <= 0 {
						continue
					}
					inv_gd := AI_MAX_TECH_DIST - aip.goal_dist[loc]
					if inv_gd <= 0 {
						continue
					}
					if loc == obj_loc {
						buf = buf_c
					} else {
						buf = buf_c * tech.aoe_ratio / 100
					}
					vmd += buf * dur * inv_gd / AI_MAX_TECH_DIST
				}
				if vmd > AI_MAX_TECH_BUF_VMD {
					vmd = AI_MAX_TECH_BUF_VMD
				}
				aip.tech_loc_score[i][obj_loc] = vmd * AI_MAX_TECH_SCORE / AI_MAX_TECH_BUF_VMD
			}
		case tech.class == TC_FLY:
			aoe_dist := tech.aoe_range * p.abi * tech.value / (100 * MAX_ABI)
			abi_c := p.abi * tech.value / 100
			for _, obj_loc := range objs {
				aoe_range := GetAoeRange (obj_loc, aoe_dist)
				vmd := 0
				for _, loc := range aoe_range {
					var abi int
					if !tech.dx[gmap[loc].dx] {
						continue
					}
					dur := tech.duration - gmap[loc].fly_weight_dur[p.opt] 
					if dur <= 0 {
						continue
					}
					inv_gd := AI_MAX_TECH_DIST - aip.goal_dist[loc]
					if inv_gd <= 0 {
						continue
					}
					if loc == obj_loc {
						abi = abi_c
					} else {
						abi = abi_c * tech.aoe_ratio / 100
					}
					vmd_local := w_map[loc] - (MAX_FLY_WEIGHT - 1) * (MAX_ABI - abi) / MAX_ABI + 1 
					if vmd_local <= 0 {
						continue
					}
					vmd += vmd_local * dur * inv_gd / AI_MAX_TECH_DIST
				}
				if vmd > AI_MAX_TECH_BUF_VMD {
					vmd = AI_MAX_TECH_BUF_VMD
				}
				aip.tech_loc_score[i][obj_loc] = vmd * AI_MAX_TECH_SCORE / AI_MAX_TECH_FLY_VMD
			}
		case tech.class == TC_DISABLE_LONG:
			aoe_dist := tech.aoe_range * p.abi / MAX_ABI
			dur := tech.duration * p.abi / MAX_ABI
			if dur < 1 {
				dur = 1
			}
			for _, obj_loc := range objs {
				aoe_range := GetAoeRange (obj_loc, aoe_dist)
				vmd := 0
				for _, loc := range aoe_range {
					if !tech.dx[gmap[loc].dx] {
						continue
					}
					dur_diff := dur - gmap[loc].disable_long_dur[p.opt]
					if dur_diff <= 0 {
						continue
					}
					inv_gd := AI_MAX_TECH_DIST - aip.goal_dist[loc]
					if inv_gd <= 0 {
						continue
					}
					vmd += dur_diff * inv_gd / AI_MAX_TECH_DIST
				}
				if vmd > AI_MAX_TECH_DL_VMD {
					vmd = AI_MAX_TECH_DL_VMD 
				}
				aip.tech_loc_score[i][obj_loc] = vmd * AI_MAX_TECH_SCORE / AI_MAX_TECH_DL_VMD
			}
		}
	}

	for i, cloc := range aip.move_loc {
		x, y := Loc2XY (cloc)
		max_score := 0
		max_score_tech := -1
		max_score_loc := -1
		for j, itech := range p.tech {
			tech := gtech[itech]
			for _, dp := range tech.rel_range {
				loc := loc_plus (x, y, dp)
				if loc == -1 { continue }
				if !tech.dx[gmap[loc].dx] { continue }
				if tech.class == TC_PEOPLE {
					ip_obj := gmap[loc].people 
					if ip_obj != -1 && people[ip_obj].opt != p.opt {
						if !tech.is_long || !gmap[loc].disable_long[people[ip_obj].opt] {
							if max_score_tech == -1 ||  aip.tech_people_score[j][ip_obj] > max_score {
								max_score_loc = loc
								max_score_tech = j
								max_score = aip.tech_people_score[j][ip_obj]
							}
						}
					}
				} else {
					if max_score_tech == -1 ||  aip.tech_loc_score[j][loc] > max_score {
						max_score_loc = loc
						max_score_tech = j
						max_score = aip.tech_loc_score[j][loc]
					}
				}
			}
		}
		att_score[i] += max_score
		att_tech[i] = max_score_tech
		att_tech_loc[i] = max_score_loc
	}

	max_score := 0
	i_max_score := -1
	for i, loc := range aip.move_loc {
		if gmap[loc].people != -1 {
			continue
		}
		score := att_score[i] * aip.fact_att + def_score[i] * aip.fact_def + loc_score[i] * aip.fact_move
		if i_max_score == -1 || score > max_score {
			max_score = score
			i_max_score = i
		}
	}

	if i_max_score == -1 {
		aip.opt_move_loc = p.loc
		aip.opt_itech = -1
		aip.opt_tech_obj = -1
	} else {
		aip.opt_move_loc = aip.move_loc[i_max_score]
		aip.opt_itech = att_tech[i_max_score]
		aip.opt_tech_obj = att_tech_loc[i_max_score]
	}
}

