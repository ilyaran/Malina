/**
 *
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ package views

import (
	"fmt"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/lang"
	"github.com/ilyaran/Malina/libraries"
	"strconv"
	"net/http"
	"html/template"
	"time"
	"strings"
	"github.com/ilyaran/Malina/berry"
)

type Parameter struct {
	OrderSelectOptions string
	InlistFields string
}

func (s *Parameter)Index(malina *berry.Malina,hasTreeList bool,w http.ResponseWriter){
	var out = `
	<!-- conf -->
	<input id="dbtable" type="hidden" value="parameter">
	<input id="Path_assets_uploads" type="hidden" value="`+app.Path_assets_uploads+`">
	`+s.InlistFields+`
	<input id="page" type="hidden" value="`+strconv.FormatInt(malina.Page,10)+`"/>
	<!-- conf -->


	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("parameters")+`</h3>
		</div>
		<div class="panel-body">
			<div class="input-group">
			  <span class="input-group-addon listing_button" id="basic-addon1"><span class="glyphicon glyphicon-send" aria-hidden="true">&nbsp;`+lang.T("Send")+`</span></span>
			  <input type="text" id="search" class="form-control" placeholder="`+lang.T("search")+`" aria-describedby="basic-addon1">
			</div>
		</div>
	</div>

	<div class="bs-example bs-example-tabs" data-example-id="togglable-tabs">
		<ul class="nav nav-tabs" id="myTabs" role="tablist">
			<li role="presentation" class="active">
				<a class="btn btn-success" href="#items_table" id="items_table-tab" role="tab" data-toggle="tab" aria-controls="home" aria-expanded="true">`+lang.T("list")+`</a>
			</li>
			<li role="presentation" class="">
				<a class="btn btn-success" href="#item_form" id="item_form-tab" role="tab" data-toggle="tab" aria-controls="profile" aria-expanded="false">`+lang.T("form")+`</a>
			</li>

			<li role="presentation" class="">
				<a>
					`+lang.T("order by")+`
					<select id="order_by" role="tab">
						`+s.OrderSelectOptions+`
					</select>
				</a>
			</li>

			<li role="presentation" class="">
				<a>`+HomeLayoutView.Per_page_select()+`</a>
			</li>


		</ul>
		<div class="tab-content" id="myTabContent">
			<div class="tab-pane fade active in" role="tabpanel" id="items_table" aria-labelledby="home-tab">
				<div class="table-responsive">

					<table class="table table-bordered">
						<thead>
							<tr>
								<th style="width:5%">Id</th>
								<th style="width:5%">`+lang.T("parent")+`</th>
								<th style="width:5%">`+lang.T("sort")+`</th>
								<th >`+lang.T("title/description")+`</th>
								`+HomeLayoutView.Table_head_last("parameter")+`

							</tr>
						</thead>
						<tbody id="listing_parameter">
							`
							if hasTreeList{
								out+=s.Listing(malina)
							}else {
								out+=s.ListingNonTree(malina)
							}
							out+=`
						</tbody>
					</table>

				</div>
			</div>
			<div class="tab-pane fade" role="tabpanel" id="item_form" aria-labelledby="profile-tab">
				<div class="table-responsive">
					<table class="table table-bordered">
						<tbody>
							<tr>
								<td >`+HomeLayoutView.Form_send_button()+`</td>
								<td >`+HomeLayoutView.ActionBar()+`</td>
							</tr>
							<tr>
								<td width="30%">`+lang.T("parent")+`<span class="error required"></span>
									<span id="parent_error"  class="error"><span>
								</td>
								<td>
									`+HomeLayoutView.Form_parameter_select(lang.T("root"),"parent")+`
								</td>
							</tr>
							`+HomeLayoutView.Form_title()+`
							`+HomeLayoutView.Form_sort()+`
							`+HomeLayoutView.Form_checkbox("enable","enable",true)+`
							`+HomeLayoutView.Form_description()+`
							<tr>
								<td colspan="2">`+HomeLayoutView.Form_send_button()+`</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
	<script src="`+app.Url_assets_home+`js/scripts.js?anticache=`+fmt.Sprintf("%d",time.Now().Unix())+`"></script>
	`

	malina.NavAuth = template.HTML(out)

	t := template.New("parameter_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, malina)
}

func (s *Parameter)Listing(malina *berry.Malina,)string{
	if len(libraries.ParameterLib.Trees.List) == 0 {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	var out, idStr string
	for _,v := range libraries.ParameterLib.Trees.List[malina.Page:malina.Per_page] {
			idStr = strconv.FormatInt(v.GetId(), 10)
			out += `
			<tr class="tree_row" data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Parameter).GetLevel()) + `" >
				<th scope="row">` + idStr + `</th>
				<td>` + fmt.Sprintf("%v", v.GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="parameter_sort_inlist" value="` + strconv.FormatInt(v.(*entity.Parameter).GetSort(), 10) + `"/></td>
				<td>
					` + strings.Repeat(`<button type="button" class="btn btn-default btn-sm" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-right" aria-hidden="true"></span></button>`, v.(*entity.Parameter).GetLevel()) + `
					`
				if v.(*entity.Parameter).Descendants != nil {
					out += `<button data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Parameter).GetLevel()) + `" type="button" class="btn btn-success btn-sm button_collapse" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-down" aria-hidden="true"></span></button>`
				}
				out += `
					<input type="text" data-item_id="` + idStr + `" class="parameter_title_inlist" value="` + v.(*entity.Parameter).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="parameter_description_inlist" style="width:100%" >` + v.(*entity.Parameter).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "parameter", v.(*entity.Parameter).GetEnable()) + `</td>
			</tr>
			`
		//}
	}
	out += `
		<tr>
			<td colspan="8">
				<nav aria-label="...">
					<ul class="pagination">
						` + malina.Paging + `
					</ul>
				</nav>
			</td>
		</tr>`

	return out
}


func (s *Parameter)ListingNonTree(malina *berry.Malina,)string{
	if malina.List == nil {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	var out, idStr string
	for _,v:=range *malina.List {
		//if s.SearchInList == "" || strings.Contains(v.(*entity.Parameter).GetTitle()+v.(*entity.Parameter).GetDescription(),s.SearchInList) {
		idStr = strconv.FormatInt(v.GetId(), 10)
		out += `
			<tr>
				<th scope="row">` + idStr + `</th>
				<td>` + fmt.Sprintf("%v", v.(*entity.Parameter).GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="parameter_sort_inlist" value="` + strconv.FormatInt(v.(*entity.Parameter).GetSort(), 10) + `"/></td>
				<td>
					<input type="text" data-item_id="` + idStr + `" class="parameter_title_inlist" value="` + v.(*entity.Parameter).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="parameter_description_inlist" style="width:100%" >` + v.(*entity.Parameter).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "parameter", v.(*entity.Parameter).GetEnable()) + `</td>
			</tr>
			`
		//}
	}
	out += `
		<tr>
			<td colspan="8">
				<nav aria-label="...">
					<ul class="pagination">
						` + malina.Paging + `
					</ul>
				</nav>
			</td>
		</tr>`

	return out
}





















