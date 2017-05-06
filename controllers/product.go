/**
 * Product Controller class.  Malina eCommerce application
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
	"Malina/views/product"
	"Malina/config"
	"fmt"
	"encoding/json"
	"Malina/language"
	"Malina/entity"
	"Malina/views"
	"Malina/libraries"
	"html/template"
	"strconv"
)

var Product = &ProductController{&CrudController{}}
type ProductController  struct {
	crud *CrudController
}

func (this *ProductController) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("product",w,r)
	model.ProductModel.Query = ``
	model.ProductModel.Where = ``
	model.ProductModel.All = 0

	switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler('a',w,r)
		case "edit" 	 : this.FormHandler('e',w,r)
		case "inlist" 	 : this.Inlist(w,r)
		case "del" 	 : this.Del(w,r)
		default		 : this.List(w,r); return // "list"
	}
	helper.SetAjaxHeaders(w)

	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *ProductController) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false,w,r)
	order := "product_price ASC"
	if library.VALIDATION.Status == 0 {
		//all := model.ProductModel.CountItems(nil,0,0,false)
		itemList := model.ProductModel.GetList("","","",pageStr,strconv.FormatInt(app.Per_page(),10),order,"")
		paging := helper.PagingLinks(model.ProductModel.All, page, app.Per_page(), app.Uri_product()+"list/?page=%d","href","a","","")


		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("product")))
		w.Write([]byte(productView.Index(itemList, paging)))
	}
}

func (this *ProductController) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, orderInt := this.crud.getAjaxList(false, w, r)
	this.crud.order_by = `product_updated DESC`
	switch orderInt {
		case 2:this.crud.order_by = "product.product_updated DESC"
		case 3:this.crud.order_by = "product.product_price ASC"
		case 4:this.crud.order_by = "product.product_price DESC"
		case 5:this.crud.order_by = "product.product_id ASC"
		case 6:this.crud.order_by = "product.product_id DESC"
	}
	price_min,price_minStr := library.VALIDATION.IsFloat64(false,"price_min",20,2,r)
	price_max,price_maxStr := library.VALIDATION.IsFloat64(false,"price_max",20,2,r)

	var categoryWhere = ""
	category_id, _ := library.VALIDATION.IsInt64(false, "category", 20, r)
	if category_id > 0{
		categoryObj,ok := library.CATEGORY.TreeMap[category_id]
		if ok {
			categoryWhere = "category_id IN (" + categoryObj.Get_descendantIdsString() + ")"
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["category"] = lang.T("category is not exist")
		}
	}
	if library.VALIDATION.Status == 0 {
		var priceInterval = ""
		if price_min > 0 && price_max <= 0 {
			priceInterval = "product_price > " + price_minStr
		}
		if price_min <= 0 && price_max > 0 {
			priceInterval = "product_price BETWEEN 0 AND " + price_maxStr
		}
		if price_min > 0 && price_max > 0 {
			priceInterval = "product_price BETWEEN "+price_minStr+" AND "+price_maxStr
		}

		var itemList = model.ProductModel.GetList(categoryWhere,priceInterval,search,pageStr,per_pageStr, this.crud.order_by,"")
		var paging = helper.PagingLinks(model.ProductModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

		helper.SetAjaxHeaders(w)
		w.Write([]byte(productView.Listing(itemList, paging)))

		return true
	}
	return false
}

func (this *ProductController) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var productObj = model.ProductModel.Get(idInt64)
		if productObj != nil {
			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(productObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["product_id"] = lang.T("not found")
		}
	}
	return false
}

func (this *ProductController) FormHandler(action byte,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64 = 0
	var productObj *entity.Product
	if action == 'e'{
		idInt64,_ = this.crud.edit(w,r)
		productObj = model.ProductModel.Get(idInt64)
		if productObj == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["product"] = lang.T("product is not exist")
		}
	}
	// validating string like this:  "/a/d/c.jpg|/fd/rt/img.png|/some/path/image.gif"
	var imgUrls = library.VALIDATION.ImgUrls("img",r)
	category_id, _ := library.VALIDATION.IsInt64(true, "category_id", 20, r)
	if category_id > 0{
		var categoryObj = model.CategoryModel.Get(category_id)
		if categoryObj == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["category_id"] = lang.T("category is not exist")
		}
	}
	price,_ := library.VALIDATION.IsFloat64(false,"price",20,2,r)
	price1,_ := library.VALIDATION.IsFloat64(false,"price1",20,2,r)
	enable := library.VALIDATION.CheckBox("enable",r)
	var title = template.HTMLEscapeString(r.FormValue("title"))
	var description = r.FormValue("description")
	code := library.VALIDATION.IsValidText(false,"code",false,12,512,`^[\w]+$`,"alpha numeric",r)

	if library.VALIDATION.Status == 0 {

		product := entity.NewProduct(idInt64,category_id,price,price1,code,title,description,enable)
		var res int64
		if action == 'a'{
			res = model.ProductModel.Add(product,imgUrls)
		}else {
			res = model.ProductModel.Edit(product,imgUrls)
		}

		if res > 0{
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}
func (this *ProductController) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.ProductModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not found")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}
func (this *ProductController) Inlist(w http.ResponseWriter, r *http.Request) {
	var columns = map[string]string{"product_price":"numeric","product_price1":"numeric","product_code":"string","product_title":"string","product_enable":"boolean"}
	this.crud.inlist(columns,r)
}