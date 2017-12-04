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
	"strconv"
	"net/http"
	"html/template"
	"time"
	"github.com/ilyaran/Malina/libraries"
	"strings"
	"github.com/ilyaran/Malina/berry"
)

type Permission struct {
	OrderSelectOptions string
	InlistFields string
}

func (s *Permission)Index(malina *berry.Malina, w http.ResponseWriter){

	out := `
	<!-- conf -->
	<input id="dbtable" type="hidden" value="permission">
	<input id="Path_assets_uploads" type="hidden" value="`+app.Path_assets_uploads+`">
	`+s.InlistFields+`
	<input id="page" type="hidden" value="`+strconv.FormatInt(malina.Page,10)+`"/>
	<!-- conf -->


	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("permissions")+`</h3>
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
				<a>`+HomeLayoutView.Form_category_select("","category_id")+`</a>
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
			<li role="presentation" class="">
				<a>
					<input type="number" id="price_min" placeholder="`+lang.T("min price")+`"/>
					<input type="number" id="price_max" placeholder="`+lang.T("max price")+`"/>
				</a>
			</li>


		</ul>
		<div class="tab-content" id="myTabContent">
			<div class="tab-pane fade active in" role="tabpanel" id="items_table" aria-labelledby="home-tab">

				<div class="table-responsive">
					<table class="table table-bordered">
						<thead>
							<tr>
								<th style="width:5%">Id</th>
								<th width="5%">`+lang.T("logo")+`</th>
								<th width="5%">`+lang.T("category")+`</th>
								<th width="5%">`+lang.T("code")+`</th>
								<th style="width:5%">`+lang.T("price")+`</th>
								<th style="width:5%">`+lang.T("price1")+`</th>
								<th>`+lang.T("title")+`</th>
								<th width="5%">`+lang.T("sold")+`</th>
								`+HomeLayoutView.Table_head_last("permission")+`
							</tr>
						</thead>
						<tbody id="listing_permission">
							`+s.Listing(malina)+`
						</tbody>
					</table>
				</div>


			</div>
			<div class="tab-pane fade" role="tabpanel" id="item_form" aria-labelledby="profile-tab">
				<div class="row">
					<div class="col-md-8">

						<div class="table-responsive">
							<table class="table table-bordered">
								<tbody>
									<tr>
										<td >`+HomeLayoutView.Form_send_button()+`</td>
										<td >`+HomeLayoutView.ActionBar()+`</td>
									</tr>
									<tr>
										<td colspan="2">`+HomeLayoutView.ImageBar()+`</td>
									</tr>
									<tr>
										<td width="30%">`+lang.T("parameter")+`<span class="error"></span>
											<span id="parameter_error"  class="error"><span>
										</td>
										<td id="parameter_list">

										</td>
									</tr>

									`+HomeLayoutView.Form_title()+`
									`+HomeLayoutView.Form_checkbox("Enable","enable",true)+`
									<tr>
										<td>`+lang.T("code")+`<span id="code_error"  class="error"><span></td>
										<td><input id="code" value=""></td>
									</tr>
									<tr>
										<td>`+lang.T("price")+`<span id="price_error"  class="error"><span></td>
										<td><input id="price" value=""></td>
									</tr>
									<tr>
										<td>`+lang.T("price1")+`<span id="price1_error"  class="error"><span></td>
										<td><input id="price1" value=""></td>
									</tr>
									`+HomeLayoutView.Form_description()+`
									`+HomeLayoutView.Form_short_description()+`
									<tr>
										<td colspan="2">`+HomeLayoutView.Form_send_button()+`</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>


					<div class="col-md-4">

						<div class="bs-example bs-example-tabs" data-example-id="togglable-tabs">
							<ul class="nav nav-tabs" id="myTabs" role="tablist">
								<li role="presentation" class="active">
									<a class="btn btn-success" href="#parameters_tab" id="parameters-tab" role="tab" data-toggle="tab" aria-controls="parameters" aria-expanded="true">`+lang.T("parameters")+`</a>
								</li>
								<li role="presentation" class="">
									<a class="btn btn-success" href="#permissions_tab" id="permissions-tab" role="tab" data-toggle="tab" aria-controls="permissions" aria-expanded="false">`+lang.T("permissions")+`</a>
								</li>
							</ul>
						</div>

						<div class="tab-content" id="myTabContent">
							<div class="tab-pane fade active in" role="tabpanel" id="items_table" aria-labelledby="home-tab">

								<div class="panel panel-default">
							<div class="panel-heading">
								<h3 class="panel-title">`+lang.T("parameters")+`</h3>
							</div>
							<div class="panel-body">
								<div class="input-group">
									<span class="input-group-addon listing_button" id="basic-addon1"><span class="glyphicon glyphicon-search" aria-hidden="true">&nbsp;`+lang.T("Send")+`</span></span>
									<input type="text" id="search" class="form-control" placeholder="`+lang.T("search")+`" aria-describedby="basic-addon1">
								</div>
								<div class="table-responsive">
									<table class="table table-bordered">
										<thead>
											<tr>
												<th style="width:5%">Id</th>
												<th>`+lang.T("title")+`</th>
												<th style="width:5%">`+lang.T("check")+`</th>
											</tr>
										</thead>
										<tbody>`

	var idStr string
	for _,v:=range libraries.ParameterLib.Trees.List{
		idStr = strconv.FormatInt(v.GetId(), 10)
		out += `
											<tr>
												<td>`+idStr+`</td>
												<td>`+strings.Repeat(`<button type="button" class="btn btn-default btn-sm" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-right" aria-hidden="true"></span></button>`, v.(*entity.Parameter).GetLevel())+v.(*entity.Parameter).GetTitle()+`</td>
												<td><input title="`+lang.T("Switch to check")+`" value="" data-item_id = "` + idStr + `" class="permission_parameter" type="checkbox" data-toggle="toggle" /></td>
											</tr>

											`

	}


	out+=`
										</tbody>
									</table>
								</div>

							</div>
						</div>


							</div>
						</div>



					</div>


				</div>
			</div>
		</div>
	</div>
	<script src="`+app.Url_assets_home+`js/scripts.js?anticache=`+fmt.Sprintf("%d",time.Now().Unix())+`"></script>
	`

	malina.NavAuth = template.HTML(out)

	t := template.New("permission_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, malina)
}

func (s *Permission)Listing(malina *berry.Malina)string{
	if malina.List == nil {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}

	var out, idStr, imgSrc string

	for _,v:=range *malina.List {
		idStr = strconv.FormatInt(v.(*entity.Permission).GetId(), 10)
		if v.(*entity.Permission).GetImg() != nil {
			imgSrc = app.Url_assets_uploads + v.(*entity.Permission).GetImg()[0]
		} else {
			imgSrc = app.Url_no_image
		}
		out += `
		<tr>
			<th scope="row">` + idStr + `</th>
			<td><img src="` + imgSrc + `" width='60' height='50' /></td>
			<td>` + fmt.Sprintf("%v", v.(*entity.Permission).GetCategory()) + `</td>
			<td><input type="text" data-item_id="` + idStr + `" class="permission_code_inlist" value="` + v.(*entity.Permission).GetCode() + `"/></td>
			<td>` + fmt.Sprintf(`<input type="number" data-item_id="` + idStr + `" class="permission_price_inlist" value="%.2f"/>`, v.(*entity.Permission).GetPrice()) + `</td>
			<td>` + fmt.Sprintf(`<input type="number" data-item_id="` + idStr + `" class="permission_price1_inlist" value="%.2f"/>`, v.(*entity.Permission).GetPrice1()) + `</td>
			<td><input type="text" data-item_id="` + idStr + `" class="permission_title_inlist" value="` + v.(*entity.Permission).GetTitle() + `"/></td>
			<td>` + fmt.Sprintf("%v", v.(*entity.Permission).GetSold()) + `</td>
			<td>` + HomeLayoutView.Table_colomun_last(idStr,"permission",v.(*entity.Permission).GetEnable())+ `</td>
		</tr>
		`
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