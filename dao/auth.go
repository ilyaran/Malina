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
 */ package dao

import (
	"github.com/ilyaran/Malina/entity"
	"net/http"
	"regexp"
	"time"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"net/smtp"
	"net"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/models"
	"io"
	"crypto/rand"
	"github.com/ilyaran/Malina/caching"
	"encoding/base64"
)

var AuthDao *Auth
func AuthDaoInit(sqlSelectFieldsAccount string){
	AuthDao = &Auth{
		chars   		:    []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"),
		chars_pincode   : []byte("0123456789"),
		Model:            models.Session{SqlSelectFieldsAccount:sqlSelectFieldsAccount},
	}
}

type Auth struct {
	chars         []byte
	chars_pincode []byte
	err           error
	Model         models.Session
}

func (s *Auth) Authentication(w http.ResponseWriter, r *http.Request) (*entity.Account,string ){
	cookieToken := s.GetCookie(app.Cookie_name,r)
	if cookieToken != "" {
		account, isSessionExist := s.Model.Get(cookieToken,s.GetIP(r),r.Header.Get("User-Agent"))
		if isSessionExist {
			return account,cookieToken
		}
	}
	s.SetSession(nil,r,w)
	return nil,cookieToken
}

func (s *Auth) GetCookie(cookieName string,r *http.Request)string {
	cookie, errCookie := r.Cookie(cookieName)
	if errCookie == nil && regexp.MustCompile(`^[a-z\-A-Z0-9_]{30,512}$`).MatchString(cookie.Value) {
		return cookie.Value
	} else {
		return ""
	}
}

func (s *Auth) generateGoToken(r *http.Request) string {
	return s.Cryptcode(fmt.Sprintf("%v%v%v%v%v", s.GetIP(r),time.Now().UTC().UnixNano(), app.Crypt_salt,s.Rand_char(64,s.chars),time.Now().Sub(caching.T0).Nanoseconds()))
}

func (s *Auth) SetSession(account *entity.Account,r *http.Request, w http.ResponseWriter) {
	cookieToken := s.Model.Add(r,account,s.generateGoToken(r),s.GetIP(r),r.Header.Get("User-Agent"))
	if cookieToken != "" {
		s.setCookie(app.Cookie_name, cookieToken,w)
	}
}
func (s *Auth) setCookie(cookieName, cookieValue string,w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   app.Cookie_expiration,
	}
	if app.Cookie_expiration > 0 {
		cookie.Expires = time.Now().Add(time.Duration(app.Cookie_expiration))
	} else {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	http.SetCookie(w, &cookie)
}

func (s *Auth) DeleteSession(sessionId string,w http.ResponseWriter) {
	if s.Model.Del(sessionId) > 0 {
		deleteCookie := http.Cookie{
			Name:     app.Cookie_name,
			Value:    "none",
			Expires:  time.Now(),
			HttpOnly: true,
			Path:     "/",
			MaxAge:   1, // 0 - the zero means eternal
		}
		http.SetCookie(w, &deleteCookie)
	}
}

func (s *Auth) Cryptcode(text string) string {
	h := sha256.New()
	h.Write([]byte(text + app.Crypt_salt))
	return hex.EncodeToString(h.Sum(nil))
}
func (s *Auth) SendEmail(subject, to, body string) {
	from := app.Admin_email
	pass := app.Admin_email_pass
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

func (s *Auth) GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	/*if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return ""
	}*/
	return ip
}

func (s *Auth) Generate_pinCode() string {
	return s.Rand_char(app.GenereatedPincodeLen,s.chars_pincode)
}
func (s *Auth) Generate_password() string {
	return s.Rand_char(app.GenereatedPasswordLen,s.chars)
}

func (s *Auth) Rand_char(length int, chars []byte ) string {
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


func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func getToken()string{
	// Example: this will give us a 44 byte, base64 encoded output
	var token, err = GenerateRandomString(32)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
	}
	return token
}
