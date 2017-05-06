/**
 * Permission model class.  Malina eCommerce application
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
	"Malina/entity"
	"database/sql"
)

const FieldsPermission = `
	permission_id,
	permission_position,
	coalesce((SELECT position_title FROM position WHERE position_id = permission_position),'user'),
	permission_data`

var PermissionModel = new(permissionModel)
type permissionModel struct {}

func (this *permissionModel)Get(id int64)*entity.Permission{
	var querySql = `SELECT `+FieldsPermission+` FROM permission
	WHERE permission_id = $1 LIMIT 1`
	row := Crud.GetRow(querySql, []interface{}{id})
	return entity.PermissionScanRow(row)
}

func (this *permissionModel)Del(id int64)int64{
	var querySql = `
	DELETE FROM permission WHERE permission_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *permissionModel)Add(permission *entity.Permission)int64{
	var querySql = `
	INSERT INTO permission (
		permission_data,
		permission_position
	)
	VALUES ($1,$2) RETURNING permission_id`
	return Crud.Insert(querySql, permission.Exec())
}

func (this *permissionModel)Edit(permission *entity.Permission)int64{
	var querySql = `
	WITH t AS (
		UPDATE sessions SET permission_data = $2
		WHERE position_id = $3
	)
	UPDATE permission SET (permission_data,permission_position) = ($2,$3)
	WHERE permission_id = $1`
	return Crud.Update(querySql, permission.ExecWithId())
}

func (this *permissionModel) GetList(search,page,perPage string,order_by string) []*entity.Permission {
	var where = this.getSearchSql(search)
	var limit = ""
	if page != "" && perPage != ""{
		limit = `LIMIT ` + perPage + ` OFFSET ` + page
	}
	if order_by == "" {order_by = "permission_id ASC"}
	var querySql = `
	SELECT ` + FieldsPermission + ` FROM permission
	` + where + `
	ORDER BY ` + order_by + ` ` + limit

	rows := Crud.GetRows(querySql, []interface{}{})
	defer rows.Close()

	var permissionList = []*entity.Permission{}

	for rows.Next() {
		permissionList = append(permissionList, entity.PermissionScanRows(rows))
	}
	return permissionList
}

func (this *permissionModel) CountItems(search string) int64{
	var querySql = `SELECT count(*) FROM permission ` + this.getSearchSql(search)
	var all int64
	row := Crud.GetRow(querySql,[]interface{}{})
	err := row.Scan(&all)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
		return -1
	}
	return all
}

func (this *permissionModel)getSearchSql(search string)string{
	if search != "" {
		return `WHERE permission_data LIKE '%`+search+`%'`
	}
	return ""
}


