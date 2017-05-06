package library

import (
	"net/http"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
	"Malina/config"
	"Malina/language"
	"github.com/haisum/recaptcha"
	"strings"
	"path"
)

var VALIDATION *Validation
type Validation struct {
	Status int
	Result map[string]string
}

func (s *Validation)IsValidText(isRquired bool,keyName string, isToLower bool,  minLength int, maxLength int, pattern,allow_symbols string, r *http.Request)string{
	var field string
	if isToLower{
		field = strings.ToLower(r.FormValue(keyName))
	}else {
		field = r.FormValue(keyName)
	}
	if field == "" {
		if isRquired{
			s.Status = 100
			s.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return ""
	}
	if len(field) < minLength {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation min length"),keyName, minLength)
		return ""
	}
	if len(field) > maxLength {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation max length"),keyName, maxLength)
		return ""
	}
	match,_ := regexp.MatchString(pattern,field)
	if !match {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation text field"),keyName,allow_symbols)
		return ""
	}
	return field
}

func (s *Validation)IsInt64(isRequired bool, keyName string,length int,r *http.Request) (int64,string) {
	var id string
	if r.Method == "GET" {
		id = r.URL.Query().Get(keyName)
	}else if r.Method == "POST" {
		id = r.FormValue(keyName)
	}
	if id == "" {
		if isRequired {
			s.Status = 100
			s.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return 0,"0"
	}
	if len(id) > length {
		s.Status = 100
		s.Result[keyName] = lang.T(`length no greater then `+strconv.Itoa(length))
		return 0,"0"
	}
	matched, _ := regexp.MatchString(`^[0-9]+$`, id)
	if !matched{
		s.Status = 100
		s.Result[keyName] = lang.T(`invalid`)
		return 0,"0"
	}
	idInt64,err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0,"0"
	}
	return idInt64,id

}

func (s *Validation)CheckBox(keyName string,r *http.Request) bool {
	if r.FormValue(keyName) != "" {
		return true
	}
	return false
}

func (s *Validation)IsFloat64(isRequired bool, keyName string,length,lengthCent int,r *http.Request) (float64,string ){
	var formData string
	if r.Method == "GET" {
		formData = r.URL.Query().Get(keyName)
	}else if r.Method == "POST" {
		formData = r.FormValue(keyName)
	}
	if formData == "" {
		if isRequired {
			s.Status = 100
			s.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return 0,"0"
	}
	if len(formData) > length {
		s.Status = 100
		s.Result[keyName] = lang.T(`length no greater then `+strconv.Itoa(length))
		return 0,"0"
	}
	matched, _ := regexp.MatchString (`^[0-9]{1,20}(\.[0-9]{0,2})?$`, formData)//(`^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$`,formData)//("[0-9]{1,20}(\\.[0-9]{0,2})?", formData)
	if !matched{
		s.Status = 100
		s.Result[keyName] = lang.T(`invalid`)
		return 0,"0"
	}
	resFloat64,err := strconv.ParseFloat(formData, 64)
	if err != nil {
		s.Status = 100
		s.Result[keyName] = lang.T(`error parse float`)
		return 0,"0"
	}
	return resFloat64,formData
}
func (s *Validation)ImgUrls(keyName string, r *http.Request)(string){
	var imgUrls = r.FormValue(keyName)
	if imgUrls != ""{
		imgArray := strings.Split(imgUrls,"|")
		imgUrls = ""
		for k,v := range imgArray{
			if !path.IsAbs(v){
				s.Status = 100
				s.Result[keyName] = lang.T("invalid image path has an image number ")+strconv.Itoa(k+1)
				break
			}else {
				imgUrls += ",'" + v + "'"
			}
		}
		if imgUrls != ``{
			return imgUrls[1:]
		}
	}
	return ""
}
func (s *Validation)ImgIds(keyName string, isRequired bool,r *http.Request)([]int64){
	var imgIds = r.FormValue(keyName)
	if imgIds != ""{
		var idsStr = strings.Split(imgIds,",") // 12,58,47,14, ...
		var idsInt64 = []int64{}
		var num = len(idsStr)
		/*if num > app.Image_upload_num(){
			num = app.Image_upload_num()
		}*/
		for n:=0; n < num; n++{
			id, err := strconv.ParseInt(idsStr[n],10,64)
			if err != nil{
				s.Status = 100
				s.Result[keyName] = fmt.Sprintf(lang.T("validation integer"),keyName)
			}else {
				idsInt64 = append(idsInt64,id)
			}
		}
		return idsInt64
	}else if isRequired {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		return nil
	}
	return nil
}

func (s *Validation)RecaptchaValid(recaptchaChannel chan bool, r *http.Request) {
	recap := recaptcha.R{Secret: app.Recaptcha_secret()}
	recaptchaChannel <- recap.Verify(*r)
}

func (s *Validation)IsEmail(isRequired bool, r *http.Request)(string){
	var email = r.FormValue("email")
	if email != "" {
		match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)
		if match {
			return strings.ToLower(email)
		}else {
			s.Status = 100
			s.Result["email"] = lang.T(`validation email`)
			return ""
		}
	}else if isRequired{
		s.Status = 100
		s.Result["email"] = lang.T("auth_email_required")
		return ""
	}
	return ""
}
func (s *Validation)IsPhone(isRequired bool, r *http.Request)(string){
	var phone = r.FormValue("phone")
	if phone != "" {
		match, _ := regexp.MatchString(`^\+[0-9]{10,16}$`, phone)
		if match {
			return phone
		}else {
			s.Status = 100
			s.Result["phone"] = lang.T(`validation phone`)
			return ""
		}
	}else if isRequired {
		s.Status = 100
		s.Result["phone"] = fmt.Sprintf(lang.T("validation required"),"phone")
		return ""
	}
	return ""
}
func (s *Validation)IsEmailOrNick(isRequired bool, r *http.Request)(string,string){
	var email_name = r.FormValue("email_nick")
	if email_name != "" {
		match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email_name)
		if match {
			return email_name,""
		}
		match, _ = regexp.MatchString(`^[a-z0-9-_]+$`, email_name)
		if match {
			return "",email_name
		}
		s.Status = 100
		s.Result["email_nick"] = lang.T(`auth_nick_or_email_valid`)
		return "",""
	}else if isRequired{
		s.Status = 100
		s.Result["email_nick"] = lang.T("auth_nick_or_email_required")
		return "",""
	}
	return "",""
}
func (s *Validation)IsEmailOrNickOrPhone(isRequired bool, r *http.Request)(string,string,string){
	var email_nick_phone = r.FormValue("email_nick_phone")
	if email_nick_phone != "" {
		match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email_nick_phone)
		if match {
			return email_nick_phone,"",""
		}
		match, _ = regexp.MatchString(`^[a-z0-9-_]+$`, email_nick_phone)
		if match {
			return "",email_nick_phone,""
		}
		match, _ = regexp.MatchString(`^\+[0-9]{10,16}$`, email_nick_phone)
		if match {
			return "","",email_nick_phone
		}
		s.Status = 100
		s.Result["email_nick_phone"] = lang.T(`auth nick or email or phone invalid`)
		return "","",""
	}else if isRequired{
		s.Status = 100
		s.Result["email_nick_phone"] = lang.T("auth nick or email or phone required")
		return "","",""
	}
	return "","",""
}
func (s *Validation)PasswordValid(isRequired bool,keyName string,hasNumber,hasUpper,hasSpecial bool,r *http.Request)(string){ //(minLength, number, upper, special bool) {
	var password = r.FormValue(keyName)
	var number, upper, special bool
	var letters int = 0
	for _, s := range password {
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
			s.Status = 100
			s.Result[keyName] = fmt.Sprintf(lang.T("validation required"),keyName)
		}
		return ""
	}
	if hasNumber && !number{
		s.Status = 100
		s.Result[keyName] = keyName + " field must contain number"
		return ""
	}
	if hasUpper && !upper{
		s.Status = 100
		s.Result[keyName] = keyName + " field must contain upper"
		return ""
	}
	if hasSpecial && !special{
		s.Status = 100
		s.Result[keyName] = keyName + " field must contain special"
		return ""
	}
	if letters < app.PasswordMinLength() {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation min length"),lang.T(keyName),app.PasswordMinLength())
		return ""
	}
	if letters > app.PasswordMaxLength() {
		s.Status = 100
		s.Result[keyName] = fmt.Sprintf(lang.T("validation max length"),lang.T(keyName),app.PasswordMaxLength())
		return ""
	}
	return password
}
func (s *Validation)ConfirmPasswordValid(password string, r *http.Request){
	//confirm password validation
	var confirm_password = r.FormValue("confirm_password")
	if confirm_password != password {
		s.Status = 100
		s.Result["confirm_password"] =  lang.T("auth_confirm_pass_valid")
	}
}
func (s *Validation)IAgree(r *http.Request) {
	if r.FormValue("i_agree") == "" {
		s.Status = 100
		s.Result["i_agree"] =  lang.T("auth_i_agree")
	}
}


