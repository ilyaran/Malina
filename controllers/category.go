/**
 * Category controller class.  Malina eCommerce application
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
	"Malina/views/category"
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
	"Malina/core"
)

var Category = &CategoryController{&CrudController{}}
type CategoryController  struct { crud *CrudController }

func (this *CategoryController) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("category",w,r)
	model.CategoryModel.Query = ``
	model.CategoryModel.Where = ``
	model.CategoryModel.All = 0
	switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler('a',w,r)
		case "edit" 	 : this.FormHandler('e',w,r)
		case "del" 	 : this.Del(w,r)
		case "inlist" 	 : this.Inlist(w,r)
		default		 : this.List(w,r); return // "list"
	}
	helper.SetAjaxHeaders(w)

	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	page, _ := this.crud.getList(false,w,r)
	order := "category_sort ASC"
	if library.VALIDATION.Status == 0 {
		itemList := model.CategoryModel.GetList("", "", "", order)
		paging := helper.PagingLinks(model.CategoryModel.All, page, app.Per_page(), app.Uri_category()+"list/?page=%d","href","a","","")

		library.CATEGORY.SetTrees(itemList)

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("category")))
		if app.Per_page()+page > int64(len(library.CATEGORY.TreeList)) {
			w.Write([]byte(categoryView.Index(library.CATEGORY.TreeList[page:len(library.CATEGORY.TreeList)], paging)))
		}else {
			w.Write([]byte(categoryView.Index(library.CATEGORY.TreeList[page:page+app.Per_page()], paging)))
		}
	}
}

func (this *CategoryController) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "category_sort ASC"
	switch order_by {
		case 2:order = "category_sort DESC"
		case 3:order = "category_id ASC"
		case 4:order = "category_id DESC"
		case 5:order = "category_title ASC"
		case 6:order = "category_title DESC"
	}
	if library.VALIDATION.Status == 0 {
		var itemList []*entity.Category
		if search == ""{
			itemList = model.CategoryModel.GetList(search, "", "", order)
			library.CATEGORY.SetTrees(itemList)
		}else {
			itemList = model.CategoryModel.GetList(search, pageStr, per_pageStr, order)
		}
		paging := helper.PagingLinks(model.CategoryModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

		helper.SetAjaxHeaders(w)
		var limit = page+per_page
		if search == "" && limit > int64(len(library.CATEGORY.TreeList)) {
			limit = int64(len(library.CATEGORY.TreeList))
		}
		if search == "" {
			w.Write([]byte(categoryView.Listing(library.CATEGORY.TreeList[page:limit], paging)))
		}else {
			w.Write([]byte(categoryView.Listing(itemList, paging)))
		}

		return true
	}
	return false
}

func (this *CategoryController) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var categoryObj = model.CategoryModel.Get(idInt64)
		if categoryObj != nil {

			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(categoryObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["category_id"] = lang.T("not found")
		}
	}
	return false
}

func (this *CategoryController) FormHandler(action byte,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64 = 0
	var categoryObj *entity.Category
	if action == 'e'{
		idInt64,_ = this.crud.edit(w,r)
		categoryObj = model.CategoryModel.Get(idInt64)
		if categoryObj == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["category"] = lang.T("category is not exist")
		}
	}
	// validating string like this:  /a/d/c.jpg|/fd/rt/img.png
	var imgUrls = library.VALIDATION.ImgUrls("img",r)
	parent, _ := library.VALIDATION.IsInt64(true, "parent", 20, r)
	if parent > 0{
		var category = model.CategoryModel.Get(parent)
		if category == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["parent"] = lang.T("parent is not exist")
		}
	}
	sort,_ := library.VALIDATION.IsInt64(false, "sort", 20, r)
	enable := library.VALIDATION.CheckBox("enable",r)
	var title = template.HTMLEscapeString(r.FormValue("title"))
	var description = r.FormValue("description")

	if library.VALIDATION.Status == 0 {
		var res int64
		category := entity.NewCategory(idInt64,parent,sort,title,description,enable)
		if action == 'e'{
			res = model.CategoryModel.Edit(category,imgUrls)
		}
		if action == 'a'{
			res = model.CategoryModel.Add(category,imgUrls)
		}
		if res > 0{
			core.MALINA.SetCategoryGlobals()
			library.VALIDATION.Status = 0
			library.VALIDATION.Result["select_parent"] = views.FACE_FORM.CategorySelect("Root","parent")
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}
func (this *CategoryController) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.CategoryModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not found")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}
func (this *CategoryController) Inlist(w http.ResponseWriter, r *http.Request) {
	var columns = map[string]string{"category_sort":"integer","category_title":"string","category_enable":"boolean"}
	this.crud.inlist(columns,r)

}