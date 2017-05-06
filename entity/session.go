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
	Id         	string	`json:"id"`
	Ip_address 	string	`json:"ip_address"`
	Account_id 	int64	`json:"account_id"`

	Email      	string	`json:"email"`
	Nick       	string	`json:"nick"`

	Data       	string	`json:"data"`

	User_agent 	string	`json:"user_agent"`
	Is_flash   	bool	`json:"is_flash"`
	Position_id   	int64	`json:"position_id"`

	Position_title  string	`json:"position_title"`

	Permission_data  string	`json:"permission_data"`

	Phone 		string  `json:"phone"`
	Balance  	int64	`json:"balance"`
}

func (s *Session)GetId()string{return s.Id}
func (s *Session)GetIp_address()string{return s.Ip_address}
func (s *Session)GetData()string{return s.Data}
func (s *Session)GetAccount_id()int64{return s.Account_id}
func (s *Session)GetEmail()string{return s.Email}
func (s *Session)GetNick()string{return s.Nick}
func (s *Session)GetPermission()string{return s.Permission_data}
func (s *Session)GetBalance()int64{return s.Balance}
func (s *Session)GetPhone()string{return s.Phone}

func (s *Session)SetData(data string){s.Data=data}

func (s *Session) ExecWithId() []interface{} {
	return []interface{}{
		s.Id,
		s.Ip_address,
		s.Account_id,

		s.Email,
		s.Nick,
		s.Phone,

		s.Data,
		s.User_agent,
		s.Is_flash,

		s.Position_id,
		s.Balance}
}
func NewSession(account *Account,id,ip_address string,data,user_agent string,is_flash bool)*Session{
	return &Session{
		Id:id,
		Ip_address:ip_address,
		Account_id:account.GetId(),

		Email:account.GetEmail(),
		Nick:account.GetNick(),
		Data:data,

		User_agent:user_agent,
		Is_flash:is_flash,
		Position_id:account.GetPosition().GetId(),
		Phone:account.GetPhone(),
		Balance:account.GetBalance(),
	}
}
func NewUnauthSession(id,ip_address string,data,user_agent string)*Session{
	return &Session{
		Id:id,
		Ip_address:ip_address,
		Data:data,
		User_agent:user_agent}
}
func SessionScanRow(row *sql.Row)*Session{
	var s = &Session{}
	err := row.Scan(
		&s.Id,
		&s.Ip_address,
		&s.Account_id,

		&s.Email,
		&s.Nick,
		&s.Phone,

		&s.Data,
		&s.User_agent,
		&s.Is_flash,

		&s.Position_id,
		&s.Balance,
		&s.Position_title,
		&s.Permission_data,

		)

	if err == sql.ErrNoRows {return nil}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}
func SessionScanRows(rows *sql.Rows)*Session{
	var s = &Session{}
	err := rows.Scan(
		&s.Id,
		&s.Ip_address,
		&s.Account_id,

		&s.Email,
		&s.Nick,
		&s.Phone,

		&s.Data,
		&s.User_agent,
		&s.Is_flash,

		&s.Position_id,
		&s.Balance,
		&s.Position_title,
		&s.Permission_data,
	)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}

	return s
}
















