/**
 * Common views functions.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package views

import (
	"Malina/language"
	"Malina/config"
)


var FACE_FORM *FaceForm = &FaceForm{}
var TABLE_FORM *TableForm= &TableForm{}

type TableForm struct {
	Inputs string
	Breadcrumb string
	Select_options_order_by string
	Head string
	Listing string
	Form string
	Out                     string
	Option1 string
	Option2 string
}
func(s *TableForm) SetNull(){
	s.Inputs =``
	s.Breadcrumb =``
	s.Select_options_order_by  =``
	s.Head  =``
	s.Listing =``
	s.Form  =``
	s.Out= ``
	s.Option1 = ``
	s.Option2 = ``
}
func(s *TableForm) IndexFrom(){
	s.Out = s.Inputs + `
	<div class="row">
	<div class="col-md-12">
		<div class="panel panel-default">
			<div class="panel-heading">
				<h3 class="panel-title">`+s.Breadcrumb +`</h3>
			</div>
		  	<div class="panel-body">
		    		`+FACE_FORM.PerPageSelectForm() + `
				&nbsp;&nbsp;
				<span>`+lang.T("order_by")+`</span>
				<select id="order_by" >
					`+s.Select_options_order_by +`
				</select>
				`+s.Option1+`
				`+s.Option2+`
				<input type="text" class="form-control" id="search" name="search"  placeholder="Search for...">

		    		<button id="get_list" class="btn btn-primary btn-lg btn-block" type="button" >`+lang.T("go")+`</button>
				<button id = "success_bar"  style="display:none;" type="button" class="btn btn-success btn-lg btn-block">`+lang.T("success")+`</button>
		  	</div>
		</div>
		<ul class="nav nav-tabs">
                	<li class="active"><a href="#home" data-toggle="tab">`+lang.T("list")+`</a></li>
                        <li class=""><a href="#form" data-toggle="tab">` + lang.T("form") + `</a></li>

		</ul>
		<div class="tab-content">
			<div class="tab-pane fade active in" id="home">
				<div class="error" id="inlist_error"></div>
				<table class="table table-striped table-bordered table-hover" id="dataTables-example">
					<thead>
						<tr>
						`+s.Head+`
						</tr>
					</thead>
					<tbody id="listing">
						`+s.Listing+`
					</tbody>
				</table>
			</div>
			<div class="tab-pane fade" id="form">
				`+s.Form+`
			</div>
		</div>
	</div>
</div>
	`
}





type FaceForm struct {}

func (s *FaceForm)Inputs(dbtable string)string{
	return `
	<input type="hidden" value="`+dbtable+`" id="dbtable"/>
	<input type="hidden" value="`+app.No_image()+`" id="url_no_image"/>
	<input type="hidden" value="`+app.Base_url()+app.Upload_path()+`" id="url_upload_path"/>
	<input type="hidden" value="`+app.Base_url()+`home/`+dbtable+`/list_ajax/" id="url_`+dbtable+`_list_ajax"/>
	<input type="hidden" value="`+app.Base_url()+`home/`+dbtable+`/get/" id="url_`+dbtable+`_get"/>
	<input type="hidden" value="`+app.Base_url()+`home/`+dbtable+`/add/" id="url_`+dbtable+`_add"/>
	<input type="hidden" value="`+app.Base_url()+`home/`+dbtable+`/edit/" id="url_`+dbtable+`_edit"/>
	<input type="hidden" value="`+app.Base_url()+`home/`+dbtable+`/del/" id="url_`+dbtable+`_del"/>

	`
}

func (s *FaceForm)BarForms()string{
	return `
	<div class="form-group" style="float:left;">
		<select id="action">
			<option value="add">`+lang.T("add")+`</option>
			<option value="get">`+lang.T("get")+`</option>
			<option value="edit">`+lang.T("edit")+`</option>
			<option value="delete">`+lang.T("delete")+`</option>
		</select>
		<span style = "display:none;" id = "item_id_bar">
			<span>`+lang.T("id")+`:</span>
			<input id = "item_id" value="" type="integer"/>
			<span id = "item_id_error" class="error"></span>
		</span>
        </div>
	<span style="float:left;" id="error" class="error"></span>`
}
func (s *FaceForm)CategorySelect(rootTile string,keyName string)string{
	var out string
	out += `&nbsp;&nbsp;
		<span>`+lang.T(keyName)+`</span>
		<select id="`+keyName+`" onchange="$('#`+keyName+`_title').html($('#`+keyName+` option:selected').text());">
		`; if rootTile != ``{
		out += `<option value="0">`+rootTile+`</option>`
	}; out += /*library.CATEGORY.SelectOptionsList +*/ `
		</select>
		&nbsp;&nbsp;
		<span id="`+keyName+`_error" class="error"></span>
		<span class="btn btn-info" id="`+keyName+`_title">Category Title</span>`
	return out
}
func (s *FaceForm)PerPageSelectForm()string{
	return lang.T("per page") + `
	<select name="per_page" id="per_page" >
		<option value="100">100</option>
		<option value="150">150</option>
		<option value="200">200</option>
		<option value="300">300</option>
		<option value="400">400</option>
		<option value="500">500</option>
	</select>`
}

func (s *FaceForm)ImageBar()string{
	return `
	<div id="roxyCustomPanel2" style="display: none;">
		<iframe src="/assets/filemanager/index.html?integration=custom&type=files&txtFieldId=image_preview" style="width:100%;height:100%" frameborder="0">
		</iframe>
	</div>

	<button id="select_image" type="button" class="btn btn-primary"><span class="glyphicon glyphicon-picture" aria-hidden="true"></span></button>
        <button  id="clean_images" type="button" class="btn btn-danger"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
	<div id="image_preview" style="background-image: url(assets/img/noimg.jpg);background-repeat: no-repeat;min-height:180px;">
        </div>
        <span id="img_error" class="error"></span>`
}

func (s *FaceForm)Description()string{
	return `
	<tr>
		<td>
			<label>`+lang.T("description")+`
			<span id="description_error"  class="error"><span></label>
		</td>
                <td>
			<div class="form-group">
				<textarea id="description" name="description" class="form-control" rows="3"></textarea>
			</div>
			<script type="text/javascript">
				var ckeditor = CKEDITOR.replace('description',{
					filebrowserBrowseUrl: '/assets/filemanager/index.html',
					filebrowserImageBrowseUrl:'/assets/filemanager/index.html?type=image',
					removeDialogTabs: 'link:upload;image:upload'
				});
			</script>
		</td>
        </tr>
	`
}
/*
<script type="text/javascript">
				var roxyFileman = 'http://roxyfilemanager.loc/assets/index.html';
				$(function(){
					CKEDITOR.replace( 'description',{
					filebrowserBrowseUrl:roxyFileman,
					filebrowserImageBrowseUrl:roxyFileman+'?type=image',
					removeDialogTabs: 'link:upload;image:upload'});
				});
			</script>
 */
func (s *FaceForm)Title()string{
	return `
	<tr>
        	<td>
        		<label>`+lang.T("title")+`<span class="error">*</span>
                        	<span id="title_error"  class="error"><span>
                        </label>
		</td>
                <td>
                        <div class="form-group">
                        	<input id="title"  class="form-control" />
                        </div>
                </td>
        </tr>`
}

func (s *FaceForm)FormField(title, keyName string, isRequired bool)string{
	o := `
	<tr>
        	<td>
        		<label>`+lang.T(keyName); if isRequired {o+=`<span class="error">*</span>`}; o+= `
                        	<span id="`+keyName+`_error"  class="error"><span>
                        </label>
		</td>
                <td>
                        <div class="form-group">
                        	<input id="`+keyName+`"  class="form-control" />
                        </div>
                </td>
        </tr>`
	return o
}

func (s *FaceForm)CheckBox(title,keyName string,isChecked bool)string{
	o := `
	<tr>
		<td>
                	<div class="form-group">
                        	<label>`+lang.T(title)+`</label>
			</div>
                </td>
                <td>
                        <div class="form-group">
                        	<input id="`+keyName+`" type="checkbox" `;if isChecked {
		o += `checked`}; o += ` class="form-control" />
                        </div>
                </td>
        </tr>`
	return o
}