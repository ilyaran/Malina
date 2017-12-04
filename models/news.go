package models

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



