package permissionView

import (
	"Malina/entity"
	"Malina/language"
	"strconv"
	"Malina/config"
	"Malina/views"
	"Malina/libraries"
)

func Index(permissionList []*entity.Permission, paging string)string{
	var inputs = `
	<input type="hidden" value="permission" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_permission_ajax()+`" id="url_permission_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_permission_get()+`" id="url_permission_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_permission_add()+`" id="url_permission_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_permission_edit()+`" id="url_permission_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_permission_del()+`" id="url_permission_del"/>`
	var out = inputs + `
<div class="row">
	<div class="col-md-12">

		<div class="panel panel-default">
			<div class="panel-heading">
				<h3 class="panel-title">Home / Permissions</h3>
			</div>
		  	<div class="panel-body">
		    		`+lang.T("per page")+`
				`+views.FACE_FORM.PerPageSelectForm()+`
				&nbsp;&nbsp;
				<span>`+lang.T("order_by")+`</span>
				<select id="order_by" >
					<option value="1">`+lang.T("id")+`&uarr;</option>
					<option value="2">`+lang.T("id")+`&darr;</option>
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
							<th>`+lang.T("role id")+`</th>
							<th>`+lang.T("role title")+`</th>
							<th>`+lang.T("data")+`</th>
							<th></th>
						</tr>

					</thead>
					<tbody id="listing">
					`
						out += Listing(permissionList,paging)+`

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

func Listing(permissionList []*entity.Permission, paging string)string{
	var out,idStr string
	if permissionList != nil && len(permissionList)>0 {
		for _, i := range permissionList {
			idStr = strconv.FormatInt(i.GetId(),10)
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td>` + strconv.FormatInt(i.GetPosition().GetId(),10) + `</td>
				<td>` + i.GetPosition().GetTitle() + `</td>
				<td>` + i.GetData()  + `</td>
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
			<div id="error" class="error"></div>
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

                <table class="table table-striped table-bordered table-hover" >
                	<tr>
                        	<td>
                                	<label>`+lang.T("role")+`
                                		<span id="position_id_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                		<select id="position_id">
                                			`+library.POSITION.SelectOptionsList +`
                                		</select>
                                        </div>
                               	</td>
                        </tr>
                     	<tr>
                        	<td>
                                	<label>`+lang.T("data")+`<span class="error">*</span>
                                		<span id="data_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="data"  class="form-control" />
                                	</div>
                               	</td>
                        </tr>

                </table>
        </div>
</div>`
	return out
}
