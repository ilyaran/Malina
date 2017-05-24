/**
 * Session entity class.  Malina eCommerce application
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

import "database/sql"

type Session struct {
	Id                 string		`json:"Id"`
	Ip_address         string		`json:"ip_address"`
	Account            *Account		`json:"account"`

	Data                  string		`json:"data"`
	Data1              string		`json:"data1"`
	Data2              string		`json:"data2"`
	Data3              string		`json:"data2"`

	User_agent         string		`json:"user_agent"`

}

func (s *Session)SetId(id string) {s.Id = id}

func (s *Session)SetData(data string) {s.Data = data}
func (s *Session)SetData1(data string) {s.Data1 = data}
func (s *Session)SetData2(data string) {s.Data2 = data}


func (s *Session)GetId() string {return s.Id}
func (s *Session)GetIp_address() string {return s.Ip_address}

func (s *Session)GetAccount() *Account {return s.Account}

func (s *Session)GetData() string {return s.Data}
func (s *Session)GetData1() string {return s.Data1}
func (s *Session)GetData2() string {return s.Data2}
func (s *Session)GetUser_agent() string {return s.User_agent}


func SessionScan(row *sql.Row, rows *sql.Rows) *Session {
	s := &Session{}
	var account_id int64
	var err error
	if row != nil {
		err = row.Scan(&s.Id,&account_id,&s.Data)
	}
	if row == nil {
		err = rows.Scan(&s.Id,&account_id,&s.Data)
	}

	if err == sql.ErrNoRows {
		//panic(err)
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	if account_id > 0 {
		s.Account = &Account{Id: account_id}
	}
	return s
}


func NewSession(id,ip string,account *Account,data,userAgent string) *Session{
	return &Session{
		Id:id,
		Ip_address:ip,
		Account:account,
		Data:data,
		User_agent:userAgent,
	}
}
