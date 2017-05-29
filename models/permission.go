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
	"github.com/ilyaran/Malina/entity"
	"database/sql"
)

const FieldsPermission = `
	permission_id,
	coalesce(permission_position,0),
	coalesce((SELECT position_title FROM position WHERE position_id = permission_position),''),
	permission_data`



var PermissionModel = new(permissionModel)

type permissionModel struct {
	Where string
	Query string
	All   int64
	Exec []interface{}
}

func (this *permissionModel)Get(id int64) *entity.Permission {
	this.Query = `
	SELECT ` + FieldsPermission + ` FROM permission
	WHERE permission_id = $1 LIMIT 1 `
	row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.PermissionScan(row,nil)
}

func (this *permissionModel)Del() int64 {
	this.Query = `DELETE FROM permission WHERE permission_id = $1`
	return Crud.Delete(this.Query, this.Exec)
}

func (this *permissionModel)Add() int64 {
	this.Query = `
	INSERT INTO permission
	(permission_data,permission_position)
	VALUES($1,$2) RETURNING permission_id`
	return Crud.Insert(this.Query, this.Exec)
}

func (this *permissionModel)Edit(permission *entity.Permission) int64 {
	this.Query = `
	UPDATE permission SET
	(permission_data,permission_position)
	= ($1,$2) WHERE permission_id = $3 `
	return Crud.Update(this.Query, this.Exec)
}

func (this *permissionModel) GetList(search, page, perPage string, order_by string) []*entity.Permission {

	this.CountItems(search)

	var limit = ""
	if page != "" && perPage != "" {
		limit = `LIMIT ` + perPage + ` OFFSET ` + page
	}
	if order_by == "" {
		order_by = "permission_id ASC"
	}
	this.Query = `
	SELECT ` + FieldsPermission + ` FROM permission
	LEFT JOIN position ON position_id = permission_position
	` + this.Where + `
	ORDER BY ` + order_by + ` ` + limit

	rows := Crud.GetRows(this.Query, []interface{}{})

	var permissionList = []*entity.Permission{}

	for rows.Next() {
		permissionList = append(permissionList, entity.PermissionScan(nil,rows))
	}
	return permissionList
}

func (this *permissionModel) CountItems(search string) {
	if search != "" {
		this.Where = `WHERE permission_data LIKE '%`+search+`%'`
	}
	this.Query = `SELECT count(*) FROM permission ` + this.Where
	row := Crud.GetRow(this.Query, []interface{}{})
	err := row.Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = 0
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}



