/**
 * Product model class.  Malina eCommerce application
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

var FieldsProduct = `
	product.product_id,
	category_id,
	product_price,
	product_price1,
	product_code,
	product_title,
	product_description,
	product_created,
	product_updated,
	product_enable,
	product_img
`

var ProductModel = new(productModel)
type productModel struct {
	Where string
	Query string
	All int64
}

func (this *productModel)Get(id int64)*entity.Product{
	this.Query = `SELECT `+FieldsProduct+` FROM product_category
	INNER JOIN product ON product.product_id = product_category.product_id
	WHERE product_category.product_id = $1 LIMIT 1`
	row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.ProductScanRow(row)
}

func (this *productModel)Del(id int64)int64{
	this.Query = `
	DELETE FROM product WHERE product_id = $1`
	return Crud.Delete(this.Query, []interface{}{id})
}

func (this *productModel)Add(product *entity.Product,imgIdsString string)int64{
	var res int64 = 0
	this.Query= `
	SELECT product_add($1,$2,$3,$4,$5,$6,$7,ARRAY[`+imgIdsString+`]::VARCHAR(255)[]);
	`
	row := Crud.GetRow(this.Query, product.Exec()).Scan(&res)
	if row == sql.ErrNoRows {
		return 0
	}
	if row != nil {
		panic(row)
		return 0
	}
	return res
}
func (this *productModel)Edit(product *entity.Product,imgIdsString string)int64{
	var res int64 = 0
	this.Query = `
	SELECT product_update($1,$2,$3,$4,$5,$6,$7,$8,ARRAY[`+imgIdsString+`]::VARCHAR(255)[]);
	`
	row := Crud.GetRow(this.Query, product.ExecWithId()).Scan(&res)
	if row == sql.ErrNoRows {
		return 0
	}
	if row != nil {
		panic(row)
		return 0
	}
	return res
}
func (this *productModel) GetList(
	categoryWhere,
	priceInterval,
	search,
	page,
	perPage,
	order_by,
	enableItems string) []*entity.Product {

	this.CountItems(categoryWhere, priceInterval, search, enableItems)
	if search != "" {
		if order_by == ""{
			order_by = "ts_rank_cd(i.search_vector, query_string) DESC"
		}else {
			order_by = order_by + ", ts_rank_cd(product.search_vector, query_string) DESC"
		}
	}
	this.Query = `
	SELECT `+FieldsProduct+` FROM product_category
	INNER JOIN product ON product.product_id = product_category.product_id
	` + this.Where + `
	ORDER BY ` + order_by + `
	LIMIT ` + perPage + ` OFFSET ` + page
	//fmt.Println(this.Query)
	rows := Crud.GetRows(this.Query, []interface{}{})
	defer rows.Close()
	var itemList = []*entity.Product{}
	for rows.Next() {
		itemList = append(itemList, entity.ProductScanRows(rows))
	}
	return itemList
}

func (this *productModel) CountItems(
	categoryWhere,
	priceInterval,
	search,
	enableItems string){
	if search != `` {
		this.Where = `INNER JOIN to_tsquery('` + search + `') AS query_string ON search_vector @@ query_string`
	}
	if enableItems != `` {
		this.Where += " WHERE " + enableItems
	}
	if categoryWhere != ""{
		if this.Where != "" {
			this.Where += " AND " + categoryWhere
		}else {
			this.Where += " WHERE " + categoryWhere
		}
	}
	if priceInterval != ""{
		if this.Where != "" {
			this.Where += " AND " + priceInterval
		}else {
			this.Where += " WHERE " + priceInterval
		}
	}
	this.Query = `SELECT count(*) FROM product_category
	INNER JOIN product ON product.product_id = product_category.product_id
	` + this.Where
	//fmt.Println(this.Query)
	row := Crud.GetRow(this.Query,[]interface{}{})
	err := row.Scan(&this.All)
	if err == sql.ErrNoRows {

	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}
