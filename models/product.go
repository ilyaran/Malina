package models

import (
	"database/sql"
)

type Product struct {

}

func (this *Product)Upsert(action byte,imgUrlsString,parameterIds string,args... interface{})int64{
	var res int64 = 0
	var err error
	if imgUrlsString!=""{
		imgUrlsString = `ARRAY[`+imgUrlsString+`]::VARCHAR(255)[]`
	}else {
		imgUrlsString="NULL"
	}
	if parameterIds!=""{
		parameterIds = `ARRAY[`+parameterIds+`]::BIGINT[]`
	}else {
		parameterIds="NULL"
	}
	if action=='a'{
		err = CrudGeneral.DB.QueryRow(`
		SELECT product_add($1,$2,$3,$4,$5,$6,$7,$8,`+imgUrlsString+`,`+parameterIds+`)
		`, args...).Scan(&res)
	}else{
		err = CrudGeneral.DB.QueryRow(`
		SELECT product_update($1,$2,$3,$4,$5,$6,$7,$8,$9,`+imgUrlsString+`,`+parameterIds+`);
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



