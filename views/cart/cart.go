package cartView

import (
	"Malina/entity"
	"Malina/views"
	"Malina/language"
	"fmt"
	"Malina/config"
)

var inputs string = `
	<input type="hidden" value="cart" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_cart_ajax()+`" id="url_cart_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_cart_get()+`" id="url_cart_get"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_cart_add()+`" id="url_cart_add"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_cart_edit()+`" id="url_cart_edit"/>
	<input type="hidden" value="`+app.Base_url()+app.Uri_cart_del()+`" id="url_cart_del"/>

`
var select_options_order_by string = `
	<option value="1">`+lang.T("price")+`&uarr;</option>
	<option value="2">`+lang.T("price")+`&darr;</option>
	<option value="3">`+lang.T("ID")+`&uarr;</option>
	<option value="4">`+lang.T("ID")+`&darr;</option>
	<option value="5">`+lang.T("title")+`&uarr;</option>
	<option value="6">`+lang.T("title")+`&darr;</option>
`

func Index(cartList []*entity.Cart, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = inputs
	views.TABLE_FORM.Breadcrumb = "Home / Carts"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(cartList,paging)
	views.TABLE_FORM.Form = Form()
	views.TABLE_FORM.IndexFrom()

	return `
	<script type="text/javascript" src="`+app.Assets_path()+`ckeditor/ckeditor.js"></script>
	` + views.TABLE_FORM.Out + views.Footer()
}
var table_head = `
	<th width="5%">ID</th>
	<th>`+lang.T("product title")+`</th>

<th></th>
<th></th>
<th></th>
<th></th>
<th>Buy now<input checked type="checkbox" onchange="if(this.checked){$('.inlist_cart_enable').bootstrapToggle('on');}else{$('.inlist_cart_enable').bootstrapToggle('off');}"></th>

	<th width="20%">

		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Listing(cartList []*entity.Cart, paging string)string{
	var out,idStr string
	if cartList != nil && len(cartList)>0 {
		for _, i := range cartList {
			idStr = i.Get_id()
			for k,v :=range i.Get_product(){
				out += `
				<tr class="even gradeA">
					<td rowspan="">` + idStr[0:12] + `</td>`

					out += `
						<td>`+v.Get_title()+`</td>
						<td>`+fmt.Sprintf("%.2f",v.Get_price())+`</td>
						<td>`+fmt.Sprintf("%.2f",v.Get_price1())+`</td>
						<td>`+fmt.Sprintf("%.2f",i.Get_quantity()[k])+`</td>
						<td>`+fmt.Sprintf("%.2f",v.Get_price() + i.Get_quantity()[k])+`</td>
						<td><input value="" data-item_id="` + idStr + `" class="inlist_cart_buy_now" type="checkbox" `; if i.Get_buy_now()[k] {out+=`checked`}; out+=` data-toggle="toggle" ></td>
						`
					out += `
					<td>
						<button data-item_id="` + idStr + `" class="btn btn-primary edit_item"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></button>
						<button data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
						<input value="0" data-item_id="` + idStr + `" class="inlist_del" type="checkbox">
					</td>
				</tr>`
			}
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
		views.FACE_FORM.BarForms() + `&nbsp;&nbsp;
                <table class="table table-striped table-bordered table-hover" >
                	<tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("product ID")+`
						<span id="product_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="_product"  value="" type="number" class="form-control" />
                                	</div>
                               	</td>
                        </tr>
                        <tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("quantity")+`
						<span id="quantity_error"  class="error"><span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                        	<input id="quantity"  value="" type="number" class="form-control" />
                                	</div>
                               	</td>
                        </tr>

                        `+views.FACE_FORM.Enable("buy_now")+`

                </table>
                <button id="submitButton" class="btn btn-primary">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}