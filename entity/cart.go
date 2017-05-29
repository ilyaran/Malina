/**
 * Cart entity class.  Malina eCommerce application
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
package entity

import (
	"time"
	"database/sql"
	"strings"
	"strconv"
)

type CartPublic struct {
	Product_id int64  	`json:"product_id"`
	Product_img []string	`json:"product_img"`
	Product_title string	`json:"product_title"`
	Product_price float64	`json:"product_price"`
	Product_quantity float64 `json:"product_quantity"`
	Subtotal float64	`json:"subtotal"`
}
type Cart struct {
	id string
	product []*Product
	quantity []float64
	buy_now [] bool
	created time.Time
}

func (s *Cart) Get_id()string { return s.id }
func (s *Cart) Get_product()[]*Product { return s.product }
func (s *Cart) Get_quantity()[]float64 { return s.quantity }
func (s *Cart) Get_buy_now()[]bool { return s.buy_now }
func (s *Cart) Get_created()time.Time { return s.created }

func NewCart(id string,product int64,quantity float64,buy_now bool) *Cart {
	return &Cart{}
}
func (this *Cart) ExecWithId() []interface{} {
	return []interface{}{
		this.id,
		this.product[0].id,
		this.quantity,
		this.buy_now,
	}
}
func CartScanRow(row *sql.Row,rows *sql.Rows) *Cart {
	var s = &Cart{product:[]*Product{}}
	var product_ids string
	var product_titles string
	var product_prices string
	var product_prices1 string
	var quantities string
	var buy_nows string

	var err error
	if row != nil && rows == nil{
		err = row.Scan(
			&s.id,
			&product_ids,
			&product_titles,
			&product_prices,
			&product_prices1,
			&quantities,
			&buy_nows,
			&s.created,
		)
	}
	if rows != nil && row == nil{
		err = rows.Scan(
			&s.id,
			&product_ids,
			&product_titles,
			&product_prices,
			&product_prices1,
			&quantities,
			&buy_nows,
			&s.created,
		)
	}

	var product_id_array []string = strings.Split(product_ids,"|")
	var product_title_array []string = strings.SplitN(product_titles,"|",len(product_id_array))
	var product_price_array []string = strings.Split(product_prices,"|")
	var product_price1_array []string = strings.Split(product_prices1,"|")
	var quantities_array []string = strings.Split(quantities,"|")
	var buy_nows_array []string = strings.Split(buy_nows,"|")
	s.quantity = []float64{}
	s.buy_now = []bool{}
	for k, v := range product_id_array{
		id,_:=strconv.ParseInt(v,10,64)
		price,_:=strconv.ParseFloat(product_price_array[k],64)
		price1,_:=strconv.ParseFloat(product_price1_array[k],64)
		s.product = append(
			s.product,
			&Product{
				id:id,
				title:product_title_array[k],
				price:price,
				price1:price1,
			})
		q,_:=strconv.ParseFloat(quantities_array[k],64)
		s.quantity = append(s.quantity,q)
		bn,_:= strconv.ParseBool(buy_nows_array[k])
		s.buy_now = append(s.buy_now,bn)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}



