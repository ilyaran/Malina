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
)

type News struct {
	Crud      CrudPostgres
	err       error

}

func (this *News)Upsert(action byte,imgUrlsString string,args... interface{})int64{
	var res int64 = 0
	if imgUrlsString!=""{
		imgUrlsString = `ARRAY[`+imgUrlsString+`]::VARCHAR(255)[]`
	}else {
		imgUrlsString="NULL"
	}
	if action=='a'{
		this.err = this.Crud.DB.QueryRow(`
		SELECT news_add($1,$2,$3,$4,$5,$6,$7,`+imgUrlsString+`)
		`, args...).Scan(&res)
	}else{
		this.err = this.Crud.DB.QueryRow(`
		SELECT news_update($1,$2,$3,$4,$5,$6,$7,$8,`+imgUrlsString+`);
		`, args...).Scan(&res)
	}
	if this.err == sql.ErrNoRows {
		return 0
	}
	if this.err != nil {
		panic(this.err)
		return 0
	}
	return res
}



