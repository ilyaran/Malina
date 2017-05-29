/**
 * Public controller.  github.com/ilyaran/Malina eCommerce application
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
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/libraries"
	"github.com/gorilla/mux"
	"github.com/ilyaran/Malina/models"
	"fmt"
	"github.com/ilyaran/Malina/helpers"
	"github.com/ilyaran/Malina/language"
	"time"
	"github.com/ilyaran/Malina/config"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"github.com/ilyaran/Malina/caching"

)

var PublicController = &publicController{crud:&CrudController{}}

type publicController struct {
	crud               *CrudController
	cart_id            string
	//CartPublicListJSON []byte
	//CartPublicList     []*entity.CartPublic
}

func (this *publicController) Index(w http.ResponseWriter, r *http.Request) {
	caching.T0 = time.Now()

	model.PublicModel.Query = ``
	model.PublicModel.TemporaryData = ``
	model.PublicModel.Where = ``
	model.PublicModel.All = 0

	library.SESSION.Authentication(w, r)

	///caching.T1 = time.Now()
	//fmt.Println("Auth After: ", caching.T1.Sub(caching.T0))

	views.LOCALS.W = w
	views.LOCALS.AccountCartList = nil
	views.LOCALS.AccountCartListJSON = nil

	// Fetch Cart cookie from database by cart_id from client cookie
	library.SESSION.GetCookie("cart_id", r)
	if library.SESSION.CookieString == `` {
		library.SESSION.SetCookie(
			"cart_id",
			library.SESSION.Cryptcode(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), library.SESSION.GetIP(r), app.Crypt_salt())),
			true,
			"/",
			app.Cart_cookie_expiration(),
			w)
	}

	switch mux.Vars(r)["action"] {
		case "crud" : this.Cart(w, r)
		case "form" : this.Order(w, r)
		case "get" : this.ProductItem(w, r)
		case "list" : this.ProductList(w, r)
		case "details" : this.CartDetails(w, r)
		case "ajax" : this.ProductAjaxList(w, r)
		case "ajax_list" :
			r.ParseForm()
			this.inlistUpdateQuery("cart","cart_quantity","integer",r)
			helper.SetAjaxHeaders(w)
			publicView.CartViewObj.CartProductsList = model.PublicModel.GetCartProducts(library.SESSION.CookieString)
			fmt.Fprintf(w, "")
		default     : this.Welcome(r)
	}

}

func (this *publicController) Order(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {



	}else {

	}
}

func (this *publicController) ProductItem(w http.ResponseWriter, r *http.Request) {
	defer func() {
		caching.T1 = time.Now()
		fmt.Println("Item Product run in: ", caching.T1.Sub(caching.T0))
		fmt.Println(`==========================================================`)
	}()
	idInt64, _ := library.VALIDATION.IsInt64(false, "id", 20, r)
	if library.VALIDATION.Status == 0 {
		publicView.ProductViewObj.Product = model.PublicModel.GetProduct(idInt64)
		if publicView.ProductViewObj.Product != nil {

			publicView.ProductViewObj.Item()
		}
	}
}

func (this *publicController) ProductList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		caching.T1 = time.Now()
		fmt.Println("Product List run in: ", caching.T1.Sub(caching.T0))
		fmt.Println(`==========================================================`)
	}()
	page,pageStr:=library.VALIDATION.IsInt64(false, "page", 20, r)
	order := "product_price ASC"
	if library.VALIDATION.Status == 0 {

		publicView.ProductViewObj.ProductList = model.PublicModel.GetList("","","",pageStr,strconv.FormatInt(app.Per_page(),10),order)
		publicView.ProductViewObj.Paging = helper.PagingLinks(model.PublicModel.All, page, app.Per_page(), "public/product/list/?page=%d","href","a","","")

		caching.T1 = time.Now()
		fmt.Println("Listing After: ", caching.T1.Sub(caching.T0))

		publicView.ProductViewObj.List()
	}
}

func (this *publicController) ProductAjaxList(w http.ResponseWriter, r *http.Request) {
	//page, pageStr, per_page, per_pageStr, search, orderInt := this.crud.getAjaxList(false, w, r)
	page, pageStr := library.VALIDATION.IsInt64(false, "page", 20, r)
	per_page, per_pageStr := library.VALIDATION.IsInt64(false, "per_page", 4, r)
	if per_page < app.Per_page() {
		per_page = app.Per_page()
		per_pageStr = strconv.FormatInt(app.Per_page(), 10)
	}
	orderInt, _ := library.VALIDATION.IsInt64(false, "order_by", 2, r)
	var search = r.FormValue("search")
	if search != "" {
		if len(search) > 64 {
			search = search[0:64]
		}
		search = strings.Trim(regexp.MustCompile(`[\W]+`).ReplaceAllString(search, "|"), "|")
	}

	this.crud.order_by = `product_price ASC`
	switch orderInt {
	case 2:this.crud.order_by = "product.product_price DESC"
	case 3:this.crud.order_by = "product.product_id ASC"
	case 4:this.crud.order_by = "product.product_id DESC"
	case 5:this.crud.order_by = "product.product_title ASC"
	case 6:this.crud.order_by = "product.product_title DESC"
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

		/*var itemList = */model.PublicModel.GetList(categoryWhere,priceInterval,search,pageStr,per_pageStr, this.crud.order_by)
		/*var paging =*/ helper.PagingLinks(model.PublicModel.All, page, per_page, "%d", "data-page", "span","", `class="paging"`)

		helper.SetAjaxHeaders(w)
		//w.Write([]byte(publicView.ProductListing(itemList, paging)))
	}
}


func (this *publicController) Welcome(r *http.Request) {
	//cache.T0 = time.Now()
	defer func() {
		caching.T1 = time.Now()
		fmt.Println("Welcome run in: ", caching.T1.Sub(caching.T0))
		fmt.Println(`==========================================================`)
	}()

	caching.T1 = time.Now()
	fmt.Println("Listing Before: ", caching.T1.Sub(caching.T0))

	if publicView.WelcomeViewObj.ProductList == nil{
		delete(caching.PublicPages,`welcome`)
		publicView.WelcomeViewObj.ProductList = model.PublicModel.GetList("", "", "", "0", "20", "product_created DESC")
	}

	caching.T1 = time.Now()
	fmt.Println("Listing After: ", caching.T1.Sub(caching.T0))

	publicView.WelcomeViewObj.Welcome()
}
func (this *publicController) CartDetails(w http.ResponseWriter, r *http.Request) {
	publicView.CartViewObj.CartProductsList = model.PublicModel.GetCartProducts(library.SESSION.CookieString)
	publicView.CartViewObj.Cart()
}
func (this *publicController) Cart(w http.ResponseWriter, r *http.Request) {
	defer func() {
		caching.T1 = time.Now()
		fmt.Println(`~~~~~~~~~~~~~Cart~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`)
		fmt.Println("Cart handler run in: ", caching.T1.Sub(caching.T0))
		fmt.Println(`~~~~~~~~~~~~~/Cart~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`)
	}()

	cart_action, _ := library.VALIDATION.IsInt64(false, "cart_action", 1, r)
	product_id, _ := library.VALIDATION.IsInt64(false, "product_id", 20, r)
	helper.SetAjaxHeaders(w)
	if product_id > 0 {
		if model.PublicModel.GetProduct(product_id) != nil {
			switch cart_action {
			case 1:
				model.PublicModel.Add_product(library.SESSION.CookieString, product_id)
			case 2:
				model.PublicModel.Del_product(library.SESSION.CookieString, product_id)
			}
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["product_id"] = lang.T("not found")
			return
		}
	}
	views.LOCALS.AccountCartListJSON, _ = json.Marshal(model.PublicModel.GetCartProducts(library.SESSION.CookieString))
	fmt.Fprintf(w, `{"Status":0,"Result":`+string(views.LOCALS.AccountCartListJSON)+`}`)
}

func (s *publicController) inlistUpdateQuery(dbtable, columnName, columnType string, r *http.Request) {
	var valueNum int = 1 // if count numbers 1,2,3 ... to "$1,$2,$3,$4,$5 ..." string
	var exec []interface{}
	var layout_item_query = `
	SELECT ` + dbtable + `_update('` + s.cart_id + `', ARRAY[ %s ], ARRAY[ %s ])`
	var value string
	var key string
	var match bool
	for _, v1 := range r.PostForm[columnName + "_inlist[]"] {
		values_arr := strings.SplitN(v1, "|", 2)
		match, _ = regexp.MatchString(`^[0-9]{1,20}$`, values_arr[0])
		if match {
			key += `,` + values_arr[0]
			// push value to exec interface array
			switch columnType {
			case "integer":
				match, _ = regexp.MatchString(`^[0-9]{1,20}$`, values_arr[1])
				if match {
					value += `,` + values_arr[1]
				}
			case "numeric":
				match, _ = regexp.MatchString(`^[0-9]{1,20}(\.[0-9]{0,2})?$`, values_arr[1])
				if match {
					value += `,` + values_arr[1]
				}
			case "boolean":
				// building "TRUE,TRUE,TRUE,FALSE, ..." string to values array: unnest(ARRAY[TRUE,TRUE,TRUE,FALSE, ...]) AS v
				valueBool, _ := strconv.ParseBool(values_arr[1])
				if valueBool {
					value += `,TRUE`
				} else {
					value += `,FALSE`
				}
			case "string":
				match, _ = regexp.MatchString(app.Pattern_title(), values_arr[0])
				if match {
					// building "$1,$2,$3,$4,$5 ..." string to values array: unnest(ARRAY[$1,$2,$3,$4,$5 ...]) AS v
					value += `,$` + strconv.Itoa(valueNum)
					exec = append(exec, values_arr[1])
				}
			}
			valueNum++
		}
	}
	if key != `` {
		model.Crud.GetRow(fmt.Sprintf(layout_item_query, key[1:], value[1:]),exec)
	}
}

