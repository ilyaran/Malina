package authView

import (
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/libraries"

)

var AuthViewObj *AuthView
func AuthViewObjInit() {
	AuthViewObj = &AuthView{Layout: views.NewLayout()}
}

type AuthView struct {
	Layout          *views.Layout
	Message 	string
	temp            string
}

func (s *AuthView)LoginForm(){
	s.temp = `
<div class="container">
	<div class="row" style="margin-top:auto;">
        <div class="col-md-4 col-md-offset-4">
        	<div class="login-panel panel panel-default">
                <div class="panel-heading">
                    <h3 class="panel-title">`+lang.T("login")+`</h3>
                </div>
                <div class="panel-body">
					<h3>`+lang.T("sign in")+`</h3>`
	
					if s.Message != `` {
						s.temp += `<h1>`+s.Message+`</h1>
						<h6><a href="`+app.Uri_auth_logout()+`" >`+lang.T("logout")+`</a></h6>`
					} else {
						s.temp += `
					<form  method="post">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["error"]; ok {s.temp += v }; s.temp += `</span>

						<input  value="`+views.LOCALS.R.FormValue("email")+`" type="text" name="email" placeholder="Your Email" required="">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["email"]; ok {s.temp += v }; s.temp += `</span>

						<input type="password" name="password" placeholder="Password" required="">
						<span class="error">`; if v,ok:= library.VALIDATION.Result["password"]; ok {s.temp += v }; s.temp += `</span>

						`; if app.Use_recaptcha_login() { s.temp += `
						<script src="https://www.google.com/recaptcha/api.js"></script>
						<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
						<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {s.temp += v }; s.temp += `<span>
						<br>
						`}; s.temp += `
						<label class="checkbox"><input type="checkbox" name="remember" checked >Remember me</label>
						<input type="submit" value="Sign In">
						<div class="forgot-grid">
							<div class="forgot">
								<a href="`+app.Uri_auth_forgot()+`" >Forgot Password?</a>
							</div>
							<div class="clearfix"> </div>
						</div>
					</form>
					<h6> Not a Member Yet? <a href="`+app.Uri_auth_register()+`">Sign Up Now</a> </h6>
					`}
					s.temp += `
				</div>
			</div>
		</div>
	</div>
</div>

	`
	s.Layout.Body2 = []byte(s.temp)
	s.Layout.WriteAuthResponse()
}

func (s *AuthView)RegisterForm(){
	s.temp = `
<div class="container">
	<div class="row" style="margin-top:auto;">
        <div class="col-md-4 col-md-offset-4">
        	<div class="login-panel panel panel-default">
                <div class="panel-heading">
                    <h3 class="panel-title">`+lang.T("register")+`</h3>
                </div>
                <div class="panel-body">
					<h3>`+lang.T("sign up")+`</h3>`

				if s.Message != `` {
					s.temp += `<h1>` + s.Message + `</h1>`
				} else {
					s.temp += `
				<form  method="post">
				<div class="form-group">
					<span class="error">`; if v,ok := library.VALIDATION.Result["error"]; ok {s.temp += v }; s.temp += `</span>

					<input value="`+views.LOCALS.R.FormValue("email")+`" type="text" name="email" placeholder="Your Email"   class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["email"]; ok {s.temp += v }; s.temp += `</span>

					<input type="password" name="password" placeholder="Password" required="" class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["password"]; ok {s.temp += v }; s.temp += `</span>

					<input type="password" name="confirm_password" placeholder="Confirm password" required=""  class="form-control">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["confirm_password"]; ok {s.temp += v }; s.temp += `</span>
					`

					if app.Use_recaptcha_register() { s.temp += `
                                    	<script src="https://www.google.com/recaptcha/api.js"></script>
					<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {s.temp += v }; s.temp += `<span>
					<br>
					`}

					s.temp += `
					<div class="signin-rit">
						<span class="agree-checkbox">
							<label class="checkbox"><input type="checkbox" name="i_agree" value="1">I agree to your <a class="w3layouts-t" href="terms.html" target="_blank">Terms of Use</a> and <a class="w3layouts-t" href="privacy.html" target="_blank">Privacy Policy</a></label>
						</span>
					</div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["i_agree"]; ok {s.temp += v }; s.temp += `</span>

					<input type="submit" value="`+lang.T("sign up")+`">
				</div>
				</form>
				`}
				s.temp += `
			</div>
		</div>
	</div>
</div>
	`
	s.Layout.Body2 = []byte(s.temp)
	s.Layout.WriteAuthResponse()
}

func (s *AuthView) ForgotForm() {
	s.temp = `
	
	 <section>
		<div id="agileits-sign-in-page" class="sign-in-wrapper">
			<div class="agileinfo_signin">
			<h3>Forgot</h3>
				` 
				if library.SESSION.SessionObj.AccountId > 0 {
				s.temp += `<h1>`+s.Message+`</h1>
				<h6><a href="`+app.Base_url()+`auth/logout/" >`+lang.T("auth_login_allready")+`</a></h6>
				`
				}else{
					s.temp += `
				<form method="post">

					<span class="error">` + s.Message + `</span>

					<input value="` + views.LOCALS.R.FormValue("email")+`" type="text" name="email" placeholder="Your Email or Nick" required="">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["email"]; ok {s.temp += v }; s.temp += `</span>

					`; if app.Use_recaptcha_forgot() { s.temp += `
                                    	<script src="https://www.google.com/recaptcha/api.js"></script>
					<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {s.temp += v }; s.temp += `<span>
					<br>
					`} 
					s.temp += `

					<input type="submit" value="Submit">

				</form>`
				}
				s.temp += `
			</div>
		</div>
	</section>
	
	`
	s.Layout.Body2 = []byte(s.temp)
	s.Layout.WriteAuthResponse()
}

func (s *AuthView)ActivationForm(){

	s.temp = `
	<!-- Activation form -->
	 <section>
		<div id="agileits-sign-in-page" class="sign-in-wrapper">
			<div class="agileinfo_signin">
				<h3>Activation</h3>

				<div class="alert alert-warning alert-dismissible" role="alert">
					<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
					<strong>` + s.Message + `</strong>
				</div>

			</div>
		</div>
	</section>
	<!-- //Activation form -->

	`
	s.Layout.Body2 = []byte(s.temp)
	s.Layout.WriteAuthResponse()
}

func (s *AuthView) ChangePassword (){
	s.temp += `
	<input type="hidden" value="login" id="actionType"/>

<div class="container">
	<div class="row" style="margin-top:auto;">
        <div class="col-md-4 col-md-offset-4">
            <div class="login-panel panel panel-default">
                <div class="panel-heading">
                    <h3 class="panel-title">`+lang.T("reset_password")+`</h3>
                </div>
                <div class="panel-body">
		` 
				if s.Message != `` {
					s.temp += `<h1>`+s.Message+`</h1>`
				}else{
					s.temp += `
					<form method="post" action="`+app.Base_url()+app.Uri_auth_change_password()+`">

						<span class="error">`; if v,ok:= library.VALIDATION.Result["no_rows"]; ok {s.temp += v }; s.temp += `<span>
						<span class="error">`; if v,ok:= library.VALIDATION.Result["error"]; ok {s.temp += v }; s.temp += `</span>

		                                <div class="form-group">
		                                    <input class="form-control" value="`+ views.LOCALS.R.FormValue("old_password")+`" placeholder="`+lang.T("auth_pass_old_title")+`" name="old_password" type="text" autofocus>
		                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["old_password"]; ok {s.temp += v }; s.temp += `<span>
		                                </div>

		                                <div class="form-group">
		                                    <input class="form-control" value="`+views.LOCALS.R.FormValue("new_password")+`" placeholder="`+lang.T("auth_pass_new_title")+`" name="new_password" type="text" autofocus>
		                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["new_password"]; ok {s.temp += v }; s.temp += `<span>
		                                </div>

		                                <div class="form-group">
		                                    <input class="form-control" value="`+views.LOCALS.R.FormValue("confirm_new_password")+`" placeholder="`+lang.T("auth_pass_new_confirm_title")+`" name="confirm_new_password" type="text" autofocus>
		                                    <span class="error">`; if v,ok:= library.VALIDATION.Result["confirm_new_password"]; ok {s.temp += v }; s.temp += `<span>
		                                </div>

		                               `; if app.Use_recaptcha_change_password() { s.temp += `
		                                        <script src="https://www.google.com/recaptcha/api.js"></script>
							<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
							<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {s.temp += v }; s.temp += `<span>
							<br>
						`}; s.temp += `

		                                <!-- Change this to a button or input when using this as a form -->
		                                <button type="submit" class="btn btn-lg btn-success btn-block">`+lang.T("send")+`</button>

					</form>`
					}
	
			s.temp += `
                    <hr>
                    <a href="`+app.Uri_auth_login() + `" class="btn btn-lg btn-primary">`+lang.T("login")+`</a>
                    <a href="`+app.Uri_auth_register() +`" class="btn btn-lg btn-default">`+lang.T("register")+`</a>
                </div>
            </div>
        </div>
    </div>
</div>
  `
	
	s.Layout.Body2 = []byte(s.temp)
	s.Layout.WriteAuthResponse()
}
