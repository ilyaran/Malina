package permissionView

import (
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/language"
	"strconv"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/libraries"
	"github.com/ilyaran/Malina/views/position"
)


var select_options_order_by string = `
	<option value="1">`+lang.T("ID")+`&uarr;</option>
	<option value="2">`+lang.T("ID")+`&darr;</option>
	<option value="3">`+lang.T("data")+`&uarr;</option>
	<option value="4">`+lang.T("data")+`&darr;</option>
`
var table_head = `
	<th width="5%">ID</th>
	<th width="5%">`+lang.T("position id")+`</th>
	<th width="10%">`+lang.T("position title")+`</th>
	<th >`+lang.T("data")+`</th>
	<th width="20%">
		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Index(permissionList []*entity.Permission, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = views.ICE_FORM.Inputs("permission")
	views.TABLE_FORM.Breadcrumb = "Home / Permissions"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(permissionList,paging)
	views.TABLE_FORM.Form = Form()
	views.TABLE_FORM.BuildIndexForm()

	return views.TABLE_FORM.Out + views.Footer()
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
				<td>`
			if _,ok := library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetDescendantIdsMap()[i.GetPosition().GetId()];
				ok || library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetParent().GetId() == 0 {
				out+=`<input data-item_id="` + idStr + `" class="inlist_permission_data" type="text" value="` + i.GetData() +`"/>
				</td>
				<td>
					<button data-item_id="` + idStr + `" class="btn btn-primary edit_item"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></button>
					<button data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
					<input value="0" data-item_id="` + idStr + `" class="inlist_del" type="checkbox">
				</td>
				`}else {
				out += i.GetData()+`</td><td><td></td>`
			}
			out+=`
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
		`<span id="select_position">
			<select id="position_id">
                        	` + positionView.GetSelectOptionsListView(library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()],"Root",false,false) +`
                        </select>
                        <span id="position_id_error"  class="error"></span>
		</span>` +
		`<br>
		<button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
                <table class="table table-striped table-bordered table-hover" >
                	<tr>
                        	<td>
                                	<div class="form-group">
                                            	<label>`+lang.T("data")+`
						<span id="data_error"  class="error"></span></label>
                                        </div>
                                </td>
                                <td>
                                	<div class="form-group">
                                		<textarea width="100%" id="data"></textarea>
                                	</div>
                               	</td>
                        </tr>
                </table>
                <button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}
