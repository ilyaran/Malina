package views

import (
	"Malina/config"
	"Malina/language"
	"Malina/libraries"
)

func Header()string  {

	return `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">

		<title>Admin</title>
		<base href="`+ app.Base_url() +`">
		<meta name="description" content="Source code generated using layoutit.com">
		<meta name="author" content="LayoutIt!">

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


	`
}

func Nav(active string)string  {
	var out = `
	<nav class="navbar navbar-default navbar-inverse navbar-fixed-top" role="navigation">
		<div class="navbar-header">
			<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
				<span class="sr-only">Toggle navigation</span>
				<span class="icon-bar"></span>
				<span class="icon-bar"></span>
				<span class="icon-bar"></span>
			</button>
			<a class="navbar-brand" href="#">` + app.Site_name() + `</a>
		</div>
		<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
			<ul class="nav navbar-nav">
				<li `; if active == "cart" {out+=`class="active"`}; out+= `>
					<a href="`+app.Base_url()+app.Uri_cart()+`list">Carts</a>
				</li>
				<li `; if active == "product" {out+=`class="active"`}; out+= `>
					<a href="`+app.Base_url()+app.Uri_product()+`list">Products</a>
				</li>
				<li `; if active == "category" {out+=`class="active"`}; out+= `>
					<a href="`+app.Base_url()+app.Uri_category()+`list">Categories</a>
				</li>
				<li class="dropdown">
					<a href="#" class="dropdown-toggle" data-toggle="dropdown">Auth<strong class="caret"></strong></a>
					<ul class="dropdown-menu">
						<li `; if active == "account" {out+=`class="active"`}; out+= `>
							<a href="`+app.Base_url()+app.Uri_account()+`list">Accounts</a>
						</li>
						<li `; if active == "position" {out+=`class="active"`}; out+= `>
							<a href="`+app.Base_url()+app.Uri_position()+`list">Positions</a>
						</li>
						<li `; if active == "permission" {out+=`class="active"`}; out+= `>
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
				`; if library.SESSION.GetSessionObj().GetAccount_id() > 0 {out+=`
				<li><a href="">` + library.SESSION.GetSessionObj().GetEmail() + `&nbsp;`+library.SESSION.GetSessionObj().GetNick() + `</a></li>
				<li><a href="` + app.Base_url() + `auth/logout/">` + lang.T("logout") + `</a></li>
				`}else {
				out += `<li><a href="` + app.Base_url() + `auth/login/">` + lang.T("login") + `</a></li>`
				}; out+=`
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
	</nav>`
	return out
}

