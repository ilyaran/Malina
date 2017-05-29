package positionView

import (
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/language"
	"strconv"
	"github.com/ilyaran/Malina/views"
	"strings"
	"github.com/ilyaran/Malina/libraries"
)

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
		Enable<input checked type="checkbox" onchange="if(this.checked){$('.inlist_enable').bootstrapToggle('on');}else{$('.inlist_enable').bootstrapToggle('off');}">
		<button id="submitInlistButton" class="btn btn-success"><span class="glyphicon glyphicon-save" aria-hidden="true"></span></button>
		Delete<input type="checkbox"  onchange="if(this.checked){$('.inlist_del').prop('checked',true);}else{$('.inlist_del').prop('checked',false);}">
	</th>
	`
func Index(positionList []*entity.Position, paging string)string{
	views.TABLE_FORM.SetNull()

	views.TABLE_FORM.Inputs = views.ICE_FORM.Inputs("position")
	views.TABLE_FORM.Breadcrumb = "Home / Positions"
	views.TABLE_FORM.Select_options_order_by = select_options_order_by
	views.TABLE_FORM.Head = table_head
	views.TABLE_FORM.Listing = Listing(positionList,paging)
	views.TABLE_FORM.Form = Form()
	views.TABLE_FORM.BuildIndexForm()

	return views.TABLE_FORM.Out + views.Footer()
}

func Listing(positionList []*entity.Position, paging string)string{
	var out,idStr string
	if positionList != nil && len(positionList)>0 {
		for _, i := range positionList {
			idStr = strconv.FormatInt(i.GetId(),10)
			var checked = ``
			if i.GetEnable() {checked = `checked`}
			out += `
			<tr class="even gradeA">
				<td>` + idStr + `</td>
				<td>` + strconv.FormatInt(i.GetParent().GetId(),10) + `</td>
				<td><input data-item_id="` + idStr + `" style="width:50px;" class="inlist_position_sort" type="number" value="` + strconv.FormatInt(i.GetSort(),10) + `"/></td>
				<td>
				` + strings.Repeat("&rarr;", i.GetLevel()) + `&nbsp;`
			if _,ok := library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetDescendantIdsMap()[i.GetId()]; ok || library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()].GetParent().GetId() == 0 {
				out += `
					<input data-item_id="` + idStr + `" class="inlist_position_title" type="text" value="` + i.GetTitle() +`"/>
				</td>
				<td>
					<input value="" data-item_id="` + idStr + `" class="inlist_position_enable" type="checkbox" `+checked+` data-toggle="toggle" >
					<button data-item_id="` + idStr + `" class="btn btn-primary edit_item"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></button>
					<button data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
					<input value="0" data-item_id="` + idStr + `" class="inlist_del" type="checkbox">
				</td>`

			}else {
				out += i.GetTitle()+`</td><td><td></td>`
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
		`Parent Position <span id="select_parent">
			<select id="parent">
                                ` + GetSelectOptionsListView(library.POSITION.TreeMap[library.SESSION.GetSessionObj().GetPositionId()],"Root",true,true) +`
                        </select>
		</span>` +
		`<br>
		<button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
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
                        `+views.ICE_FORM.CheckBox("enable","enable",true)+`
                        `+views.ICE_FORM.Title()+`

                </table>
                <button class="btn btn-primary submitButton">`+lang.T("Send")+`</button>
        </div>
</div>`
	return out
}

func GetSelectOptionsListView(currentPosition *entity.Position, rootTitle string, isIncludedOwnId, hasRoot bool)string{

	if isIncludedOwnId || currentPosition.GetParent().GetId() == 0 {
		currentPosition.GetDescendantIdsMap()[currentPosition.GetId()]=false

	}
	positionOptionsList := library.POSITION.BuildSelectOptionsView(nil, currentPosition.GetDescendantIdsMap(), 0)
	if currentPosition.GetParent().GetId() == 0 {
		if hasRoot {
			positionOptionsList = `<option value="0">`+rootTitle+`</option>` + positionOptionsList
		}
	}
	return positionOptionsList
}
