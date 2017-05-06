package authView

import (
	"net/http"
	"Malina/config"
	"Malina/language"
	"Malina/libraries"
)

func LoginForm(r *http.Request)string{
	var out string = ""
	out += `
<div class="container">
	<div class="row" style="margin-top:auto;">
        	<div class="col-md-4 col-md-offset-4">
        		<div class="login-panel panel panel-default">
                    		<div class="panel-heading">
                        		<h3 class="panel-title">`+lang.T("login")+`</h3>
                    		</div>
                    		<div class="panel-body">
					<h3>`+lang.T("sign in")+`</h3>
					`; if v,ok:= library.VALIDATION.Result["success"]; ok {
					out += `<h1>`+v+`</h1>
					<h6><a href="`+app.Base_url()+`auth/logout/" >`+lang.T("logout")+`</a></h6>
					`;
					}else{out += `
					<form  method="post">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["error"]; ok {out += v }; out += `</span>

						<input  value="`+r.FormValue("email_nick_phone")+`" type="text" name="email_nick_phone" placeholder="Your Email (ilyaran@mail.ru) or Nick (ilyasvetozar) or Phone (+77058436633)" required="">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["email_nick_phone"]; ok {out += v }; out += `</span>

						<input type="password" name="password" placeholder="Password" required="">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["password"]; ok {out += v }; out += `</span>

						`; if app.Use_recaptcha_login() { out += `
						<script src="https://www.google.com/recaptcha/api.js"></script>
						<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
						<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {out += v }; out += `<span>
						<br>
						`}; out += `
						<label class="checkbox"><input type="checkbox" name="remember" checked >Remember me</label>
						<input type="submit" value="Sign In">
						<div class="forgot-grid">
							<div class="forgot">
								<a href="`+app.Base_url()+`auth/forgot/" >Forgot Password?</a>
							</div>
							<div class="clearfix"> </div>
						</div>
					</form>
					<h6> Not a Member Yet? <a href="`+app.Base_url()+`auth/register/">Sign Up Now</a> </h6>
					`}; out += `
				</div>
			</div>
		</div>
	</div>
</div>

	`

	return out
}
