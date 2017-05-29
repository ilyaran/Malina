/**
 * Permissions controller class.  github.com/ilyaran/Malina eCommerce application
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
	"github.com/ilyaran/Malina/config"
	"fmt"
	"encoding/json"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/views"
	"html/template"
	"github.com/ilyaran/Malina/libraries"
	"github.com/ilyaran/Malina/views/permission"
	"github.com/ilyaran/Malina/core"
	"strconv"
)

var Permission = &permission{&CrudController{}}
type permission  struct {
	crud *CrudController
}

func (this *permission) Index(w http.ResponseWriter, r *http.Request) {

	if this.crud.hasPermission("permission","permission_id",w,r) {

		model.PositionModel.Query = ``
		model.PositionModel.Where = ``
		model.PositionModel.All = 0
		model.PositionModel.Exec = nil

		model.PermissionModel.Query = ``
		model.PermissionModel.Where = ``
		model.PermissionModel.All = 0
		model.PermissionModel.Exec = nil

		switch this.crud.action {
			case "list_ajax" : if this.AjaxList(w,r) {return}
			case "inlist"    : this.Inlist(w, r)
			case "get" 	 : if this.Get(w,r) {return}
			case "add" 	 : this.FormHandler(this.crud.action,w,r)
			case "edit" 	 : this.FormHandler(this.crud.action,w,r)
			case "del" 	 : this.Del(w,r)
			default		 : this.List(w,r); return // "list"
		}
	}

	helper.SetAjaxHeaders(w)
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *permission) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false,w,r)
	order := "permission_id ASC"
	if library.VALIDATION.Status == 0 {

		itemList := model.PermissionModel.GetList("", pageStr, strconv.FormatInt(app.Per_page(), 10), order)
		paging := helper.PagingLinks(model.PermissionModel.All, page, app.Per_page(), app.Uri_account()+"list/?page=%d","href","a","","")

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("permissions")))
		w.Write([]byte(permissionView.Index(itemList, paging)))
	}
}

func (this *permission) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "permission_id ASC"
	switch order_by {
	case 2:order = "permission_id DESC"
	case 3:order = "permission_data ASC"
	case 4:order = "permission_data DESC"
	}
	if library.VALIDATION.Status == 0 {

		itemList := model.PermissionModel.GetList(search, pageStr, per_pageStr, order)
		paging := helper.PagingLinks(model.PermissionModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

		helper.SetAjaxHeaders(w)
		w.Write([]byte(permissionView.Listing(itemList, paging)))
		return true
	}
	return false
}

func (this *permission) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	var permissionObj = this.CheckIfExistsItem(idInt64)
	if library.VALIDATION.Status == 0 {
		out, _ := json.Marshal(permissionObj)
		fmt.Fprintf(w, string(out))
		return true
	}
	return false
}

func (this *permission) FormHandler(action string,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64
	var permissionObj *entity.Permission
	if action == "edit"{
		idInt64,_ = this.crud.edit(w,r)
		permissionObj = this.CheckIfExistsItem(idInt64)
	}
	positionId, _ := library.VALIDATION.IsInt64(true, "position_id", 20, r)
	var positionObj *entity.Position = PositionController.CheckIfExistsItem("position_id",positionId,false)

	var data = template.HTMLEscapeString(r.FormValue("data"))
	if library.VALIDATION.Status == 0 {
		var res int64
		if action == "edit"{
			permissionObj.SetData(data)
			permissionObj.SetPosition(positionObj)
			model.PermissionModel.Exec = []interface{}{data,positionId,idInt64}
			res = model.PermissionModel.Edit(permissionObj)
		}
		if action == "add"{
			//permissionObj = entity.NewPermission(data,positionId)
			model.PermissionModel.Exec = []interface{}{data,positionId}
			res = model.PermissionModel.Add()
		}
		if res > 0{
			core.MALINA.SetPositionGlobals()
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}

func (this *permission) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	this.CheckIfExistsItem(idInt64)
	if library.VALIDATION.Status == 0 {
		model.PermissionModel.Exec = []interface{}{idInt64}
		res := model.PermissionModel.Del()
		if res > 0 {
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T("server error")
		}
	}
}

func (this *permission) Inlist(w http.ResponseWriter, r *http.Request) {
	var columns = map[string]string{"permission_data":"string"}
	this.crud.inlist(columns, r)
}

func (this *permission) CheckIfExistsItem (idInt64 int64) *entity.Permission {
	var permissionObj *entity.Permission
	if idInt64 > 0{
		permissionObj = model.PermissionModel.Get(idInt64)
		if permissionObj == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["id"] = lang.T("permission doesn't exist")
		}else {
			if library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetParent().GetId() > 0 {
				if _, ok := library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetDescendantIdsMap()[permissionObj.GetPosition().GetId()]; ok {

					return permissionObj

				}else {
					library.VALIDATION.Status = 160
					library.VALIDATION.Result["id"] = lang.T("no permission")
				}
			}
		}
	}
	return nil
}