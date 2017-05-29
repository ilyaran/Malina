package views

import (
	"net/http"
	"github.com/ilyaran/Malina/entity"
	"fmt"
	"time"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/libraries"
)

var LOCALS = &Locals{}
type Locals struct {
	W                   http.ResponseWriter
	R 		    *http.Request

	CurrentPage 	    string
	AccountCartList     []*entity.CartPublic
	AccountCartListJSON []byte
}

type Layout struct {
	Header1 []byte
	Nav1     []byte
	Nav2     []byte
	Nav3     []byte
	Body1    []byte
	Body2    []byte
	Body3    []byte
	Footer1  []byte
}

func (s *Layout) WriteAuthResponse() {
	LOCALS.W.Write(s.Header1)
	LOCALS.W.Write(s.Nav1)
	s.SetAuthNav2()
	LOCALS.W.Write(s.Nav2)
	LOCALS.W.Write(s.Nav3)
	LOCALS.W.Write(s.Body2)
	LOCALS.W.Write(s.Footer1)
}

func NewAuthLayout() *Layout {
	s := &Layout{}
	s.SetStaticParts()
	return s
}

func NewLayout() *Layout {
	s := &Layout{}
	s.SetStaticParts()
	library.SESSION.GetSessionObj()
	return s
}

func (s *Layout) WriteResponse() {

	LOCALS.W.Write(s.Header1)
	LOCALS.W.Write(s.Nav1)

	s.SetNav2()
	LOCALS.W.Write(s.Nav2)

	LOCALS.W.Write(s.Nav3)
	LOCALS.W.Write(s.Body1)
	LOCALS.W.Write(s.Body2)
	LOCALS.W.Write(s.Body3)
	LOCALS.W.Write(s.Footer1)

}


func(s *Layout) SetAuthNav2(){
	var out = `
<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
	<ul class="nav navbar-nav">
	<li>
		<a href="`+app.Uri_auth_login()+`">Sign In</a>
	</li>
	<li>
		<a href="`+app.Uri_auth_register()+`">Sign Up</a>
	</li>
	<li>
		<a href="`+app.Uri_auth_forgot()+`">Forgot</a>
	</li>
	</ul>

	<ul class="nav navbar-nav navbar-right">
	`

	s.Nav2 = []byte(out)
}

func(s *Layout) SetNav2(){
	var out = `
<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
	<ul class="nav navbar-nav">
	<li `; if LOCALS.CurrentPage == "cart" {out+=`class="active"`}; out+= `>
		<a href="`+app.Base_url()+app.Uri_cart()+`list">Carts</a>
	</li>
	<li `; if LOCALS.CurrentPage == "product" {out+=`class="active"`}; out+= `>
		<a href="`+app.Base_url()+app.Uri_product()+`list">Products</a>
	</li>
	<li `; if LOCALS.CurrentPage == "category" {out+=`class="active"`}; out+= `>
		<a href="`+app.Base_url()+app.Uri_category()+`list">Categories</a>
	</li>
	<li class="dropdown">
		<a href="#" class="dropdown-toggle" data-toggle="dropdown">Auth<strong class="caret"></strong></a>
		<ul class="dropdown-menu">
			<li `; if LOCALS.CurrentPage == "account" {out+=`class="active"`}; out+= `>
				<a href="`+app.Base_url()+app.Uri_account()+`list">Accounts</a>
			</li>
			<li `; if LOCALS.CurrentPage == "position" {out+=`class="active"`}; out+= `>
				<a href="`+app.Base_url()+app.Uri_position()+`list">Positions</a>
			</li>
			<li `; if LOCALS.CurrentPage == "permission" {out+=`class="active"`}; out+= `>
				<a href="`+app.Base_url()+app.Uri_permission()+`list">Permissions</a>
			</li>
			<li class="divider"></li>
			<li><a href="#">Separated link</a></li>
			<li class="divider"></li>
			<li><a href="#">One more separated link</a></li>
		</ul>
	</li>
	</ul>

	<ul class="nav navbar-nav navbar-right">
	`; if library.SESSION.GetSessionObj().AccountId > 0 {out+=`
		<li><a href="">` + library.SESSION.GetSessionObj().GetEmail() + `</a></li>
		<li><a href="` + app.Base_url() + `auth/logout/">` + lang.T("logout") + `</a></li>
	`}else {
		out += `<li><a href="` + app.Base_url() + `auth/login/">` + lang.T("login") + `</a></li>`
	}

	s.Nav2 = []byte(out)
}
func(s *Layout) SetStaticParts(){
	s.Header1 = []byte(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">

		<title>Title</title>
		<base href="`+ app.Base_url() +`">
		<meta name="description" content="Source code generated using layoutit.com">
		<meta name="author" content="LayoutIt!">
		<link rel="apple-touch-icon" sizes="76x76" href="assets/apple-icon-76x76.png">
		<link href="`+ app.Assets_backend_path() +`css/bootstrap.min.css" rel="stylesheet">
		<link href="`+ app.Assets_backend_path() +`css/style.css" rel="stylesheet">

		<link href='http://fonts.googleapis.com/css?family=Open+Sans' rel='stylesheet' type='text/css' />
		<link href="`+ app.Assets_backend_path() +`css/bootstrap-toggle.css" rel="stylesheet">
  	</head>

	<body>

		<!-- Alert Dialog -->
		<div id="alert_dialog" title="Alert message" style="display: none">
		    <div class="ui-dialog-content ui-widget-content">
			<p>
			    <span class="ui-icon ui-icon-alert" style="float: left; margin: 0 7px 20px 0"></span>
			    <label id="alert_message">Alert!</label>
			</p>
		    </div>
		</div>
		<!-- /Alert Dialog -->

		<div id="loader" class="loader" style="background-color: black;position: absolute;top:30%; left: 50%; ;z-index: 1500;display: none;">
		    <img src="`+app.Assets_path()+`img/ajax_loader_blue_256.gif" />
		</div>

		<div class="container-fluid">
			<div class="row">
				<div class="col-md-12">


	`)

	//**************************************************************************************
	//**************************************************************************************
	//**************************************************************************************

	s.Nav1 = []byte(`
			<nav class="navbar navbar-default navbar-inverse navbar-fixed-top" role="navigation">
				<div class="navbar-header">
					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
						<span class="sr-only">Toggle navigation</span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="` + app.Base_url() + `">` + app.Site_name() + `</a>
				</div>
		`)

	s.Nav3 = []byte(`
						<li class="dropdown">
							<a href="#" class="dropdown-toggle" data-toggle="dropdown">Dropdown<strong class="caret"></strong></a>
							<ul class="dropdown-menu">
								<li><a href="#">Action</a></li>
								<li><a href="#">Another action</a></li>
								<li><a href="#">Something else here</a></li>
								<li class="divider"></li>
								<li><a href="#">Separated link</a></li>
							</ul>
						</li>
					</ul>
				</div>
			</nav>
	`)

	//**************************************************************************************
	//**************************************************************************************
	//**************************************************************************************

	s.Footer1 = []byte(`

				</div>
			</div>
		</div>
	    	<script src="`+ app.Assets_backend_path() +`js/jquery.min.js"></script>
	<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
  <script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
	    	<script src="`+ app.Assets_backend_path() +`js/bootstrap.min.js"></script>
	    	<script src="`+ app.Assets_backend_path() +`js/bootstrap-toggle.js"></script>
	    	<script src="`+app.Assets_path()+`js/home.js?`+fmt.Sprintf("%v",time.Now())+`"></script>
	</body>
</html>`)
}

