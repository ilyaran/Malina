package library

import (
	"io"
	"net/http"
	"fmt"
	"time"
	"regexp"
	"crypto/sha256"
	"encoding/hex"
	"net/smtp"
	"net"
	"Malina/models"
	"Malina/entity"
	"Malina/config"
	"crypto/rand"
)

var SESSION *Session
type Session struct {
	cookieString string
	SessionObj   *entity.Session
}
func (s *Session)GetSessionObj() *entity.Session {
	return s.SessionObj
}
func (s *Session) Authentication(w http.ResponseWriter, r *http.Request) {
	s.cookieString = s.GetCookie(app.Cookie_name(),r)
	if s.cookieString != "" {
		s.SessionObj = model.Session.Get(s.cookieString, s.GetIP(r))
		if s.SessionObj != nil {
			return
		}
	}
	s.SessionObj = s.SetUnauthSession("unauth",s.GetIP(r),r.Header.Get("User-Agent"),w)
}

func (s *Session) GetCookie(cookieName string,r *http.Request) string {
	cookie, errCookie := r.Cookie(cookieName)
	if errCookie == nil && regexp.MustCompile(`^[a-z\.\-\_0-9]{30,512}$`).MatchString(cookie.Value) {
		return cookie.Value
	}
	return ``
}
func (s *Session) SetUserDataToSession(sessionId,data string)bool{
	var session = model.Session.Get(sessionId, "noip")
	if session != nil {
		session.SetData(data)
		if model.Session.Update(session) > 0{
			return true
		}
	}
	return false
}
func (s *Session) SetSession(account *entity.Account,data string,is_flash bool,ip_address, user_agent string, w http.ResponseWriter)*entity.Session {
	var sessionId string = s.Cryptcode(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(), ip_address))
	var session = entity.NewSession(account,sessionId,ip_address,data,user_agent,is_flash)
	return s.addSessionToDb(session,w)
}
func (s *Session) SetUnauthSession(data string, ip_address, user_agent string, w http.ResponseWriter)*entity.Session {
	var sessionId string = s.Cryptcode(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(), ip_address))
	var session = entity.NewUnauthSession(sessionId,ip_address,data,user_agent)
	//fmt.Println(ip_address)
	return s.addSessionToDb(session,w)
}
func (s *Session) SetCookie(name,value string,httpOnly bool,path string,max_age int, w http.ResponseWriter){
	cookie := http.Cookie{
		Name: name,
		Value: value,
		HttpOnly: httpOnly,
		Path: path,
		MaxAge:max_age,
	}
	if app.Cookie_expiration() > 0 {
		cookie.Expires = time.Now().Add(time.Duration(max_age))
	} else {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	http.SetCookie(w, &cookie)
}
func (s *Session) addSessionToDb(session *entity.Session,w http.ResponseWriter) *entity.Session {
	if model.Session.Add(session){
		cookie := http.Cookie{
			Name: app.Cookie_name(),
			Value: session.GetId(),
			HttpOnly: true,
			Path: "/",
			MaxAge:app.Cookie_expiration(),
		}
		if app.Cookie_expiration() > 0 {
			cookie.Expires = time.Now().Add(time.Duration(app.Cookie_expiration()))
		} else {
			// Set it to the past to expire now.
			cookie.Expires = time.Unix(1, 0)
		}
		http.SetCookie(w, &cookie)
	}

	return session
}
func (s *Session) DeleteSession(sessionId string,w http.ResponseWriter) {
	if model.Session.Del(sessionId) {
		deleteCookie := http.Cookie{
			Name: app.Cookie_name(),
			Value: "none",
			Expires: time.Now(),
			HttpOnly: true,
			Path: "/",
			MaxAge:1, // 0 - the zero means eternal
		}
		http.SetCookie(w, &deleteCookie)
	}
}
func (s *Session) Cryptcode(text string) string {
	h := sha256.New()
	h.Write([]byte(text + app.Crypt_salt()))
	return hex.EncodeToString(h.Sum(nil))
}
func (s *Session) SendEmail(subject, to, body string) {
	from := app.Admin_email()
	pass := app.Admin_email_pass()
	//to := "foobarbazz@mailinator.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + " \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return
	}
}

func (s *Session) GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	/*if !IsIP(ip) {
		return ""
	}*/
	return ip
}
func (s *Session) IsIP(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}
func (s *Session) Generate_pinCode(length int) string {
	var StdChars = []byte("0123456789")
	return s.Rand_char(length, StdChars)
}
func (s *Session) Generate_password(length int) string {
	var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
	return s.Rand_char(length, StdChars)
}
func (s *Session) Rand_char(length int, chars []byte) string {
	new_pword := make([]byte, length)
	random_data := make([]byte, length + (length / 4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c % clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
	panic("unreachable")
}


