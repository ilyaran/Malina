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
	"Malina/entity"
	"Malina/language"
	"strconv"
	"fmt"
	"Malina/config"
	"Malina/views"
	"Malina/libraries"
)

func Index(accountList []*entity.Account, paging string)string{
	var inputs = `
	<input type="hidden" value="account" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_account_ajax()+`" id="url_account_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_account_get()+`" id="url_account_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_account_add()+`" id="url_account_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_account_edit()+`" id="url_account_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_account_del()+`" id="url_account_del"/>`
	var out = inputs+ `
<div class="row">
	<div class="col-md-12">

		<div class="panel panel-default">
			<div class="panel-heading">
				<h3 class="panel-title">Home / Accounts</h3>
			</div>
		  	<div class="panel-body">
		    		`+lang.T("per page")+`
				`+views.FACE_FORM.PerPageSelectForm()+`
				&nbsp;&nbsp;
				<span>`+lang.T("order_by")+`</span>
				<select id="order_by" >
					<option value="1">`+lang.T("last_logged")+`&darr;</option>
					<option value="2">`+lang.T("last_logged")+`&uarr;</option>
					<option value="3">`+lang.T("nick")+`&uarr;</option>
					<option value="4">`+lang.T("nick")+`&darr;</option>
					<option value="5">`+lang.T("fist_name")+`&uarr;</option>
					<option value="6">`+lang.T("fist_name")+`&darr;</option>
					<option value="7">`+lang.T("last_name")+`&uarr;</option>
					<option value="8">`+lang.T("last_name")+`&darr;</option>
					<option value="9">`+lang.T("updated")+`&uarr;</option>
					<option value="10">`+lang.T("updated")+`&darr;</option>
				</select>

				<input type="text" class="form-control" id="search" name="search"  placeholder="Search for...">

		    		<button id="get_list" class="btn btn-primary" type="button" >`+lang.T("go")+`</button>

		  	</div>
		</div>

		<ul class="nav nav-tabs">
                	<li class="active"><a href="#home" data-toggle="tab">`+lang.T("list")+`</a></li>
                        <li class=""><a href="#form" data-toggle="tab">` + lang.T("form") + `</a></li>
		</ul>
		<div class="tab-content">
			<div class="tab-pane fade active in" id="home">
				<table class="table table-striped table-bordered table-hover" id="dataTables-example">
					<thead>
						<tr>
							<th>ID</th>
							<th>`+lang.T("avatar")+`</th>
							<th>`+lang.T("nick")+`</th>
							<th>`+lang.T("email")+`</th>
							<th>`+lang.T("position")+`</th>
							<th>`+lang.T("banned")+`</th>
							<th>`+lang.T("last IP")+`</th>
							<th>`+lang.T("last logged")+`</th>
							<th>`+lang.T("created")+`</th>
							<th>`+lang.T("action")+`</th>
						</tr>

					</thead>
					<tbody id="listing">
					`
					out += Listing(accountList,paging)+`

					</tbody>
				</table>
			</div>
			<div class="tab-pane fade" id="form">
				`+Form()+`

			</div>
		</div>
	</div>
</div>

`+ views.Footer()

	return out
}

func Listing(accountList []*entity.Account, paging string)string{
	var out,idStr string
	if accountList != nil && len(accountList)>0 {
		for _, i := range accountList {
			idStr = strconv.FormatInt(i.GetId(),10)
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td><img src="` + i.GetImg() + `" width='60' height='50' /></td>
				<td>` + i.GetNick() + `</td>
				<td>` + i.GetEmail()+`</td>
				<td>` + i.GetPosition().GetTitle() + `</td>
				<td>` + i.GetBan_reason()  + `</td>
				<td>` + i.GetLast_ip()  + `</td>
				<td>` + fmt.Sprintf("%v",i.GetLast_logged())  + `</td>
				<td>` + fmt.Sprintf("%v",i.GetCreated())  + `</td>
				<td>
					<a data-item_id="` + idStr + `" class="btn btn-primary edit_item"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></a>
					<a data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
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
	<div class="col-md-12">
		<div class="form-group">

			<select id="action">
				<option value="add">`+lang.T("add")+`</option>
				<option value="get">`+lang.T("get")+`</option>
				<option value="edit">`+lang.T("edit")+`</option>
				<option value="delete">`+lang.T("delete")+`</option>
			</select>
                        <span style = "display:none;" id = "item_id_bar">
                        	<span>`+lang.T("id")+`:</span>
                        	<input id = "item_id" value="" type="integer"/>
                        	<span id = "item_id_error" class="error"></span>
                        </span>
                        <button id = "success_bar" style="display:none;" type="button" class="btn btn-success btn-lg">Success</button>
                </div>
                <button id="submitButton" class="btn btn-primary">`+lang.T("Send")+`</button>
		<div id="error" class="error"></div>
                <table class="table table-striped table-bordered table-hover" >
                	<tr>
                        	<td>
                                	<label>`+lang.T("nick")+`<span class="error">*</span>
                                		<span id="nick_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="nick"  class="form-control" type="text" pattern="[a-z0-9-_]+"/>
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("email")+`
                                	<span id="email_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="email" type="email" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("phone")+`
                                	<span id="phone_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="phone"  type="tel" pattern="^\+\d{28}-\d{3}-\d{4}$" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("new password")+`<span class="error noneed">*</span>
                                	<span id="newpass_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="newpass"  class="form-control">
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("role")+`
                                		<span id="position_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                		<select id="position">
                                			`+library.POSITION.SelectOptionsList +`
                                		</select>
                                        </div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("ban")+`
                                		<span id="ban_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="ban" type="checkbox"  class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<label>`+lang.T("ban reason")+`
                                		<span id="ban_reason_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="ban_reason"  class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                </table>
        </div>
</div>`
	return out
}