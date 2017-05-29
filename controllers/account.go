/**
 * Account controller class.  Malina eCommerce application
 *
 *
 * @author		John Aran Ilyas Aranzhanovich Toxanbayev)
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
	"github.com/ilyaran/Malina/views/account"
	"github.com/ilyaran/Malina/helpers"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/config"
	"fmt"
	"encoding/json"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/entity"
	"html/template"
	"github.com/ilyaran/Malina/libraries"
)

var AccountController = &accountController{&CrudController{}}

type accountController struct{ crud *CrudController }

func (this *accountController) Index(w http.ResponseWriter, r *http.Request) {
	this.crud.hasPermission("account", "account_id", w, r)
	if library.VALIDATION.Status == 0 {

		model.AccountModel.Query = ``
		model.AccountModel.Where = ``
		model.AccountModel.All = 0

		switch this.crud.action {
		case "list_ajax":
			if this.AjaxList(w, r) {
				return
			}
		case "inlist":
			this.Inlist(w, r)
		case "get":
			if this.Get(w, r) {
				return
			}
		case "add":
			this.FormHandler('a', w, r)
		case "edit":
			this.FormHandler('e', w, r)
		case "del":
			this.Del(w, r)
		default:
			this.List(w, r); return // "list"
		}

	}
	helper.SetAjaxHeaders(w)
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}

func (this *accountController) List(w http.ResponseWriter, r *http.Request) {
	page, pageStr := this.crud.getList(false, w, r)
	order := "account_last_logged DESC"
	if library.VALIDATION.Status == 0 {

		accountView.AccountViewObj.AccountList = model.AccountModel.GetList("", pageStr, strconv.FormatInt(app.Per_page(), 10), order, library.POSITION.TreeMap)
		accountView.AccountViewObj.Paging = helper.PagingLinks(model.AccountModel.All, page, app.Per_page(), app.Uri_account()+"list/?page=%d", "href", "a", "", "")

		accountView.AccountViewObj.Index()
	}
}

func (this *accountController) AjaxList(w http.ResponseWriter, r *http.Request) bool {
	page, pageStr, per_page, per_pageStr, search, order_by := this.crud.getAjaxList(false, w, r)
	var order = "account_last_logged DESC"
	switch order_by {
	case 2:
		order = "account_last_logged ASC"
	case 3:
		order = "account_updated DESC"
	case 4:
		order = "account_updated ASC"
	case 5:
		order = "account_created ASC"
	case 6:
		order = "account_created DESC"
	}

	if library.VALIDATION.Status == 0 {

		itemList := model.AccountModel.GetList(search, pageStr, per_pageStr, order, library.POSITION.TreeMap)
		paging := helper.PagingLinks(model.AccountModel.All, page, per_page, "%d", "data-page", "span", "", `class="paging"`)

		helper.SetAjaxHeaders(w)

		w.Write([]byte(accountView.Listing(itemList, paging)))
		return true
	}
	return false
}

func (this *accountController) Get(w http.ResponseWriter, r *http.Request) bool {
	idInt64, _ := this.crud.get(w, r)
	if library.VALIDATION.Status == 0 {
		var accountObj = model.AccountModel.Get(idInt64)
		if accountObj != nil {
			accountObj.SetPassword("")

			helper.SetAjaxHeaders(w)

			out, _ := json.Marshal(accountObj)
			fmt.Fprintf(w, string(out))
			return true
		} else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["item_id"] = lang.T("not found")
		}
	}
	return false
}
func (this *accountController) FormHandler(action byte, w http.ResponseWriter, r *http.Request) {
	var idInt64 int64
	var account *entity.Account = nil
	if action == 'e' {
		idInt64, _ = this.crud.edit(w, r)
		account = model.AccountModel.Get(idInt64)
		if account == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["item_id"] = lang.T("account not found")
		}
	}
	var exec = []interface{}{}
	var q = ``
	var qvalues = ``
	var n int = 1

	position, _ := library.VALIDATION.IsInt64(false, "position", 10, r)
	if position > 0 {
		if _, ok := library.POSITION.TreeMap[position]; ok {
			q += `,account_position`
			qvalues += `,$` + strconv.Itoa(n)
			exec = append(exec, position)
			n++
		} else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["position"] = lang.T("position is not exist")
		}
	}

	var nick = library.VALIDATION.IsValidText(false, "nick", true, 1, 255, `^[a-z0-9-_]+$`, ` a-z 0-9 - _ `, r)
	if nick != `` {
		model.AccountModel.Where = ` account_nick = $1`
		if model.AccountModel.CheckDetail([]interface{}{nick}) && (account != nil && account.Nick != nick) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["nick"] = lang.T("exists allready")
		} else {
			q += `,account_nick`
			qvalues += `,$` + strconv.Itoa(n)
			exec = append(exec, nick)
			n++
		}
	}
	var email = library.VALIDATION.IsEmail(false, r)
	if email != `` {
		model.AccountModel.Where = ` account_email = $1`
		if model.AccountModel.CheckDetail([]interface{}{email}) && (account != nil && account.Email != email) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("exists allready")
		} else {
			q += `,account_email`
			qvalues += `,$` + strconv.Itoa(n)
			exec = append(exec, email)
			n++
		}
	}

	var phone = library.VALIDATION.IsPhone(false, "phone", r)
	if phone != `` {
		model.AccountModel.Where = ` account_phone = $1`
		if model.AccountModel.CheckDetail([]interface{}{phone}) && (account != nil && account.Phone != phone) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["phone"] = lang.T("exists allready")
		} else {
			q += `,account_phone`
			qvalues += `,$` + strconv.Itoa(n)
			exec = append(exec, phone)
			n++
		}

	}

	if email+nick+phone == "" {
		library.VALIDATION.Status = 130
		library.VALIDATION.Result["error"] = lang.T("email or nick or phone required")
	}


	if library.VALIDATION.CheckBox("ban", r) {
		exec = append(exec, true)
	}else {
		exec = append(exec, false)
	}
	q += `,account_ban`
	qvalues += `,$` + strconv.Itoa(n)
	n++

	var ban_reason = template.HTMLEscapeString(r.FormValue("ban_reason")) //library.Validation.IsValidText(false,"ban_reason",false,255,`[\w\s-_]+`,`words, spaces, underscore`,r)
	if ban_reason != `` {
		q += `,account_ban_reason`
		qvalues += `,$` + strconv.Itoa(n)
		exec = append(exec, ban_reason)
		n++
	}

	var newpass string
	if action == 'a' {
		newpass = library.VALIDATION.PasswordValid(true, "newpass", false, false, false, r)
	} else {
		newpass = library.VALIDATION.PasswordValid(false, "newpass", false, false, false, r)
	}
	if newpass != `` {
		q += `,account_password`
		qvalues += `,$` + strconv.Itoa(n)
		exec = append(exec, library.SESSION.Cryptcode(newpass))
		n++
	}

	if library.VALIDATION.Status == 0 && q != `` {
		var res int64
		if action == 'e' {
			// edit
			exec = append(exec, idInt64)
			res = model.Crud.Update(`
			UPDATE account SET
			(`+ q[1:]+ `) = (`+ qvalues[1:]+ `) WHERE account_id = $`+ strconv.Itoa(n), exec)
		}
		if action == 'a' {
			// add
			res = model.Crud.Insert(`
			INSERT INTO account (`+ q[1:]+ `) VALUES (`+ qvalues[1:]+ `) RETURNING account_id`, exec)
			library.VALIDATION.Result["id"] = strconv.FormatInt(res, 10)
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		} else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T(`server error`)
		}
	}
}

func (this *accountController) Del(w http.ResponseWriter, r *http.Request) {
	idInt64, _ := this.crud.del(w, r)
	if library.VALIDATION.Status == 0 {
		res := model.AccountModel.Del(idInt64)
		if res == 0 {
			library.VALIDATION.Status = 250
			library.VALIDATION.Result["error"] = lang.T("not found")
		}
		if res > 0 {
			library.VALIDATION.Status = 0
		}
	}
}

func (this *accountController) Inlist(w http.ResponseWriter, r *http.Request) {
	var columns = map[string]string{"account_nick": "string", "account_ban": "boolean", "account_email": "email", "account_phone": "phone", "account_password": "password"}
	this.crud.inlist(columns, r)
}
