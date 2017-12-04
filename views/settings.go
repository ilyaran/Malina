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
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/lang"
	"net/http"
	"html/template"
)

type Settings struct {
	Out template.HTML
}

func (s *Settings)Index(w http.ResponseWriter){

	s.Out = template.HTML(`
	<div class="panel panel-default">
		<div class="panel-heading">
			<h3 class="panel-title">`+lang.T("Settings")+`</h3>
		</div>
		<div class="panel-body">

		</div>
	</div>

	<div class="bs-example bs-example-tabs" data-example-id="togglable-tabs">
		<ul class="nav nav-tabs" id="myTabs" role="tablist">
			<li role="presentation" class="active">
				<a href="#items_table" id="home-tab" role="tab" data-toggle="tab" aria-controls="home" aria-expanded="true">`+lang.T("list")+`</a>
			</li>

		</ul>
		<div class="tab-content" id="myTabContent">
			<div class="tab-pane fade active in" role="tabpanel" id="items_table" aria-labelledby="home-tab">
				<div class="row">
					<div class="col-md-6">
						<div class="table-responsive">

						<table class="table table-bordered">
							<thead>
								<tr>
									<th>`+lang.T("Name")+`</th>
									<th>`+lang.T("Value")+`</th>
									<th>`+lang.T("Action")+`</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<th scope="row">Base_url</th>
									<td><input type="text" value="`+app.Base_url+`"></td>
									<td>
										<button class="btn btn-primary btn-sm">Send</button>
										<button class="btn btn-success btn-sm">Reset</button>
									</td>
								</tr>
							</tbody>
						</table>

						</div>
					</div>
					<div class="col-md-6">
						<div class="table-responsive">

						<table class="table table-bordered">
							<thead>
								<tr>
									<th>`+lang.T("Name")+`</th>
									<th>`+lang.T("Value")+`</th>
									<th>`+lang.T("Action")+`</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<th scope="row">Base_url</th>
									<td><input type="text" value="`+app.Base_url+`"></td>
									<td>
										<button class="btn btn-primary btn-sm">Send</button>
										<button class="btn btn-success btn-sm">Reset</button>
									</td>
								</tr>
							</tbody>
						</table>

						</div>
					</div>

				</div>
			</div>

		</div>
	</div>`)


	t := template.New("Settings_index")
	t, _ = t.Parse(HomeLayoutView.Layout)

	t.Execute(w, s)
}
