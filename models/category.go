/**
 * Category model class.  Malina eCommerce application
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

const FieldsCategory = `
	category_id,
	COALESCE(category_parent,0),
	category_sort,

	category_title,
	category_description,
	category_lang,

	category_enable,
	category_quantity,
	array_to_string(category_img, '|', '*')
	`

var CategoryModel = new(categoryModel)
type categoryModel struct {
	Where string
	Query string
	All int64
}

func (this *categoryModel)Get(id int64)*entity.Category{
	var querySql = `SELECT `+FieldsCategory+` FROM category
	WHERE category_id = $1 LIMIT 1`
	row := Crud.GetRow(querySql, []interface{}{id})
	return entity.CategoryScan(row,nil)
}

func (this *categoryModel)Del(id int64)int64{
	var querySql = `
	DELETE FROM category WHERE category_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *categoryModel)Add(category *entity.Category,imgIdsString string)int64{
	var querySql = `
	INSERT INTO category (
		category_parent,
		category_sort,
		category_title,
		category_description,
		category_lang,
		category_enable,
		category_img
	)
	VALUES ($1,$2,$3,$4,$5,$6,ARRAY[`+imgIdsString+`]::VARCHAR[]) RETURNING category_id `
	return Crud.Insert(querySql, category.Exec())
}

func (this *categoryModel)Edit(category *entity.Category,imgIdsString string)int64{

	var querySql = `
	UPDATE category SET (
		category_parent,
		category_sort,
		category_title,
		category_description,
		category_lang,
		category_enable,
		category_img
	) = ($2,$3,$4,$5,$6,$7,ARRAY[`+imgIdsString+`]::VARCHAR[])
	WHERE category_id = $1`
	return Crud.Update(querySql, category.ExecWithId())
}

func (this *categoryModel) GetList(search,pageStr, perPageStr string,order_by string) []*entity.Category {

	this.CountItems(search)

	this.Query = `
	SELECT ` + FieldsCategory + ` FROM category
	` + this.Where + `
	ORDER BY ` + order_by

	if pageStr != "" && perPageStr != ""{
		this.Query += ` LIMIT ` + perPageStr + ` OFFSET ` + pageStr
	}

	rows := Crud.GetRows(this.Query, []interface{}{})
	defer rows.Close()

	var categoryList = []*entity.Category{}

	for rows.Next() {
		categoryList = append(categoryList, entity.CategoryScan(nil,rows))
	}
	return categoryList
}

func (this *categoryModel) CountItems(search string) {
	if search != "" {
		this.Where = `WHERE category_title LIKE '%`+search+`%'
		OR category_description LIKE '%`+search+`%'`
	}

	this.Query = `SELECT count(*) FROM category ` + this.Where

	err := Crud.GetRow(this.Query,[]interface{}{}).Scan(&this.All)
	if err == sql.ErrNoRows {
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}




/*
import (
	"github.com/ilyaran/Malina/entity"
	"database/sql"
)

const FieldsCategory = `
	category_id,
	COALESCE(category_parent,0),
	category_sort,

	category_title,
	category_description,
	category_lang,

	category_enable,
	category_quantity,
	array_to_string(category_img, '|', '*')
	`

var CategoryModel = new(categoryModel)
type categoryModel struct {
	Where string
	Query string
	All int64
	Exec []interface{}
}*/

/*
func (this *categoryModel)Get(id int64)*entity.Category{
	var querySql = `SELECT `+FieldsCategory+` FROM category
	WHERE category_id = $1 LIMIT 1`
	Row := Crud.GetRow(querySql, []interface{}{id})
	return entity.CategoryScanRow(Row)
}

func (this *categoryModel)Del(id int64)int64{
	var querySql = `
	DELETE FROM category WHERE category_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}
*/
/*
func ( this *categoryModel ) Get ( id int64) *entity.Category{
	this.Query = `SELECT `+FieldsCategory+` FROM category WHERE id = ? LIMIT 1 `
	Row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.CategoryScan(Row,nil)
}
func (this *categoryModel)Del() int64 {
	this.Query = `DELETE FROM category WHERE category_id = $1 LIMIT 1`
	return Crud.Delete(this.Query, this.Exec)
}
func (this *categoryModel)Add(categoryParentId int64) int64 {
	this.Query = `
	INSERT INTO category SET `
	if categoryParentId > 0 {
		this.Query += `parent = ?,`
	}
	this.Query += `
		title = ?,
		sort = ?,
		enable = ?`
	return Crud.Insert(this.Query, this.Exec)
}

func (this *categoryModel)Edit(categoryParentId int64) int64 {
	this.Query = `
	UPDATE category SET `
	if categoryParentId > 0 {
		this.Query += `category_parent = $1,`
	}
	this.Query += `
		category_title = $2,
		category_sort = $3,
		category_enable = $4 WHERE category_id = $5 `
	return Crud.Update(this.Query, this.Exec)
}


func (this *categoryModel) GetList(search,page,perPage string,order_by string) []*entity.Category {

	this.CountItems(search)

	var limit = ""
	if page != "" && perPage != ""{
		limit = ` LIMIT ` + perPage + ` OFFSET ` + page
	}

	if order_by == "" {order_by = "sort ASC"}
	this.Query = `
	SELECT ` + FieldsCategory + ` FROM category
	` + this.Where + `
	ORDER BY ` + order_by + ` ` + limit

	rows := Crud.GetRows(this.Query, []interface{}{})

	var categoryList = []*entity.Category{}

	for rows.Next() {
		categoryList = append(categoryList, entity.CategoryScan(nil,rows))
	}

	return categoryList
}

func (this *categoryModel) CountItems(search string) {
	if search != `` {
		this.Where = ` WHERE title LIKE '%` + search + `%'`
	}

	this.Query = `SELECT count(*) FROM category ` + this.Where

	Row := Crud.GetRow(this.Query,[]interface{}{})
	err := Row.Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = 0
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}*/

/*
func (this *categoryModel)Add(category *entity.Category,imgIdsString string)int64{
	var querySql = `
	INSERT INTO category (
		category_parent,
		category_sort,
		category_title,
		category_description,
		category_lang,
		category_enable,
		category_img
	)
	VALUES ($1,$2,$3,$4,$5,$6,ARRAY[`+imgIdsString+`]::VARCHAR[]) RETURNING category_id `
	return Crud.Insert(querySql, category.Exec())
}

func (this *categoryModel)Edit(category *entity.Category,imgIdsString string)int64{

	var querySql = `
	UPDATE category SET (
		category_parent,
		category_sort,
		category_title,
		category_description,
		category_lang,
		category_enable,
		category_img
	) = ($2,$3,$4,$5,$6,$7,ARRAY[`+imgIdsString+`]::VARCHAR[])
	WHERE category_id = $1`
	return Crud.Update(querySql, category.ExecWithId())
}*/
/*

func (this *categoryModel) GetList(search,pageStr, perPageStr string,order_by string) []*entity.Category {

	this.CountItems(search)

	this.Query = `
	SELECT ` + FieldsCategory + ` FROM category
	` + this.Where + `
	ORDER BY ` + order_by

	if pageStr != "" && perPageStr != ""{
		this.Query += ` LIMIT ` + perPageStr + ` OFFSET ` + pageStr
	}

	rows := Crud.GetRows(this.Query, []interface{}{})
	defer rows.Close()

	var categoryList = []*entity.Category{}

	for rows.Next() {
		categoryList = append(categoryList, entity.CategoryScanRows(rows))
	}
	return categoryList
}

func (this *categoryModel) CountItems(search string) {
	if search != "" {
		this.Where = `WHERE category_title LIKE '%`+search+`%'
		OR category_description LIKE '%`+search+`%'`
	}

	this.Query = `SELECT count(*) FROM category ` + this.Where

	err := Crud.GetRow(this.Query,[]interface{}{}).Scan(&this.All)
	if err == sql.ErrNoRows {
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}


*/
