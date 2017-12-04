package models

import (
	"github.com/ilyaran/Malina/entity"
	"database/sql"
	"github.com/ilyaran/Malina/app"
	"net/http"
)

type Session struct {
	SqlSelectFieldsAccount string
}

// return Account and bool. Bool means that session is found, false not found
func(s *Session)Get(sessionId, ip, user_agent string)(*entity.Account,bool){
	account := new(entity.Account)
	row := CrudGeneral.DB.QueryRow(`
		WITH t AS (
			DELETE FROM Session WHERE
			date_part('epoch',CURRENT_TIMESTAMP)::bigint - session_created > `+app.Session_expirationStr+`
		)
		SELECT `+s.SqlSelectFieldsAccount+` FROM session
		LEFT JOIN account ON account_id = session_account
		WHERE session_id = $1 LIMIT 1`, sessionId)
	if account.Scanning(row,nil) == 'o' {
		if account.Id > 0 {
			go s.UpdateAccount(ip,account.Id)
			return account,true
		}
		return nil,true
	}

	return nil,false
}

func(s *Session)Add(r *http.Request, account *entity.Account,sessionId,ip,userAgent string)string{
	var err error
	if account == nil {
		err = CrudGeneral.DB.QueryRow(`
		INSERT INTO Session
		(
			session_id,
			session_ip,
			session_user_agent
		)VALUES (uuid_generate_v4()||'`+`_`+sessionId+`',$1,$2) RETURNING session_id`,ip,userAgent).Scan(&sessionId)
	}else {
		err = CrudGeneral.DB.QueryRow(`
		INSERT INTO Session
		(
			session_id,
			session_ip,
			session_user_agent,

			session_account,
			session_email,
			session_nick,
			session_phone,
			session_token
		)
		SELECT uuid_generate_v4()||'`+`_`+sessionId+`',$1,$2,account_id,account_email,account_nick,account_phone,account_token
		FROM account
		WHERE account_id = $3
		RETURNING session_id`, ip,userAgent,account.Id).Scan(&sessionId)
	}

	if err == sql.ErrNoRows {
		panic(err)
		return ""
	}
	if err != nil {
		panic(err)
		return ""
	}
	return sessionId
}
func(s *Session)UpdateAccount(ip string,accountId int64)int64{

	return CrudGeneral.Update(`
	UPDATE account
		SET
			account_last_logged = now(),
			account_last_ip = $1
	WHERE account_id = $2 `,ip,accountId)
}

func(s *Session)Del(sessionId string)int64{

	return CrudGeneral.Delete(`DELETE FROM Session WHERE session_id = $1 `,sessionId)
}