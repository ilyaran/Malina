package productView

import (
	"strconv"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/libraries"
	"fmt"
)

var inputs string = `
	<input type="hidden" value="product" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_product_ajax()+`" id="url_product_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_product_get()+`" id="url_product_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_product_add()+`" id="url_product_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_product_edit()+`" id="url_product_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_product_del()+`" id="url_product_del"/>

`
var select_options_order_by string = `
	<option value="1">`+lang.T("updated")+`&darr;</option>
	<option value="2">`+lang.T("updated")+`&uarr;</option>
	<option value="3">`+lang.T("price")+`&uarr;</option>
	<option value="4">`+lang.T("price")+`&darr;</option>
	<option value="5">`+lang.T("id")+`&uarr;</option>
	<option value="6">`+lang.T("id")+`&darr;</option>`

func Index(productList []*entity.Product, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = inputs
	views.TABLE_FORM.Breadcrumb = "Home / Products"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(productList,paging)
	views.TABLE_FORM.Form = Form()
	views.TABLE_FORM.Option1 = views.ICE_FORM.CategorySelect("All","category")
	views.TABLE_FORM.Option2 = `<input type="number" id="price_min" placeholder="Min price"/><input type="number" placeholder="Max price" id="price_max"/>`

	views.TABLE_FORM.BuildIndexForm()
	return `
	<script type="text/javascript" src="`+app.Assets_path()+`ckeditor/ckeditor.js"></script>
	` + views.TABLE_FORM.Out + views.Footer()
}
var table_head = `
	<th width="5%">ID</th>
	<th width="5%">`+lang.T("logo")+`</th>
	<th >`+lang.T("title")+`</th>
	<th width="5%">`+lang.T("price")+`</th>
	<th width="5%">`+lang.T("price1")+`</th>
	<th width="5%">`+lang.T("code")+`</th>
	<th width="20%">
		Enable<input checked type="checkbox" onchange="if(this.checked){$('.inlist_product_enable').bootstrapToggle('on');}else{$('.inlist_product_enable').bootstrapToggle('off');}">
		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Listing(productList []*entity.Product, paging string)string{
	var out,idStr,imgSrc string
	if productList != nil && len(productList)>0 {
		for _, i := range productList {
			idStr = strconv.FormatInt(i.Get_id(),10)
			if i.Get_img() !=nil {imgSrc = i.Get_img()[0]} else {imgSrc = app.No_image()}
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td><img src="` + imgSrc + `" width='60' height='50' /></td>
				<td>
					<div>
						<input data-item_id="` + idStr + `" class="inlist_product_title" type="text" value="` + i.Get_title() +`"/>
					</div>
					<div>`; if v,ok := library.CATEGORY.TreeMap[i.Get_category_id()]; ok {out += v.Get_path()}; out += `</div>
				</td>
				<td>
					<input data-item_id="` + idStr + `" class="inlist_product_price" type="text" value="` + fmt.Sprintf("%.2f",i.Get_price()) +`"/>
				</td>
				<td>
					<input data-item_id="` + idStr + `" class="inlist_product_price1" type="text" value="` + fmt.Sprintf("%.2f",i.Get_price1()) +`"/>
				</td>
				<td>
					<input data-item_id="` + idStr + `" class="inlist_product_code" type="text" value="` + i.Get_code() +`"/>
				</td>
				<td>
					<input value="" data-item_id="` + idStr + `" class="inlist_product_enable" type="checkbox" `; if i.Get_enable() {out+=`checked`}; out+=` data-toggle="toggle" >
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
		views.ICE_FORM.BarForms() + `&nbsp;&nbsp;` +
		`<span id="select_category">` + views.ICE_FORM.CategorySelect("","category_id") + `</span>` +
		`<br>` +
		views.ICE_FORM.ImageBar() + `
		<button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
                <table class="table table-striped table-bordered table-hover" >
                	<tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("price")+`
						<span id="price_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="price"  value="" type="number" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("price 1")+`
						<span id="price1_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="price1"  value="" type="number" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("code")+`
						<span id="code_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="code"  value="" type="text" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        `+views.ICE_FORM.CheckBox("enable","enable",true)+`
                        `+views.ICE_FORM.Title()+`
                        `+views.ICE_FORM.Description()+`
                </table>
                <button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}