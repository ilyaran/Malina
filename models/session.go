package model

import (
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/config"
	"strconv"
	"database/sql"
)

var Session *SessionModel = new(SessionModel)
type SessionModel struct {
	Where string
	Query string
	Exec  []interface{}
	All   int64
	Row   *sql.Row
}
func (this *SessionModel) Get(sessionId string){
	this.Row = Crud.GetRow(`
	WITH t AS (
		DELETE FROM session WHERE
		date_part('epoch',CURRENT_TIMESTAMP)::bigint - session.session_timestamp > `+strconv.FormatInt(app.Session_expiration(),10)+`
	)
	SELECT session_id,session_data,session_data1,session_data2,
		coalesce(session_account,0),
		coalesce(session_email,''),
		coalesce(session_nick,''),
		coalesce(session_phone,''),
		session_provider,
		session_token,
		coalesce(session_position,0)
	FROM session WHERE session_id = $1 LIMIT 1
	`, []interface{}{sessionId})
}

func (this *SessionModel) Add(sess *entity.Session) int64 {
	// ********** FUTURE TASK ************
	// should create insert if position != nil etc.
	if sess.AccountId > 0{
		var exec = []interface{}{sess.Id,sess.Data,sess.AccountId,sess.Email,sess.Ip_address}
		this.Query = `
		WITH t AS (
			UPDATE account SET account_last_ip = $5, account_last_logged = now()
			WHERE account_id = $3
		)
		INSERT INTO session (session_id,session_data,session_account,session_email)
		VALUES ($1,$2,$3,$4,$5) RETURNING 1`
		return Crud.Insert(this.Query, exec)
	}
	this.Query = `INSERT INTO session (session_id,session_data) VALUES ($1,$2) RETURNING 1`
	return Crud.Insert(this.Query, []interface{}{sess.Id,sess.Data})
}

func (this *SessionModel) Del(id string) bool {
	this.Query = `DELETE FROM session WHERE session_id = $1`
	if Crud.Delete(this.Query, []interface{}{id}) > 0{
		return true
	}
	return false
}

func (this *SessionModel) Update(sess *entity.Session) int64 {
	this.Query = `
		UPDATE session SET
			session_ip = $1,
			session_data = $2
		WHERE session_id = $3
	`
	return 0//Crud.Update(this.Query,[]interface{}{sess.GetIp_address(),sess.GetAccount().GetId(),sess.GetData(),sess.GetId()})
}


