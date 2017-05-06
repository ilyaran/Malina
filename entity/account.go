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
	Id           int64 	`json:"id"`
	Position     *Position 	`json:"position"`
	Email        string	`json:"email"`
	Password     string	`json:"password"`
	Phone        string	`json:"phone"`
	Provider     string	`json:"provider"`
	Token        string	`json:"token"`
	Ban          int	`json:"banned"`
	Ban_reason   string	`json:"ban_reason"`
	Newpass      string	`json:"newpass"`
	Newpass_key  string	`json:"newpass_key"`
	Newpass_time time.Time	`json:"newpass_time"`
	Last_ip      string	`json:"last_ip"`
	Last_logged  time.Time	`json:"last_logged"`
	Created      time.Time	`json:"created"`
	Updated      time.Time	`json:"updated"`

	Nick         string	`json:"nick"`
	First_name   string	`json:"first_name"`
	Last_name    string	`json:"last_name"`
	Birth        time.Time	`json:"birth"`
	State        string	`json:"state"`
	City         string	`json:"city"`
	Skype        string	`json:"skype"`
	Steam_id     string	`json:"steam_id"`
	Balance      int64	`json:"balance"`

	Img          string	`json:"img"`
}
func (this *Account) Exec() []interface{} {
	return []interface{}{
		this.Email,
		this.Nick,
		this.Password,
		this.Ban,
		this.Ban_reason,
		this.Position.GetId(),
	}
}
func (this *Account) ExecWithId() []interface{} {
	if this.Newpass != "" {
		return []interface{}{
			this.Id,
			this.Email,
			this.Nick,
			this.Ban,

			this.Ban_reason,
			this.Position.GetId(),
			this.Password,
		}
	}
	return []interface{}{
		this.Id,
		this.Email,
		this.Nick,
		this.Ban,
		this.Ban_reason,
		this.Position.GetId(),
	}
}

func (s *Account)GetId() int64 {return s.Id}
func (s *Account)GetPosition() *Position {return s.Position}
func (s *Account)GetEmail() string {return s.Email}
func (s *Account)GetNick()string  {return s.Nick}
func (s *Account)GetPassword() string {return s.Password}
func (s *Account)GetNew_password() string {return s.Newpass}
func (s *Account)GetNew_pass_key() string {return s.Newpass_key}
func (s *Account)GetNew_pass_time() time.Time {return s.Newpass_time}
func (s *Account)GetPhone() string {return s.Phone}
func (s *Account)GetProvider() string {return s.Provider}
func (s *Account)GetToken() string {return s.Token}
func (s *Account)GetBanned() int {return s.Ban
}
func (s *Account)GetBan_reason()string  {return s.Ban_reason}
func (s *Account)GetLast_ip() string {return s.Last_ip}
func (s *Account)GetLast_logged() time.Time {return s.Last_logged}

func (s *Account)GetCreated() time.Time {return s.Created}
func (s *Account)GetUpdated() time.Time {return s.Updated}
func (s *Account)GetFirst_name() string {return s.First_name}
func (s *Account)GetLast_name() string {return s.Last_name}
func (s *Account)GetBirth() time.Time {return s.Birth}
func (s *Account)GetState() string {return s.State}
func (s *Account)GetCity() string {return s.City}
func (s *Account)GetSkype() string {return s.Skype}
func (s *Account)GetSteam_id() string {return s.Steam_id}
func (s *Account)GetBalance() int64 {return s.Balance}
func (s *Account)GetImg() string {return s.Img}

func (s *Account)SetId(id int64) {s.Id=id}
func (s *Account)SetPass(pass string) {s.Password=pass}
func (s *Account)SetNewPass(newpass string) {s.Newpass=newpass}

func AccountScanRow(row *sql.Row) *Account {
	s := &Account{Position:&Position{}}
	err := row.Scan(
		&s.Id,&s.Email,&s.Nick, &s.Password, &s.Phone,
		&s.First_name,&s.Last_name,&s.Img, &s.Provider,
		&s.Token,&s.Ban,&s.Ban_reason,&s.Newpass,
		&s.Newpass_key,&s.Newpass_time,&s.Last_ip,&s.Last_logged,
		&s.Created,&s.Updated,&s.Birth,&s.State,&s.City,
		&s.Skype,&s.Steam_id,&s.Position.Id,&s.Position.Title,&s.Balance)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}

	return s
}
func AccountScanRows(rows *sql.Rows) *Account {
	s := &Account{Position:&Position{}}
	err := rows.Scan(
		&s.Id,&s.Email,&s.Nick, &s.Password, &s.Phone,
		&s.First_name,&s.Last_name,&s.Img, &s.Provider,
		&s.Token,&s.Ban,&s.Ban_reason,&s.Newpass,
		&s.Newpass_key,&s.Newpass_time,&s.Last_ip,&s.Last_logged,
		&s.Created,&s.Updated,&s.Birth,&s.State,&s.City,
		&s.Skype,&s.Steam_id,&s.Position.Id,&s.Position.Title,&s.Balance)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}

/*func NewAccount(id int64, email, name, newpass string) *Account {
	return &Account{Id:id, Email:email, Nick:name, Newpass:newpass, Position:&Position{Id:int64(1),Title:"user"}}
}*/
func NewAccount(id int64, email string, position int64, nick, ban_reason, newpass string, ban int) *Account {
	return &Account{Id:id, Email:email, Nick:nick, Newpass:newpass,Ban_reason:ban_reason,Ban:ban,Position:&Position{Id:position}}
}








