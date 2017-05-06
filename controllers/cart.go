/**
 * Cart controller class.  Malina eCommerce application
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
	"Malina/views/cart"
	"Malina/helpers"
	"Malina/models"
	"Malina/config"
	"fmt"
	"encoding/json"
	"Malina/language"
	"Malina/entity"
	"Malina/views"
	"Malina/libraries"
	"strconv"
)

var Cart = &cartController{&CrudController{}}
type cartController  struct {
	crud *CrudController

}

func (this *cartController) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("cart",w,r)

	model.CartModel.Query = ``
	model.CartModel.Where = ``
	model.CartModel.All = 0

	switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler('a',w,r)
		case "edit" 	 : this.FormHandler('e',w,r)
		case "del" 	 : this.Del(w,r)
		default		 : this.List(w,r); return // "list"
	}
	helper.SetAjaxHeaders(w)

	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *cartController) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false,w,r)
	if library.VALIDATION.Status == 0 {

		itemList := model.CartModel.GetList("","",pageStr,strconv.FormatInt(app.Per_page(),10),"cart_created DESC")
		paging := helper.PagingLinks(model.CartModel.All, page, app.Per_page(), "%d", "data-page", "span","", `class="paging"`)

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("cart")))
		w.Write([]byte(cartView.Index(itemList, paging)))
	}
}

func (this *cartController) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, orderInt := this.crud.getAjaxList(false, w, r)
	this.crud.order_by = `cart_created DESC`
	switch orderInt {
		case 2:this.crud.order_by = "cart_created ASC"
		case 3:this.crud.order_by = "cart_product ASC"
		case 4:this.crud.order_by = "cart_product DESC"
		case 5:this.crud.order_by = "cart_id ASC"
		case 6:this.crud.order_by = "cart_id ASC"
	}
	_, product_idStr := library.VALIDATION.IsInt64(true, "cart_product", 20, r)
	if library.VALIDATION.Status == 0 {
		itemList := model.CartModel.GetList(search,product_idStr,pageStr,per_pageStr, this.crud.order_by)
		paging := helper.PagingLinks(model.CartModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

		helper.SetAjaxHeaders(w)
		w.Write([]byte(cartView.Listing(itemList, paging)))
		return true
	}
	return false
}

func (this *cartController) Get(w http.ResponseWriter, r *http.Request) bool{
	_,idInt64Str := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var cartObj = model.CartModel.Get(idInt64Str)
		if cartObj != nil {
			helper.SetAjaxHeaders(w)
			out, _ := json.Marshal(cartObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["cart_id"] = lang.T("not found")
		}
	}
	return false
}
func (this *cartController) FormHandler(action byte,w http.ResponseWriter, r *http.Request) {
	var idString string
	if action == 'e'{
		idString = library.VALIDATION.IsValidText(true,"cart_id",false,12,128,`^[0-9a-zA-Z]+$`,"alpha numeric",r)
		var cart = model.CartModel.Get(idString)
		if cart == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["cart_id"] = lang.T("not found")
		}
	}
	var buy_now = library.VALIDATION.CheckBox("buy_now",r)
	quantity, _ := library.VALIDATION.IsFloat64(false, "quantity", 20,2, r)
	product_id,_ := library.VALIDATION.IsInt64(true, "product_id", 20, r)
	var product = model.ProductModel.Get(product_id)
	if product == nil {
		library.VALIDATION.Status = 100
		library.VALIDATION.Result["product_id"] = lang.T("not found")
	}

	if library.VALIDATION.Status == 0 {
		var res int64
		cart := entity.NewCart(idString, product_id, quantity, buy_now)
		res = model.CartModel.Add(cart)
		if res > 0{
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}

func (this *cartController) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.CartModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not find")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}



