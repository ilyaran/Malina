package authView

import (
	"net/http"
	"Malina/language"
	"Malina/config"
	"Malina/libraries"
)

func ChangePassword(r *http.Request)string{

	var out string = ""

	out += `

<input type="hidden" value="login" id="actionType"/>

         <div class="container">
        <div class="row" style="margin-top:auto;">
            <div class="col-md-4 col-md-offset-4">
                <div class="login-panel panel panel-default">
                    <div class="panel-heading">
                        <h3 class="panel-title">`+lang.T("reset_password")+`</h3>
                    </div>
                    <div class="panel-body">

		`; if v,ok := library.VALIDATION.Result["success"]; ok {
		out += `<h1>`+v+`</h1>`;
	}else{

		out += `

                    	<form method="post" action="`+app.Base_url()+app.Uri_auth_change_password()+`">

				<span class="error">`; if v,ok:= library.VALIDATION.Result["no_rows"]; ok {out += v }; out += `<span>
				<span class="error">`; if v,ok:= library.VALIDATION.Result["error"]; ok {out += v }; out += `</span>

                                <div class="form-group">
                                    <input class="form-control" value="`+r.FormValue("old_password")+`" placeholder="`+lang.T("auth_pass_old_title")+`" name="old_password" type="text" autofocus>
                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["old_password"]; ok {out += v }; out += `<span>
                                </div>

                                <div class="form-group">
                                    <input class="form-control" value="`+r.FormValue("new_password")+`" placeholder="`+lang.T("auth_pass_new_title")+`" name="new_password" type="text" autofocus>
                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["new_password"]; ok {out += v }; out += `<span>
                                </div>

                                <div class="form-group">
                                    <input class="form-control" value="`+r.FormValue("confirm_new_password")+`" placeholder="`+lang.T("auth_pass_new_confirm_title")+`" name="confirm_new_password" type="text" autofocus>
                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["confirm_new_password"]; ok {out += v }; out += `<span>
                                </div>

                               `; if app.Use_recaptcha_change_password() { out += `
                                    	<script src="https://www.google.com/recaptcha/api.js"></script>
					<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {out += v }; out += `<span>
					<br>
				`}; out += `

                                <!-- Change this to a button or input when using this as a form -->
                                <button type="submit" class="btn btn-lg btn-success btn-block">`+lang.T("send")+`</button>

			</form>

			`}; out += `

                        <br>
                        <a href="`+app.Base_url()+app.Uri_auth_login() + `" class="btn btn-lg btn-primary">`+lang.T("login")+`</a>
                        <a href="`+app.Base_url()+app.Uri_auth_register() +`" class="btn btn-lg btn-default">`+lang.T("register")+`</a>
                    </div>
                </div>
            </div>
        </div>
    </div>



	`

	return out }

