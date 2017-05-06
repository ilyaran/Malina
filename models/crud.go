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
	db          *sql.DB
}
/*
if err := db.QueryRow("SELECT pin FROM pins").Scan(&pin); err != nil {
	http.Error(w, "database error", http.StatusInternalServerError)
	return
}
 */

func (this *CrudPostgres) GetRows(querySql string, exec []interface{}) *sql.Rows {
	if this.db != nil {
		rows, err := this.db.Query(querySql, exec...)
		if err != nil {
			panic(err)
			return nil
		}
		return rows
	}
	return nil
}
func (this *CrudPostgres) GetRow(querySql string, exec []interface{}) *sql.Row {
	if this.db != nil {
		row := this.db.QueryRow(querySql, exec...)
		return row
	}
	return nil
}

func (this *CrudPostgres) Insert(querySql string, exec []interface{})int64{
	if this.db != nil {
		var lastInsertId int64
		err := this.db.QueryRow(querySql, exec...).Scan(&lastInsertId)
		if err == sql.ErrNoRows{
			return -1
		}
		if err != nil {
			panic(err)
			return -2
		}
		return lastInsertId
	}
	return -4
}

func (this *CrudPostgres) Update(querySql string, exec []interface{})int64{
	if this.db != nil {
		stmt, err := this.db.Prepare(querySql)
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
	return -4

}

func (this *CrudPostgres) Delete(querySql string, exec []interface{})int64{
	if this.db != nil {
		stmt, err := this.db.Prepare(querySql)
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
	return -4
}

