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
)

var PermissionController *Permission
func PermissionControllerInit(){
	PermissionController = &Permission{base:&base{
		dbtable:"permission",
		item:&entity.Permission{},
		selectSqlFieldsDefault 	:
		`permission_id,
		 permission_code,
		 coalesce(permission_category,0),
		 permission_parameter,
		 permission_img,
		 permission_title,
		 coalesce(permission_description,''),
		 permission_price,
		 permission_price1,
		 permission_price2,
		 coalesce(permission_quantity,0),
		 permission_sold,
		 permission_views,
		 permission_created,
		 permission_updated,
		 permission_enable`,
		}}

	//helpers.GenerateControllerFields(PermissionController.base.dbtable,PermissionController.base.item)

	// model init
	PermissionController.model = models.Permission{}
	// end model init

	// orders init
	PermissionController.base.orderList =[][2]string{
		[2]string{"ORDER BY permission_updated DESC",lang.T("updated")+`&darr;`},
		[2]string{"ORDER BY permission_updated ASC",lang.T("updated")+`&uarr;`},
		[2]string{"ORDER BY permission_created DESC",lang.T("created")+`&darr;`},
		[2]string{"ORDER BY permission_created ASC",lang.T("created")+`&uarr;`},
		[2]string{"ORDER BY permission_price ASC",lang.T("price")+`&uarr;`},
		[2]string{"ORDER BY permission_price DESC",lang.T("price")+`&darr;`},
	}

	if PermissionController.base.orderList==nil || len(PermissionController.base.orderList) < 0{
		panic("order list is not init")
	}
	PermissionController.base.orderListLength = int64(len(PermissionController.base.orderList))
	// end orders init

	PermissionController.base.searchSqlTemplate=`permission_title LIKE '%~%' OR permission_description LIKE '%~%' OR permission_code LIKE '%~%'`

	PermissionController.base.inlistSqlFields = map[string]byte{
		"permission_title"     :'s',
		"permission_price"     :'n',
		"permission_price1"    :'n',
		"permission_code"      :'s',
		"permission_enable"    :'b',
	}

	// view init
	PermissionController.view = views.Permission{}
	for k,_ := range PermissionController.base.inlistSqlFields{
		PermissionController.view.InlistFields += `<input class="inlist_fields" type="hidden" value="`+k+`">`
	}
	for k,v := range PermissionController.base.orderList{
		PermissionController.view.OrderSelectOptions += `<option value="`+strconv.Itoa(k)+`">`+v[1]+`</option>`
	}
	// end view init
}

type Permission struct {
	base     *base
	model			models.Permission
	view            views.Permission
}

func(s *Permission)Index(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	malina.Controller = s
	malina.TableSql = s.base.dbtable
	s.base.index(malina,w,r,"")
}

func(s *Permission)GetList(malina *berry.Malina, w http.ResponseWriter, r *http.Request, condition string){

	category_id,category_idStr := filters.IsUint(malina,false,"category_id",20,r)
	price_min,price_minStr := filters.IsUint(malina,false,"price_min",20,r)
	price_max,price_maxStr := filters.IsUint(malina,false,"price_max",20,r)

	if malina.Status > 0 {
		return
	}

	if category_id > 0 {
		if v,ok:=libraries.CategoryLib.Trees.Map[category_id];ok{
			models.CrudGeneral.WhereAnd(malina,"permission_category IN ","("+category_idStr+v.(*entity.Category).DescendantsString+")")
		}
	}
	if price_min > 0 && price_max == 0 {
		models.CrudGeneral.WhereAnd(malina,"permission_price >",price_minStr)
	}else if price_min == 0 && price_max > 0 {
		models.CrudGeneral.WhereAnd(malina,"permission_price <",price_maxStr)
	}else if price_min > 0 && price_max > 0 {
		models.CrudGeneral.WhereAnd(malina,"permission_price >",price_minStr)
		models.CrudGeneral.WhereAnd(malina,"permission_price <",price_maxStr)
	}

	if malina.Department == "public" {
		models.CrudGeneral.WhereAnd(malina,"permission_enable =","TRUE")
	}

	s.base.getList(malina,w,r)

	if malina.Device == "browser" {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

			if malina.Department == "public" {

			}else if malina.Department == "home"{
				malina.Paging = s.base.paging(malina,"","")
				w.Write([]byte(s.view.Listing(malina)))
			}

		}else {
			if malina.Department == "public" {

			}else if malina.Department == "home"{
				malina.Paging = s.base.paging(malina,"",app.Url_home_permission_list)
				s.view.Index(malina,w)
			}
		}
	}else {
		out, _ := json.Marshal(malina.List)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(out)
	}

}


/*func (s *Permission) FormHandler(malina *berry.Malina,action byte,w http.ResponseWriter,r *http.Request) {

	form := s.base.formHandler(malina,action,r)

	if malina.Status > 0 {
		return
	}

	var res int64
	if action == 'a' {
		res = s.model.Upsert(action,form.imgUrls,form.category_id,form.title,form.description,form.code,form.price,form.price1,form.enable)
	} else {
		res = s.model.Upsert(action,form.imgUrls,malina.IdInt64,form.category_id,form.title,form.description,form.code,form.price,form.price1,form.enable)
	}
	if res > 0 {
		malina.Status = http.StatusOK
		models.CrudGeneral.WhereAnd(malina, "permission_id =","$1")
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
}*/
func (s *Permission) FormHandler(malina *berry.Malina,action byte,w http.ResponseWriter,r *http.Request) {
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

	parameterIds:=filters.IdsList(malina,"parameter",0,false,r)
	fmt.Println(parameterIds)

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
		res = s.model.Upsert(action,imgUrls,parameterIds,category_id,title,description,short_description,code,price,price1,enable)
	} else {
		res = s.model.Upsert(action,imgUrls,parameterIds,malina.IdInt64,category_id,title,description,short_description,code,price,price1,enable)
	}
	if res > 0 {
		malina.Status = http.StatusOK
		models.CrudGeneral.WhereAnd(malina, "permission_id =","$1")
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


func (s *Permission) ResultHandler(malina *berry.Malina,args ... interface{}){

}



func(s *Permission)Default(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	if malina.Action==""{
		s.GetList(malina,w,r,"")
	}
}














