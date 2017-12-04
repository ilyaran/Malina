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
	"github.com/ilyaran/Malina/libraries"
	"strconv"
	"github.com/ilyaran/Malina/entity"
	"strings"
)


func HomeLayoutViewInit(){
	HomeLayoutView =&HomeLayout{}
	HomeLayoutView.SetLayout()


}
var HomeLayoutView *HomeLayout
type HomeLayout struct {

	Layout string
}


func (s *HomeLayout)SetLayout(){
	s.Layout = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Admin</title>
	<base href="`+ app.Base_url +`"/>
    <meta name="description" content="">

    <link href="`+app.Url_assets_home+`css/bootstrap.min.css" rel="stylesheet">
    <link href="`+app.Url_assets_home+`css/style.css" rel="stylesheet">
	<link href="`+app.Url_assets_home+`css/bootstrap-toggle.css" rel="stylesheet">

	<script src="http://code.jquery.com/jquery-1.11.1.min.js"></script>
	<script src="http://code.jquery.com/ui/1.11.1/jquery-ui.min.js"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.11.1/themes/smoothness/jquery-ui.css" />

  </head>
  <body>
	<div id="loader" class="loader" style=" opacity:0.6;background-color: rgba(255,245,215,1);position: fixed;width:100%; height: 100%;z-index: 1500;display:none;">
		    <img src="`+app.Url_assets+`img/ajax_loader_blue_256.gif" style="position: fixed;top:40%; left: 40%;" />
	</div>
    <div class="container-fluid">
	<div class="row">
		<div class="col-md-12">

			<!-- nav -->
			<nav class="navbar navbar-default navbar-inverse navbar-fixed-top" role="navigation">
				<div class="navbar-header">

					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
						 <span class="sr-only">Toggle navigation</span><span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span>
					</button> <a class="navbar-brand" href="`+app.Base_url+`">`+app.Site_name+`</a>
				</div>

				<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
					<ul class="nav navbar-nav navbar-right" id="nav_bar">
						<li class="active"><a href="`+app.Base_url+`home/index">`+lang.T("Dashboard")+`</a></li>
						<li><a href="`+app.Base_url+`home/order">`+lang.T("orders")+`</a></li>
						<li><a href="`+app.Base_url+`home/cart">`+lang.T("carts")+`</a></li>
						<li><a href="`+app.Base_url+`home/product">`+lang.T("products")+`</a></li>
						<li><a href="`+app.Base_url+`home/category">`+lang.T("categories")+`</a></li>
						<li><a href="`+app.Base_url+`home/parameter">`+lang.T("parameters")+`</a></li>
						<li><a href="`+app.Base_url+`home/news">`+lang.T("news")+`</a></li>
						<li><a href="`+app.Base_url+`home/page">`+lang.T("pages")+`</a></li>
						<li><a href="`+app.Base_url+`home/settings">`+lang.T("settings")+`</a></li>
						<li class="dropdown">
							<a href="#" class="dropdown-toggle" data-toggle="dropdown">`+lang.T("Users")+`<strong class="caret"></strong></a>
							<ul class="dropdown-menu">
								<li><a href="`+app.Base_url+`home/account">`+lang.T("Accounts")+`</a></li>
								<li><a href="`+app.Base_url+`home/account/ban">`+lang.T("Banned")+`</a></li>
								<li><a href="`+app.Base_url+`home/role">`+lang.T("Roles")+`</a></li>
								<li class="divider"></li>
								<li><a href="`+app.Base_url+`home/permission">`+lang.T("Permissions")+`</a></li>
								<li class="divider"></li>
								<li><a href="`+app.Base_url+`home/activation">`+lang.T("Activations")+`</a></li>
							</ul>
						</li>
						<li><a >   </a></li>
					</ul>
				</div>
			</nav>
			<!-- end nav -->


			<!-- content -->

			<div class="row">
				<div class="col-md-12">

					{{ .Out }}

				</div>
			</div>

			<!-- end content -->

		</div>
	</div>
</div>


<!-- Modal block -->

<div id="modal_dialog_block">

	<div  id="modal_dialog" class="modal fade" tabindex="-1" role="dialog">
	  <div class="modal-dialog" role="document">
		<div class="modal-content">
		  <div class="modal-header">
			<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
			<h4 id="modal_title" class="modal-title"></h4>
		  </div>
		  <div class="modal-body">
			<p id="modal_body"></p>
		  </div>
		  <div class="modal-footer">
			<button id="modal_button1" type="button" class="btn btn-default" data-dismiss="modal"></button>
			<button id="modal_button2" type="button" class="btn btn-primary"></button>
		  </div>
		</div><!-- /.modal-content -->
	  </div><!-- /.modal-dialog -->
	</div><!-- /.modal -->

</div>

<!-- end Modal block -->


	<script src="`+app.Url_assets_home+`js/bootstrap.min.js"></script>
    <script src="`+ app.Url_assets_home +`js/bootstrap-toggle.js"></script>

  </body>
</html>`

}
func (s *HomeLayout)Form_sort()string{
	return `
	<tr>
        <td>
            <div class="form-group">
                <label>`+lang.T("sort")+`<span id="sort_error"  class="error"><span></label>
            </div>
        </td>
        <td>
            <div class="form-group">
                <input id="sort"  value="100" type="number" class="form-control" />
            </div>
        </td>
    </tr>
	`
}
func (s *HomeLayout) Form_title()string{
	return `
	<tr>
		<td width="30%">
			`+lang.T("title")+`<span class="error required">*</span>
         	<span id="title_error"  class="error"><span>
		</td>
		<td><input id="title" value=""></td>
	</tr>`
}
func (s *HomeLayout) Form_description()string{
	return `
	<tr>
		<script type="text/javascript" src="`+app.Url_assets+`ckeditor/ckeditor.js"></script>
		<td>`+lang.T("description")+`<span id="description_error"  class="error"><span></td>
		<td><textarea id="description" style="width:100%" rows="6"></textarea></td>
		<script type="text/javascript">
			var ckeditor = CKEDITOR.replace('description',{
				filebrowserBrowseUrl: '`+app.Url_assets+`filemanager/index.html',
				filebrowserImageBrowseUrl:'`+app.Url_assets+`filemanager/index.html?type=image',
				removeDialogTabs: 'link:upload;image:upload'
			});
		</script>
	</tr>`
}
func (s *HomeLayout) Form_description_simple()string{
	return `
	<tr>
		<td>`+lang.T("description")+`<span id="description_error"  class="error"><span></td>
		<td><textarea id="description" style="width:100%" rows="5"></textarea></td>
	</tr>`
}
func (s *HomeLayout) Form_short_description()string{
	return `
	<tr>
		<td>`+lang.T("short description")+` length:`+strconv.Itoa(app.Short_description_len)+` symbols <span id="short_description_error"  class="error"><span></td>
		<td><textarea id="short_description" style="width:100%" rows="5"></textarea></td>
	</tr>`
}

func (s *HomeLayout) Form_checkbox(title,keyName string,isChecked bool)string{
	o := `
	<tr>
		<td>
			<div class="form-group">
				<label>`+lang.T(title)+`</label>
			</div>
		</td>
        <td>
            <div class="form-group">
                <input id="`+keyName+`" type="checkbox" `; if isChecked {o += `checked`}; o += ` class="form-control" />
            </div>
        </td>
    </tr>`
	return o
}

func (s *HomeLayout) Form_send_button()string{
	return `<button class="btn btn-primary form_send">`+lang.T("Send")+`</button>`
}

func (s *HomeLayout) Table_head_last(dbtable string)string{
	return `
	<th style="width:20%">
		Enable<input title="`+lang.T("Check to be enable all items")+`" checked type="checkbox" onchange="if(this.checked){$('.`+dbtable+`_enable_inlist').bootstrapToggle('on');}else{$('.`+dbtable+`_enable_inlist').bootstrapToggle('off');}"/>
		Delete<input title="`+lang.T("Check to delete all items")+`" type="checkbox"  onchange="if(this.checked){$('.`+dbtable+`_del_item_inlist').prop('checked',true);}else{$('.`+dbtable+`_del_item_inlist').prop('checked',false);}"/>
		<button title="`+lang.T("Click to send update in list")+`" id="send_inlist_`+dbtable+`" class="btn btn-success btn-sm"><span class="glyphicon glyphicon-send" aria-hidden="true"></span></button>
	</th>`
}

func (s *HomeLayout) Table_colomun_last(idStr, dbtable string,enable bool)string{
	out := `<input title="`+lang.T("Switch Enable")+`" value="" data-item_id = "` + idStr + `" class="`+dbtable+`_enable_inlist" type="checkbox" `
	if enable {
		out += `checked`
	}
	out += ` data-toggle="toggle" />
	<button title="`+lang.T("Click to edit item")+`" data-item_id="` + idStr + `" class="btn btn-primary edit_item_`+dbtable+`"><span class="glyphicon glyphicon-edit" aria-hidden="true"></span></button>
	<button title="`+lang.T("Click to delete item")+`" data-item_id="` + idStr + `" class="btn btn-danger del_item_`+dbtable+`"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
	<input title="`+lang.T("Check to delete item")+`" value="0" data-item_id="` + idStr + `" class="`+dbtable+`_del_item_inlist" type="checkbox">`
	return out
}

func (s *HomeLayout) Form_parameter_select(rootTile string,keyName string)string{
	var out string
	out += `
		&nbsp;&nbsp;
		<span>`+lang.T(keyName)+`</span>
		<select id="`+keyName+`" onchange="$('#`+keyName+`_title').html($('#`+keyName+` option:selected').text());">
		`
			if rootTile != ``{
				out += `<option value="0">`+rootTile+`</option>`
			}
			for _,v:=range libraries.ParameterLib.Trees.List {
				out +=`<option value="`+strconv.FormatInt(v.(*entity.Parameter).GetId(),10)+`">` + strings.Repeat("&rarr;", v.(*entity.Parameter).GetLevel())+v.(*entity.Parameter).GetTitle()+`</option>`
			}
		out += `
		</select>
		&nbsp;&nbsp;
		<span id="`+keyName+`_error" class="error"></span>
		<span class="btn btn-info" id="`+keyName+`_title"></span>`
	return out
}

func (s *HomeLayout) Form_category_select(rootTile string,keyName string)string{
	var out string
	out += `
		&nbsp;&nbsp;
		<span>`+lang.T(keyName)+`</span>
		<select id="`+keyName+`" onchange="$('#`+keyName+`_title').html($('#`+keyName+` option:selected').text());">
		`
	if rootTile != ``{
		out += `<option value="0">`+rootTile+`</option>`
	}
	for _,v:=range libraries.CategoryLib.Trees.List {
		out +=`<option value="`+strconv.FormatInt(v.(*entity.Category).GetId(),10)+`">` + strings.Repeat("&rarr;", v.(*entity.Category).GetLevel())+v.(*entity.Category).GetTitle()+`</option>`
	}
	out += `
		</select>
		&nbsp;&nbsp;
		<span id="`+keyName+`_error" class="error"></span>
		<span class="btn btn-info" id="`+keyName+`_title"></span>`
	return out
}


func (s *HomeLayout) Per_page_select()string{
	var out string
	var step = (app.Per_page_max -1-app.Per_page)/app.Per_page_select_options_widget_num
	for i:=int64(app.Per_page); i<(app.Per_page_max -step); i+= step {
		out += `<option value="`+strconv.FormatInt(i,10)+`">`+strconv.FormatInt(i,10)+`</option>`
	}
	out += `<option value="`+strconv.FormatInt(app.Per_page_max-1,10)+`">`+strconv.FormatInt(app.Per_page_max-1,10)+`</option>`

	return lang.T("per page") + `<select name="per_page" id="per_page" >`+out+`</select>`
}

func (s *HomeLayout)ActionBar()string{
	return `
	<div class="form-group" style="float:left;">
		`+lang.T("action")+`
		<select id="action" onchange='
			if ($(this).val() == "add") {
				$("#item_id_bar").hide();
					$(".required").html("*");
				} else{
					$("#item_id_bar").show();
					$(".required").html("");
			}'>
			<option value="add">`+lang.T("add")+`</option>
			<option value="get">`+lang.T("get")+`</option>
			<option value="edit">`+lang.T("edit")+`</option>
			<option value="del">`+lang.T("delete")+`</option>
		</select>
		<span style = "display:none;" id = "item_id_bar">
			<span>`+lang.T("id")+`:</span>
			<input id = "item_id" value="" type="integer"/>
		</span>
    </div>`
}

func (s *HomeLayout)ImageBar()string{
	return `
	<div id="roxyCustomPanel2" style="display: none;">
		<iframe src="`+app.Url_assets+`filemanager/index.html?integration=custom&type=files&txtFieldId=image_preview" style="width:100%;height:100%" frameborder="0">
		</iframe>
	</div>
	<button id="select_image" type="button" class="btn btn-primary"><span class="glyphicon glyphicon-picture" aria-hidden="true"></span></button>
    <button  id="clean_image_preview" type="button" class="btn btn-danger"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
	<div id="image_preview" style="background-image: url(`+app.Url_assets+`/img/noimg.jpg);background-repeat: no-repeat;min-height:180px;"></div>
    <span id="error_image_preview" class="error"></span>`
}






























