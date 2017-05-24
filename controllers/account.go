/**
 * Account controller class.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      il.aranov@gmail.com
 * @link
 * @github      https://github.com/ilyaran/Malina
 * @license	MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package controller

import (
	"net/http"
	"strconv"
	"Malina/views/account"
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
)
var Account = &account{&CrudController{}}
type account  struct { crud *CrudController }

func (this *account) Index(w http.ResponseWriter, r *http.Request) {
	action := this.crud.authAdmin("account",w,r)
	if action!=""{
		switch action {
		case "ajax_list" : if this.AjaxList(w,r) {return}
		case "get" 	 : if this.Get(w,r) {return}
		case "add" 	 : this.FormHandler('a',w,r)
		case "edit" 	 : this.FormHandler('e',w,r)
		case "del" 	 : this.Del(w,r)
		default		 : this.List(w,r); return // "list"
		}
		helper.SetAjaxHeaders(w)
	}
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *account) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false,w,r)
	order := "account_last_logged DESC"
	if library.VALIDATION.Status == 0 {

		all := model.AccountModel.CountItems("")
		paging := helper.PagingLinks(all, page, app.Per_page(), app.Uri_account()+"list/?page=%d","href","a","","")
		itemList := model.AccountModel.GetList("", pageStr, strconv.FormatInt(app.Per_page(), 10), order)

		w.Write([]byte(views.Header()))
		w.Write([]byte(views.Nav("account")))
		w.Write([]byte(accountView.Index(itemList, paging)))
	}
}

func (this *account) AjaxList(w http.ResponseWriter, r *http.Request) bool{
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "account_last_logged DESC"
	switch order_by {
	case 2:order = "account_last_logged ASC"
	case 3:order = "account_nick ASC"
	case 4:order = "account_nick DESC"
	case 5:order = "account_fist_name ASC"
	case 6:order = "account_fist_name DESC"
	case 7:order = "account_last_name ASC"
	case 8:order = "account_last_name DESC"
	case 9:order = "account_updated ASC"
	case 10:order = "account_updated DESC"
	}
	if library.VALIDATION.Status == 0 {
		all := model.AccountModel.CountItems(search)
		paging := helper.PagingLinks(all, page, per_page, "%d", "data-page", "span","", `class="paging"`)
		itemList := model.AccountModel.GetList(search, pageStr, per_pageStr, order)
		helper.SetAjaxHeaders(w)
		w.Write([]byte(accountView.Listing(itemList, paging)))
		return true
	}
	return false
}

func (this *account) Get(w http.ResponseWriter, r *http.Request) bool{
	idInt64,_ := this.crud.get(w,r)
	if library.VALIDATION.Status == 0 {
		var accountObj = model.AccountModel.Get(idInt64)
		if accountObj != nil {
			accountObj.SetPass("")

			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(accountObj)
			fmt.Fprintf(w, string(out))
			return true
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["account_id"] = lang.T("not found")
		}
	}
	return false
}
func (this *account) FormHandler(action byte,w http.ResponseWriter, r *http.Request) {
	var idInt64 int64 = 0
	if action == 'e'{
		idInt64,_ = this.crud.edit(w,r)
	}
	var email = library.VALIDATION.IsEmail(false, r)
	if email != ""{
		var account = model.AccountModel.GetByEmailByNickByPhoneByPassword(email,"", "","")
		if account != nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("auth_email_exist_allready")
		}
	}

	position,_ := library.VALIDATION.IsInt64(true,"position",10,r)
	var nick  = library.VALIDATION.IsValidText(false,"nick",true,1,255,`^[a-z0-9-_]+$`,` a-z 0-9 - _ `,r)
	if nick != ""  {
		var account = model.AccountModel.GetByEmailByNickByPhoneByPassword("",nick, "","")
		if account != nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["nick"] = lang.T("auth_nick_exist_allready")
		}
	}
	var phone  = library.VALIDATION.IsValidText(false,"phone",true,10,16,`^\+[0-9]+$`,` a-z 0-9 - _ `,r)
	if phone != ""  {
		var account = model.AccountModel.GetByEmailByNickByPhoneByPassword("","", phone,"")
		if account != nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["phone"] = lang.T("auth phone exist allready")
		}
	}
	if email == "" && nick == "" && phone == ""{
		library.VALIDATION.Status = 130
		library.VALIDATION.Result["error"] = lang.T("auth email or nick or phone required")
	}

	var ban_reason  = template.HTMLEscapeString(r.FormValue("ban_reason")) //library.Validation.IsValidText(false,"ban_reason",false,255,`[\w\s-_]+`,`words, spaces, underscore`,r)
	var newpass string
	if action=='a'{
		newpass = library.VALIDATION.PasswordValid(true,"newpass",false,false,false,r)
	}else {
		newpass = library.VALIDATION.PasswordValid(false,"newpass",false,false,false,r)
	}
	var ban int = 0
	if library.VALIDATION.CheckBox("ban",r){
		ban = 1

	}
	if library.VALIDATION.Status == 0 {
		var res int64
		account := entity.NewAccount(idInt64,email,position,nick,ban_reason, newpass,ban)
		if account.GetNew_password() != ""{
			account.SetPass(library.SESSION.Cryptcode(account.GetNew_password()))
			account.SetNewPass("")
		}
		if action == 'e'{ // edit
			res = model.AccountModel.Edit(account)
		}
		if action == 'a'{ // add
			res = model.AccountModel.Add(account)
			library.VALIDATION.Result["id"] = strconv.FormatInt(res,10)
		}
		if res > 0{
			library.VALIDATION.Status = 0
		}else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}

func (this *account) Del(w http.ResponseWriter, r *http.Request) {
	idInt64,_ := this.crud.del(w,r)
	if library.VALIDATION.Status == 0 {
		res := model.AccountModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not find")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}

