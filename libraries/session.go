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
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/config"
	"crypto/rand"
)

var SESSION *Session
type Session struct {

	CookieString string
	SessionObj   *entity.Session
	StdChars     []byte
}
func (s *Session)GetSessionObj() *entity.Session {
	return s.SessionObj
}
func (s *Session) Authentication(w http.ResponseWriter, r *http.Request) {
	s.GetCookie(app.Cookie_name(),r) //fmt.Println("Cookie: ",s.CookieString)
	if s.CookieString != `` {
		model.Session.Get(s.CookieString)
		if s.SessionObj.ScanRow(model.Session.Row) {
			if s.SessionObj.AccountId > 0 {
				go model.Crud.Update(`
				UPDATE account SET
					account_last_logged = now(),
					account_last_ip = $1
				WHERE account_id = $2
				`, []interface{}{s.GetIP(r), s.SessionObj.AccountId})
			}
			return
		}
	}
	s.SetCookieString()
	s.SessionObj.SetBaseDataToObject(s.CookieString,s.GetIP(r),r.Header.Get("User-Agent"))
	if model.Crud.Insert(`insert into session
		(session_id,session_ip,session_agent,session_data) values ($1,$2,$3,'unauth')
		returning 1`,
		[]interface{}{s.CookieString,s.SessionObj.Ip_address,s.SessionObj.User_agent}) > 0 {

		s.SetSessionCookie(w)
	}
}

func (s *Session) GetCookie(cookieName string,r *http.Request) {
	cookie, errCookie := r.Cookie(cookieName)
	if errCookie == nil && regexp.MustCompile(`^[a-zA-Z0-9]{30,255}$`).MatchString(cookie.Value) {
		s.CookieString = cookie.Value
	} else {
		s.CookieString = ``
	}
}
func (s *Session) SetSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name: app.Cookie_name(),
		Value: s.CookieString,
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

func (s *Session)SetCookieString(){s.CookieString = s.Cryptcode(fmt.Sprintf("%v%v%v", time.Now().UTC().UnixNano(), app.Crypt_salt(),s.Rand_char(12,s.StdChars)))}

func (s *Session) SetSession(accountId int64, data string, w http.ResponseWriter) {
	s.SetCookieString()
	if model.Crud.Insert(`
		INSERT INTO session
		(
			session_data,
			session_id,
			session_account,

			session_email,
			session_nick,
			session_phone,

			session_provider,
			session_token,
			session_position
		)
		SELECT $1,$2,$3,account_email,account_nick,account_phone,account_provider,account_token,account_position
		FROM account
		WHERE account_id = $3
		returning 1 `,
		[]interface{}{data, s.CookieString, accountId}) > 0 {
		s.SetSessionCookie(w)
	}
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
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return ""
	}
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

	return s.Rand_char(length, s.StdChars)
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


