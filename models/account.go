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
	"github.com/ilyaran/Malina/entity"
	"database/sql"
)

const  FieldsAccount = `
	account_id,
	coalesce(account_position,0),
	coalesce(account_email,''),
	coalesce(account_nick,''),
	account_password,
	coalesce(account_phone,''),
	account_provider,
	account_token,
	account_ban,
	account_ban_reason,
	account_newpass,
	account_newpass_key,
	account_newpass_time,
	account_last_ip,
	account_last_logged,
	account_created,
	account_updated
	`

var AccountModel *accountModel = new(accountModel)
type accountModel struct {
	Where string
	Query string
	All int64
}

func (this *accountModel)CheckDetail(exec []interface{})bool{

	this.Query  = `SELECT count(*) FROM account WHERE `+this.Where
	err := Crud.GetRow(this.Query , exec).Scan(&this.All)
	if err!=nil{
		panic(err)
	}
	if this.All > 0 {
		this.All = 0
		return true
	}
	return false
}

func (this *accountModel)Get(id int64)*entity.Account{
	this.Query = `
	SELECT `+FieldsAccount+`
	FROM account
	LEFT JOIN position ON position_id = account_id
	WHERE account_id = $1 LIMIT 1`
	row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.AccountScan(row,nil)
}

func (this *accountModel)Del(id int64)int64{
	var querySql = `
	DELETE FROM account WHERE account_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *accountModel)Add(q string,exec []interface{})int64{
	this.Query = `
	INSERT INTO account `+q
	return Crud.Insert(this.Query, exec)
}
func (this *accountModel)Edit(q string,exec []interface{})int64{
	this.Query = `UPDATE account SET `+q
	return Crud.Update(this.Query, exec)
}

func (this *accountModel) GetList(search,page,perPage string,order_by string,TreeMap map[int64]*entity.Position) []*entity.Account {

	this.CountItems(search)

	if order_by == "" {
		order_by = "account_last_logged DESC"
	}
	this.Query = `
	SELECT ` + FieldsAccount + `
	FROM account
	LEFT JOIN position ON position_id = account_id
	` + this.Where + `
	ORDER BY ` + order_by + `
	LIMIT ` + perPage + ` OFFSET ` + page

	rows := Crud.GetRows(this.Query, []interface{}{})

	var accountList = []*entity.Account{}

	var account *entity.Account
	for rows.Next() {
		account = entity.AccountScan(nil,rows)
		if account.GetPosition() != nil && account.GetPosition().GetId() > 0 {
			if v,ok := TreeMap[account.GetPosition().GetId()]; ok {
				account.Position = v
			}
		}
		accountList = append(accountList, account)
	}
	return accountList
}

func (this *accountModel) CountItems(search string){
	if search != "" {
		this.Where = ` WHERE account_email LIKE '%` + search + `%'
		OR account_nick LIKE '%` + search + `%'
		OR account_phone LIKE '%` + search + `%' `
		this.Query = `SELECT count(*) FROM account ` + this.Where
		row := Crud.GetRow(this.Query, []interface{}{})
		err := row.Scan(&this.All)
		if err == sql.ErrNoRows {
			this.All = 0
		}
		if err != nil {
			panic(err)
			this.All = -1
		}
	}
}
