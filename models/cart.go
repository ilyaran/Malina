/**
 * Cart model class.  Malina eCommerce application
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

const  FieldsCart = `
	cart_id,
	cart_product,
	product_title,

	product_price,
	product_price1,
	cart_quantity,

	cart_buy_now,
	cart_created`

var CartModel *cartModel = new(cartModel)
type cartModel struct {
	Where string
	Query string
	All int64
}


func (this *cartModel)Get(id string)*entity.Cart{
	this.Query  = `SELECT `+FieldsCart+` FROM cart
	INNER JOIN product ON product_id = cart_product
	WHERE cart_id = $1 LIMIT 1`
	row := Crud.GetRow(this.Query , []interface{}{id})
	return entity.CartScanRow(row,nil)
}

func (this *cartModel)Del(id int64)int64{
	this.Query  = `
	DELETE FROM cart WHERE cart_id = $1`
	return Crud.Delete(this.Query , []interface{}{id})
}

func (this *cartModel)Add(cart *entity.Cart)int64{
	this.Query  = `
	INSERT INTO cart (cart_id,cart_product,cart_quantity,buy_now)
	VALUES ($1,$2,$3,$4) RETURNING cart_id`
	return Crud.Insert(this.Query , cart.ExecWithId())
}

func (this *cartModel)Edit(cart *entity.Cart)int64{
	this.Query  = `
	UPDATE cart SET (cart_product,cart_quantity,buy_now) = ($2,$3,$4)
	WHERE cart_id = $1
	`
	return Crud.Update(this.Query , cart.ExecWithId())
}

func (this *cartModel) GetList(searchStr, product_idStr,page,perPage string,order_by string) []*entity.Cart {

	this.CountItems(searchStr,product_idStr)

	this.Query = `
	SELECT
		cart_id,
		array_to_string(array_agg(cart_product), '|', '*'),
		array_to_string(array_agg(product_title),'|', '*'),

		array_to_string(array_agg(product_price),'|', '*'),
		array_to_string(array_agg(product_price1),'|', '*'),
		array_to_string(array_agg(cart_quantity),'|', '*'),

		array_to_string(array_agg(cart_buy_now),'|', '*'),
		cart_created
	FROM cart  `+ this.Where + `
	INNER JOIN product ON product_id = cart_product
	GROUP BY cart_id, cart_created
	ORDER BY ` + order_by + `
	LIMIT ` + perPage + ` OFFSET ` + page

	rows := Crud.GetRows(this.Query , []interface{}{})
	defer rows.Close()

	var cartList = []*entity.Cart{}
	for rows.Next() {
		cartList = append(cartList, entity.CartScanRow(nil,rows))
	}
	return cartList
}

func (this *cartModel) CountItems(searchStr, product_idStr string){
	if searchStr != ``{
		this.Where = `WHERE (product_title LIKE '%` + searchStr + `%' OR cart_id LIKE '%` + searchStr + `%')`
	}
	if product_idStr != `` {
		if this.Where != `` {
			this.Where += ` OR cart_product = ` + product_idStr
		}else {
			this.Where = `WHERE cart_product = ` + product_idStr
		}
	}

	this.Query  = `SELECT count(*) FROM cart `+ this.Where

	row := Crud.GetRow(this.Query ,[]interface{}{})
	err := row.Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = 0
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}

