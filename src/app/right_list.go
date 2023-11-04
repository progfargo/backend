package app

import (
	"fmt"
)

type Right struct {
	PageName string
	FuncName string
	Exp      string
}

func NewRight(pageName, funcName, exp string) *Right {
	rv := new(Right)
	rv.PageName = pageName
	rv.FuncName = funcName
	rv.Exp = exp

	return rv
}

type RightTab struct {
	Title string
	List  []*Right
	rMap  map[string]*Right
}

func NewRightTab(title string) *RightTab {
	rv := new(RightTab)

	rv.Title = title
	rv.List = make([]*Right, 0, 100)
	rv.rMap = make(map[string]*Right)

	return rv
}

func (rt *RightTab) Add(r *Right) {
	if _, ok := rt.rMap[MakeKey(r.PageName, r.FuncName)]; ok {
		panic("this item already exists. name: " + MakeKey(r.PageName, r.FuncName))
	}

	rt.List = append(rt.List, r)
	rt.rMap[MakeKey(r.PageName, r.FuncName)] = r
}

func (rt *RightTab) GetRight(pageName, funcName string) *Right {
	if _, ok := rt.rMap[MakeKey(pageName, funcName)]; !ok {
		panic(fmt.Sprintf("can not find right. page name: %s function name: %s", pageName, funcName))
	}

	return rt.rMap[MakeKey(pageName, funcName)]
}

type rightList struct {
	List []*RightTab
	rMap map[string]bool
}

func NewRightList() *rightList {
	rv := new(rightList)

	rv.List = make([]*RightTab, 0, 100)
	rv.rMap = make(map[string]bool)
	return rv
}

func (rl *rightList) Add(name string, rt *RightTab) {
	if _, ok := rl.rMap[name]; ok {
		panic("This right tab already exists.")
	}

	rl.List = append(rl.List, rt)
	rl.rMap[name] = true
}

func (rl *rightList) GetRightMap() map[string]*Right {
	rightMap := make(map[string]*Right, 100)

	var key string
	for _, rightTab := range rl.List {
		for _, right := range rightTab.List {
			key = MakeKey(right.PageName, right.FuncName)
			rightMap[key] = right
		}
	}

	return rightMap
}

func MakeKey(pageName, funcName string) string {
	return pageName + ":" + funcName
}

var UserRightList *rightList
var UserRightMap map[string]*Right

//app.UserRightList is being set in role_lib page to be able to translate explanations.
