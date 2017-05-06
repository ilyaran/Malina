/**
 * Account model class.  Malina eCommerce application
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
	"Malina/entity"
	"database/sql"
	"Malina/config"
)

var FieldsAccount = `
	account_id,
	account_email,
	account_nick,
	account_password,
	account_phone,
	account_fist_name,
	account_last_name,
	coalesce('`+app.Base_url()+app.Upload_path()+`' || account_img,'`+app.No_image()+`'),
	account_provider,
	account_token,
	account_banned,
	account_ban_reason,
	account_newpass,
	account_newpass_key,
	account_newpass_time,
	account_last_ip,
	account_last_logged,
	account_created,
	account_updated,
	account_birth,
	account_state,
	account_city,
	account_skype,
	account_steam_id,
	account_position,
	coalesce((SELECT position_title FROM position WHERE position_id = account.account_position),'user'),
	account_balance `
/*
date_trunc('hour', account_last_logged),
	date_trunc('hour', account_created),
	date_trunc('hour', account_updated),
	date_trunc('day', account_birth),
 */
var AccountModel *accountModel = new(accountModel)
type accountModel struct {}

func (this *accountModel)GetByEmailByNickByPhoneByPassword(email, nick, phone, password string)*entity.Account{
	var where string
	var exec []interface{}
	if email != "" {
		where = "account_email = $1 AND account_password = $2"
		exec = []interface{}{email,password}
	}else if nick != "" {
		where = "account_nick = $1 AND account_password = $2"
		exec = []interface{}{nick,password}
	}else if phone != "" {
		where = "account_phone = $1 AND account_password = $2"
		exec = []interface{}{nick,password}
	}
	if password == ""{
		if email != "" {
			where = "account_email = $1"
			exec = []interface{}{email}
		}else if nick != "" {
			where = "account_nick = $1"
			exec = []interface{}{nick}
		}else if phone != "" {
			where = "account_phone = $1"
			exec = []interface{}{nick}
		}else {
			return nil
		}
	}
	var querySql = `SELECT `+FieldsAccount+` FROM account WHERE `+where+` LIMIT 1`
	row := Crud.GetRow(querySql, exec)
	return entity.AccountScanRow(row)
}

func (this *accountModel)Get(id int64)*entity.Account{
	var querySql = `SELECT `+FieldsAccount+` FROM account
	WHERE account_id = $1 LIMIT 1`
	row := Crud.GetRow(querySql, []interface{}{id})
	return entity.AccountScanRow(row)
}

func (this *accountModel)Del(id int64)int64{
	var querySql = `
	WITH t AS (
		DELETE FROM sessions WHERE account_id = $1
	)
	DELETE FROM account WHERE account_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *accountModel)Add(account *entity.Account)int64{
	var querySql = `
	INSERT INTO account (
		account_email,
		account_nick,
		account_password,
		account_banned,
		account_ban_reason,
		account_position
	)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING account_id`
	return Crud.Insert(querySql, account.Exec())
}

func (this *accountModel)Edit(account *entity.Account)int64{
	var querySql = `
	WITH t AS (
		UPDATE sessions SET (position_id,position_title,permission_data)=
		($6,
		coalesce((SELECT position_title FROM position WHERE position_id = $6),'user'),
		coalesce((SELECT permission_data FROM permission WHERE permission_position = $6),'user')
		)
		WHERE account_id = $1
	)
	UPDATE account
	SET (
		account_email,
		account_nick,
		account_banned,
		account_ban_reason,
		account_position,
		account_password,
		account_updated

	) = ($2,$3,$4,$5,$6,$7,now())
	WHERE account_id = $1`
	return Crud.Update(querySql, account.ExecWithId())
}

func (this *accountModel) GetList(search,page,perPage string,order_by string) []*entity.Account {
	var where = this.getSearchSql(search)
	if order_by == "" {
		order_by = "account_last_logged DESC"
	}
	var querySql = `
	SELECT ` + FieldsAccount + `
	FROM account
	` + where + `
	ORDER BY ` + order_by + `
	LIMIT ` + perPage + ` OFFSET ` + page

	rows := Crud.GetRows(querySql, []interface{}{})
	defer rows.Close()

	var accountList = []*entity.Account{}

	for rows.Next() {
		accountList = append(accountList, entity.AccountScanRows(rows))
	}
	return accountList
}

func (this *accountModel) CountItems(search string) int64{
	var querySql = `SELECT count(*) FROM account ` + this.getSearchSql(search)
	var all int64
	row := Crud.GetRow(querySql,[]interface{}{})
	err := row.Scan(&all)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
		return -1
	}
	return all
}

func (this *accountModel)getSearchSql(search string)string{
	if search != "" {
		return `WHERE account_email LIKE '%`+search+`%'
		OR account_nick LIKE '%`+search+`%'
		OR account_fist_name LIKE '%`+search+`%'
		OR account_last_name LIKE '%`+search+`%'
		OR account_phone LIKE '%`+search+`%'
		OR account_state LIKE '%`+search+`%'
		OR account_city LIKE '%`+search+`%'
		OR account_steam_id LIKE '%`+search+`%'`
	}
	return ""
}
