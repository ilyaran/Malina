package authView

import (
	"net/http"
	"Malina/language"
	"Malina/config"
	"Malina/libraries"
)

func RegisterForm(r *http.Request)string{
	var out string = ""

	out += `
<div class="container">
	<div class="row" style="margin-top:auto;">
        	<div class="col-md-4 col-md-offset-4">
        		<div class="login-panel panel panel-default">
                    		<div class="panel-heading">
                        		<h3 class="panel-title">`+lang.T("register")+`</h3>
                    		</div>
                    		<div class="panel-body">
				<h3>`+lang.T("sign up")+`</h3>
				`; if v,ok:= library.VALIDATION.Result["success"]; ok {
				if library.VALIDATION.Status == 50 {
					out += `<h1>` + v + `</h1>
					<h6><a href="` + app.Base_url() + `auth/logout/" >` + lang.T("logout") + `</a></h6>`;
				}else if library.VALIDATION.Status == 0 {
					out += `<h1>` + v + `</h1>
					<h6><a href="` + app.Base_url() + `auth/login/" >` + lang.T("login") + `</a></h6>`;
				}

				}else{out += `
				<form  method="post">
				<div class="form-group">
					<span class="error">`; if v,ok := library.VALIDATION.Result["error"]; ok {out += v }; out += `</span>

					<input value="`+r.FormValue("email")+`" type="text" name="email" placeholder="Your Email"   class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["email"]; ok {out += v }; out += `</span>

					<input type="password" name="password" placeholder="Password" required="" class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["password"]; ok {out += v }; out += `</span>

					<input type="password" name="confirm_password" placeholder="Confirm password" required=""  class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["confirm_password"]; ok {out += v }; out += `</span>


					`; if app.Use_recaptcha_register() { out += `
                                    	<script src="https://www.google.com/recaptcha/api.js"></script>
					<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {out += v }; out += `<span>
					<br>
					`}; out += `

					<div class="signin-rit">
						<span class="agree-checkbox">
							<label class="checkbox"><input type="checkbox" name="i_agree" value="1">I agree to your <a class="w3layouts-t" href="terms.html" target="_blank">Terms of Use</a> and <a class="w3layouts-t" href="privacy.html" target="_blank">Privacy Policy</a></label>
						</span>
					</div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["i_agree"]; ok {out += v }; out += `</span>

					<input type="submit" value="`+lang.T("sign up")+`">
				</div>
				</form>
				`}; out += `
			</div>
		</div>
	</div>
</div>`
	return out
}
