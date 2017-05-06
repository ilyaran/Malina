/**
 * Authorization controller class.  Malina eCommerce application
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
package controller

import (
	"net/http"
	"fmt"
	"regexp"
	"Malina/language"
	"Malina/config"
	"Malina/models"
	"github.com/gorilla/mux"
	"strings"
	"Malina/views/auth"
	"Malina/views"
	"Malina/libraries"
	"Malina/views/errors"
	"time"
)
var AuthController = &authController{&CrudController{}}
type authController struct{crud *CrudController }

func (this *authController) Index (w http.ResponseWriter, r *http.Request) {
	this.crud.auth("",w,r)
	this.crud.dbtable = "account"

	vars := mux.Vars(r)
	action := vars["action"]

	switch action {
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
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		library.VALIDATION.Status = 90
		library.VALIDATION.Result["success"] = lang.T("auth_login_allready")
	}else if r.Method == "POST" {
		if model.AuthModel.CountLoginAttempts(library.SESSION.GetIP(r)) >= app.Login_attempts() {
			library.VALIDATION.Status = 120
			library.VALIDATION.Result["error"] = lang.T("auth exceed max attempts")
		}else {
			var recaptchaChannel chan bool
			if app.Use_recaptcha_login() {
				recaptchaChannel = make(chan bool)
				go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
			}
			var email,nick,_ = library.VALIDATION.IsEmailOrNickOrPhone(true,r)
			pass := library.VALIDATION.PasswordValid(true,"password",false,false,false, r)
			if app.Use_recaptcha_login() {
				if !(<-recaptchaChannel) {
					library.VALIDATION.Status = 100
					library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
				}
			}
			if library.VALIDATION.Status == 0 {
				var cryptedPass string = library.SESSION.Cryptcode(pass)
				var account = model.AccountModel.GetByEmailByNickByPhoneByPassword(email,nick,"",cryptedPass)
				if account != nil {
					if account.GetBanned() == 0{
						library.VALIDATION.Status = 0
						library.VALIDATION.Result["success"] = lang.T("success login")
						model.Session.Del(library.SESSION.GetSessionObj().GetId())
						if r.FormValue("remember") != "" {
							library.SESSION.SetSession(account, "logged", false, library.SESSION.GetIP(r), r.Header.Get("User-Agent"), w)
						} else {
							library.SESSION.SetSession(account, "logged", true, library.SESSION.GetIP(r), r.Header.Get("User-Agent"), w)
						}
						// Clear login attempts
						model.AuthModel.ClearLoginAttempts(library.SESSION.GetSessionObj().GetIp_address())

						http.Redirect(w, r, "/", 301)
					}else {
						w.Write([]byte(errors.ErrorView("Access is denied","You are banned")))
						return
					}
				} else {
					library.VALIDATION.Status = 100
					library.VALIDATION.Result["error"] = lang.T("auth email or nick or phone and password are not exist")
					// Increase login attempts
					model.AuthModel.IncreaseLoginAttempts(library.SESSION.GetSessionObj().GetIp_address())
				}
			} else {
				// Increase login attempts
				model.AuthModel.IncreaseLoginAttempts(library.SESSION.GetSessionObj().GetIp_address())
			}
		}
	}
	w.Write([]byte(views.Header()))
	w.Write([]byte(authView.LoginForm(r)))
	w.Write([]byte(views.Footer()))
}

func (this *authController) Logout(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		library.SESSION.DeleteSession(library.SESSION.GetSessionObj().GetId(),w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (this *authController)  DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		library.SESSION.DeleteSession(library.SESSION.GetSessionObj().GetId(),w)
		model.AccountModel.Del(library.SESSION.GetSessionObj().GetAccount_id())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}else {
		http.Redirect(w, r, "/"+app.Uri_auth_login(), http.StatusTemporaryRedirect)
	}
}

func (this *authController)  Register(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		library.VALIDATION.Status = 50
		library.VALIDATION.Result["success"] = lang.T("auth_logout_first")
	}else if r.Method == "POST" {
		var recaptchaChannel chan bool
		if app.Use_recaptcha_register() {
			recaptchaChannel = make(chan bool)
			go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
		}
		var email = library.VALIDATION.IsEmail(true,r)

		var account = model.AccountModel.GetByEmailByNickByPhoneByPassword(email,"","","")
		if account != nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email"] = lang.T("auth_email_exist_allready")
		}
		var pass = library.VALIDATION.PasswordValid(true,"password",false,false,false, r)
		library.VALIDATION.ConfirmPasswordValid(pass, r)
		library.VALIDATION.IAgree(r)
		if app.Use_recaptcha_register() {
			if !(<-recaptchaChannel) {
				library.VALIDATION.Status = 100
				library.VALIDATION.Result["recaptcha"] = lang.T("wrong_recaptcha")
			}
		}
		if library.VALIDATION.Status == 0 {
			var activation_key = library.SESSION.Cryptcode(fmt.Sprintf(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(), library.SESSION.GetIP(r))))

			if model.AuthModel.AddActivation(email,"","",pass,activation_key, library.SESSION.GetIP(r)) {
				if email != "" {
					var emailContent string = fmt.Sprintf(lang.T("auth_activate_content"),
						app.Base_url(),
						app.Base_url() + app.Uri_auth_activation() + "?activation_key=" + activation_key,
						app.Email_activation_expire() / int64(3600),
						email,
						pass,
						app.Base_url(),
					)
					fmt.Println(app.Base_url() + app.Uri_auth_activation() + "?activation_key=" + activation_key)

					go library.SESSION.SendEmail(app.Site_name() + lang.T("auth_activate_subject"), email, emailContent)
				}
				library.VALIDATION.Result["success"] = lang.T("auth_success_reg_check")
			} else {
				library.VALIDATION.Result["error"] = lang.T("server error")
			}
		}
	}
	w.Write([]byte(views.Header()))
	w.Write([]byte(authView.RegisterForm(r)))
	w.Write([]byte(views.Footer()))
}

func (this *authController) Activation(w http.ResponseWriter, r *http.Request) {
	var msg string
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		msg = lang.T("auth_logout_first")
	}
	msg = lang.T("auth_activation_incorrect_code")
	activation_key := r.URL.Query().Get("activation_key")
	match, err := regexp.MatchString("^[a-zA-Z0-9_]{30,512}$", activation_key)
	if err == nil && match {
		var account = model.AuthModel.GetActivation(activation_key)
		if account != nil {
			account.SetPass(library.SESSION.Cryptcode(account.GetNew_password()))
			var newpass = account.GetNew_password()
			account.SetNewPass("")
			last_user_id := model.AccountModel.Add(account)
			if last_user_id>0{
				account.SetId(last_user_id)
			}
			library.SESSION.SetSession(account,"activated",false,library.SESSION.GetIP(r),r.Header.Get("User-Agent"),w)
			var emailContent string = fmt.Sprintf(
				lang.T("auth_account_content"),
				app.Base_url(),
				account.GetEmail(),
				newpass,
				app.Base_url()+app.Uri_cabinet(),
				app.Base_url(),
			)
			go library.SESSION.SendEmail(app.Site_name()+lang.T("auth_account_subject"), account.GetEmail(), emailContent)
			msg = lang.T("auth_activation_success") + `<br><a href="` + app.Base_url() + app.Uri_auth_login() + `">` + lang.T("login") + `</a>`
			model.Session.Del(library.SESSION.GetSessionObj().GetId())
			model.Session.Del(activation_key)
		}
	}
	w.Write([]byte(views.Header()))
	w.Write([]byte(authView.ActivationForm(msg)))
	w.Write([]byte(views.Footer()))
}

func (this *authController)  Forgot(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
		library.VALIDATION.Status = 0
		library.VALIDATION.Result["success"] = lang.T("auth_login_allready")
	}else if r.Method == "POST" {
		var recaptchaChannel chan bool
		if app.Use_recaptcha_forgot() {
			recaptchaChannel = make(chan bool)
			go library.VALIDATION.RecaptchaValid(recaptchaChannel, r)
		}

		email, nick, phone := library.VALIDATION.IsEmailOrNickOrPhone(true, r)
		account := model.AccountModel.GetByEmailByNickByPhoneByPassword(email, nick, phone, "")
		if account == nil {
			library.VALIDATION.Status = 100
			library.VALIDATION.Result["email_nick_phone"] = lang.T("auth nick or email or phone not exist")
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
	w.Write([]byte(views.Header()))
	w.Write([]byte(authView.ForgotForm( r)))
	w.Write([]byte(views.Footer()))
}

func (this *authController)  ChangePassword(w http.ResponseWriter, r *http.Request) {
	if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
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

	if oldPass != "" {
		user := model.AccountModel.GetByEmailByNickByPhoneByPassword(library.SESSION.GetSessionObj().GetEmail(),library.SESSION.GetSessionObj().GetNick(),library.SESSION.GetSessionObj().GetPhone(), library.SESSION.Cryptcode(oldPass))
		if user == nil {
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
		var cryptedOldPass string = library.SESSION.Cryptcode(oldPass)
		var cryptedNewPass string = library.SESSION.Cryptcode(newPass)
		res, ok := model.AuthModel.ChangePass(library.SESSION.GetSessionObj().GetEmail(),
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
