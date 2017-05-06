package positionView

import (
	"Malina/entity"
	"Malina/language"
	"strconv"
	"Malina/config"
	"Malina/views"
	"strings"
	"Malina/libraries"
)

func Index(positionList []*entity.Position, paging string)string{
	var inputs = `
	<input type="hidden" value="position" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_position_ajax()+`" id="url_position_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_position_get()+`" id="url_position_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_position_add()+`" id="url_position_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_position_edit()+`" id="url_position_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_position_del()+`" id="url_position_del"/>`
	var out = inputs + `
<div class="row">
	<div class="col-md-12">

		<div class="panel panel-default">
			<div class="panel-heading">
				<h3 class="panel-title">Home / Positions</h3>
			</div>
		  	<div class="panel-body">
		    		`+lang.T("per page")+`
				`+views.FACE_FORM.PerPageSelectForm()+`
				&nbsp;&nbsp;
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
							<th>`+lang.T("title")+`</th>
							<th></th>
						</tr>

					</thead>
					<tbody id="listing">
					`
						out += Listing(positionList,paging)+`

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

func Listing(positionList []*entity.Position, paging string)string{
	var out,idStr string
	if positionList != nil && len(positionList)>0 {
		for _, i := range positionList {
			idStr = strconv.FormatInt(i.GetId(),10)
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td>` + strings.Repeat("&rarr;", i.GetLevel()) + `&nbsp;` + i.GetTitle()  + `</td>
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
                                	<label>`+lang.T("parent")+`
                                		<span id="parent_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                		<select id="parent">
                                			<option value="0">Root</option>
                                			`+library.POSITION.SelectOptionsList +`
                                		</select>
                                        </div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("sort")+`
						<span id="sort_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="sort"  value="100" type="number" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
			<tr>
                        	<td>
                                	<label>`+lang.T("enable")+`</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="enable" type="checkbox" checked class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                	<tr>
                        	<td>
                                	<label>`+lang.T("title")+`<span class="error">*</span>
                                		<span id="title_error"  class="error"><span>
                                	</label>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="title"  class="form-control" />
                                	</div>
                               	</td>
                        </tr>

                </table>
        </div>
</div>`
	return out
}