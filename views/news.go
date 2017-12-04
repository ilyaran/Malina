package views


/*

import (
	"fmt"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/lang"
	"strconv"
	"net/http"
	"html/template"
	"time"
	"strings"
)

type News struct {
	List   *[app.Per_page_max]entity.Scanable
	Paging string
	NavAuth template.HTML
	OrderSelectOptions string
	InlistFields string
}

func (s *News)Index(w http.ResponseWriter){

	s.NavAuth = template.HTML(`
	<!-- conf -->
	<input id="dbtable" type="hidden" value="news">
	<input id="Path_assets_uploads" type="hidden" value="`+app.Path_assets_uploads+`">
	`+s.InlistFields+`
	<!-- conf -->


	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("newss")+`</h3>
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
								<th width="5%">`+lang.T("logo")+`</th>
								<th width="5%">`+lang.T("views")+`</th>

								<th width="5%">`+lang.T("likes")+`</th>
								<th width="15%">`+lang.T("title")+`</th>
								<th>`+lang.T("description")+`</th>
								`+HomeLayoutView.Table_head_last("news")+`
							</tr>
						</thead>
						<tbody id="listing_news">
							`+s.Listing()+`
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
								<td colspan="2">`+HomeLayoutView.ImageBar()+`</td>
							</tr>


							`+HomeLayoutView.Form_title()+`
							`+HomeLayoutView.Form_checkbox("Enable","enable",true)+`

							`+HomeLayoutView.Form_description()+`
							`+HomeLayoutView.Form_short_description()+`
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
	`)


	t := template.New("news_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, s)
}

func (s *News)Listing()string{
	if s.List[0].GetId() == -1 {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}

	var out, idStr, imgSrc string

	for i := 0; s.List[i].GetId() != -1; i++ {
		idStr = strconv.FormatInt(s.List[i].GetId(), 10)
		if s.List[i].(*entity.News).GetImg() != nil {
			imgSrc = app.Url_assets_uploads + s.List[i].(*entity.News).GetImg()[0]
		} else {
			imgSrc = app.Url_no_image
		}
		out += `
		<tr>
			<th scope="row">` + idStr + `</th>
			<td><img src="` + imgSrc + `" width='60' height='50' /></td>
			<td>` + fmt.Sprintf("%v", s.List[i].(*entity.News).GetViews()) + `</td>

			<td>` + fmt.Sprintf("%v", s.List[i].(*entity.News).GetLike()) + `</td>
			<td>` + s.List[i].(*entity.News).GetTitle() + `</td>
			<td>` + strings.Replace(s.List[i].(*entity.News).GetDescription(),"<","",-1) + `</td>
			<td>` + HomeLayoutView.Table_colomun_last(idStr,"news",s.List[i].(*entity.News).GetEnable())+ `</td>
		</tr>
		`
	}
	out += `
		<tr>
			<td colspan="8">
				<nav aria-label="...">
					<ul class="pagination">
						` + s.Paging + `
					</ul>
				</nav>
			</td>
		</tr>`

	return out
}*/
