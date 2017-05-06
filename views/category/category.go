package categoryView

import (
	"strconv"
	"Malina/entity"
	"Malina/views"
	"Malina/language"
	"Malina/config"
	"strings"
)

var inputs string = `
	<input type="hidden" value="category" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_category_ajax()+`" id="url_category_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_category_get()+`" id="url_category_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_category_add()+`" id="url_category_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_category_edit()+`" id="url_category_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_category_del()+`" id="url_category_del"/>

`
var select_options_order_by string = `
	<option value="1">`+lang.T("sort")+`&uarr;</option>
	<option value="2">`+lang.T("sort")+`&darr;</option>
	<option value="3">`+lang.T("ID")+`&uarr;</option>
	<option value="4">`+lang.T("ID")+`&darr;</option>
	<option value="5">`+lang.T("title")+`&uarr;</option>
	<option value="6">`+lang.T("title")+`&darr;</option>
`
var table_head = `
	<th width="5%">ID</th>
	<th width="5%">`+lang.T("parent")+`</th>
	<th width="5%">`+lang.T("sort")+`</th>
	<th >`+lang.T("title")+`</th>
	<th width="20%">
		Enable<input checked type="checkbox" onchange="if(this.checked){$('.inlist_category_enable').bootstrapToggle('on');}else{$('.inlist_category_enable').bootstrapToggle('off');}">
		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Index(categoryList []*entity.Category, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = inputs
	views.TABLE_FORM.Breadcrumb = "Home / Categories"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(categoryList,paging)
	views.TABLE_FORM.Form = Form()
	views.TABLE_FORM.IndexFrom()

	return `
	<script type="text/javascript" src="`+app.Assets_path()+`ckeditor/ckeditor.js"></script>
	` + views.TABLE_FORM.Out + views.Footer()
}

func Listing(categoryList []*entity.Category, paging string)string{
	var out,idStr string
	if categoryList != nil && len(categoryList)>0 {
		for _, i := range categoryList {
			idStr = strconv.FormatInt(i.Get_id(),10)
			var imgSrc string = app.Base_url()+"assets/img/noimg.jpg"
			if len(i.Get_img())>0{imgSrc = i.Get_img()[0]}
			var checked = ``
			if i.Get_enable() {checked = `checked`}
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td>` + strconv.FormatInt(i.Get_parent(),10) + `</td>
				<td><input data-item_id="` + idStr + `" style="width:50px;" class="inlist_category_sort" type="number" value="` + strconv.FormatInt(i.Get_sort(),10) + `"/></td>
				<td>` + strings.Repeat("&rarr;", i.Get_level()) + `
					&nbsp;<img src="` + imgSrc + `" width='60' height='50' />
					&nbsp;<input data-item_id="` + idStr + `" class="inlist_category_title" type="text" value="` + i.Get_title() +`"/>
					(`+ strconv.FormatInt(i.Get_quantity(),10) +`)
				</td>
				<td>
					<input value="" data-item_id="` + idStr + `" class="inlist_category_enable" type="checkbox" `+checked+` data-toggle="toggle" >
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
		views.FACE_FORM.BarForms() + `&nbsp;&nbsp;` +
		`<span id="select_parent">` + views.FACE_FORM.CategorySelect("Root","parent") + `</span>` +
		`<br>` +
		views.FACE_FORM.ImageBar() + `
                <table class="table table-striped table-bordered table-hover" >
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
                        `+views.FACE_FORM.Enable("enable")+`
                        `+views.FACE_FORM.Title()+`
                        `+views.FACE_FORM.Description()+`
                </table>
                <button id="submitButton" class="btn btn-primary">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}
