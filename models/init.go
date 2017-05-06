/**
 * Initialize model package.  Malina eCommerce application
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
	"Malina/config"
)

func init()  {
	db, err := sql.Open("postgres", "postgres://" + app.DB_USER + ":" + app.DB_PASSWORD + "@" + app.HOST + "/" + app.DB_NAME + "?sslmode=disable")
	if err != nil {
		Crud = nil
		return
	}
	if db.Ping() != nil {
		Crud = nil
		return
	}
	Crud = &CrudPostgres{db}
	if Crud != nil {

	}
}

func CloseDb(){
	Crud.db.Close()
}
