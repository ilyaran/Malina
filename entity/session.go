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

import (
	"database/sql"
)

type Session struct {
	Id         	string
	Ip_address 	string

	AccountId      	int64

	Email    	string
	Nick     	string
	Phone    	string

	Position 	int64

	Provider     	string
	Token    	string

	Data  		string
	Data1 		string
	Data2 		string

	User_agent 	string

	err error
}

func (s *Session) SetId(id string) { s.Id = id }

func (s *Session) SetData(data string)  { s.Data = data }
func (s *Session) SetData1(data string) { s.Data1 = data }
func (s *Session) SetData2(data string) { s.Data2 = data }

func (s *Session) GetId() string         { return s.Id }
func (s *Session) GetIp_address() string { return s.Ip_address }
/*func (s *Session) GetAccount() *Account {
	return &Account{Id:s.AccountId,Email:s.Email,Nick:s.Nick,Phone:s.Phone,Provider:s.Provider,Token:s.Token,
		Position:&Position{Id:s.Position}}
}*/
func (s *Session) GetEmail() string      { return s.Email}
func (s *Session) GetAccountId() int64   { return s.AccountId}
func (s *Session) GetPositionId() int64  { return s.Position}
func (s *Session) GetNick() string       { return s.Nick }
func (s *Session) GetPhone() string      { return s.Phone }
func (s *Session) GetProvider() string   { return s.Provider }
func (s *Session) GetToken() string      { return s.Token }
func (s *Session) GetData() string       { return s.Data }
func (s *Session) GetData1() string      { return s.Data1 }
func (s *Session) GetData2() string      { return s.Data2 }
func (s *Session) GetUser_agent() string { return s.User_agent }

func (s *Session) SetBaseDataToObject(id,ip,user_agent string) {
	s.Id = id
	s.Ip_address=ip
	s.User_agent=user_agent
	s.AccountId=0

	s.Email=``
	s.Nick=``
	s.Phone=``

	s.Position=0

	s.Provider=``
	s.Token=``

	s.Data =``
	s.Data1 =``
	s.Data2 =``

	s.err = nil
}
func (s *Session) ScanRow(row *sql.Row) bool {

	s.err = row.Scan(&s.Id, &s.Data, &s.Data1, &s.Data2,
		&s.AccountId,
		&s.Email,
		&s.Nick,
		&s.Phone,
		&s.Provider,
		&s.Token,
		&s.Position,
	)

	if s.err == sql.ErrNoRows {
		panic(s.err)
		return false
	}
	if s.err != nil {
		panic(s.err)
		return false
	}
	return true
}

func SessionScan(row *sql.Row, rows *sql.Rows) *Session {
	s := &Session{}
	var err error
	if row != nil {
		err = row.Scan(&s.Id, &s.Data, &s.Data1, &s.Data2,
			&s.AccountId,
			&s.Email,
			&s.Nick,
			&s.Phone,
			&s.Provider,
			&s.Token,
			&s.Position,
		)
	} else {
		err = rows.Scan(&s.Id, &s.Data, &s.Data1, &s.Data2,
			&s.AccountId,
			&s.Email,
			&s.Nick,
			&s.Phone,
			&s.Provider,
			&s.Token,
			&s.Position,
		)
	}
	if err == sql.ErrNoRows {
		//panic(err)
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}

	//fmt.Println(s)
	return s
}

func NewSession(id string, data,data1, data2 string,
		accountId int64, email,nick,phone string,
		provider,token string) *Session {
	return &Session{
		Id:id,Data:data,Data1:data1,Data2:data2,
		AccountId:accountId,Email:email,Nick:nick,Phone:phone,
		Provider:provider,Token:token}
}
