/**
 * Position controller class.  github.com/ilyaran/Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package controller

import (
	"net/http"
	"github.com/ilyaran/Malina/helpers"
	"github.com/ilyaran/Malina/models"
	"fmt"
	"encoding/json"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/views"
	"html/template"
	"github.com/ilyaran/Malina/libraries"
	"github.com/ilyaran/Malina/views/position"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/core"
)

var PositionController = &positionController{crud: &CrudController{}}
type positionController struct {
	crud *CrudController
}

func (this *positionController) Index(w http.ResponseWriter, r *http.Request) {
	this.crud.hasPermission("position", "position_id", w, r)
	if library.VALIDATION.Status == 0 {
		model.PositionModel.Query = ``
		model.PositionModel.Where = ``
		model.PositionModel.All = 0
		model.PositionModel.Exec = nil

		switch this.crud.action {
		case "list_ajax":
			if this.AjaxList(w, r) {
				return
			}
		case "inlist":
			this.Inlist(w, r)
		case "get":
			if this.Get(w, r) {
				return
			}
		case "add":
			this.FormHandler(this.crud.action, w, r)
		case "edit":
			this.FormHandler(this.crud.action, w, r)
		case "del":
			this.Del(w, r)
		default:
			this.List(w, r); return // "list"
		}
	}
	helper.SetAjaxHeaders(w)
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *positionController) List(w http.ResponseWriter, r *http.Request) {
	page, _ := this.crud.getList(false,w,r)
	order := "position_sort ASC"
	if library.VALIDATION.Status == 0 {

		l := model.PositionModel.GetList("", "", "", order,library.POSITION.TreeMap)
		p := helper.PagingLinks(model.PositionModel.All, page, app.Per_page(), app.Uri_position()+"list/?page=%d","href","a","","")

		library.POSITION.SetTrees(l)

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("positionController")))

		if app.Per_page()+page > int64(len(library.POSITION.TreeList)) {
			w.Write([]byte(positionView.Index(library.POSITION.TreeList[page:len(library.POSITION.TreeList)],p)))
		}else {
			w.Write([]byte(positionView.Index(library.POSITION.TreeList[page:page+app.Per_page()],p)))
		}
	}
}

func (this *positionController) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "position_sort ASC"
	switch order_by {
	case 2:order = "position_sort DESC"
	case 3:order = "position_id ASC"
	case 4:order = "position_id DESC"
	case 5:order = "position_title ASC"
	case 6:order = "position_title DESC"
	}
	if library.VALIDATION.Status == 0 {
		var itemList []*entity.Position
		if search == ""{
			itemList = model.PositionModel.GetList(search, "", "", order,library.POSITION.TreeMap)
			library.POSITION.SetTrees(itemList)
		}else {
			itemList = model.PositionModel.GetList(search, pageStr, per_pageStr, order,library.POSITION.TreeMap)
		}
		paging := helper.PagingLinks(model.PositionModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

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

func (this *positionController) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	var positionObj = this.CheckIfExistsItem("id",idInt64,false)
	if positionObj != nil {
		helper.SetAjaxHeaders(w)
		positionObj.Parent = &entity.Position{Id:positionObj.GetParent().GetId()}
		positionObj.PermissionsMap = nil
		out, _ := json.Marshal(struct {
			ID int64 `json:"id"`
			Title string `json:"title"`
			Parent *entity.Position `json:"parent"`
			Enable bool `json:"enable"`
			Sort int64 `json:"sort"`
		}{positionObj.Id,positionObj.Title,positionObj.Parent,positionObj.Enable,positionObj.Sort})
		fmt.Fprintf(w, string(out))
		return true
	}
	return false
}

func (this *positionController) FormHandler(action string,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64
	if action == "edit"{
		idInt64,_ = this.crud.edit(w,r)
		this.CheckIfExistsItem("id",idInt64,false)
	}
	parentPositionId, _ := library.VALIDATION.IsInt64(true, "parent", 20, r)
	this.CheckIfExistsItem("parent",parentPositionId,true)

	sort,_ := library.VALIDATION.IsInt64(false, "sort", 20, r)
	enable := library.VALIDATION.CheckBox("enable",r)
	var title = template.HTMLEscapeString(r.FormValue("title"))
	if title == ``{
		library.VALIDATION.Status = 100
		library.VALIDATION.Result["title"] = lang.T("title is required")
	}
	if library.VALIDATION.Status == 0 {
		var res int64

		if action == "edit"{
			if parentPositionId > 0 {
				model.PositionModel.Exec = []interface{}{parentPositionId,title,sort,enable,idInt64}
			}else {
				model.PositionModel.Exec = []interface{}{title,sort,enable,idInt64}
			}
			res = model.PositionModel.Edit(parentPositionId)
		}
		if action == "add"{
			if parentPositionId > 0 {
				model.PositionModel.Exec = []interface{}{parentPositionId,title,sort,enable}
			}else {
				model.PositionModel.Exec = []interface{}{title,sort,enable}
			}
			res = model.PositionModel.Add(parentPositionId)
		}
		if res > 0{
			core.MALINA.SetPositionGlobals()
			library.VALIDATION.Status = 0
			library.VALIDATION.Result["parent"] = positionView.GetSelectOptionsListView(library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()],"Root",true,true)
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}

func (this *positionController) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	this.CheckIfExistsItem("id",idInt64,false)
	if library.VALIDATION.Status == 0 {
		model.PositionModel.Exec = []interface{}{idInt64}
		res := model.PositionModel.Del()
		if res > 0 {
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T("server error")
		}
	}
}
func (this *positionController) Inlist(w http.ResponseWriter, r *http.Request) {
	var columns = map[string]string{"position_title":"string", "position_sort":"integer", "position_enable":"boolean"}
	this.crud.inlist(columns, r)
}

func (this *positionController) CheckIfExistsItem (keyName string,idInt64 int64, isParent bool) *entity.Position {
	var positionObj *entity.Position
	if idInt64 > 0 {
		var ok bool
		if positionObj,ok = library.POSITION.TreeMap[idInt64]; ok {
			if library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetParent().GetId() > 0 {
				if _, ok := library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetDescendantIdsMap()[positionObj.GetParent().GetId()]; ok {
					return positionObj
				}else if isParent && positionObj.GetId() == library.SESSION.SessionObj.GetPositionId() {
					return positionObj
				}else if positionObj.GetParent().GetId() == library.SESSION.SessionObj.GetPositionId() {
					return positionObj
				}
				library.VALIDATION.Status = 160
				library.VALIDATION.Result[keyName] = lang.T("no permission")
			}
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result[keyName] = lang.T("position doesn't exist")
		}
	}
	return nil
}
