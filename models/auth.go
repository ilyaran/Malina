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
	"fmt"
	"database/sql"
	"Malina/entity"
	"Malina/config"
	"strconv"
)

var AuthModel = new(authModel)
type authModel struct{}

func (this *authModel)GetActivation(activation_key string) *entity.Account{

	var id int64
	var user_name string
	var user_password string
	var user_email string

	var exec = []interface{}{activation_key}
	var querySql = `
	SELECT id, user_name, user_password, user_email
	FROM delete_expired_sessions(`+fmt.Sprintf("%v",app.Session_expiration())+`), user_temp
	WHERE activation_key = $1 LIMIT 1`

	err := Crud.GetRow(querySql, exec).Scan(&id, &user_name, &user_password, &user_email)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
		return nil
	default:
		this.RemoveFromActivateTable(id)
		return entity.NewAccount(0,user_email,1,user_name,"",user_password,0)
	}
	return nil
}

func (this *authModel) RemoveFromActivateTable(id int64) {
	var exec = []interface{}{id}
	var querySql = `DELETE FROM user_temp WHERE id = $1`
	Crud.Delete(querySql, exec)
}

func (this *authModel)DeleteSession(cookieValue string) bool {
	var exec = []interface{}{cookieValue}
	var querySql = `DELETE FROM sessions WHERE id = $1`
	if Crud.Delete(querySql, exec) > 0 {
		return true
	}
	return false
}
func (this *authModel)ForgotPass(email, cryptedNewPass string) (string, bool) {
	var exec = []interface{}{cryptedNewPass, email}
	var querySql = `
		UPDATE users SET user_password = $1
		WHERE user_email = $2
	`
	if Crud.Update(querySql, exec) > 0 {
		return email, true
	}
	return "", false
}
func (this *authModel) DeleteUserByEmailPassword(email,password string) {
	var exec []interface{}
	var querySql string
	if password != ""{
		exec = []interface{}{email,password}
		querySql = `
		DELETE FROM users
		WHERE user_email = $1 AND user_password = $2
		`
	}else {
		exec = []interface{}{email}
		querySql = `
		DELETE FROM users
		WHERE user_email = $1
		`
	}
	Crud.Delete(querySql, exec)
}

func (this *authModel) ChangePass(email, cryptedOldPass, cryptedNewPass string) (string, bool) {
	var userId uint64
	var exec = []interface{}{email, cryptedOldPass}
	var querySql = `
	SELECT user_id FROM users
	WHERE user_email = $1 AND user_password = $2 LIMIT 1
	`
	err := Crud.GetRow(querySql, exec).Scan(&userId)
	if err == sql.ErrNoRows {
		return "no_rows", false
	} else if err != nil {
		panic(err)
		return "server_error", false
	} else {
		var exec = []interface{}{cryptedNewPass, email, cryptedOldPass}
		var querySql = `
		UPDATE users SET user_password = $1
		WHERE user_email = $2 AND user_password = $3
		`
		if Crud.Update(querySql, exec) > 0 {
			return "", true
		}
	}
	return "", false
}

func (this *authModel) AddActivation(email,nick,phone,pass,activation_key, ip_address string) bool {
	var exec = []interface{}{email,nick,phone,pass,activation_key,ip_address}
	var querySql = `
		INSERT INTO user_temp
		(user_email,user_name,user_phone,user_password,activation_key,last_ip)
		VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	if Crud.Insert(querySql, exec) > 0 {
		return true
	}
	return false
}

func (this *authModel) CountLoginAttempts(ip_address string) int {
	var exec = []interface{}{ip_address}
	var querySql = `SELECT COUNT(*) FROM login_attempts WHERE ip_address = $1`
	var count int
	err := Crud.GetRow(querySql, exec).Scan(&count)
	if err == sql.ErrNoRows {
		return -110
	} else if err != nil {
		panic(err)
		return -100
	}
	return count
}

func (this *authModel) IncreaseLoginAttempts(ip_address string) bool {
	var exec = []interface{}{ip_address}
	var querySql = `
		WITH t as (
			DELETE FROM login_attempts
			WHERE attempt_time + ` + strconv.Itoa(app.Login_attempts_period()) + ` < date_part('epoch',CURRENT_TIMESTAMP)::BIGINT
		)
		INSERT INTO login_attempts
		(ip_address) VALUES ( $1 )`
	if Crud.Insert(querySql, exec) > 0 {
		return true
	} //
	return false
}

func (this *authModel) ClearLoginAttempts(ip_address string) bool {
	var exec = []interface{}{ip_address}
	var querySql = `
	DELETE FROM login_attempts
	WHERE ip_address = $1
	`
	if Crud.Delete(querySql, exec) > 0 {
		return true
	}
	return false
}

func (this *authModel) CheckIsAllreadySentToActivation(email,nick,phone string) bool {
	var where string
	var exec []interface{}
	if email != "" {
		where = "user_email = $1"
		exec = []interface{}{email}
	}else if nick != "" {
		where = "user_name = $1"
		exec = []interface{}{nick}
	}else if nick != "" {
		where = "user_phone = $1"
		exec = []interface{}{phone}
	}
	var querySql = `SELECT count(*) FROM user_temp WHERE `+where
	var c int64 = 0
	err := Crud.GetRow(querySql, exec).Scan(&c)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(err)
		return false
	}
	return c > 0
}
/*

var querySql = `
		WITH t as (
			DELETE FROM login_attempts
			WHERE attempt_time +86400 < date_part('epoch',CURRENT_TIMESTAMP)::BIGINT
		)
		INSERT INTO login_attempts
		(ip_address)
		SELECT $1 FROM t
			WHERE (
				SELECT COUNT(*)
				from  login_attempts
				WHERE ip_address = $1
				) < $2`

 */
