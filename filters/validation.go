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
 */ package filters

import (
	"net/http"
	"strconv"
	"regexp"
	"github.com/ilyaran/Malina/lang"
	"fmt"
	"strings"
	"path"
	"net/url"
	"github.com/ilyaran/Malina/berry"
	"unicode"
	"github.com/ilyaran/Malina/app"
	"github.com/haisum/recaptcha"
)

var(	
	ErrorCode int64 = 406
)

func IsUint(malina *berry.Malina,isRequired bool, keyName string,length int,r *http.Request) (int64,string) {
	var id string
	if r.Method == "GET" {
		id = r.URL.Query().Get(keyName)
	}else if r.Method == "POST" {
		id = r.FormValue(keyName)
	}
	if id == "" {
		if isRequired {
			malina.Status = ErrorCode
			malina.Result[keyName] = lang.T("required")
		}
		return 0,"0"
	}
	if len(id) > length {
		malina.Status = ErrorCode
		malina.Result[keyName] = lang.T(`length no greater then:`)+strconv.Itoa(length)
		return 0,"0"
	}
	matched, _ := regexp.MatchString(`^[0-9]+$`, id)
	if !matched{
		malina.Status = ErrorCode
		malina.Result[keyName] = lang.T(`invalid`)
		return 0,"0"
	}
	idInt64,err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0,"0"
	}
	return idInt64,id
}
func CheckBox(malina *berry.Malina,keyName string,r *http.Request) bool {
	if r.FormValue(keyName) != "" {
		return true
	}
	return false
}

func IsFloat64(malina *berry.Malina,isRequired bool, keyName string,length,lengthCent int,r *http.Request) (float64,string){
	var formData string
	if r.Method == "GET" {
		formData = r.URL.Query().Get(keyName)
	}else if r.Method == "POST" {
		formData = r.FormValue(keyName)
	}
	if formData == "" {
		if isRequired {
			malina.Status = ErrorCode
			malina.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return 0.00,"0.00"
	}
	if len(formData) > length {
		malina.Status = ErrorCode
		malina.Result[keyName] = lang.T(`length no greater then `+strconv.Itoa(length))
		return 0.00,"0.00"
	}
	matched, _ := regexp.MatchString (`^[0-9]{1,20}(\.[0-9]{0,2})?$`, formData)//(`^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$`,formData)//("[0-9]{1,20}(\\.[0-9]{0,2})?", formData)
	if !matched{
		malina.Status = ErrorCode
		malina.Result[keyName] = lang.T(`invalid`)
		return 0.00,"0.00"
	}
	resFloat64,err := strconv.ParseFloat(formData, 64)
	if err != nil {
		malina.Status = ErrorCode
		malina.Result[keyName] = lang.T(`error parse float`)
		return 0.00,"0.00"
	}
	return resFloat64,formData
}

func IsValidText(malina *berry.Malina,isRquired bool,keyName string, isToLower bool,  minLength int, maxLength int, pattern,allow_symbols string, r *http.Request)(string){
	var field string
	if isToLower{
		field = strings.ToLower(r.FormValue(keyName))
	}else {
		field = r.FormValue(keyName)
	}
	if field == "" {
		if isRquired{
			malina.Status = ErrorCode
			malina.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return ""
	}
	if len(field) < minLength {
		malina.Status = ErrorCode
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation min length"),keyName, minLength)
		return ""
	}
	if len(field) > maxLength {
		malina.Status = ErrorCode
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation max length"),keyName, maxLength)
		return ""
	}
	match,_ := regexp.MatchString(pattern,field)
	if !match {
		malina.Status = ErrorCode
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation text field"),keyName,allow_symbols)
		return ""
	}
	return field
}
func ImgUrls(malina *berry.Malina,returnPath bool,keyName string, r *http.Request)(string,string,string){
	if r.FormValue(keyName) != ""{
		var imgAllUrls,imgUrls,imgNames = "","",""
		for k,v := range strings.Split(r.FormValue(keyName),"|"){
			if path.IsAbs(v) {
				if returnPath{
					imgAllUrls += ",'" + v + "'"
				}else {
					imgAllUrls += ",'" + path.Base(v) + "'"
				}
				imgNames   += ",'" + path.Base(v) + "'"
			}else {
				base, err := url.Parse(v)
				if err == nil && base.IsAbs(){
					imgAllUrls += ",'" + v + "'"
					imgUrls    += ",'" + v + "'"
				}else {
					malina.Status = ErrorCode
					malina.Result[keyName] = lang.T("invalid image url :")+strconv.Itoa(k+1)
					break
				}
			}
		}
		if imgAllUrls != ``{
			imgAllUrls =  imgAllUrls[1:]
			if imgUrls != ``{
				imgUrls =  imgUrls[1:]
			}
			if imgNames != ``{
				imgNames =  imgNames[1:]
			}
			return imgAllUrls,imgUrls,imgNames
		}
	}
	return "","",""
}


func IsEmail(malina *berry.Malina,isRequired bool, keyName string,r *http.Request)(string){
	var email = r.FormValue(keyName)
	if email != "" {
		if len(email)>255 {
			malina.Status = 406
			malina.Result[keyName] = lang.T(`length no greater then `)+`255`
			return ""
		}
		match, _ := regexp.MatchString(app.Pattern_email, email)
		if match {
			return strings.ToLower(email)
		}else {
			malina.Status = 406
			malina.Result[keyName] = keyName+" "+lang.T(`invalid`)
			return ""
		}
	}else if isRequired{
		malina.Status = 406
		malina.Result[keyName] = keyName+" "+lang.T("required")
		return ""
	}
	return ""
}

func IsPhone(malina *berry.Malina,isRequired bool,keyName string, r *http.Request)(string){
	var phone = r.FormValue(keyName)
	if phone != "" {
		match, _ := regexp.MatchString(app.Pattern_phone, phone)
		if match {
			return phone
		}else {
			malina.Status = 406
			malina.Result[keyName] = lang.T(`invalid `)+keyName + `, example: +77058436633`
			return ""
		}
	}else if isRequired {
		malina.Status = 406
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		return ""
	}
	return ""
}

func PasswordValid(malina *berry.Malina,isRequired bool,keyName string,hasNumber,hasUpper,hasSpecial bool,r *http.Request)(string){ //(minLength, number, upper, special bool) {
	var number, upper, special bool
	var letters int = 0
	for _, s := range r.FormValue(keyName) {
		switch {
		case unicode.IsNumber(s):
			number = true
			letters++
		case unicode.IsUpper(s):
			upper = true
			letters++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
			letters++
		case unicode.IsLetter(s) || s == ' ':
			letters++
			//default:
			//return false, false, false, false
		}
	}
	if letters == 0 {
		if isRequired {
			malina.Status = 406
			malina.Result[keyName] = lang.T("validation required")
		}
		return ""
	}
	if hasNumber && !number{
		malina.Status = 406
		malina.Result[keyName] = keyName + " field must contain number"
		return ""
	}
	if hasUpper && !upper{
		malina.Status = 406
		malina.Result[keyName] = keyName + " field must contain upper"
		return ""
	}
	if hasSpecial && !special{
		malina.Status = 406
		malina.Result[keyName] = keyName + " field must contain special"
		return ""
	}
	if letters < app.PasswordMinLength {
		malina.Status = 406
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation min length"),lang.T(keyName), app.PasswordMinLength)
		return ""
	}
	if letters > app.PasswordMaxLength {
		malina.Status = 406
		malina.Result[keyName] = fmt.Sprintf(lang.T("validation max length"),lang.T(keyName), app.PasswordMaxLength)
		return ""
	}
	return r.FormValue(keyName)
}
func ConfirmPasswordValid(malina *berry.Malina,keyName,password string, r *http.Request){
	if r.FormValue(keyName) != password {
		malina.Status = 406
		malina.Result[keyName] =  lang.T("auth_confirm_pass_valid")
	}
}


func IdsList(malina *berry.Malina,keyName string,length int, isRequired bool,r *http.Request)string{
	var ids = r.FormValue(keyName)
	if ids != ""{
		if length < 1 {
			match, _ := regexp.MatchString(`^([0-9]+)(\,[0-9]+)*?$`, ids)
			if match {
				return ids
			} else {
				malina.Status = 406
				malina.Result[keyName] = lang.T("invalid ids")
			}
		}else if len(ids) < length {
			match, _ := regexp.MatchString(`^([0-9]+)(\,[0-9]+){0,`+strconv.Itoa(length-1)+`}$`, ids)
			if match {
				return ids
			} else {
				malina.Status = 406
				malina.Result[keyName] = lang.T("invalid ids")
			}
		}else{
			malina.Status = 406
			malina.Result[keyName] = fmt.Sprintf(lang.T("too long"),keyName)
		}
	}else if isRequired {
		malina.Status = 406
		malina.Result[keyName] = fmt.Sprintf(lang.T("required"),keyName)
	}
	return ""
}


func RecaptchaValid(malina *berry.Malina, r *http.Request) {
	recap := recaptcha.R{Secret: app.Recaptcha_secret}
	malina.ChannelBool <- recap.Verify(*r)
}