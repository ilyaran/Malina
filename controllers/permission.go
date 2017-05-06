package controller

import (
	"net/http"
	"Malina/helpers"
	"Malina/models"
	"Malina/config"
	"fmt"
	"encoding/json"
	"Malina/language"
	"Malina/entity"
	"Malina/views"
	"html/template"
	"Malina/libraries"
	"Malina/views/permission"
	"strconv"
)

var Permission = &permission{&CrudController{}}
type permission  struct {
	crud *CrudController
}

func (this *permission) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("permission",w,r)
	if action!=""{
		switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler(action,w,r)
		case "edit" 	 : this.FormHandler(action,w,r)
		case "del" 	 : this.Del(w,r)
		default		 : this.List(w,r); return // "list"
		}
		helper.SetAjaxHeaders(w)
	}
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *permission) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false,w,r)
	order := "permission_id ASC"
	if library.VALIDATION.Status == 0 {
		all := model.PermissionModel.CountItems("")
		paging := helper.PagingLinks(all, page, app.Per_page(), app.Uri_account()+"list/?page=%d","href","a","","")
		itemList := model.PermissionModel.GetList("", pageStr, strconv.FormatInt(app.Per_page(), 10), order)

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("permissions")))
		w.Write([]byte(permissionView.Index(itemList, paging)))
	}
}

func (this *permission) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "permission_id ASC"
	switch order_by {
	case 2:order = "permission_id ASC"
	}
	if library.VALIDATION.Status == 0 {
		all := model.PermissionModel.CountItems(search)
		paging := helper.PagingLinks(all, page, per_page, "%d", "data-page", "span","", `class="paging"`)
		itemList := model.PermissionModel.GetList(search, pageStr, per_pageStr, order)
		helper.SetAjaxHeaders(w)
		w.Write([]byte(permissionView.Listing(itemList, paging)))
		return true
	}
	return false
}

func (this *permission) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var permissionObj = model.PermissionModel.Get(idInt64)
		if permissionObj != nil {

			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(permissionObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["permission_id"] = lang.T("not found")
		}
	}
	return false
}

func (this *permission) FormHandler(action string,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64
	var permission *entity.Permission
	if action == "edit"{
		idInt64,_ = this.crud.edit(w,r)
		if idInt64 > 0{
			permission = model.PermissionModel.Get(idInt64)
			if permission == nil {
				library.VALIDATION.Status = 100
				library.VALIDATION.Result["permission_id"] = lang.T("permission is not exist")
			}
		}
	}
	positionId, _ := library.VALIDATION.IsInt64(true, "position_id", 20, r)
	if positionId > 0{
		var position = model.PositionModel.Get(positionId)
		if position == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["position_id"] = lang.T("position is not exist")
		}
	}
	var data = template.HTMLEscapeString(r.FormValue("data"))
	if library.VALIDATION.Status == 0 {
		var res int64
		if action == "edit"{
			permission = entity.NewPermission(idInt64,data,positionId)
			res = model.PermissionModel.Edit(permission)
		}
		if action == "add"{
			permission = entity.NewPermission(0,data,positionId)
			res = model.PermissionModel.Add(permission)
		}
		if res > 0{
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}
func (this *permission) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.PermissionModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not found")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}


