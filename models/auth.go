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
 */ package models

import (
	"database/sql"
	"strconv"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/app"
)

var AuthModel = new(authModel)
type authModel struct{

}
// param must have format example: CheckDetails("email = ?","ilyaran@mail.ru")
func (this *authModel)CheckDetails(param, value string)int64{
	var all int64
	query := `SELECT count(*) FROM account WHERE `+param+` LIMIT 1`
	err := CrudGeneral.DB.QueryRow(query,  value).Scan(&all)
	if err == sql.ErrNoRows {
		panic(err)
		return -1
	}
	if err != nil {
		panic(err)
		return -2
	}
	return all
}
func (this *authModel) AddActivation(email,pass,activation_key, ip_address string) bool {
	query := `
		INSERT INTO activation
		(activation_email,activation_password,activation_key,activation_last_ip)
		VALUES
		($1,$2,$3,$4)
		ON CONFLICT (activation_email) DO UPDATE
        SET activation_password = $2,
        activation_key = $3,
        activation_last_ip = $4,
        activation_created = date_part('epoch',CURRENT_TIMESTAMP)::bigint

		RETURNING activation_created`
	//fmt.Println(query)
	if CrudGeneral.Insert(query,  email,pass,activation_key,ip_address) > 0 {
		return true
	}
	return false
}

func (s *authModel)GetActivation(activation_key string) *entity.Account{

	var user_password string
	var user_email string

	err := CrudGeneral.DB.QueryRow(`
	WITH t AS (
		DELETE FROM activation WHERE
		date_part('epoch',CURRENT_TIMESTAMP)::bigint - activation_created > `+strconv.FormatInt(app.Email_activation_expire,10)+`
	)
	SELECT activation_email, activation_password
	FROM activation WHERE activation_key = $1 LIMIT 1
	`,  activation_key).Scan(&user_email,&user_password)

	if err == sql.ErrNoRows{
		//panic(err)
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	go s.RemoveFromActivateTable(user_email)
	return &entity.Account{Email:user_email,Password:user_password}
}


func (this *authModel) RemoveFromActivateTable(user_email string) {
	CrudGeneral.Delete(`DELETE FROM activation WHERE activation_email = $1`,  user_email)
}
func (this *authModel)AddAccount(q string,exec... interface{})int64{
	return CrudGeneral.Insert(`INSERT INTO account `+q +` RETURNING account_id`, exec...)
}

func (this *authModel)DeleteSession(cookieValue string) bool {
	if CrudGeneral.Delete(`DELETE FROM session WHERE session_id = $1`,  cookieValue) > 0 {
		return true
	}
	return false
}
func (this *authModel) CountLoginAttempts(ip_address string) int64{
	var all int64
	err := CrudGeneral.DB.QueryRow(`SELECT COUNT(*) FROM attempt WHERE attempt_ip = $1`,  ip_address).Scan(&all)
	if err == sql.ErrNoRows {
		all = 0
	} else if err != nil {
		panic(err)
		all = -1
	}
	return all
}
func (this *authModel) IncreaseLoginAttempts(ip_address string) bool {
	query := `
	WITH t AS (
		DELETE FROM attempt
		WHERE attempt_time + ` + strconv.Itoa(app.Login_attempts_period) + ` < date_part('epoch',CURRENT_TIMESTAMP)::bigint
	)
	INSERT INTO attempt (attempt_ip) VALUES ($1) RETURNING attempt_id`
	if CrudGeneral.Insert(query,  ip_address) > 0 {
		return true
	}
	return false
}
func (this *authModel) ClearLoginAttempts(ip_address string) bool {
	if CrudGeneral.Delete(`DELETE FROM attempt WHERE attempt_ip = $1`,  ip_address) > 0 {
		return true
	}
	return false
}
func (this *authModel)ForgotPass(email, cryptedNewPass string) (string, bool) {
	if CrudGeneral.Update(`
		UPDATE account SET account_password = $1 WHERE account_email = $2
	`, cryptedNewPass, email) > 0 {
		return email, true
	}
	return "", false
}
func (this *authModel) DeleteUserByEmailPassword(email,password string) {
	if password != ""{
		CrudGeneral.Delete(`DELETE FROM account WHERE account_email = $1 AND account_password = $2`, email,password)
		return
	}else {
		CrudGeneral.Delete(`DELETE FROM account WHERE account_email = $1`, email)
		return
	}
}

func (this *authModel) ChangePass(email, cryptedOldPass, cryptedNewPass string) (string, bool) {
	var userId int64
	query := `SELECT account_id FROM account
	WHERE account_email = $1 AND account_password = $2 LIMIT 1`

	err := CrudGeneral.DB.QueryRow(query,  email, cryptedOldPass).Scan(&userId)
	if err == sql.ErrNoRows {
		return "no_rows", false
	} else if err != nil {
		panic(err)
		return "server_error", false
	} else {
		query := `
		UPDATE account SET account_password = $1
		WHERE account_email = $2 AND account_password = $3
		`
		if CrudGeneral.Update(query,  cryptedNewPass, email, cryptedOldPass) > 0 {
			return ``, true
		}
	}
	return ``, false
}
func (this *authModel)GetAccountByEmail(id int64,email,sqlSelectFields string)*entity.Account{
	var query string
	var row *sql.Row
	if email != "" {
		query = `
	SELECT `+sqlSelectFields+`
	FROM account
	WHERE account_email = $1 LIMIT 1`
		row = CrudGeneral.DB.QueryRow(query,  email)
	}
	if id > 0 {
		query = `
	SELECT `+sqlSelectFields+`
	FROM account
	WHERE account_id = $1 LIMIT 1`
		row = CrudGeneral.DB.QueryRow(query,  id)
	}

	account:=&entity.Account{}
	if account.Scanning(row,nil) == 'o'{
		return account
	}
	return nil
}

func (this *authModel) CheckIsAllreadySentToActivation(email,nick,phone string) bool {
	var all int64
	var where string
	var exec []interface{}
	if email != "" {
		where = "activation_email = $1"
		exec =  []interface{}{email}
	}else if nick != "" {
		where = "activation_nick = $1"
		exec =  []interface{}{nick}
	}

	query := `SELECT count(*) FROM activation WHERE ` + where

	err := CrudGeneral.DB.QueryRow(query, exec).Scan(&all)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(err)
		return false
	}
	return all > 0
}



