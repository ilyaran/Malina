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
	"github.com/gorilla/mux"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/filters"
	"regexp"
	"github.com/ilyaran/Malina/models"
	"strconv"
	"fmt"
	"github.com/ilyaran/Malina/helpers"
	"github.com/ilyaran/Malina/dao"
	"encoding/json"
	"github.com/ilyaran/Malina/lang"
	"strings"
	"html/template"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/libraries"

	"github.com/ilyaran/Malina/berry"
)

func getNewMalina(w http.ResponseWriter, r *http.Request) *berry.Malina {
	malina := new(berry.Malina)
	malina.CurrentAccount,malina.SessionId = dao.AuthDao.Authentication(w,r)
	malina.Status = 0
	malina.Result = map[string]interface{}{}
	malina.Device = "browser"
	if malina.CurrentAccount!=nil{

	}

	return malina
}
func HomeDepartment (w http.ResponseWriter, r *http.Request){
	malina := getNewMalina(w,r)
	malina.ControllerName = mux.Vars(r)["controller"]
	malina.Department = "home"
	if malina.CurrentAccount!=nil && malina.CurrentAccount.Role == app.Admin_role_id {
		switch malina.ControllerName {
			case "product"      : ProductController.Index(malina,w,r)
			case "category"     : CategoryController.Index(malina,w,r)
			case "parameter"    : ParameterController.Index(malina,w,r)
			case "cart"         :
			case "settings"   	: SettingsController.Index(w,r)
			case "account"      :
			case "role"         :
			case "permission"   :
			case "activation"   :
			default:HomeController.Index(w,r)
		}
	}else{
		http.Redirect(w,r, "/", 301)
	}
}
func PublicDepartment (w http.ResponseWriter, r *http.Request){
	malina := getNewMalina(w,r)
	malina.ControllerName = mux.Vars(r)["controller"]
	malina.Department = "public"
	switch malina.ControllerName {
	case "product"		:	ProductController.Index(malina,w,r)
	case "category"		:	CategoryController.Index(malina,w,r)
	case "parameter"	:	ParameterController.Index(malina,w,r)
	case "cart"			:
	default				: 	WelcomeController.Index(malina,w,r)
	}
}
func CabinetDepartment (w http.ResponseWriter, r *http.Request){
	malina := getNewMalina(w,r)
	malina.Department = "cabinet"
	malina.Action = mux.Vars(r)["action"]
	if malina.CurrentAccount != nil {
		switch malina.Action {
		case "cart"		:
		default			:	CabinetController.Index(malina,w,r)
		}
	}else {
		http.Redirect(w,r, "/", 301)
	}
}
func AuthDepartment(w http.ResponseWriter, r *http.Request){
	malina := getNewMalina(w,r)
	malina.Action=mux.Vars(r)["action"]
	switch malina.Action {
	case "login" 			: 	AuthController.Login(malina,w,r)
	case "forgot" 			:	AuthController.Forgot(malina,w,r)
	case "register" 		:	AuthController.Register(malina,w,r)
	case "change_password" 	:	AuthController.ChangePassword(malina,w,r)
	case "activation" 		:	AuthController.Activation(malina,w,r)
	case "logout" 			:	AuthController.Logout(malina,w,r)
	default 				:	AuthController.Login(malina,w,r)
	}
}
type base struct {
	dbtable            		string

	orderList       		[][2]string
	orderListLength 		int64
	searchSqlTemplate 		string
	inlistSqlFields     	map[string]byte

	selectSqlFieldsDefault 	string
	item 					entity.Scanable
	itemHierarchical 		entity.Hierarchical

}
func(s *base) index(malina *berry.Malina, w http.ResponseWriter, r *http.Request, condition string){
	malina.Action = mux.Vars(r)["action"]
	switch malina.Action {
		case "list"     :   malina.Controller.GetList(malina, w, r, condition)
		case "get"      :   s.getItem(malina,w, r,condition)
		case "add"      :   malina.Controller.FormHandler(malina,'a',w, r)
		case "edit"     :   malina.Controller.FormHandler(malina,'e',w, r)
		case "del"      :   s.del(malina,r,condition)
		case "inlist"   :   s.inlist(malina,r,condition)

		default			: 	malina.Controller.Default(malina, w, r)
	}
	if malina.Status > 0 {
		out, err := json.Marshal(malina.Result)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(int(malina.Status))
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(out))
	}
}

func(s *base) getList(malina *berry.Malina,w http.ResponseWriter, r *http.Request, args... interface{}){
	s.listingPrepare(malina,w, r)

	// count all items
	//s.start_count_items(malina)
	malina.ChannelBool = make(chan bool)
	go models.CrudGeneral.Count(malina,"", args...)

	malina.SelectSql = s.selectSqlFieldsDefault
	models.CrudGeneral.GetItems(malina,s.item,"",args...)

	// sync all items
	<- malina.ChannelBool
}


func(s *base)paging(malina *berry.Malina,class, url string)string{
	if url == "" {
		return helpers.PagingLinks(malina.All, malina.Page, malina.Per_page, "%d", "data-page", "span", class, `class="page_`+s.dbtable+`"`)
	}
	return helpers.PagingLinks(malina.All, malina.Page, malina.Per_page, url,"href","a",class,"")
}


func(s *base)listingPrepare(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	malina.Page, malina.PageStr = filters.IsUint(malina,false, "page", 20, r)
	malina.Per_page,_ = filters.IsUint(malina,false, "per_page", 20, r)
	malina.Url += "&page=%d"

	// per_page handle
	// if no any per page data then set default per page
	if malina.Per_page == 0 {
		malina.Per_page = app.Per_page
	}else {
		if malina.Per_page < app.Per_page  {
			malina.Per_page = app.Per_page
		}
		if malina.Per_page > app.Per_page_max-1 {
			malina.Per_page = app.Per_page_max -1
		}
	}
	malina.Url += "&per_page="+strconv.FormatInt(malina.Per_page,10)
	// end per_page handle

	malina.Order_index, _ = filters.IsUint(malina,false, "order_by", 2, r)
	if malina.Order_index > s.orderListLength {
		malina.OrderBySql = s.orderList[0][0]
	}else {
		malina.OrderBySql = s.orderList[malina.Order_index][0]
		malina.Url += "&order_by="+strconv.FormatInt(malina.Order_index,10)
	}


	if r.Method == "POST" {
		malina.Search = r.FormValue("search")
	}else if r.Method == "GET" {
		malina.Search = r.URL.Query().Get("search")
	}
	if malina.Search != "" {
		if len(malina.Search) > 64 {
			malina.Search = malina.Search[0:64]
		}
		malina.Search = regexp.MustCompile(`[\W]`).ReplaceAllString(malina.Search, "")
		malina.Search = regexp.MustCompile(`[~]`).ReplaceAllString(s.searchSqlTemplate, malina.Search)
		if malina.WhereSql == ``{
			malina.WhereSql += ` WHERE `+ malina.Search+` `
		}else {
			malina.WhereSql +=  ` AND (`+ malina.Search+`) `
		}
		malina.Url += "&search="+malina.Search
	}

	malina.LimitSql = "LIMIT "+strconv.FormatInt(malina.Per_page,10)+" OFFSET "+ malina.PageStr
}

func(s *base)getItem(malina *berry.Malina,w http.ResponseWriter, r *http.Request, condition string){
	malina.IdInt64, malina.IdStr = filters.IsUint(malina,false, "id", 20, r)
	if malina.IdInt64 > 0 {

		models.CrudGeneral.WhereAnd(malina,s.dbtable+"_id =", malina.IdStr)
		if malina.CurrentAccount != nil && malina.CurrentAccount.Role != app.Admin_role_id {
			models.CrudGeneral.WhereAnd(malina,s.dbtable+"_enable =","TRUE")
		}
		malina.SelectSql = s.selectSqlFieldsDefault
		models.CrudGeneral.GetItem(malina,s.item,"")

		if malina.Item != nil {
			malina.Status = http.StatusOK
			malina.Result["item"] = malina.Item
		}else {
			malina.Status = http.StatusNotFound
			malina.Result["item_id"] = lang.T("not found")
		}
	}
}

func(s *base)del(malina *berry.Malina,r *http.Request, condition string){
	malina.IdInt64, malina.IdStr = filters.IsUint(malina,false, "id", 20, r)
	if malina.IdInt64 > 0{


		res:=models.CrudGeneral.Delete("DELETE FROM "+s.dbtable+" WHERE "+s.dbtable+"_id = "+malina.IdStr)
		if  res > 0 {
			malina.Status = http.StatusOK
			malina.Controller.ResultHandler(malina)
		}else if res == 0{
			malina.Status = http.StatusNotFound
			malina.Result["item_id"] = s.dbtable+"_id:"+ malina.IdStr+" - " + lang.T("not found")
		}else {
			malina.Status = http.StatusInternalServerError
			malina.Result["error"] = lang.T("server error")
		}
	}
}

func(s *base)isExistItem(malina *berry.Malina,id int64, querySql string,isViaCrudGeneralChannel, hasNull bool) {
	if id > 0 {
		if isViaCrudGeneralChannel {
			if querySql == "" {
				go models.CrudGeneral.Count(malina,"SELECT count(*) FROM "+s.dbtable+" WHERE "+s.dbtable+"_id = $1", id)
			}else {
				go models.CrudGeneral.Count(malina,querySql, id)
			}
		} else {
			if querySql == "" {
				go models.CrudGeneral.Count(malina,"SELECT count(*) FROM "+s.dbtable+" WHERE "+s.dbtable+"_id = $1", id)
			}else {
				go models.CrudGeneral.Count(malina, querySql, id)
			}
		}
	} else{
		if hasNull {
			if isViaCrudGeneralChannel {
				malina.ChannelBool <- true
			} else {
				malina.ChannelBool <- true
			}
		} else {
			if isViaCrudGeneralChannel {
				malina.ChannelBool <- false
			} else {
				malina.ChannelBool <- false
			}
		}
	}
}

func(s *base) setTree(malina *berry.Malina,treeSetter libraries.TreeSetter){
	// set public trees
	malina.SelectSql = s.selectSqlFieldsDefault
	malina.WhereSql=""
	malina.OrderBySql = s.orderList[0][0]
	malina.LimitSql=""
	models.CrudGeneral.WhereAnd(malina, s.dbtable+"_enable =","TRUE")
	models.CrudGeneral.GetItems(malina,s.item,"")

	treeSetter.SetTrees(malina.List)
	treeSetter.SetPublicTrees()
	// end set public trees

	// set home trees
	malina.WhereSql=""
	models.CrudGeneral.GetItems(malina,s.item,"")

	treeSetter.SetTrees(malina.List)
	// end set home trees
}

func(s *base)inlist(malina *berry.Malina,r *http.Request, condition string){
	if s.dbtable == "role" || s.dbtable == "permission" {


	}
	r.ParseForm()

	tx, err := models.CrudGeneral.DB.Begin()
	if err != nil {
		return
	}

	for columnName, columnType := range s.inlistSqlFields {
		var query, values, ids string
		var valueNum = 1  // count numbers 1,2,3 ... to "$1,$2,$3,$4,$5 ..."
		var args []interface{}
		for _, v := range r.PostForm[columnName+"_inlist[]"] {
			query = `UPDATE ` + s.dbtable + `
				SET ` + columnName + ` = data_table.v
				FROM (
				       SELECT
					 unnest(ARRAY[ %s ]) AS v,
					 unnest(ARRAY[ %s ]) AS i
				     )AS data_table
				WHERE ` + s.dbtable + `_id = data_table.i `

			// The id-value string - 25|89 or 25|253.56 or 25|Some Title
			// is splitting to array [25,89] or [25,253.56] or [25,Some Title]
			id_value := strings.SplitN(v, "|", 2)
			match, _ := regexp.MatchString(`^[0-9]{1,20}$`, id_value[0])

			if match {
				switch columnType {
				case 'i': // i - integer
					match, _ = regexp.MatchString(`^[0-9]+$`, id_value[1])
					if match {
						values += `,` + id_value[1]
						ids += `,` + id_value[0]
						valueNum++
					}
				case 'n': // n - numeric
					match, _ = regexp.MatchString(`^[0-9]{1,20}(\.[0-9]{0,2})?$`, id_value[1])
					if match {
						values += `,` + id_value[1]
						ids += `,` + id_value[0]
						valueNum++
					}

				case 's': // s - string
					if id_value[1] != `` {
						// building "$1,$2,$3,$4,$5 ..." string to values array: unnest(ARRAY[$1,$2,$3,$4,$5 ...]) AS v
						values += `,$` + strconv.Itoa(valueNum)
						args = append(args, template.HTMLEscapeString(id_value[1]))
						ids += `,` + id_value[0]
						valueNum++
					}
				case 'b': // b - boolean
					// building "TRUE,TRUE,TRUE,FALSE, ..." string to values array: unnest(ARRAY[TRUE,TRUE,TRUE,FALSE, ...]) AS v
					valueBool, _ := strconv.ParseBool(id_value[1])
					if valueBool {
						values += `,TRUE`
					} else {
						values += `,FALSE`
					}
					ids += `,` + id_value[0]
					valueNum++
				case 'P': // P - phone
					match, _ = regexp.MatchString(app.Pattern_phone, id_value[1])
					if match {
						values += `,$` + strconv.Itoa(valueNum)
						args = append(args, id_value[1])
						ids += `,` + id_value[0]
						valueNum++
					}
				case 'p': // p - password
					if len(id_value[1])>5 && len(id_value[1])<128 {
						values += `,$` + strconv.Itoa(valueNum)
						args = append(args, dao.AuthDao.Cryptcode(id_value[1]))
						ids += `,` + id_value[0]
						valueNum++
					}
				}
			}
		}
		if ids != `` {
			stmt, err := tx.Prepare(fmt.Sprintf(query, values[1:], ids[1:]))
			if err != nil {
				return
			}
			defer stmt.Close()

			if _, err := stmt.Exec(args...); err != nil {
				tx.Rollback()
				return
			}
		}
	}

	// delete in list
	var del_ids string
	for _, v := range r.PostForm[s.dbtable+"_del_item_inlist[]"] {
		ids := strings.SplitN(v, "|", 2)
		match, _ := regexp.MatchString(`^[0-9]{1,20}$`, ids[0])
		if match {
			valueBool, _ := strconv.ParseBool(ids[1])
			if valueBool {
				del_ids += `,` + ids[0]
			}
		}
	}
	if del_ids != `` {
		del_q := fmt.Sprintf(`DELETE FROM `+ s.dbtable+ ` WHERE `+ s.dbtable+ `_id = ANY(ARRAY[%s])`, del_ids[1:])

		stmt, err := tx.Prepare(del_q)
		if err != nil {
			return
		}
		defer stmt.Close()

		if _, err := stmt.Exec(); err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	malina.Status = 200
	malina.Result["success"] = "ok"
}








