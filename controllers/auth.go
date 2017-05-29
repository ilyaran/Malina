/**
 * Authorization controller class.  github.com/ilyaran/Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package controller

import (
	"net/http"
	"fmt"
	"regexp"
	"github.com/ilyaran/Malina/language"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/libraries"
	"github.com/gorilla/mux"
	"strings"
	"time"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/views/auth"
)

var AuthController = &authController{&CrudController{}}
type authController struct{crud *CrudController }

func (this *authController) Index (w http.ResponseWriter, r *http.Request) {

	model.AuthModel.Query, model.AuthModel.Where = ``, ``
	model.AuthModel.All = 0

	library.SESSION.Authentication(w, r)

	library.VALIDATION.Status = 0
	library.VALIDATION.Result = map[string]string{}

	views.LOCALS.W = w
	views.LOCALS.R = r
	views.LOCALS.CurrentPage = ``

	authView.AuthViewObj.Message = ``

	switch mux.Vars(r)["action"] {
		case "register" 	: this.Register(w,r)
		case "logout" 		: this.Logout(w,r)
		case "delete_account" 	: this.DeleteAccount(w,r)
		case "forgot" 		: this.Forgot(w,r)
		case "activation" 	: this.Activation(w,r)
		case "change_password" 	: this.ChangePassword(w,r)
		default :
			this.Login(w,r)
	}
}

func (this *authController) Login(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId > 0 {

		authView.AuthViewObj.Message = lang.T("auth_login_allready")

	} else if r.Method == "POST" {

		model.AuthModel.CountLoginAttempts(library.SESSION.GetIP(r))

		if model.AuthModel.All >= app.Login_attempts() {
			library.VALIDATION.Status = 120
			library.VALIDATION.Result["error"] = lang.T("auth exceed max attempts")
		}else {
			var recaptchaChannel chan bool
			if app.Use_recaptcha_login() {
				recaptchaChannel = make(chan bool)
				go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
			}
			var email = library.VALIDATION.IsEmail(true, r)
			var pass = library.VALIDATION.PasswordValid(true,"password",false,false,false, r)
			if app.Use_recaptcha_login() {
				if !(<-recaptchaChannel) {
					library.VALIDATION.Status = 100
					library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
				}
			}
			if library.VALIDATION.Status == 0 {
				var accountObj = model.AuthModel.GetAccountByEmail(email)
				if accountObj != nil {
					if accountObj.GetBan() {
						authView.AuthViewObj.Message = "Access is denied. You are banned"
					}else {
						if library.SESSION.Cryptcode(pass) == accountObj.GetPassword() {
							library.VALIDATION.Status = 0
							authView.AuthViewObj.Message = lang.T("success login")

							go model.AuthModel.ClearLoginAttempts(library.SESSION.SessionObj.Ip_address)

							model.Session.Del(library.SESSION.SessionObj.Id)

							library.SESSION.SetSession(accountObj.Id,"logged_in" ,w)


							http.Redirect(w, r, "/", 301)

						} else {
							library.VALIDATION.Status = 100
							library.VALIDATION.Result["error"] = lang.T("sign in failed your password is an incorrect")
							// Increase login attempts
							model.AuthModel.IncreaseLoginAttempts(library.SESSION.GetSessionObj().GetIp_address())
						}
					}
				} else {
					library.VALIDATION.Status = 100
					library.VALIDATION.Result["email"] = lang.T("email not found")
				}
			}
		}
	}

	authView.AuthViewObj.LoginForm()
}

func (this *authController) Logout(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId >0 {
		library.SESSION.DeleteSession(library.SESSION.GetSessionObj().GetId(),w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}else {
		http.Redirect(w, r, "/"+app.Uri_auth_login(), http.StatusTemporaryRedirect)
	}
}

func (this *authController)  DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId > 0 {
		library.SESSION.DeleteSession(library.SESSION.SessionObj.Id,w)
		model.AccountModel.Del(library.SESSION.SessionObj.AccountId)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}else {
		http.Redirect(w, r, "/"+app.Uri_auth_login(), http.StatusTemporaryRedirect)
	}
}

func (this *authController)  Register(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId > 0 {

		authView.AuthViewObj.Message = lang.T("auth_logout_first")

	}else if r.Method == "POST" {

		var recaptchaChannel chan bool
		if app.Use_recaptcha_register() {
			recaptchaChannel = make(chan bool)
			go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
		}

		var email = library.VALIDATION.IsEmail(true, r)
		model.AuthModel.CheckDetails(`account_email = $1`, email)
		if model.AuthModel.All > 0 {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("email exists allready")
		}

		var pass = library.VALIDATION.PasswordValid(true, "password", false, false, false, r)
		library.VALIDATION.ConfirmPasswordValid(pass, r)
		if app.Use_recaptcha_register() {
			if !(<-recaptchaChannel) {
				library.VALIDATION.Status = 100
				library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
			}
		}

		if library.VALIDATION.Status == 0 {

			var activation_key = library.SESSION.Cryptcode(fmt.Sprintf(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(), library.SESSION.GetIP(r))))

			if model.AuthModel.AddActivation(email, pass, activation_key, library.SESSION.GetIP(r)) {

				/*var emailContent string = fmt.Sprintf(lang.T("auth_activate_content"),
					app.Base_url(),
					app.Base_url()+app.Uri_auth_activation()+"?activation_key="+activation_key,
					app.Email_activation_expire()/int64(3600),
					email,
					pass,
					app.Base_url(),
				)*/

				fmt.Println(app.Base_url() + app.Uri_auth_activation() + "?activation_key=" + activation_key)

				//go library.SESSION.SendEmail(app.Site_name()+lang.T("auth_account_subject"), email, emailContent)

				library.VALIDATION.Status = 0

				authView.AuthViewObj.Message = lang.T("auth_success_reg_check")
			} else {
				library.VALIDATION.Status = 30
				library.VALIDATION.Result["error"] = lang.T("server error")
			}
		}
	}

	authView.AuthViewObj.RegisterForm()
}

func (this *authController) Activation(w http.ResponseWriter, r *http.Request) {

	if library.SESSION.SessionObj.AccountId > 0 {

		authView.AuthViewObj.Message = lang.T("auth_logout_first")

	}else {
		activation_key := r.URL.Query().Get("activation_key")
		match, err := regexp.MatchString("^[a-zA-Z0-9_]{30,255}$", activation_key)
		if err == nil && match {
			var account= model.AuthModel.GetActivation(activation_key)
			if account != nil {
				account.Password = library.SESSION.Cryptcode(account.GetNewpass())
				last_user_id := model.AuthModel.AddAccount(`(account_email,account_password)VALUES($1,$2)`, []interface{}{account.GetEmail(), account.GetPassword()})
				if last_user_id > 0 {
					account.SetId(last_user_id)

					model.Session.Del(library.SESSION.SessionObj.Id)
					//fmt.Println(library.SESSION.SessionObj.Id)
					library.SESSION.SetSession(account.Id, "activated",w)
					/*var emailContent string = fmt.Sprintf(
						lang.T("auth_account_content"),
						app.Base_url(),
						accountController.GetEmail(),
						accountController.GetNewpass(),
						app.Base_url()+app.Uri_cabinet(),
						app.Site_name(),
					)
					go library.SESSION.SendEmail(app.Site_name()+lang.T("auth_account_subject"), accountController.GetEmail(), emailContent)
					*/
					//views.LOCALS.CurrentAccount = account

					authView.AuthViewObj.Message = lang.T("auth_activation_success") + `<br><a href="` + app.Base_url() + `">` + lang.T("home") + `</a>`
				}
			} else {
				authView.AuthViewObj.Message = lang.T("auth_activation_incorrect_code")
			}
		}
	}
	authView.AuthViewObj.ActivationForm()
}

func (this *authController)  Forgot(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId > 0 {
		library.VALIDATION.Status = 0
		library.VALIDATION.Result["success"] = lang.T("auth_login_allready")
	}else if r.Method == "POST" {
		var recaptchaChannel chan bool
		if app.Use_recaptcha_forgot() {
			recaptchaChannel = make(chan bool)
			go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
		}

		email := library.VALIDATION.IsEmail(true, r)
		var account = model.AuthModel.GetAccountByEmail(email)
		if account == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("email doesn't exist")
		}

		if app.Use_recaptcha_forgot() {
			if !(<-recaptchaChannel) {
				library.VALIDATION.Status = 100
				library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
			}
		}
		if library.VALIDATION.Status == 0 {
			var newPass = library.SESSION.Generate_password(app.PasswordMinLength())
			// generate and crypt password
			var cryptedNewPass string = library.SESSION.Cryptcode(newPass)
			var email string = strings.ToLower(r.FormValue("email"))
			emailOrError, ok := model.AuthModel.ForgotPass(email, cryptedNewPass)
			if ok {
				var emailContent string = fmt.Sprintf(lang.T("auth_forgot_password_content"),
					app.Base_url(),
					email,
					newPass,
					app.Admin_email(),
					app.Base_url(),
				)
				go library.SESSION.SendEmail(lang.T("auth_forgot_password_subject"), emailOrError, emailContent)
				library.VALIDATION.Status = 0
				authView.AuthViewObj.Message = lang.T("auth_new_pass")
			} else {
				library.VALIDATION.Status = 300
				library.VALIDATION.Result["server_error"] = lang.T("server error")
			}
		}
	}

	authView.AuthViewObj.ForgotForm()
}

func (this *authController)  ChangePassword(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountId > 0 {
		library.VALIDATION.Result["success"] = lang.T("auth_reset_first")
		//w.Write([]byte(authView.ChangePassword(r)))
		return
	}
	//recaptcha
	var recaptchaChannel chan bool
	if app.Use_recaptcha_change_password() {
		recaptchaChannel = make(chan bool)
		go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
	}

	var oldPass string = library.VALIDATION.PasswordValid(true,"old_password",false,false,false, r)
	var newPass string = library.VALIDATION.PasswordValid(true,"new_password",false,false,false, r)
	library.VALIDATION.ConfirmPasswordValid(newPass,r)
	var cryptedOldPass string = library.SESSION.Cryptcode(oldPass)
	if oldPass != "" {
		cryptedOldPass = library.SESSION.Cryptcode(oldPass)

		var accountObj = model.AccountModel.Get(library.SESSION.GetSessionObj().GetAccountId())
		if accountObj.Password != cryptedOldPass {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["old_password"] = lang.T("auth_incorrect_password")
		}
	}

	if app.Use_recaptcha_change_password() {
		if !(<-recaptchaChannel) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
		}
	}
	if library.VALIDATION.Status  == 0 {
		//var cryptedOldPass string = library.SESSION.Cryptcode(oldPass)
		var cryptedNewPass string = library.SESSION.Cryptcode(newPass)
		res, ok := model.AuthModel.ChangePass(library.SESSION.SessionObj.Email, cryptedOldPass, cryptedNewPass)

		if ok {
			library.VALIDATION.Status = 0
			authView.AuthViewObj.Message = lang.T("auth_pass_success_changed")
		} else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["error"] = lang.T(res)
		}
	}

	authView.AuthViewObj.ChangePassword()
}










/*

var AuthController = &authController{&CrudController{}}

type authController struct { crud *CrudController }

func (this *authController) Index(w http.ResponseWriter, r *http.Request) {

	model.AuthModel.Query = ``
	model.AuthModel.Where = ``
	model.AuthModel.All = 0

	library.SESSION.Authentication(w, r)

	library.VALIDATION.Status = 0
	library.VALIDATION.Result = map[string]string{}

	switch mux.Vars(r)["action"] {
		case "register":this.Register(w, r)
		case "logout":this.Logout(w, r);return
		case "delete_account":this.DeleteAccount(w, r)
		case "forgot":this.Forgot(w, r)
		case "activation":this.Activation(w, r);return
		case "change_password":this.ChangePassword(w, r)
		case "login": this.Login(w, r)
		//default: this.Form(r)
	}

	helper.SetAjaxHeaders(w)
	out, _ := json.Marshal(library.VALIDATION)
	fmt.Fprintf(w, string(out))
}
*/
/*func (this *authController) Form(r *http.Request) {

	publicView.AuthViewObj.Form()
}*//*

func (this *authController) Login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(library.SESSION.SessionObj)
	if library.SESSION.SessionObj.AccountController != nil {
		library.VALIDATION.Status = 50
		library.VALIDATION.Result["allready_logged_in"] = lang.T("auth_login_allready")
		return
	}
	model.AuthModel.CountLoginAttempts(library.SESSION.GetIP(r))
	if model.AuthModel.All >= app.Login_attempts() {
		library.VALIDATION.Status = 120
		library.VALIDATION.Result["error"] = lang.T("auth exceed max attempts")
		return
	}
	var recaptchaChannel chan bool
	if app.Use_recaptcha_login() {
		recaptchaChannel = make(chan bool)
		go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
	}
	var email = library.VALIDATION.IsEmail(true, r)
	var pass = library.VALIDATION.PasswordValid(true, "password", false, false, false, r)
	if app.Use_recaptcha_login() {
		if !(<-recaptchaChannel) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
		}
	}
	if library.VALIDATION.Status == 0 {
		var cryptedPass string = library.SESSION.Cryptcode(pass)
		var accountObj = model.AuthModel.GetAccountByEmail(email)
		if accountObj != nil {
			//fmt.Println(accountObj)
			if accountObj.Ban {
				library.VALIDATION.Status = 190
				library.VALIDATION.Result["ban"] = lang.T("access is denied you are banned")
				return
			} else {
				if cryptedPass == accountObj.GetPassword() {
					library.VALIDATION.Status = 0
					library.VALIDATION.Result["success"] = lang.T("success login")

					go model.AuthModel.ClearLoginAttempts(library.SESSION.SessionObj.Ip_address)

					go model.Session.Del(library.SESSION.SessionObj.Id)

					go library.SESSION.SetSession(accountObj, "logged_in", library.SESSION.GetIP(r), r.Header.Get("User-Agent"), w)

					return
				} else {
					library.VALIDATION.Status = 100
					library.VALIDATION.Result["error"] = lang.T("auth_login_email_or_password_incorrect")
					// Increase login attempts
					go model.AuthModel.IncreaseLoginAttempts(library.SESSION.GetSessionObj().GetIp_address())
					return
				}
			}
		}else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("auth_email_not_exist")
		}
	}
}

func (this *authController) Activation(w http.ResponseWriter, r *http.Request) {
	var msg string
	if library.SESSION.SessionObj.AccountController != nil {
		msg = lang.T("auth_logout_first")
	} else {
		activation_key := r.URL.Query().Get("activation_key")
		match, err := regexp.MatchString("^[a-zA-Z0-9_]{30,512}$", activation_key)
		if err == nil && match {
			var accountController = model.AuthModel.GetActivation(activation_key)
			if accountController != nil {
				accountController.SetPassword(library.SESSION.Cryptcode(accountController.GetNewpass()))
				last_user_id := model.AuthModel.AddAccount(`email = ? , password = ?`, []interface{}{accountController.GetEmail(), accountController.GetPassword()})
				if last_user_id > 0 {
					accountController.SetId(last_user_id)
				}
				model.Session.Del(library.SESSION.SessionObj.Id)
				//fmt.Println(library.SESSION.SessionObj.Id)
				library.SESSION.SetSession(accountController, "activated", library.SESSION.GetIP(r), r.Header.Get("User-Agent"), w)
				*/
/*var emailContent string = fmt.Sprintf(
					lang.T("auth_account_content"),
					app.Base_url(),
					accountController.GetEmail(),
					accountController.GetNewpass(),
					app.Base_url()+app.Uri_cabinet(),
					app.Site_name(),
				)
				go library.SESSION.SendEmail(app.Site_name()+lang.T("auth_account_subject"), accountController.GetEmail(), emailContent)
				*//*

				publicView.PublicLocalsObj.CurrentAccount = accountController
				msg = lang.T("auth_activation_success") + `<br><a href="` + app.Base_url() + app.Uri_auth_login() + `">` + lang.T("login") + `</a>`
			} else {
				msg = lang.T("auth_activation_incorrect_code")
			}
		}
	}
	publicView.WelcomeViewObj.Layout.Body1 = []byte(`<h1>` + msg + `</h1>`)

}

func (this *authController) Register(w http.ResponseWriter, r *http.Request) {

	if library.SESSION.SessionObj.AccountController != nil {
		library.VALIDATION.Status = 50
		library.VALIDATION.Result["success"] = lang.T("logout first")
		return
	}

	var recaptchaChannel chan bool
	if app.Use_recaptcha_register() {
		recaptchaChannel = make(chan bool)
		go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
	}
	var email = library.VALIDATION.IsEmail(true, r)
	model.AuthModel.CheckDetails(`email = ?`, email)
	if model.AuthModel.All > 0 {
		library.VALIDATION.Status = 100
		library.VALIDATION.Result["email"] = lang.T("email exists allready")
	}
	var pass = library.VALIDATION.PasswordValid(true, "password", false, false, false, r)
	library.VALIDATION.ConfirmPasswordValid(pass, r)
	//library.VALIDATION.IAgree(r)
	if app.Use_recaptcha_register() {
		if !(<-recaptchaChannel) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
		}
	}
	if library.VALIDATION.Status == 0 {

		var activation_key = library.SESSION.Cryptcode(fmt.Sprintf(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(), library.SESSION.GetIP(r))))

		if model.AuthModel.AddActivation(email, pass, activation_key, library.SESSION.GetIP(r)) {

			*/
/*var emailContent string = fmt.Sprintf(lang.T("auth_activate_content"),
				app.Base_url(),
				app.Base_url()+app.Uri_auth_activation()+"?activation_key="+activation_key,
				app.Email_activation_expire()/int64(3600),
				email,
				pass,
				app.Base_url(),
			)*//*

			fmt.Println(app.Base_url() + app.Uri_auth_activation() + "?activation_key=" + activation_key)

			//go library.SESSION.SendEmail(app.Site_name()+lang.T("auth_account_subject"), email, emailContent)

			library.VALIDATION.Status = 0
			library.VALIDATION.Result["success"] = lang.T("auth_success_reg_check")
		} else {
			library.VALIDATION.Status = 30
			library.VALIDATION.Result["error"] = lang.T("server error")
		}
	}
}
func (this *authController) Logout(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountController != nil {
		library.SESSION.DeleteSession(library.SESSION.SessionObj.Id, w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (this *authController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountController != nil {
		library.SESSION.DeleteSession(library.SESSION.SessionObj.Id, w)
		model.AccountModel.Del(library.SESSION.SessionObj.AccountController.Id)
	}
	http.Redirect(w, r, "/"+app.Uri_auth_login(), http.StatusTemporaryRedirect)
}

func (this *authController) Forgot(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.SessionObj.AccountController != nil {
		library.VALIDATION.Status = 50
		library.VALIDATION.Result["success"] = lang.T("auth_login_allready")
		return
	}

	var recaptchaChannel chan bool
	if app.Use_recaptcha_forgot() {
		recaptchaChannel = make(chan bool)
		go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
	}

	email := library.VALIDATION.IsEmail(true, r)
	var accountController = model.AuthModel.GetAccountByEmail(email)
	if accountController == nil {
		library.VALIDATION.Status = 100
		library.VALIDATION.Result["email"] = lang.T("email doesn't exist")
	}

	if app.Use_recaptcha_forgot() {
		if !(<-recaptchaChannel) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
		}
	}
	if library.VALIDATION.Status == 0 {
		var newPass = library.SESSION.Generate_password(app.PasswordMinLength())
		// generate and crypt password
		var cryptedNewPass string = library.SESSION.Cryptcode(newPass)
		var email string = strings.ToLower(r.FormValue("email"))
		emailOrError, ok := model.AuthModel.ForgotPass(email, cryptedNewPass)
		if ok {
			var emailContent string = fmt.Sprintf(lang.T("auth_forgot_password_content"),
				app.Base_url(),
				email,
				newPass,
				app.Admin_email(),
				app.Base_url(),
			)
			go library.SESSION.SendEmail(lang.T("auth_forgot_password_subject"), emailOrError, emailContent)
			library.VALIDATION.Status = 0
			library.VALIDATION.Result["success"] = lang.T("auth_new_pass")
		} else {
			library.VALIDATION.Status = 300
			library.VALIDATION.Result["server_error"] = lang.T("server error")
		}
	}

}

func (this *authController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount().GetId() > 0 {
		library.VALIDATION.Result["success"] = lang.T("auth_reset_first")
		//w.Write([]byte(authView.ChangePassword(r)))
		return
	}
	//recaptcha
	var recaptchaChannel chan bool
	if app.Use_recaptcha_change_password() {
		recaptchaChannel = make(chan bool)
		go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
	}

	var oldPass string = library.VALIDATION.PasswordValid(true, "old_password", false, false, false, r)
	var newPass string = library.VALIDATION.PasswordValid(true, "new_password", false, false, false, r)
	library.VALIDATION.ConfirmPasswordValid(newPass, r)

	if oldPass != "" {
		user := model.AuthModel.GetAccountByEmail(library.SESSION.GetSessionObj().GetAccount().GetEmail())
		if user != nil {
			if user.GetPassword() != library.SESSION.Cryptcode(oldPass) {
				library.VALIDATION.Status = 100
				library.VALIDATION.Result["old_password"] = lang.T("incorrect old password")
			}
		} else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("does not exist")
		}
	}

	if app.Use_recaptcha_change_password() {
		if !(<-recaptchaChannel) {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
		}
	}
	if library.VALIDATION.Status == 0 {
		var cryptedOldPass string = library.SESSION.Cryptcode(oldPass)
		var cryptedNewPass string = library.SESSION.Cryptcode(newPass)
		res, ok := model.AuthModel.ChangePass(library.SESSION.GetSessionObj().GetAccount().GetEmail(),
			cryptedOldPass, cryptedNewPass)

		if ok {
			library.VALIDATION.Status = 0
			library.VALIDATION.Result["success"] = lang.T("auth_pass_success_changed")
		} else {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result[res] = lang.T(res)
		}
	}

	//w.Write([]byte(authView.ChangePassword(r)))
}
*/
