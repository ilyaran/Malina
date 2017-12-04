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
	"fmt"
)

type Account struct {
}

func (this *Account)Upsert(action byte, img string,args... interface{})int64{
	var res int64 = 0
	var err error
	if img !=""{
		img = `ARRAY[`+ img +`]::VARCHAR(255)[]`
	}else {
		img ="NULL"
	}
	if action=='a'{
		err = CrudGeneral.DB.QueryRow(`
		SELECT account_add($1,$2,$3,$4,$5,$6,$7,`+ img+`)
		`, args...).Scan(&res)
	}else{
		err = CrudGeneral.DB.QueryRow(`
		SELECT account_update($1,$2,$3,$4,$5,$6,$7,$8,`+ img+`);
		`, args...).Scan(&res)
	}
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
		return 0
	}
	return res
}

func (s *Account)CheckDetail(where string ,value interface{})bool{
	var count int64
	err := CrudGeneral.DB.QueryRow(`SELECT count(*) FROM account WHERE `+where , value).Scan(&count)
	if err!=nil{
		panic(err)
		return false
	}
	if count > 0 {
		count = 0
		return true
	}
	return false
}
func (s *Account)EditByObject(a *entity.Account)int64{

	query := `
		UPDATE account SET
			account_email=$1,
			account_nick=$2,
			account_password=$3,

			account_updated=now(),
			account_state=`;if a.State==0 {query +="NULL"}else{query +=strconv.FormatInt(a.State,10)}; query +=`,

			account_game=$4,

			account_skype=$5,

			account_trade=$6,
			account_steam=$7,
			account_first_name=$8,

			account_last_name=$9,
			account_sex=$10

		WHERE account_id = $11
		`
	return CrudGeneral.Update(query,
		a.Email,a.Nick,a.Password,
		a.Game, a.Skype, a.Trade,
		a.Steam, a.FirstName,a.LastName,
		a.Sex,a.Id,
	)
}

func (this *Account) GetListPublic(sqlSelectFields,join, condition, search string,page,perPage int64,order_by string) ([]*entity.Account,int64) {
	query := ``
	where := ``
	var all int64

	if search != "" {
		where = ` WHERE account_email LIKE '%` + search + `%'
		OR account_nick LIKE '%` + search + `%'
		OR account_first_name LIKE '%` + search + `%'
		OR account_last_name LIKE '%` + search + `%'
		OR account_steam LIKE '%` + search + `%'
		OR account_skype LIKE '%` + search + `%'
		OR account_trade LIKE '%` + search + `%'`
	}
	if condition != "" && where != `` {
		where += condition
	}else if condition != "" && where == `` {
		where = ` WHERE `+condition
	}
	where = join + where

	//***************** count ************************
	err := CrudGeneral.DB.QueryRow(`SELECT count(*) FROM account ` + where).Scan(&all)
	if err == sql.ErrNoRows {
		all = 0
		return nil,all
	}
	if err != nil {
		panic(err)
		all = -1
		return nil,all
	}
	//*****************************************

	if order_by == "" {
		order_by = "account_last_logged DESC"
	}
	query = `
		SELECT ` + sqlSelectFields + `
		FROM account
		` + where + `
		 ` + order_by + `
		LIMIT $1 OFFSET $2`
	fmt.Println(query)
	rows,err := CrudGeneral.DB.Query(query,perPage,page)
	if err!=nil{panic(err)}
	defer rows.Close()

	accountList:=make([]*entity.Account,0,perPage)
	for rows.Next() {
		item:=new(entity.Account)
		if item.Scanning(nil, rows) == 'o' {
			accountList=append(accountList,item)
		}
	}

	return accountList,all
}

func (this *Account)AccountAccount(account_id1,account_id2 int64,status string)(int64){

	if account_id1==account_id2{
		return 0
	}
	// 10 - friend request, 20 - it means a friend
	if status == "friend_request_add" {
		var all int64
		CrudGeneral.DB.QueryRow(`
		SELECT count(*) FROM account_account
		WHERE account_id1 = $1 AND account_id2 = $2 LIMIT 1
		`,account_id1,account_id2).Scan(&all)

		if all != 0 {
			return -10
		}
		return CrudGeneral.Insert(`
			INSERT INTO account_account
			(account_id1,account_id2,status)
			VALUES ($1,$2,$3)
			RETURNING account_id1`,  account_id2, account_id1, 10)
	}
	if status == "friend_request_accept" {
		return CrudGeneral.Update(`
		UPDATE account_account SET status = $1 , created = now()
		WHERE account_id1 = $2 AND account_id2 = $3
	`, 20, account_id1, account_id2)
	}
	if status == "friend_request_remove" || status == "friend_remove" {
		return CrudGeneral.Delete(`
		DELETE FROM account_account
		WHERE account_id1 = $1 AND account_id2 = $2
	`, account_id1, account_id2)
	}
	return 0
}









