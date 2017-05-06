/**
 * Position controller class.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package controller

import (
	"net/http"
	"Malina/helpers"
	"Malina/models"
	"fmt"
	"encoding/json"
	"Malina/language"
	"Malina/entity"
	"Malina/views"
	"html/template"
	"Malina/libraries"
	"Malina/views/position"
	"Malina/config"
	"Malina/core"
)

var Position = &position{&CrudController{}}
type position  struct {
	crud *CrudController
}

func (this *position) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("position",w,r)

	switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler(action,w,r)
		case "edit" 	 : this.FormHandler(action,w,r)
		case "del" 	 : this.Del(w,r)
		default		 : this.List(w,r); return // "list"
	}
	helper.SetAjaxHeaders(w)

	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *position) List(w http.ResponseWriter, r *http.Request) {
	page, _ := this.crud.getList(false,w,r)
	order := "position_id ASC"
	if library.VALIDATION.Status == 0 {
		all := model.PositionModel.CountItems("")
		p := helper.PagingLinks(all, page, app.Per_page(), app.Uri_position()+"list/?page=%d","href","a","","")
		l := model.PositionModel.GetList("", "", "", order)
		library.POSITION.SetTrees(l)
		library.POSITION.SetSelectOptionsView(library.POSITION.TreeList,-1,2)
		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("position")))
		if app.Per_page()+page > int64(len(library.POSITION.TreeList)) {
			w.Write([]byte(positionView.Index(library.POSITION.TreeList[page:len(library.POSITION.TreeList)],p)))
		}else {
			w.Write([]byte(positionView.Index(library.POSITION.TreeList[page:page+app.Per_page()],p)))
		}
	}
}

func (this *position) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, _ := this.crud.getAjaxList(false, w, r)
	var order = "position_sort ASC"
	if library.VALIDATION.Status == 0 {
		all := model.PositionModel.CountItems(search)
		paging := helper.PagingLinks(all, page, per_page, "%d", "data-page", "span","", `class="paging"`)
		var itemList []*entity.Position
		if search == ""{
			itemList = model.PositionModel.GetList(search, "", "", order)
			library.POSITION.SetTrees(itemList)
			library.POSITION.SetSelectOptionsView(library.POSITION.TreeList,-1,2)
		}else {
			itemList = model.PositionModel.GetList(search, pageStr, per_pageStr, order)
		}

		helper.SetAjaxHeaders(w)
		var limit = page+per_page
		if search == "" && limit > int64(len(library.POSITION.TreeList)) {
			limit = int64(len(library.POSITION.TreeList))
		}
		if search == "" {
			w.Write([]byte(positionView.Listing(library.POSITION.TreeList[page:limit], paging)))
		}else {
			w.Write([]byte(positionView.Listing(itemList, paging)))
		}

		return true
	}
	return false
}

func (this *position) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var positionObj = model.PositionModel.Get(idInt64)
		if positionObj != nil {

			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(positionObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["position_id"] = lang.T("not found")
		}
	}
	return false
}

func (this *position) FormHandler(action string,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64
	if action == "edit"{
		idInt64,_ = this.crud.edit(w,r)
	}
	parent, _ := library.VALIDATION.IsInt64(true, "parent", 20, r)
	if parent > 0{
		var position = model.PositionModel.Get(parent)
		if position == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["parent"] = lang.T("position is not exist")
		}
	}
	sort,_ := library.VALIDATION.IsInt64(false, "sort", 20, r)
	enable := library.VALIDATION.CheckBox("enable",r)
	var title = template.HTMLEscapeString(r.FormValue("title"))
	if library.VALIDATION.Status == 0 {
		var res int64
		if action == "edit"{
			position := entity.NewPosition(idInt64,parent,title,sort,enable)
			res = model.PositionModel.Edit(position)
		}
		if action == "add"{
			position := entity.NewPosition(0,parent,title,sort,enable)
			res = model.PositionModel.Add(position)
		}
		if res > 0{
			core.MALINA.SetGlobals()
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}
func (this *position) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.PositionModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not found")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}

