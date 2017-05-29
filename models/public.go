/**
 * Public model class.  Malina eCommerce application
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
	"github.com/ilyaran/Malina/entity"
	"strings"
)

var PublicModel *publicModel = new(publicModel)
type publicModel struct {
	Where string
	Query string
	TemporaryData string
	All int64
}
func (this *publicModel)GetProduct(id int64)*entity.Product{
	this.Query = `SELECT `+FieldsProduct+` FROM product_category
	INNER JOIN product ON product.product_id = product_category.product_id
	WHERE product_enable = TRUE AND product_category.product_id = $1 LIMIT 1`
	row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.ProductScanRow(row)
}
func (this *publicModel) GetList(
	categoryWhere,
	priceInterval,
	search,
	page,
	perPage,
	order_by string) []*entity.Product {

	this.CountItems(categoryWhere, priceInterval, search)
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

func (this *publicModel) CountItems(
	categoryWhere,
	priceInterval,
	search string){
	if search != `` {
		this.Where = `INNER JOIN to_tsquery('` + search + `') AS query_string ON search_vector @@ query_string`
	}
	this.Where += " WHERE product_enable = true"

	if categoryWhere != ""{
		this.Where += " AND " + categoryWhere
	}
	if priceInterval != ""{
		this.Where += " AND " + priceInterval
	}

	this.Query = `SELECT count(*) FROM product_category
	INNER JOIN product ON product.product_id = product_category.product_id
	` + this.Where
	//fmt.Println(this.Query)
	row := Crud.GetRow(this.Query,[]interface{}{})
	err := row.Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = -2
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}

func (this *publicModel)Add_product(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_add_product($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}
func (this *publicModel)Del_product(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_remove_product($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}

func (this *publicModel)GetCartProducts(cart_id string)[]*entity.CartPublic{
	this.Query  = `
		SELECT * FROM cart_get_products($1)
	`
	this.TemporaryData = ``
	var rows = Crud.GetRows(this.Query , []interface{}{cart_id})
	defer rows.Close()
	var cartList = []*entity.CartPublic{}
	var item *entity.CartPublic
	for rows.Next() {
		item = &entity.CartPublic{}
		err := rows.Scan(&item.Product_id,&this.TemporaryData,&item.Product_title,&item.Product_price,&item.Product_quantity,&item.Subtotal)
		if err != nil {
			panic(err)
		}
		if this.TemporaryData != `` {
			item.Product_img = strings.Split(this.TemporaryData,"|")
		}
		cartList = append(cartList, item)
	}
	return cartList
}

func (this *publicModel)CountCartProducts(cart_id string){
	this.Query  = `SELECT sum(cart_quantity) FROM cart WHERE cart_id = $1`
	err := Crud.GetRow(this.Query , []interface{}{cart_id}).Scan(&this.All)
	if err != nil {
		this.All = 0
	}
}


func (this *publicModel)Get_saved_products(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}
func (this *publicModel)Get_total_amount(cart_id string, product_id int64)float64{

	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	var total float64
	err := Crud.GetRow(this.Query , []interface{}{cart_id,product_id}).Scan(total)
	if err == sql.ErrNoRows {
		return -1.00
	}
	if err != nil {
		panic(err)
		return -2.00
	}
	return total
}
func (this *publicModel)Move_product_to_cart(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}
func (this *publicModel)Remove_product(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}
func (this *publicModel)Save_product_for_later(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}
func (this *publicModel)Cart_update(cart_id string, product_id int64)*sql.Row{
	this.Query  = `
		SELECT cart_get_products($1,$2)
	`
	return Crud.GetRow(this.Query , []interface{}{cart_id,product_id})
}








