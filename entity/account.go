/**
 * Account entity class.  Malina eCommerce application
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
package entity

import (
	"time"
	"database/sql"
)

type Account struct {
	Id           int64         `json:"Id"`
	Position     *Position     `json:"position"`
	Nick         string        `json:"nick"`
	Email        string        `json:"email"`
	Password     string        `json:"password"`
	Phone        string        `json:"phone"`

	Provider     string        `json:"provider"`
	Token        string        `json:"token"`
	Ban          bool           `json:"ban"`

	Ban_reason   string        `json:"ban_reason"`
	Newpass      string        `json:"newpass"`
	Newpass_key  string        `json:"newpass_key"`

	Newpass_time time.Time     `json:"newpass_time"`
	Last_ip      string        `json:"last_ip"`
	Last_logged  time.Time     `json:"last_logged"`

	Created      time.Time     `json:"created"`
	Updated      time.Time     `json:"updated"`
}
func(s *Account)SetId(id int64){s.Id = id}
func(s *Account)SetPassword(pass string){s.Password = pass}
func(s *Account)SetNewPassword(pass string){s.Password = pass}

func(s *Account)GetId () int64 {return s.Id}
func(s *Account)GetPosition () *Position {return s.Position}
func(s *Account)GetEmail ()string {return s.Email}
func(s *Account)GetPassword () string {return s.Password}
func(s *Account)GetPhone () string {return s.Phone}
func(s *Account)GetProvider () string {return s.Provider}
func(s *Account)GetToken() string {return s.Token}
func(s *Account)GetBan () bool {return s.Ban}
func(s *Account)GetBan_reason () string{return s.Ban_reason}
func(s *Account)GetNewpass () string {return s.Newpass}
func(s *Account)GetNewpass_key () string {return s.Newpass_key}
func(s *Account)GetNewpass_time ()time.Time {return s.Newpass_time}
func(s *Account)GetLast_ip () string{return s.Last_ip}
func(s *Account)GetLast_logged () time.Time {return s.Last_logged}
func(s *Account)GetCreated () time.Time {return s.Created}
func(s *Account)GetUpdated () time.Time {return s.Updated}
func(s *Account)GetNick () string{return s.Nick}

func(s *Account)SetEmail (email string) {s.Email=email}
func(s *Account)SetNick (nick string) {s.Nick=nick}
func(s *Account)SetPhone (phone string) {s.Phone=phone}


func AccountScan(row *sql.Row,rows *sql.Rows) *Account {
	s := &Account{}
	var positionId 		int64

	var err error
	if row!=nil{
		err = row.Scan(
			&s.Id,
			&positionId,
			&s.Nick,

			&s.Password,
			&s.Email,
			&s.Ban,

			&s.Ban_reason,
			&s.Newpass,
			&s.Newpass_key,

			&s.Newpass_time,
			&s.Last_ip,
			&s.Last_logged,

			&s.Created,
			&s.Updated,
			&s.Phone,

			&s.Provider,
			&s.Token)
	}else {
		err = rows.Scan(
			&s.Id,
			&positionId,
			&s.Nick,

			&s.Password,
			&s.Email,
			&s.Ban,

			&s.Ban_reason,
			&s.Newpass,
			&s.Newpass_key,

			&s.Newpass_time,
			&s.Last_ip,
			&s.Last_logged,

			&s.Created,
			&s.Updated,
			&s.Phone,

			&s.Provider,
			&s.Token)
	}
	if err == sql.ErrNoRows {
		return nil
	}

	if positionId > 0 {
		s.Position = &Position{Id:positionId}
	}
	if err != nil {
		panic(err)
		return nil
	}

	return s
}

func NewAccount(Email string,NewPassword string) *Account {
	return &Account{Email:Email,Newpass:NewPassword}
}





