/**
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package controllers

import (
	"net/http"
	"fmt"
	"time"
	"regexp"
	"github.com/ilyaran/Malina/lang"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/filters"
	"github.com/ilyaran/Malina/dao"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/berry"
	"strings"
)

var AuthController = &Auth{view:publicView.Auth{}}
type Auth struct {

	view publicView.Auth

}
func(s *Auth)Index(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	s.view.Index(malina,"Welcome",w)
}

func(s *Auth)Login(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	var message string
	if malina.CurrentAccount != nil {

		message = lang.T("auth_login_allready")

	} else if r.Method == "POST" {

		attemptCounts:=models.AuthModel.CountLoginAttempts(dao.AuthDao.GetIP(r))

		if attemptCounts >= app.Login_attempts {
			malina.Status = 406
			malina.Result["error"] = lang.T("auth exceed max attempts")
		}else {

			var email = filters.IsEmail(malina,true,"email", r)
			var pass = filters.PasswordValid(malina,true,"password",false,false,false, r)

			if malina.Status == 0 {
				//fmt.Println("Login method ",email)
				var accountObj = models.AuthModel.GetAccountByEmail(0,email,AccountController.base.selectSqlFieldsDefault)
				if accountObj != nil {
					if accountObj.Ban {
						message = "Access is denied. You are banned"
					}else {
						if dao.AuthDao.Cryptcode(pass) == accountObj.Password {
							malina.Status = 0
							message = lang.T("success login")

							models.AuthModel.ClearLoginAttempts(dao.AuthDao.GetIP(r))

							//fmt.Println(dao.AuthDao.SessionObj.Id)
							dao.AuthDao.Model.Del(malina.SessionId)
							dao.AuthDao.SetSession(accountObj,r,w)

							http.Redirect(w, r, app.Url_cabinet_profile, 301 )
							return
						}
					}
				}
				malina.Status = 401
				malina.Result["error"] = lang.T("sign in failed your password is an incorrect")
				// Increase login attempts
				models.AuthModel.IncreaseLoginAttempts(dao.AuthDao.GetIP(r))
			}
		}
	}
	if malina.Status>0{
		for _,v:=range malina.Result {
			message += `<p>`+v.(string)+`</p>`
		}
		s.view.Index(malina,message,w)
		return
	}
	s.view.Login(malina,w,r)
}

func(s *Auth)Register(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	var message string
	if malina.CurrentAccount != nil {
		message = lang.T("auth_logout_first")
	}else if r.Method == "POST" {

		if app.Use_recaptcha_register {
			malina.ChannelBool = make(chan bool)
			go filters.RecaptchaValid(malina, r)
		}
		var email = filters.IsEmail(malina,true ,"email",r)
		malina.All = models.AuthModel.CheckDetails(`account_email = $1`, email)
		if malina.All > 0 {
			malina.Status = 404
			malina.Result["email"] = "email "+lang.T("exists allready")
		}
		var pass = filters.PasswordValid(malina,true, "password", false, false, false, r)
		filters.ConfirmPasswordValid(malina,"confirm_password",pass, r)
		if app.Use_recaptcha_register {
			if !(<-malina.ChannelBool) {
				malina.Status = 406
				malina.Result["recaptcha"] = lang.T("wrong_recaptcha")
			}
		}
		if malina.Status == 0 {

			var activation_key = dao.AuthDao.Cryptcode(fmt.Sprintf(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt, dao.AuthDao.GetIP(r))))
			if models.AuthModel.AddActivation(email, pass, activation_key, dao.AuthDao.GetIP(r)) {
				var emailContent string = fmt.Sprintf(lang.T("auth_activate_content"),
					app.Site_name,
					app.Url_auth_activation+`?activation_key=`+activation_key,
					app.Email_activation_expire/int64(3600),
					email,
					pass,
					app.Site_name,
				)
				fmt.Println(app.Url_auth_activation + "?activation_key=" + activation_key)
				go dao.AuthDao.SendEmail(app.Site_name+" "+lang.T("auth_activate_subject"), email, emailContent)

				malina.Status = 0
				message = lang.T("auth_success_reg_check")
				//s.view.RegistrationResult(malina,email,w)
				s.view.Index(malina,message,w)
				return
			} else {
				malina.Status = 500
				malina.Result["error"] = lang.T("server error")
			}
		}
	}

	//if malina.Status > 0 {
		//for _,v:=range malina.Result {
		//	message += `<p>`+v.(string)+`</p>`
		//}
		//s.view.AlertNotice(malina,message,w)
		//return
	//}

	s.view.Register(malina,w,r)
}

func(s *Auth)Activation(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	var message string
	if malina.CurrentAccount != nil {
		message = lang.T("auth_logout_first")
	}else {
		activation_key := r.URL.Query().Get("activation_key")
		match, err := regexp.MatchString("^[a-zA-Z0-9_]{30,255}$", activation_key)
		if err == nil && match {
			account := models.AuthModel.GetActivation(activation_key)
			//fmt.Println(account)
			if account != nil {
				account.Password = dao.AuthDao.Cryptcode(account.Newpass)
				//fmt.Println(account)
				last_user_id := models.AuthModel.AddAccount(`(account_email,account_password)VALUES($1,$2)`, account.Email, account.Password)
				if last_user_id > 0 {
					account.SetId(last_user_id)
					dao.AuthDao.Model.Del(malina.SessionId)
					dao.AuthDao.SetSession(account,r,w)
					var emailContent string = fmt.Sprintf(
						lang.T("auth_account_content"),
						app.Base_url,
						account.Email,
						account.Newpass,
						app.Url_cabinet,
						app.Site_name,
					)
					go dao.AuthDao.SendEmail(app.Site_name+" "+lang.T("auth_account_subject"), account.Email, emailContent)
					malina.CurrentAccount = account
					message = lang.T("auth_activation_success") + `<br><a href="` + app.Base_url + `">` + lang.T("home") + `</a>`
					s.view.Index(malina,message,w)
					return
				}
			} else {
				message = lang.T("auth_activation_incorrect_code")
			}
		}else {
			message = "Код активации "+lang.T("invalid")
		}
	}

	if malina.Status > 0 {
		for _,v:=range malina.Result {
			message += `<p>`+v.(string)+`</p>`
		}
	}
	s.view.Index(malina,message,w)
}

func(s *Auth)Logout(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	if malina.CurrentAccount != nil {
		dao.AuthDao.DeleteSession(malina.SessionId,w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (s *Auth)Forgot(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	var message string
	if malina.CurrentAccount != nil {
		malina.Status = 0
		//malina.Result["success"] = lang.T("auth_login_allready")
		message = lang.T("auth_login_allready")
	}else if r.Method == "POST" {
		if app.Use_recaptcha_forgot {
			malina.ChannelBool = make(chan bool)
			go filters.RecaptchaValid(malina, r)
		}

		email := filters.IsEmail(malina,true,"email", r)
		if email != "" {
			var account = models.AuthModel.GetAccountByEmail(0,email,AccountController.base.selectSqlFieldsDefault)
			if account == nil {
				malina.Status = 404
				malina.Result["email"] = lang.T("auth_email_not_exist")
			}
		}
		if app.Use_recaptcha_forgot {
			if !(<-malina.ChannelBool) {
				malina.Status = 406
				malina.Result["recaptcha"] = lang.T("wrong_recaptcha")
			}
		}
		if malina.Status == 0 {
			var newPass = dao.AuthDao.Generate_password()
			// generate and crypt password
			var cryptedNewPass string = dao.AuthDao.Cryptcode(newPass)
			var email string = strings.ToLower(r.FormValue("email"))
			emailOrError, ok := models.AuthModel.ForgotPass(email, cryptedNewPass)
			if ok {
				var emailContent string = fmt.Sprintf(lang.T("auth_forgot_password_content"),
					app.Base_url,
					email,
					newPass,
					app.Url_support,
					app.Base_url,
				)

				go dao.AuthDao.SendEmail(lang.T("auth_forgot_password_subject"), emailOrError, emailContent)

				malina.Status = 0
				s.view.ForgotResult(malina,message,w)
				return

			} else {
				malina.Status = 500
				malina.Result["error"] = lang.T("server error")
			}
		}
		if malina.Status>0{
			for _,v:=range malina.Result {
				message += `<p>`+v.(string)+`</p>`
			}
		}
	}

	s.view.ForgotForm(malina,message,w,r)
}

func (s *Auth)  ChangePassword(malina *berry.Malina,w http.ResponseWriter, r *http.Request) {
	var message string
	if malina.CurrentAccount == nil {
		message = lang.T("auth first")
		s.view.ChangePassword(malina,message,w,r)
		return
	}
	if r.Method == "POST" {
		var oldPass string = filters.PasswordValid(malina,true, "old_password", false, false, false, r)
		var newPass string = filters.PasswordValid(malina,true, "new_password", false, false, false, r)
		filters.ConfirmPasswordValid(malina,"confirm_new_password",newPass, r)
		var cryptedOldPass string = dao.AuthDao.Cryptcode(oldPass)
		if oldPass != "" {
			cryptedOldPass = dao.AuthDao.Cryptcode(oldPass)
			var accountObj = models.AuthModel.GetAccountByEmail(malina.CurrentAccount.Id,"",AccountController.base.selectSqlFieldsDefault)
			if accountObj.Password != cryptedOldPass {
				malina.Status = 406
				malina.Result["old_password"] = lang.T("incorrect password")
			}
		}

		if malina.Status == 0 {
			var cryptedOldPass string = dao.AuthDao.Cryptcode(oldPass)
			var cryptedNewPass string = dao.AuthDao.Cryptcode(newPass)
			res, ok := models.AuthModel.ChangePass(malina.CurrentAccount.Email, cryptedOldPass, cryptedNewPass)

			if ok {
				malina.Status = 0
				message = lang.T("password success changed")
			} else {
				malina.Status = 406
				malina.Result["error"] = lang.T(res)
			}
		}
	}
	s.view.ChangePassword(malina,message,w,r)
}

