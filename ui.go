package ftbt

type SelObjStatus int
const (
	SOS_PEOPLE SelObjStatus = 0
	SOS_MOVE_OBJ SelObjStatus = 1
	SOS_TECH_OBJ SelObjStatus = 2
)

/* all UI event is driven before data update */
type UINote interface {
	/* DoTech triger before HPDown, Loc.. */
	DoTech (itech, p_src, p_obj, l_obj int)
	HPDown (ip, loc, hp, dhp int)
	PeopleIn (ip, loc int)
	PeopleOut (ip, loc int)
	/* PeopleMove trigger after PeopleStep */
	PeopleMove (ip, loc_start, loc_end int)
	PeopleStep (ip, loc, dir int)
	/* att_def: 0: att; 1: def */
	LocBuf (loc, att_def, val, dur int)
	LocFly (loc, weight, dur int)
	LocDisableLong (loc, dur int)

	LoadMapDone ()
	GameStart ()
	/* win_opt == -1: 退出*/ 
	GameOver (win_opt int) 
}

type UISelect interface {
	TurnStart (opt int)
      /* return -1 to quit, -2 to end turn (only for SOS_PEOPLE) 
	* valid_loc is only valid in SOS_MOVE_OBJ status
	* */
	SelObj (status SelObjStatus, itech int, valid_loc []int) int 
      /* return -1 to give up tech use */
	SelTech (ip int) int
	/* return true to confirm */
	Confirm (msg string) bool 
}

type UI struct {
	note UINote
	sel [2]UISelect
}
