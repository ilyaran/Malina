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
	"Malina/entity"
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
	return entity.CategoryScanRow(row)
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


