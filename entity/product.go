/**
 * Product entity class.  Malina eCommerce application
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
	"github.com/ilyaran/Malina/config"
	"encoding/json"
	"strings"
)

type Product struct {
	id                int64
	category_id       int64
	price             float64
	price1            float64
	quantity	  float64
	code              string
	title             string
	description       string
	created           time.Time
	updated           time.Time
	enable            bool
	img               []string
	sold          	  int64
	watch 			int64
	like 		int64
}
func NewProduct(
	id int64,
	category_id int64,
	price float64,
	price1 float64,
	code string,
	title string,
	description string,
	enable bool) *Product {

	return &Product{
		id:id,
		category_id:category_id,
		price:price,
		price1:price1,
		code:code,
		title:title,
		description:description,
		enable:enable}
}
func (s *Product) Exec() []interface{} {
	return []interface{}{
		s.category_id,
		s.title,
		s.description,
		s.code,
		s.price,
		s.price1,
		s.enable,
	}
}
func (s *Product) ExecWithId() []interface{} {
	return []interface{}{
		s.id,
		s.category_id,
		s.title,
		s.description,
		s.code,
		s.price,
		s.price1,
		s.enable,
	}
}
func ProductScanRow(row *sql.Row) *Product {
	var s = &Product{}
	var images_sting string
	err := row.Scan(
		&s.id,
		&s.category_id,
		&s.price,
		&s.price1,
		&s.code,
		&s.title,
		&s.description,
		&s.created,
		&s.updated,
		&s.enable,
		&images_sting,
	)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	if images_sting != `{}` {
		s.img = strings.Split(images_sting[1:len(images_sting)-1],",")
	}
	return s
}
func ProductScanRows(rows *sql.Rows) *Product {
	var s = &Product{}
	var images_sting string
	err := rows.Scan(
		&s.id,
		&s.category_id,
		&s.price,
		&s.price1,
		&s.code,
		&s.title,
		&s.description,
		&s.created,
		&s.updated,
		&s.enable,
		&images_sting,
	)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	if  images_sting != `{}` {
		s.img = strings.Split(images_sting[1:len(images_sting)-1],",")
	}
	return s
}

func (this *Product)  Get_id() int64 {
	return this.id
}
func (this *Product)  Get_category_id() int64 {
	return this.category_id
}
func (this *Product)  Get_price() float64 {
	return this.price
}
func (this *Product)  Get_price1() float64 {
	return this.price1
}
func (this *Product)  Get_quantity() float64 {
	return this.quantity
}
func (this *Product)  Total() float64 {
	return this.price * this.quantity
}
func (this *Product)  Get_code() string {
	return this.code
}
func (this *Product)  Get_title() string {
	return this.title
}
func (this *Product)  Get_description() string {
	return this.description
}
func (this *Product)  Get_created() time.Time {
	return this.created
}
func (this *Product)  Get_updated() time.Time {
	return this.updated
}
func (this *Product)  Get_enable() bool {
	return this.enable
}
func (this *Product)  Get_img() []string {
	return this.img
}
func (this *Product)  Get_logo() string {
	if len(this.img) > 0 {
		return app.Base_url() + app.Upload_path() + this.img[0]
	}
	return app.No_image()
}

//**********************************
func (t *Product) JsonEncode() string {
	b, err := json.Marshal(t.GetJsonAdaptedType())
	if err != nil {
		return ""
		//log.Fatal(err)
	}
	//log.Printf("%s", b)
	return string(b)
}
func (t *Product) JsonDecode(jsonString []byte) *Product {
	o := new(Product)
	err := json.Unmarshal(jsonString, &o)
	if err != nil {
		//log.Fatal(err)
		return nil
	}
	//log.Printf("%#v", t)
	return o
}
func (t *Product) GetJsonAdaptedType() *JsonProduct {
	return &JsonProduct{
		t.id,
		t.category_id,
		t.price,
		t.price1,
		t.code,
		t.title,
		t.description,
		t.created,
		t.updated,
		t.enable,
		t.img,
	}
}

type JsonProduct struct {
	Id                int64                `json:"id"`
	Category_id       int64        `json:"category_id"`
	Price             float64                `json:"price"`
	Price1            float64                `json:"price1"`
	Code              string                `json:"code"`
	Title             string                `json:"title"`
	Description       string        `json:"description"`
	Created           time.Time        `json:"created"`
	Updated           time.Time        `json:"updated"`
	Enable            bool                `json:"enable"`
	Img               []string                `json:"img"`
}

func (t *Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.GetJsonAdaptedType())
}

func (t *Product) UnmarshalJSON(b []byte) error {
	temp := &JsonProduct{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	t.id = temp.Id
	t.category_id = temp.Category_id
	t.price = temp.Price
	t.price1 = temp.Price1
	t.code = temp.Code
	t.title = temp.Title
	t.description = temp.Description
	t.created = temp.Created
	t.updated = temp.Updated
	t.enable = temp.Enable
	t.img = temp.Img

	return nil
}


//*********** collection
func GetProductCollectionFromJsonString(jsonString []byte) []*Product {

	ic := new(ProductCollection)
	err := ic.FromJson(jsonString)
	if err != nil {
		panic(err)
		return nil
	}
	var productList = []*Product{}
	for _, v := range ic.Pool {
		productList = append(productList, &Product{
			id:v.Id,
			category_id       :v.Category_id,
			price             :v.Price,
			price1            :v.Price1,
			code              :v.Code,
			title             :v.Title,
			description       :v.Description,
			created           :v.Created,
			updated           :v.Updated,
			enable            :v.Enable,
			img               :v.Img,
		})
	}

	return productList
}

type ProductCollection struct {
	Pool []*JsonProduct
}

func (mc *ProductCollection) FromJson(jsonStr []byte) error {
	var data = &mc.Pool
	b := jsonStr
	return json.Unmarshal(b, data)
}
