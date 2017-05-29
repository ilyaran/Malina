/**
 * CRUD Postgresql class.  Malina eCommerce application
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
	"database/sql"
	_ "github.com/lib/pq"
)

var Crud *CrudPostgres

type CrudPostgres struct {
	DB *sql.DB
}

/*
if err := DB.QueryRow("SELECT pin FROM pins").Scan(&pin); err != nil {
	http.Error(w, "database error", http.StatusInternalServerError)
	return
}
 */

func (this *CrudPostgres) GetRows(querySql string, exec []interface{}) *sql.Rows {
	rows, err := this.DB.Query(querySql, exec...)
	if err != nil {
		panic(err)
		return nil
	}
	return rows

	return nil
}
func (this *CrudPostgres) GetRow(querySql string, exec []interface{}) *sql.Row {
	row := this.DB.QueryRow(querySql, exec...)
	return row

	return nil
}

func (this *CrudPostgres) Insert(querySql string, exec []interface{}) int64 {
	var lastInsertId int64
	err := this.DB.QueryRow(querySql, exec...).Scan(&lastInsertId)
	if err == sql.ErrNoRows {
		return -1
	}
	if err != nil {
		panic(err)
		return -2
	}
	return lastInsertId

}

func (this *CrudPostgres) Update(querySql string, exec []interface{}) int64 {
	stmt, err := this.DB.Prepare(querySql)
	if err != nil {
		panic(err)
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(exec...)
	if err != nil {
		panic(err)
		return -2
	}

	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
		return -3
	}
	return affect
}

func (this *CrudPostgres) Delete(querySql string, exec []interface{}) int64 {
	stmt, err := this.DB.Prepare(querySql)
	if err != nil {
		panic(err)
		return -1
	}
	res, err := stmt.Exec(exec...)
	if err != nil {
		panic(err)
		return -2
	}
	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
		return -3
	}
	return affect

}
