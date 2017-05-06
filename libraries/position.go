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
	SelectOptionsList string
}

func (this *PositionLibrary) SetTrees(positionList []*entity.Position){
	if positionList != nil && len(positionList) > 0 {
		this.TreeList = []*entity.Position{}
		this.TreeMap = map[int64]*entity.Position{}
		var anonymous func(int64)
		var level int = 0
		anonymous = func(parent int64) {
			for i := 0; i < len(positionList); i++ {
				if positionList[i].GetParent().GetId() == parent {
					positionList[i].SetLevel(level)
					this.TreeMap[positionList[i].GetId()] = positionList[i]
					this.TreeList = append(this.TreeList, positionList[i])
					level ++
					anonymous(positionList[i].GetId())
				}
			}
			level --
		}
		anonymous(0)

	}
}
func (this *PositionLibrary) SetSelectOptionsView(treeList []*entity.Position, disable, selected int64){
	this.SelectOptionsList = ""
	const dis string = ` disabled="disabled"`
	const sel string = ` selected="selected"`
	for _, v := range treeList {
		this.SelectOptionsList += `<option`
		if v.GetId() == disable {
			this.SelectOptionsList += dis
		}
		if v.GetId() == selected {
			this.SelectOptionsList += sel
		}
		this.SelectOptionsList += ` value="` + fmt.Sprintf("%d",v.GetId()) + `">` + strings.Repeat("&rarr;", v.GetLevel()) + v.GetTitle() + `</option>`
	}

}


