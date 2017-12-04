package caching

import (
	"time"

	"fmt"
	"github.com/ilyaran/Malina/entity"
)

var T0 time.Time




var StatesMap map[int64]*entity.State
var StatesList []*entity.State
var StatesSelectOptionsView string
func SetStatesSelectOptionsView(statesList []*entity.State) string{
	StatesMap = make(map[int64]*entity.State,len(statesList))
	StatesList = []*entity.State{}
	StatesMap[0]=&entity.State{}
	StatesSelectOptionsView = `<select id="state" name="state">`
	for _,v := range statesList{
		StatesList = append(StatesList, v)
		StatesMap[v.GetId()] = v
		StatesMap[v.GetId()].SetFlag()
		StatesSelectOptionsView += fmt.Sprintf(`<option data-phone="%s" value="%d">%s</option>`,v.GetPhone(),v.GetId(),v.GetTitle())
	}
	StatesSelectOptionsView += `</select>`
	return StatesSelectOptionsView
}


var PublicPages = map[string][]byte{}

type Page struct {
	content []byte
	timestamp time.Time
}


func Get(keyName string)[]byte{

	if v,ok := PublicPages[keyName]; ok{
		return v
	}
	return nil
}

func Set(keyName string, cacheData []byte){
	PublicPages[keyName]=cacheData
}




























