/**
 * Authorization model class.  Malina eCommerce application
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
package model

import (
	"database/sql"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/config"
	"strconv"
)

var AuthModel = new(authModel)
type authModel struct{
	Where string
	Query string
	All int64
}
// param must have format example: CheckDetails("email = ?","ilyaran@mail.ru")
func (this *authModel)CheckDetails(param, value string){
	this.Query = `SELECT count(*) FROM account WHERE `+param+` LIMIT 1`
	err := Crud.GetRow(this.Query, []interface{}{value}).Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = -1
	}
	if err != nil {
		this.All = -2
	}
}
func (this *authModel) AddActivation(email,pass,activation_key, ip_address string) bool {
	this.Query = `
		INSERT INTO activation
		(activation_email,activation_password,activation_key,activation_last_ip)
		VALUES
		($1,$2,$3,$4)
		RETURNING activation_id`
	//fmt.Println(this.Query)
	if Crud.Insert(this.Query, []interface{}{email,pass,activation_key,ip_address}) > 0 {
		return true
	}
	return false
}

func (s *authModel)GetActivation(activation_key string) *entity.Account{

	var id int64
	var user_name sql.NullString
	var user_password string
	var user_email sql.NullString

	err := Crud.GetRow(`
	WITH t AS (
		DELETE FROM activation WHERE
		date_part('epoch',CURRENT_TIMESTAMP)::bigint - activation_created > `+strconv.FormatInt(app.Email_activation_expire(),10)+`
	)
	SELECT activation_id, coalesce(activation_nick,''), activation_password, coalesce(activation_email,'')
	FROM activation WHERE activation_key = $1 LIMIT 1
	`, []interface{}{activation_key}).Scan(&id, &user_name, &user_password, &user_email)

	if err == sql.ErrNoRows{
		//panic(err)
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	go s.RemoveFromActivateTable(id)
	return entity.NewAccount(user_email.String,user_password)
}


func (this *authModel) RemoveFromActivateTable(id int64) {
	this.Query = `DELETE FROM activation WHERE activation_id = $1`
	Crud.Delete(this.Query, []interface{}{id})
}
func (this *authModel)AddAccount(q string,exec []interface{})int64{
	this.Query = `INSERT INTO account `+q +` RETURNING account_id`
	return Crud.Insert(this.Query, exec)
}

func (this *authModel)DeleteSession(cookieValue string) bool {
	this.Query = `DELETE FROM session WHERE session_id = $1`
	if Crud.Delete(this.Query, []interface{}{cookieValue}) > 0 {
		return true
	}
	return false
}
func (this *authModel) CountLoginAttempts(ip_address string) {
	this.Query = `SELECT COUNT(*) FROM attempt WHERE attempt_ip = $1`
	err := Crud.GetRow(this.Query, []interface{}{ip_address}).Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = 0
	} else if err != nil {
		panic(err)
		this.All = -1
	}
}
func (this *authModel) IncreaseLoginAttempts(ip_address string) bool {
	this.Query = `
	WITH t AS (
		DELETE FROM attempt
		WHERE attempt_time + ` + strconv.Itoa(app.Login_attempts_period()) + ` < date_part('epoch',CURRENT_TIMESTAMP)::bigint
	)
	INSERT INTO attempt (attempt_ip)VALUES($1) RETURNING attempt_id`
	if Crud.Insert(this.Query, []interface{}{ip_address}) > 0 {
		return true
	}
	return false
}
func (this *authModel) ClearLoginAttempts(ip_address string) bool {
	this.Query = `DELETE FROM attempt WHERE attempt_ip = $1`
	if Crud.Delete(this.Query, []interface{}{ip_address}) > 0 {
		return true
	}
	return false
}
func (this *authModel)ForgotPass(email, cryptedNewPass string) (string, bool) {
	this.Query = `
		UPDATE account SET account_password = $1 WHERE account_email = $2
	`
	if Crud.Update(this.Query, []interface{}{cryptedNewPass, email}) > 0 {
		return email, true
	}
	return "", false
}
func (this *authModel) DeleteUserByEmailPassword(email,password string) {
	var exec []interface{}
	if password != ""{
		exec = []interface{}{email,password}
		this.Query = `DELETE FROM account WHERE account_email = $1 AND account_password = $2`
	}else {
		exec = []interface{}{email}
		this.Query = `DELETE FROM account WHERE account_email = $1`
	}
	Crud.Delete(this.Query, exec)
}

func (this *authModel) ChangePass(email, cryptedOldPass, cryptedNewPass string) (string, bool) {
	var userId int64
	this.Query = `SELECT account_id FROM account
	WHERE account_email = $1 AND account_password = $2 LIMIT 1`

	err := Crud.GetRow(this.Query, []interface{}{email, cryptedOldPass}).Scan(&userId)
	if err == sql.ErrNoRows {
		return "no_rows", false
	} else if err != nil {
		panic(err)
		return "server_error", false
	} else {
		this.Query = `
		UPDATE account SET account_password = $1
		WHERE account_email = $2 AND account_password = $3
		`
		if Crud.Update(this.Query, []interface{}{cryptedNewPass, email, cryptedOldPass}) > 0 {
			return ``, true
		}
	}
	return ``, false
}
func (this *authModel)GetAccountByEmail(email string)*entity.Account{
	this.Query = `
	SELECT `+FieldsAccount+`
	FROM account
	LEFT JOIN position ON position_id = account_position
	WHERE account_email = $1 LIMIT 1`
	row := Crud.GetRow(this.Query, []interface{}{email})
	return entity.AccountScan(row,nil)
}

func (this *authModel) CheckIsAllreadySentToActivation(email,nick,phone string) bool {
	var exec []interface{}
	if email != "" {
		this.Where = "activation_email = $1"
		exec = []interface{}{email}
	}else if nick != "" {
		this.Where = "activation_nick = $1"
		exec = []interface{}{nick}
	}
	this.Query = `SELECT count(*) FROM activation WHERE ` + this.Where

	err := Crud.GetRow(this.Query, exec).Scan(&this.All)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(err)
		return false
	}
	return this.All > 0
}

