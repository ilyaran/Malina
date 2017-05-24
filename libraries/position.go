package library

import (
	"fmt"
	"Malina/entity"
	"strings"
)

var POSITION *PositionLibrary

type PositionLibrary struct {
	TreeList []*entity.Position
	TreeMap  map[int64]*entity.Position
	//SelectOptionsList string
}

func (this *PositionLibrary) SetTrees(positionList []*entity.Position){
	if positionList != nil && len(positionList) > 0 {
		this.TreeList = []*entity.Position{}
		this.TreeMap = map[int64]*entity.Position{}
		var closure func(int64)
		var level int = 0
		closure = func(parent int64) {
			for i := 0; i < len(positionList); i++ {
				if positionList[i].GetParent().GetId() == parent {
					positionList[i].SetLevel(level)
					this.TreeMap[positionList[i].GetId()] = positionList[i]
					this.TreeList = append(this.TreeList, positionList[i])
					level ++
					closure(positionList[i].GetId())
				}
			}
			level --
		}
		closure(0)
		if len(this.TreeList) > 0 {
			for _, v := range this.TreeList {
				v.Set_descendants(this.TreeList)
			}
		}
		//fmt.Println(this.TreeMap)
		//library.POSITION.SetSelectOptionsView(positionList,-1,-1)
		//this.SetSelectOptionsView(positionList,-1,-1)

	}
}


func (this *PositionLibrary) BuildSelectOptionsView(disable map[int64]bool, enable map[int64]bool, selected int64)string{
	var out string
	const dis string = ` disabled="disabled"`
	const sel string = ` selected="selected"`
	for _, v := range this.TreeList {
		out += `<option`
		if SESSION.SessionObj.Account.Position.GetParent().GetId() > 0 {
			if disable != nil {
				if _, ok := disable[v.GetId()]; ok {
					out += dis
				}
			} else if enable != nil {
				if _, ok := enable[v.GetId()]; !ok {
					out += dis
				}
			}
		}
		if v.GetId() == selected {
			out += sel
		}
		out += ` value="` + fmt.Sprintf("%d",v.GetId()) + `">` + strings.Repeat("&rarr;", v.GetLevel()) + v.GetTitle() + `</option>`
	}
	//this.SelectOptionsList = out
	return out
}


