package models

import (
	"database/sql"
	"fmt"
	"github.com/ilyaran/Malina/app"
	"errors"
	"github.com/ilyaran/Malina/berry"
	"github.com/ilyaran/Malina/entity"
)

const TAG = "CrudPostgres"

var CrudGeneral *CrudPostgres
type CrudPostgres struct {
	DB          			*sql.DB
}

func (s *CrudPostgres) SetDBConnection() error{
	// database connection
	var err error
	s.DB, err = sql.Open("postgres", "postgres://" + app.DB_user+ ":" + app.DB_password+ "@" + app.DB_host+ "/" + app.DB_name+ "?sslmode=disable")
	if err != nil{
		err = errors.New("Fail Database initialize")
		panic(err)
		return err
	}
	err = s.DB.Ping()
	if err != nil {
		err = errors.New("Fail Database ping")
		panic(err)
		return err
	}
	return nil
	// end database connection
}

func (this *CrudPostgres)Select(malina *berry.Malina,selectSql string){
	malina.SelectSql = selectSql
}
func (this *CrudPostgres)Table(malina *berry.Malina,tableSql string){
	malina.TableSql = tableSql
}
func (this *CrudPostgres)Join(malina *berry.Malina,pre,table, on string){
	malina.JoinSql += pre+" JOIN "+table+" ON "+on
}
func (this *CrudPostgres)WhereAnd(malina *berry.Malina,key string,value string){
	if malina.WhereSql == ``{
		malina.WhereSql += ` WHERE `
	}else {
		malina.WhereSql +=  ` AND `
	}
	malina.WhereSql += key+value
}
func (this *CrudPostgres)Limit(malina *berry.Malina,perPage, page int64){
	if perPage>0 {
		malina.LimitSql = fmt.Sprintf(` LIMIT %d OFFSET %d `, perPage, page)
	}else {
		malina.LimitSql = ``
	}
}
func (this *CrudPostgres)OrderBy(malina *berry.Malina,key, dir string){
	if key!="" && dir!="" {
		malina.OrderBySql = " ORDER BY "+key+" "+dir+" "
	}else {
		malina.OrderBySql = ""
	}
}

func (s *CrudPostgres)GetItem(malina *berry.Malina,item entity.Scanable,query string, args ...interface{}){
	var row *sql.Row
	if query != ""{
		row = s.DB.QueryRow(query, args...)
	}else {
		if malina.SelectSql == ""{
			malina.SelectSql = "*"
		}
		row = s.DB.QueryRow(`
		SELECT ` + malina.SelectSql + ` FROM ` + malina.TableSql + ` `+ malina.JoinSql + ` `+malina.WhereSql+" LIMIT 1", args...)


	}
	malina.Item = item.Item(row)
}

func (s *CrudPostgres)GetItems(malina *berry.Malina,item entity.Scanable,query string, args ...interface{}){var rows *sql.Rows
	var err error
	if query != ""{
		rows,err = s.DB.Query(query, args...)
	}else {
		if malina.SelectSql ==""{
			malina.SelectSql = "*"
		}
		rows,err = s.DB.Query(`SELECT ` + malina.SelectSql + ` FROM ` + malina.TableSql + ` `+ malina.JoinSql + ` `+malina.WhereSql + ` `+malina.OrderBySql + ` `+malina.LimitSql, args...)
	}
	//fmt.Println(`SELECT ` + malina.SelectSql + ` FROM ` + malina.TableSql + ` `+ malina.JoinSql + ` `+malina.WhereSql + ` `+malina.OrderBySql + ` `+malina.LimitSql, args)

	if err != nil {
		panic(err)
		return
	}
	defer rows.Close()

	malina.List = item.Appending(rows,malina.Per_page)
}

func (s *CrudPostgres) Count(malina *berry.Malina,query string, args... interface{}) {
	var err error
	if query != ""{
		err = s.DB.QueryRow(query, args...).Scan(&malina.All)
	}else {
		err = s.DB.QueryRow(`SELECT count(*) FROM ` + malina.TableSql+ ` `+malina.JoinSql + ` `+malina.WhereSql, args...).Scan(&malina.All)
		//fmt.Println(`SELECT count(*) FROM ` + malina.TableSql+ ` `+malina.JoinSql + ` `+malina.WhereSql, args)
	}
	if err != nil {
		malina.All = 0
		malina.ChannelBool <- false
		panic(err)
	}
	if malina.All == 0 {
		malina.ChannelBool <- false
	}else {
		malina.ChannelBool <- true
	}

}

func (this *CrudPostgres) Insert(querySql string, args ... interface{}) int64 {
	var lastInsertId int64
	err := this.DB.QueryRow(querySql, args...).Scan(&lastInsertId)
	if err == sql.ErrNoRows {
		return -10
	}
	if err != nil {
		panic(err)
		return -20
	}
	return lastInsertId

}

func (this *CrudPostgres) Update(querySql string, args ... interface{}) int64 {
	stmt, err := this.DB.Prepare(querySql)
	if err != nil {
		panic(err)
		return -10
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
		return -20
	}

	affected, err := result.RowsAffected()
	if err != nil {
		//panic(this.Err)
		return -30
	}
	return affected
}
func (this *CrudPostgres) Delete(querySql string, args ... interface{}) int64 {
	stmt, err := this.DB.Prepare(querySql)
	if err != nil {
		panic(err)
		return -10
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
		return -20
	}
	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
		return -30
	}
	return affected

}
















