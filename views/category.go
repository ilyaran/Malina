package views


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

type Category struct {
	OrderSelectOptions string
	InlistFields string
}

func (s *Category)Index(malina *berry.Malina,hasTreeList bool,w http.ResponseWriter){
	var out = `
	<!-- conf -->
	<input id="dbtable" type="hidden" value="category">
	<input id="Path_assets_uploads" type="hidden" value="`+app.Path_assets_uploads+`">
	`+s.InlistFields+`
	<input id="page" type="hidden" value="`+strconv.FormatInt(malina.Page,10)+`"/>
	<!-- conf -->


	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("categories")+`</h3>
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
								`+HomeLayoutView.Table_head_last("category")+`

							</tr>
						</thead>
						<tbody id="listing_category">
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
								<td colspan="2">`+HomeLayoutView.ImageBar()+`</td>
							</tr>
							<tr>
								<td width="30%">`+lang.T("parent")+`<span class="error required"></span>
									<span id="parent_error"  class="error"><span>
								</td>
								<td>
									`+HomeLayoutView.Form_category_select(lang.T("root"),"parent")+`
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

	t := template.New("category_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, malina)
}

func (s *Category)Listing(malina *berry.Malina,)string{
	if len(libraries.CategoryLib.Trees.List) == 0 {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	var out, idStr string
	for _,v := range libraries.CategoryLib.Trees.List[malina.Page:malina.Per_page] {
		idStr = strconv.FormatInt(v.GetId(), 10)
		out += `
			<tr class="tree_row" data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Category).GetLevel()) + `" >
				<th scope="row">` + idStr + `</th>
				<td>` + fmt.Sprintf("%v", v.GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="category_sort_inlist" value="` + strconv.FormatInt(v.(*entity.Category).GetSort(), 10) + `"/></td>
				<td>
					` + strings.Repeat(`<button type="button" class="btn btn-default btn-sm" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-right" aria-hidden="true"></span></button>`, v.(*entity.Category).GetLevel()) + `
					`
		if v.(*entity.Category).Descendants != nil {
			out += `<button data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Category).GetLevel()) + `" type="button" class="btn btn-success btn-sm button_collapse" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-down" aria-hidden="true"></span></button>`
		}
		out += `
					<input type="text" data-item_id="` + idStr + `" class="category_title_inlist" value="` + v.(*entity.Category).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="category_description_inlist" style="width:100%" >` + v.(*entity.Category).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "category", v.(*entity.Category).GetEnable()) + `</td>
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


func (s *Category)ListingNonTree(malina *berry.Malina,)string{
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
		//if s.SearchInList == "" || strings.Contains(v.(*entity.Category).GetTitle()+v.(*entity.Category).GetDescription(),s.SearchInList) {
		idStr = strconv.FormatInt(v.GetId(), 10)
		out += `
			<tr>
				<th scope="row">` + idStr + `</th>
				<td>` + fmt.Sprintf("%v", v.(*entity.Category).GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="category_sort_inlist" value="` + strconv.FormatInt(v.(*entity.Category).GetSort(), 10) + `"/></td>
				<td>
					<input type="text" data-item_id="` + idStr + `" class="category_title_inlist" value="` + v.(*entity.Category).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="category_description_inlist" style="width:100%" >` + v.(*entity.Category).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "category", v.(*entity.Category).GetEnable()) + `</td>
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






/*


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
)

type Category struct {
	List           		*[app.Per_page_max]entity.Scanable
	Paging         		string
	ListIndexStart 		int64
	ListindexEnd   		int64
	SearchInList   		string

	NavAuth                	template.HTML
	OrderSelectOptions 	string
	InlistFields       	string


}

func (s *Category)Index(hasTreeList bool,w http.ResponseWriter){
	var out = `
	<!-- conf -->
	<input id="dbtable" type="hidden" value="category">
	<input id="Path_assets_uploads" type="hidden" value="`+app.Path_assets_uploads+`">
	`+s.InlistFields+`
	<!-- conf -->


	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("categories")+`</h3>
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
								<th style="width:5%">`+lang.T("parent")+`</th>
								<th style="width:5%">`+lang.T("sort")+`</th>
								<th >`+lang.T("title/description")+`</th>
								`+HomeLayoutView.Table_head_last("category")+`

							</tr>
						</thead>
						<tbody id="listing_category">
							`
	if hasTreeList{
		out+=s.Listing()
	}else {
		out+=s.ListingNonTree()
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
								<td colspan="2">`+HomeLayoutView.ImageBar()+`</td>
							</tr>
							<tr>
								<td width="30%">`+lang.T("parent")+`<span class="error required"></span>
									<span id="parent_error"  class="error"><span>
								</td>
								<td>
									`+HomeLayoutView.Form_category_select(lang.T("root"),"parent")+`
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

	s.NavAuth = template.HTML(out)
	t := template.New("category_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, s)
}

func (s *Category)Listing()string{
	if len(libraries.CategoryLib.Tree.List) == 0 {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	var out, idStr,imgSrc string

	for _,v := range libraries.CategoryLib.Tree.List[s.ListIndexStart:s.ListindexEnd] {
		idStr = strconv.FormatInt(v.GetId(), 10)
		if v.(*entity.Category).GetImg() != nil {
			imgSrc = app.Url_assets_uploads + v.(*entity.Category).GetImg()[0]
		} else {
			imgSrc = app.Url_no_image
		}
		out += `
			<tr class="tree_row" data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Category).GetLevel()) + `" >
				<th scope="row">` + idStr + `</th>
				<td><img src="` + imgSrc + `" width='60' height='50' /></td>
				<td>` + fmt.Sprintf("%v", v.GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="category_sort_inlist" value="` + strconv.FormatInt(v.(*entity.Category).GetSort(), 10) + `"/></td>
				<td>
					` + strings.Repeat(`<button type="button" class="btn btn-default btn-sm" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-right" aria-hidden="true"></span></button>`, v.(*entity.Category).GetLevel()) + `
					`
		if v.(*entity.Category).Descendants != nil {
			out += `<button data-item_id="` + idStr + `" data-item_parent="` + strconv.FormatInt(v.GetParent(), 10) + `" data-item_level="` + strconv.Itoa(v.(*entity.Category).GetLevel()) + `" type="button" class="btn btn-success btn-sm button_collapse" aria-label="Left Align"><span class="glyphicon glyphicon-arrow-down" aria-hidden="true"></span></button>`
		}
		out += `
					<input type="text" data-item_id="` + idStr + `" class="category_title_inlist" value="` + v.(*entity.Category).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="category_description_inlist" style="width:100%" >` + v.(*entity.Category).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "category", v.(*entity.Category).GetEnable()) + `</td>
			</tr>
			`
		//}
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
}


func (s *Category)ListingNonTree()string{
	if s.List[0].GetId() == -1 {
		return `
			<tr>
				<td colspan="8">
					<h2>` + lang.T("no items") + `</h2>
				</td>
			</tr>`
	}
	var out, idStr,imgSrc string
	for i := 0; s.List[i].GetId() != -1; i++ {
		//if s.SearchInList == "" || strings.Contains(v.(*entity.Category).GetTitle()+v.(*entity.Category).GetDescription(),s.SearchInList) {
		idStr = strconv.FormatInt(s.List[i].GetId(), 10)
		if s.List[i].(*entity.Product).GetImg() != nil {
			imgSrc = app.Url_assets_uploads + s.List[i].(*entity.Product).GetImg()[0]
		} else {
			imgSrc = app.Url_no_image
		}
		out += `
			<tr>
				<th scope="row">` + idStr + `</th>
				<td><img src="` + imgSrc + `" width='60' height='50' /></td>
				<td>` + fmt.Sprintf("%v", s.List[i].(*entity.Category).GetParent()) + `</td>
				<td><input type="text"  data-item_id="` + idStr + `" class="category_sort_inlist" value="` + strconv.FormatInt(s.List[i].(*entity.Category).GetSort(), 10) + `"/></td>
				<td>
					<input type="text" data-item_id="` + idStr + `" class="category_title_inlist" value="` + s.List[i].(*entity.Category).GetTitle() + `"/>
					<br/><textarea data-item_id="` + idStr + `" class="category_description_inlist" style="width:100%" >` + s.List[i].(*entity.Category).GetDescription() + `</textarea>
				</td>
				<td >` + HomeLayoutView.Table_colomun_last(idStr, "category", s.List[i].(*entity.Category).GetEnable()) + `</td>
			</tr>
			`
		//}
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
}






















*/
