package model

import (
	"fmt"
	"Malina/entity"
	"Malina/config"
	"database/sql"
)

const FieldsSession = `
	ip_address,
	account_id,
	email,

	nick,
	phone,
	data,

	user_agent,
	is_flash,
	position_id,

	balance,
	position_title,
	permission_data`

var Session *SessionModel = new(SessionModel)
type SessionModel struct {}

func (this *SessionModel) Add(sess *entity.Session) bool {
	var querySql = `
		INSERT INTO sessions
		(id,`+ FieldsSession +`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
		coalesce((SELECT position_title FROM position WHERE position_id = $10),''),
		coalesce((SELECT permission_data FROM permission WHERE permission.permission_position = $10),'')
		)
		returning id`
	//fmt.Println(querySql)
	var id string = ""
	err := Crud.GetRow(querySql, sess.ExecWithId()).Scan(&id)
	if err == sql.ErrNoRows{
		return false
	}
	if err != nil {
		panic(err)
		return false
	}
	return true
}
func (this *SessionModel) Update(sess *entity.Session) int64 {
	var querySql = `
		UPDATE sessions
		SET (`+ FieldsSession +`)
		=
		($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
		coalesce((SELECT position_title FROM position WHERE position_id = $10),''),
		coalesce((SELECT permission_data FROM permission WHERE permission.permission_position = $10),'')
		)
		WHERE id = $1
	`
	return Crud.Update(querySql,sess.ExecWithId())
}
func (this *SessionModel) Del(id string) bool {
	var exec = []interface{}{id}
	var querySql = `DELETE FROM sessions WHERE id = $1`
	if Crud.Delete(querySql, exec) > 0{
		return true
	}
	return false
}
func (this *SessionModel) Get(sessionId,ip_address string) *entity.Session{
	var querySql string = `
	WITH t AS(
		SELECT delete_expired_sessions(`+fmt.Sprintf("%v",app.Session_expiration())+`)
	)`
	if ip_address != "noip"{
		querySql += `
		, t1 AS(
			SELECT id,` + FieldsSession + `
			FROM sessions
			WHERE id = $1 LIMIT 1
		),t2 AS (
			UPDATE account SET (account_last_logged, account_last_ip) = (now(), $2)
			WHERE account_id = (SELECT t1.account_id FROM t1)
		)
		SELECT * FROM t1
		`
	}else {
		querySql += `
		SELECT id,`+ FieldsSession +`
		FROM sessions
		WHERE id = $1 LIMIT 1
		`
	}
	row := Crud.GetRow(querySql, []interface{}{sessionId,ip_address})
	return entity.SessionScanRow(row)
}

/*
WITH t AS(
			SELECT id,` + FieldsSession + `
			FROM delete_expired_sessions(` + fmt.Sprintf("%v", app.Session_expiration()) + `), sessions
			WHERE id = $1 LIMIT 1
		), t1 AS (
			UPDATE sessions SET timestamp = (date_part('epoch'::text, now()))::bigint
			WHERE id = $1
		),t2 AS (
			UPDATE account SET (account_last_logged, account_last_ip) = (now(), $2)
			WHERE account_id = (SELECT t.account_id FROM t)
		)
		SELECT * FROM t
 */
