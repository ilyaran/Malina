package authView

import (
	"net/http"
	"Malina/language"
	"Malina/config"
	"Malina/libraries"
)

func ForgotForm(r *http.Request)string{
	var out string = ""

	out += `
	<!-- sign Up form -->
	 <section>
		<div id="agileits-sign-in-page" class="sign-in-wrapper">
			<div class="agileinfo_signin">
			<h3>Forgot</h3>
				`; if v,ok:= library.VALIDATION.Result["success"]; ok {
				out += `<h1>`+v+`</h1>
				<h6><a href="`+app.Base_url()+`auth/logout/" >`+lang.T("auth_login_allready")+`</a></h6>
				`;
				}else{out += `
				<form method="post">

					<span class="error">`; if v,ok:= library.VALIDATION.Result["error"]; ok {out += v }; out += `</span>

					<input value="`+r.FormValue("email_nick")+`" type="text" name="email_nick" placeholder="Your Email or Nick" required="">
					<span class="error">`; if v,ok:= library.VALIDATION.Result["email_nick"]; ok {out += v }; out += `</span>

					`; if app.Use_recaptcha_forgot() { out += `
                                    	<script src="https://www.google.com/recaptcha/api.js"></script>
					<div id="recaptcha" class="g-recaptcha" data-sitekey="`+app.Recaptcha_public()+`"></div>
					<span class="error">`; if v,ok:= library.VALIDATION.Result["recaptcha"]; ok {out += v }; out += `<span>
					<br>
					`}; out += `

					<input type="submit" value="Submit">

				</form>

				`}; out += `

			</div>
		</div>
	</section>
	<!-- //sign Up form -->
	`

	return out
}

