package core

import (
	"Malina/libraries"
	"Malina/models"
	"Malina/config"
	"Malina/language"
)

var MALINA *Malina
type Malina struct {}


func (s *Malina) SetGlobals(){

	library.POSITION = &library.PositionLibrary{}
	library.CATEGORY = &library.CategoryLibrary{Id:-1,Parent:-1}
	library.SESSION = &library.Session{}
	library.VALIDATION = &library.Validation{}
	library.UPLOAD = &library.Upload{}

	s.SetPositionGlobals()
	s.SetCategoryGlobals()

	for _,v := range app.Language_list(){
		if i,ok := lang.FullLanguageList[v];ok{
			lang.LangSelectOptions += `<option value="` + v + `" class="in-of">` + i + `</option>`
		}
	}

}
func (s *Malina) SetPositionGlobals(){
	//set globals for category
	var positionList = model.PositionModel.GetList("","","","position_sort ASC")
	library.POSITION.SetTrees(positionList)
	library.POSITION.SetSelectOptionsView(positionList,-1,-1)

}
func (s *Malina) SetCategoryGlobals(){
	//set globals for category
	var categoryList = model.CategoryModel.GetList("","","","category_sort ASC")
	library.CATEGORY.SetTrees(categoryList)
}


























