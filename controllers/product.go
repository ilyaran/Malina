/**
 *
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ package controllers

import (
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/filters"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/lang"
	"github.com/ilyaran/Malina/berry"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/libraries"
	"html/template"
	"fmt"
	"github.com/ilyaran/Malina/views/publicView"
)

var ProductController *Product
func ProductControllerInit(){
	ProductController = &Product{base:&base{
		dbtable:"product",
		item:&entity.Product{},
		selectSqlFieldsDefault 	:
		`product_id,
		 product_code,
		 coalesce(product_category,0),
		 product_parameter,
		 product_img,
		 product_title,
		 coalesce(product_description,''),
		 product_price,
		 product_price1,
		 product_price2,
		 coalesce(product_quantity,0),
		 product_sold,
		 product_views,
		 product_created,
		 product_updated,
		 product_enable`,
		}}

	//helpers.GenerateControllerFields(ProductController.base.dbtable,ProductController.base.item)

	// model init
	ProductController.model = models.Product{}
	// end model init

	// orders init
	ProductController.base.orderList =[][2]string{
		{"ORDER BY product_created DESC",lang.T("created")+`&darr;`},
		{"ORDER BY product_created ASC",lang.T("created")+`&uarr;`},
		{"ORDER BY product_updated DESC",lang.T("updated")+`&darr;`},
		{"ORDER BY product_updated ASC",lang.T("updated")+`&uarr;`},
		{"ORDER BY product_price DESC",lang.T("price")+`&darr;`},
		{"ORDER BY product_price ASC",lang.T("price")+`&uarr;`},
	}

	if ProductController.base.orderList==nil || len(ProductController.base.orderList) < 0{
		panic("order list is not init")
	}
	ProductController.base.orderListLength = int64(len(ProductController.base.orderList))
	// end orders init

	ProductController.base.searchSqlTemplate=`product_title LIKE '%~%' OR product_description LIKE '%~%' OR product_code LIKE '%~%'`

	ProductController.base.inlistSqlFields = map[string]byte{
		"product_title"     :'s',
		"product_price"     :'n',
		"product_price1"    :'n',
		"product_code"      :'s',
		"product_enable"    :'b',
	}

	// view init
	ProductController.view = views.Product{}
	for k,_ := range ProductController.base.inlistSqlFields{
		ProductController.view.InlistFields += `<input class="inlist_fields" type="hidden" value="`+k+`">`
	}
	for k,v := range ProductController.base.orderList{
		ProductController.view.OrderSelectOptions += `<option value="`+strconv.Itoa(k)+`">`+v[1]+`</option>`
	}

	ProductController.viewPublic=publicView.Product{}
	// end view init

}

type Product struct {
	base     *base
	model			models.Product
	view            views.Product
	viewPublic      publicView.Product
}

func(s *Product)Index(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	malina.Controller = s
	malina.TableSql = s.base.dbtable
	s.base.index(malina,w,r,"")
}

func(s *Product)GetList(malina *berry.Malina, w http.ResponseWriter, r *http.Request, condition string){

	category_id,category_idStr := filters.IsUint(malina,false,"category_id",20,r)
	price_min,price_minStr := filters.IsUint(malina,false,"price_min",20,r)
	price_max,price_maxStr := filters.IsUint(malina,false,"price_max",20,r)

	if malina.Status > 0 {
		return
	}

	var category *entity.Category
	if category_id > 0 {
		if malina.Department=="public"{
			if v,ok := libraries.CategoryLib.MapPublic[category_id]; ok {
				category=v
				models.CrudGeneral.WhereAnd(malina,"product_category IN ","("+category_idStr+v.DescendantsString+")")
				malina.Url+="&category_id="+category_idStr
				malina.CategoryParameters=v.ParametersListView
			}
		}
		if malina.Department=="home"{
			if v,ok := libraries.CategoryLib.Trees.Map[category_id]; ok {
				models.CrudGeneral.WhereAnd(malina,"product_category IN ","("+category_idStr+v.(*entity.Category).DescendantsString+")")
				malina.Url+="&category_id="+category_idStr
			}
		}
	}
	if price_min > 0 && price_max == 0 {
		models.CrudGeneral.WhereAnd(malina,"product_price >",price_minStr)
		malina.Url+="&price_min="+price_minStr
	}else if price_min == 0 && price_max > 0 {
		models.CrudGeneral.WhereAnd(malina,"product_price <",price_maxStr)
		malina.Url+="&price_max="+price_maxStr
	}else if price_min > 0 && price_max > 0 {
		models.CrudGeneral.WhereAnd(malina,"product_price >",price_minStr)
		models.CrudGeneral.WhereAnd(malina,"product_price <",price_maxStr)
		malina.Url+="&price_max="+price_maxStr+"&price_min="+price_minStr
	}
	if malina.Department == "public" {
		models.CrudGeneral.WhereAnd(malina,"product_enable =","TRUE")
	}

	s.base.getList(malina,w,r)

	if malina.Device == "browser" {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

			if malina.Department == "public" {
				malina.Paging = s.base.paging(malina,"","")
				w.Write([]byte(s.viewPublic.Listing(malina,category)))
			}else if malina.Department == "home"{
				malina.Paging = s.base.paging(malina,"","")
				w.Write([]byte(s.view.Listing(malina)))
			}

		}else {
			if malina.Department == "public" {
				malina.Paging = s.base.paging(malina,"",app.Url_public_product_list+"?"+malina.Url[1:])
				s.viewPublic.Index(malina,category,w)
			}else if malina.Department == "home"{
				malina.Paging = s.base.paging(malina,"",app.Url_home_product_list+"?"+malina.Url[1:])
				s.view.Index(malina,w)
			}
		}
	}else {
		out, _ := json.Marshal(malina.List)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(out)
	}

}

func (s *Product) FormHandler(malina *berry.Malina,action byte,w http.ResponseWriter,r *http.Request) {
	// <- 1
	malina.ChannelBool = make(chan bool)
	if action=='e'{
		malina.IdInt64, malina.IdStr = filters.IsUint(malina,true, "id", 20, r)
		go s.base.isExistItem(malina, malina.IdInt64,"",false,false)
	}
	// end <- 1

	imgUrls, _, _ := filters.ImgUrls(malina,false, "img", r)

	category_id,_ := filters.IsUint(malina,true,"category_id",20,r)
	if category_id > 0 && malina.Status == 0 {
		if _,ok:=libraries.CategoryLib.Trees.Map[category_id];!ok{
			malina.Status = http.StatusNotAcceptable
			malina.Result["category_id"]=lang.T("not found")
		}
	}
	price,_ := filters.IsFloat64(malina,false,"price",20,2,r)
	price1,_ := filters.IsFloat64(malina,false,"price1",20,2,r)
	code := filters.IsValidText(malina,false,"code",false,1,255,`^[\w]+$`,"alpha numeric",r)
	enable := filters.CheckBox(malina,"enable",r)

	title := template.HTMLEscapeString(r.FormValue("title"))
	if title==""{
		malina.Status = http.StatusNotAcceptable
		malina.Result["title"]=lang.T("required")
	}else{
		if len(title) > app.Title_max_len{
			malina.Status = http.StatusNotAcceptable
			malina.Result["title"]=lang.T("length must have no more then:")+strconv.Itoa(app.Title_max_len)+" symbols"
		}
	}
	description := r.FormValue("description")
	short_description := template.HTMLEscapeString(r.FormValue("short_description"))
	if short_description!="" && len(short_description) > app.Short_description_max_len {
		malina.Status = http.StatusNotAcceptable
		malina.Result["short_description"]=lang.T("length must have no more then:")+strconv.Itoa(app.Short_description_max_len)+" symbols"
	}

	parameters :=filters.IdsList(malina,"parameters",0,false,r)
	fmt.Println(parameters)
	related_products:=filters.IdsList(malina,"related_products",0,false,r)
	fmt.Println(related_products)
	// <- 1
	//waiting for item exist result
	if action == 'e' {
		if !<-malina.ChannelBool {
			malina.Status = http.StatusNotFound
			malina.Result["item_id"] = lang.T("not found")
		}
	}
	// end <- 1

	if malina.Status > 0 {
		return
	}

	var res int64
	if action == 'a' {
		res = s.model.Upsert(action,imgUrls, parameters,category_id,title,description,short_description,code,price,price1,enable)
	} else {
		res = s.model.Upsert(action,imgUrls, parameters,malina.IdInt64,category_id,title,description,short_description,code,price,price1,enable)
	}
	if res > 0 {
		malina.Status = http.StatusOK
		models.CrudGeneral.WhereAnd(malina, "product_id =","$1")
		malina.SelectSql=s.base.selectSqlFieldsDefault
		if action == 'a'{
			models.CrudGeneral.GetItem(malina,s.base.item,"",res)
		} else {
			models.CrudGeneral.GetItem(malina,s.base.item,"",malina.IdInt64)
		}
		malina.Result["item"] =  malina.Item
	} else {
		malina.Status = http.StatusInternalServerError
		malina.Result["error"] = lang.T(`server error`)
	}
}

func (s *Product) Item(malina *berry.Malina, w http.ResponseWriter, r *http.Request){

	malina.IdInt64, malina.IdStr = filters.IsUint(malina,true, "id", 20, r)

	if malina.Status > 0 {
		return
	}
	if malina.Department == "public" {
		models.CrudGeneral.WhereAnd(malina,"product_enable =","TRUE")
	}

	models.CrudGeneral.WhereAnd(malina, "product_id =","$1")
	malina.SelectSql=s.base.selectSqlFieldsDefault
	models.CrudGeneral.GetItem(malina,s.base.item,"",malina.IdInt64)


	if malina.Device == "browser" {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

			if malina.Department == "public" {

			}else if malina.Department == "home"{

			}

		}else {
			if malina.Department == "public" {
				s.viewPublic.Single(malina,malina.Item.(*entity.Product),w)
			}else if malina.Department == "home"{

			}
		}
	}else {
		out, _ := json.Marshal(malina.Item)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(out)
	}



}


func (s *Product) ResultHandler(malina *berry.Malina,args ... interface{}){

}



func(s *Product)Default(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	if malina.Action==""{
		s.GetList(malina,w,r,"")
	}
	if malina.Action=="item"{
		s.Item(malina,w,r)
	}
}














