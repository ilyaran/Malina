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
	"database/sql"
	"github.com/ilyaran/Malina/caching"
	"github.com/ilyaran/Malina/dao"

	"github.com/ilyaran/Malina/views/publicView"
)

var AccountController *Account
func AccountControllerInit()string{


	AccountController = &Account{base:&base{dbtable:"account", item:&entity.Account{},}}

	//helpers.GenerateControllerFields(AccountController.base.dbtable,AccountController.base.item)

	// model init
	AccountController.model = models.Account{}
	// end model init

	// set SQL fields
	for _,v := range app.AccountSelectSqlFieldsList {
		AccountController.base.selectSqlFieldsDefault += ","+v
	}
	AccountController.base.selectSqlFieldsDefault = AccountController.base.selectSqlFieldsDefault[1:]



	// orders init
	AccountController.base.orderList =[][2]string{
		{"ORDER BY account_created DESC",lang.T("created")+`&darr;`},
		{"ORDER BY account_created ASC",lang.T("created")+`&uarr;`},
		{"ORDER BY account_updated DESC",lang.T("updated")+`&darr;`},
		{"ORDER BY account_updated ASC",lang.T("updated")+`&uarr;`},
	}

	if AccountController.base.orderList==nil || len(AccountController.base.orderList) < 0{
		panic("order list is not init")
	}
	AccountController.base.orderListLength = int64(len(AccountController.base.orderList))
	// end orders init

	AccountController.base.searchSqlTemplate=` account_nick LIKE '%~%'
												OR account_email LIKE '%~%'
												OR account_first_name LIKE '%~%'
												OR account_last_name LIKE '%~%'																							OR account_email LIKE '%~%'
												OR account_phone LIKE '%~%' `

	AccountController.base.inlistSqlFields = map[string]byte{
		"account_first_name"     :'s',
		"account_last_name"     :'s',
		"account_nick"     :'s',
		"account_email"     :'s',
		"account_ban"    :'b',
	}

	// view init
	AccountController.view = views.Account{}
	for k,_ := range AccountController.base.inlistSqlFields{
		AccountController.view.InlistFields += `<input class="inlist_fields" type="hidden" value="`+k+`">`
	}
	for k,v := range AccountController.base.orderList{
		AccountController.view.OrderSelectOptions += `<option value="`+strconv.Itoa(k)+`">`+v[1]+`</option>`
	}
	AccountController.viewPublic=publicView.Account{}
	// end view init

	return AccountController.base.selectSqlFieldsDefault
}

type Account struct {
	ArraySelectSqlFields []string
	base     		*base
	model			models.Account
	view            views.Account
	viewPublic publicView.Account
}

func(s *Account)Index(malina *berry.Malina,department, device string, w http.ResponseWriter, r *http.Request){
	malina.Controller = s
	malina.Department = department
	malina.Device = device
	malina.TableSql = s.base.dbtable
	s.base.index(malina,w,r,"")
}

func(s *Account)GetList(malina *berry.Malina, w http.ResponseWriter, r *http.Request, condition string){


	if malina.Status > 0 {
		return
	}

	if malina.Department == "public"{
		models.CrudGeneral.WhereAnd(malina,"account_enable =","TRUE")
	}

	s.base.getList(malina,w,r)
	//fmt.Println(malina.Order_bySql)
	// if the user agent is a browser
	if malina.Device == "browser" {
		// if via ajax
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

			//helpers.SetAjaxHeader(w)

			//if public
			if malina.Department == "public"{
				//s.view.Paging = s.base.paging(malina,"","")
				//w.Write([]byte(public.AccountView.Listing()))

			// if home admin
			}else {
				malina.Paging = s.base.paging(malina,"","")
				w.Write([]byte(s.view.Listing(malina)))
			}

		// if non ajax req
		}else {
			//if public
			if malina.Department == "public"{
				//public.AccountView.Paging = s.base.paging(malina,app.Url_public_account_list)
				//public.AccountView.Index()
			// if home admin
			}else {
				malina.Paging = s.base.paging(malina,"",app.Url_home_account_list)
				s.view.Index(malina,w)
			}
		}

	// if the user agent is a mobile or another device
	}else {
		out, _ := json.Marshal(malina.List)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(out)
	}

}


func (s *Account) FormHandler(malina *berry.Malina,action byte,w http.ResponseWriter,r *http.Request) {
	if malina.Department=="cabinet"{
		s.UpdateAccountColumn(malina,w, r)
		return
	}


	if malina.Status > 0 {
		return
	}

	var res int64
	if action == 'a' {
		//res = s.model.Upsert(action,form.imgUrls,form.category_id,form.title,form.description,form.code,form.price,form.price1,form.enable)
	} else {
		//res = s.model.Upsert(action,form.imgUrls,malina.IdInt64,form.category_id,form.title,form.description,form.code,form.price,form.price1,form.enable)
	}
	if res > 0 {
		malina.Status = http.StatusOK
		models.CrudGeneral.WhereAnd(malina, "account_id =","$1")
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



func (s *Account) ResultHandler(malina *berry.Malina,args ... interface{}){

}

func(s *Account)Default(malina *berry.Malina, w http.ResponseWriter, r *http.Request){

	if malina.Action == "form" {
		s.AccountForm(malina,w, r)
		return
	}
	if malina.Action == "phone" {

		return
	}
	if malina.Action == "img" {
		s.AccountImg(malina,w, r)
		return
	}



	s.viewPublic.Index(malina,w)
}


func (s *Account) AccountForm(malina *berry.Malina,w http.ResponseWriter, r *http.Request) {

	if r.Method=="GET"{


	}else if r.Method=="POST"{

		first_name :=filters.IsValidText(malina,false, "first_name", false, 1, 64, `^[а-яА-Я\sa-zA-Z0-9_\-]+$`, `позволительно только следующие символы: а-я А-Я a-z A-Z 0-9 _ - и пробел`, r)
		if first_name != `` {malina.CurrentAccount.FirstName=first_name}

		last_name:=filters.IsValidText(malina,false, "last_name", false, 1, 64, `^[а-яА-Я\sa-zA-Z0-9_\-]+$`, `позволительно только следующие символы: а-я А-Я a-z A-Z 0-9 _ - и пробел`, r)
		if last_name != `` {malina.CurrentAccount.LastName=last_name}

		nick:=filters.IsValidText(malina,false, "nick", true, 1, 8, `^[a-zA-Z0-9]+$`, `позволительно только следующие символы: a-z A-Z 0-9`, r)
		if nick!=`` {
			if malina.CurrentAccount.Nick != nick {
				if AccountController.model.CheckDetail(`account_nick = $1`, nick) {
					malina.Status = 404
					malina.Result["nick"] = "nick " + lang.T("exists allready")
				} else {
					malina.CurrentAccount.Nick = nick
				}
			}
		}
		state, _ :=filters.IsUint(malina,false, "state", 3, r)
		if state>0 {
			if _,ok := caching.StatesMap[state]; ok {
				malina.CurrentAccount.State = state
			}else {
				malina.Status = 404
				malina.Result["state"] = "state "+lang.T("not found")
			}
		}else {
			malina.CurrentAccount.State = 0
		}

		/*city :=filters.IsValidText(malina,false, "city", false, 1, 64, `^[а-яА-Я\sa-zA-Z0-9_\-]+$`, `позволительно только следующие символы: а-я А-Я a-z A-Z 0-9 _ - и пробел`, r)
		if city != `` { malina.CurrentAccount.City=city }*/

		var skype =filters.IsValidText(malina,false, "skype", true, 1, 8, `^[a-zA-Z0-9]+$`, `позволительно только следующие символы: a-z A-Z 0-9`, r)
		if skype != `` {
			if malina.CurrentAccount.Skype != skype {
				if AccountController.model.CheckDetail(`account_skype = $1`, skype) {
					malina.Status = 404
					malina.Result["skype"] = "skype " + lang.T("exists allready")
				} else {
					malina.CurrentAccount.Skype = skype
				}
			}
		}
		var email =filters.IsEmail(malina,true ,"email",r)
		if email != `` {
			if malina.CurrentAccount.Email != email{
				if AccountController.model.CheckDetail(`account_email = $1`, email) {
					malina.Status = 404
					malina.Result["email"] = "email " + lang.T("exists allready")
				}else {
					malina.CurrentAccount.Email = email
				}
			}
		}

		/*var phone =filters.IsValidText(malina,false, "phone", true, 5, 16, `^[\+]?[0-9]+$`, `позволительно только : +77058436633`, r)
		if phone != `` {
			if malina.CurrentAccount.Phone != phone {
				model.AuthModel.CheckDetails(`account_phone = $1`, phone)
				if model.AuthModel.All > 0 {
					malina.Status = 404
					malina.Result["phone"] = "phone " + lang.T("exists allready")
				} else {
					malina.CurrentAccount.Phone = phone
				}
			}
		}*/

		var oldPass =filters.PasswordValid(malina,false, "old_password", false, false, false, r)
		if oldPass != `` {
			if malina.CurrentAccount.Password != dao.AuthDao.Cryptcode(oldPass){
				malina.Status = 404
				malina.Result["old_password"] = "старый пароль не правильный"
			}
			var new_pass =filters.PasswordValid(malina,false, "new_password", false, false, false, r)

			filters.ConfirmPasswordValid(malina,"confirm_new_password",new_pass, r)

			if malina.Status == 0 {
				malina.CurrentAccount.Password = dao.AuthDao.Cryptcode(new_pass)
			}
		}

		if r.FormValue("sex") == "male" {
			malina.CurrentAccount.Sex=false
		}else {
			malina.CurrentAccount.Sex=true
		}

		if malina.Status == 0 {
			if AccountController.model.EditByObject(malina.CurrentAccount) > 0 {
				//msg="Успешно обновлен"
			} else {
				//msg="ошибка сервера"
			}
			/*malina.Status = 200
			malina.Result["success"] = msg*/
		}
	}
	//s.viewPublic.AccountForm(malina,msg,w)

}


func (s *Account) UpdateAccountColumn(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	key := r.FormValue("key")
	value := filters.IsValidText(malina,true, "value", true, 1, 128, `^[a-zA-Z0-9_\-]+$`, `позволительно только следующие символы: a-z A-Z 0-9 _ -`, r)
	if malina.Status == 0 {
		if key == "trade"{
			s.UpdateColumn(malina,"account_trade",value)
		}
	}
}

func (s *Account) UpdateColumn(malina *berry.Malina,columnName string,value interface{}) bool {

	if models.CrudGeneral.Update(`
			UPDATE account SET `+columnName+`= $1 WHERE account_id = $2
		`,value, malina.CurrentAccount.Id)>0{
		malina.Status=200
		malina.Result["ok"]="ok"
		return true
	}
	malina.Status=500
	malina.Result["error"]=lang.T("server error")
	return false
}


func (s *Account) AccountImg(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	avatar,_ := libraries.UploadLib.Img(malina,"avatar",w,r)
	/*if avatar == "" {
		malina.Status = 200
		malina.Result["account_img"] = app.Url_no_avatar

		return
	}*/
	if malina.Status > 0 {
		return
	}

	sqlAvatar:=sql.NullString{String:avatar,Valid:true}
	if avatar=="null" || avatar==""{sqlAvatar.Valid=false}
	var oldAvatar string
	err := models.CrudGeneral.DB.QueryRow(`
			UPDATE "account" x
			SET    "account_img" = $1
			FROM  (SELECT account_id, "account_img" FROM "account" WHERE account_id = $2 FOR UPDATE) y
			WHERE  x.account_id = y.account_id
			RETURNING COALESCE(y."account_img",'')
		`,sqlAvatar,malina.CurrentAccount.Id).Scan(&oldAvatar)
	if err==sql.ErrNoRows{
		//panic(err)
		malina.Status = 404
		malina.Result["error"] = lang.T("not found")
		return
	}
	if err!=nil{
		panic(err)
		malina.Status = 500
		malina.Result["error"] = lang.T("server error")
		return
	}
	libraries.UploadLib.DelFile(app.Root_path+"/"+app.Path_assets_uploads, oldAvatar)

	malina.Status = 200
	if avatar=="null" || avatar==""{
		malina.Result["account_img"] = app.Url_no_image
	}else{
		malina.Result["account_img"] = app.Url_assets_uploads + avatar
	}

}


















