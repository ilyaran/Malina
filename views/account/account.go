/**
 * Account views functions.  Malina eCommerce application
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
package accountView

import (
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/language"
	"strconv"
	"fmt"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/libraries"
	"github.com/ilyaran/Malina/views/position"
)


var AccountViewObj *AccountView
func AccountViewObjInit() {
	AccountViewObj = &AccountView{
		Layout: views.NewLayout(),
		TableForm:&views.TableForm{
			Inputs:views.ICE_FORM.Inputs("account"),
			Breadcrumb : "Home / Accounts",
			Select_options_order_by : select_options_order_by,
			Head : table_head,
		},
	}
}

type AccountView struct {
	Layout          *views.Layout
	TableForm       *views.TableForm
	AccountList     []*entity.Account
	Paging          string
	Form string

}

func (s *AccountView)Index(){

	s.TableForm.Form = Form()
	s.TableForm.Listing = Listing(s.AccountList,s.Paging)
	s.TableForm.BuildIndexForm()
	s.Layout.Body2 = []byte(s.TableForm.Out)

	s.Layout.WriteResponse()
}

var select_options_order_by string = `
	<option value="1">`+lang.T("last logged")+`&darr;</option>
	<option value="2">`+lang.T("last logged")+`&uarr;</option>
	<option value="3">`+lang.T("updated")+`&darr;</option>
	<option value="4">`+lang.T("updated")+`&uarr;</option>
	<option value="5">`+lang.T("created")+`&uarr;</option>
	<option value="6">`+lang.T("created")+`&darr;</option>
	`

func Index(accountList []*entity.Account, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = views.ICE_FORM.Inputs("account")
	views.TABLE_FORM.Breadcrumb = "Home / Accounts"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(accountList,paging)
	views.TABLE_FORM.Form = Form()

	views.TABLE_FORM.BuildIndexForm()

	return views.TABLE_FORM.Out + views.Footer()
}
var table_head = `
	<th>ID</th>

	<th>`+lang.T("phone")+`</th>
	<th >`+lang.T("email")+`</th>

	<th>`+lang.T("nick")+`</th>
	<th>`+lang.T("new pass")+`</th>
	<th>`+lang.T("role")+`</th>

	<th>`+lang.T("provider")+`</th>
	<th>`+lang.T("token")+`</th>

	<th>`+lang.T("last ip")+`</th>
	<th>`+lang.T("last login")+`</th>
	<th>`+lang.T("created")+`</th>
	<th>`+lang.T("updated")+`</th>
	<th>
		Ban<input type="checkbox" onchange="if(this.checked){$('.inlist_ban').bootstrapToggle('on');}else{$('.inlist_ban').bootstrapToggle('off');}">
		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Listing(accountList []*entity.Account, paging string)string{
	var out,idStr string
	if accountList != nil && len(accountList)>0 {
		for _, i := range accountList {
			idStr = strconv.FormatInt(i.GetId(),10)
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>

				<td>
					<div>
						<input data-item_id="` + idStr + `" class="inlist_account_phone" type="text" value="` + i.GetPhone() +`"/>
					</div>
				</td>
				<td>
					<div>
						<input data-item_id="` + idStr + `" class="inlist_account_email" type="text" value="` + i.GetEmail() +`"/>
					</div>
				</td>
				<td>
					<div>
						<input data-item_id="` + idStr + `" class="inlist_account_nick" type="text" value="` + i.GetNick() +`"/>
					</div>
				</td>
				<td>
					<div>
						<input data-item_id="` + idStr + `" class="inlist_account_password" type="text" value=""/>
					</div>

				</td>
				<td>
					<div>
						`; if i.GetPosition()!=nil { out += i.GetPosition().GetTitle() }; out += `
					</div>
				</td>

				<td>
					<div>` + i.GetProvider() +`</div>
				</td>
				<td>
					<div>` + i.GetToken() +`</div>
				</td>

				<td>` + fmt.Sprintf("%v",i.GetLast_ip()) + `</td>
				<td>` + fmt.Sprintf("%v",i.GetLast_logged()) + `</td>
				<td>` + fmt.Sprintf("%v",i.GetCreated()) + `</td>
				<td>` + fmt.Sprintf("%v",i.GetUpdated()) + `</td>

				<td>
					<input value="" data-item_id="` + idStr + `" class="inlist_account_ban" type="checkbox" `; if i.GetBan() {out+=`checked`}; out+=` data-toggle="toggle" >
					<button data-item_id="` + idStr + `" class="btn btn-primary edit_item"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></button>
					<button data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
					<input value="0" data-item_id="` + idStr + `" class="inlist_del" type="checkbox">
				</td>
			</tr>`
		}
		out += `
			<tr>
				<td colspan="9">
					<nav aria-label="...">
						<ul class="pagination">
							`+paging+`
						</ul>
					</nav>
				</td>
			</tr>`
	}else {
		out += `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	return out
}

func Form()string{
	out := `
<div class="row">
	<div class="col-md-12">` +
		views.ICE_FORM.BarForms() + `<br>
		<button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
                <table class="table table-striped table-bordered table-hover" >

                        `+views.ICE_FORM.FormField("nick","nick",false)+`
                        `+views.ICE_FORM.FormField("email","email",false)+`
                        `+views.ICE_FORM.FormField("phone","phone",false)+`
                        `+views.ICE_FORM.FormField("new password","newpass",false)+`
                        <tr>
                        	<td>
                                	<label>`+lang.T("role")+`
                                		<span id="position_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                		<select id="position">
                                		<>
                        			` + positionView.GetSelectOptionsListView(library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()],"Nonposition",false,true) +`
                        			</select>
                        			<span id="position_error"  class="error"></span>

                                        </div>
                               	</td>
                        </tr>
                        `+views.ICE_FORM.CheckBox("ban","ban",false)+`
                        `+views.ICE_FORM.FormField("ban reason","ban_reason",false)+`
                </table>
                <button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}