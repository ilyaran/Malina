package entity

import (
	"time"
	"database/sql"
	"github.com/lib/pq"
)

type Product struct {
	Id int64			`json:"id"`
	Code string			`json:"code"`
	Category int64		`json:"category"`
	Parameter []int64	`json:"parameter"`
	Img []string		`json:"img"`
	Title string		`json:"title"`
	Description string	`json:"description"`
	Price float64		`json:"price"`
	Price1 float64		`json:"price1"`
	Price2 float64		`json:"price2"`
	Quantity float64	`json:"quantity"`
	Sold int64			`json:"sold"`
	Views int64			`json:"views"`
	Created time.Time	`json:"created"`
	Updated time.Time	`json:"updated"`
	Enable bool			`json:"enable"`
}

func (s *Product) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &Product{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *Product) Item(row *sql.Row) Scanable {
	item := Product{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}


func (s *Product) GetId() int64 {return s.Id}
func (s *Product) SetId(id int64) {s.Id = id}

func (s *Product) GetCode() string {return s.Code}
func (s *Product) GetCategory() int64 {return s.Category}
func (s *Product) GetParameter() []int64 {return s.Parameter}
func (s *Product) GetImg() []string {return s.Img}
func (s *Product) GetTitle() string {return s.Title}
func (s *Product) GetDescription() string {return s.Description}
func (s *Product) GetPrice() float64 {return s.Price}
func (s *Product) GetPrice1() float64 {return s.Price1}
func (s *Product) GetQuantity() float64 {return s.Quantity}
func (s *Product) GetSold() int64 {return s.Sold}
func (s *Product) GetViews() int64 {return s.Views}
func (s *Product) GetCreated() time.Time {return s.Created}
func (s *Product) GetUpdated() time.Time {return s.Updated}
func (s *Product) GetEnable() bool {return s.Enable}



func (s *Product) Scanning(row *sql.Row,rows *sql.Rows)byte {
	var err error
	if row!=nil {
		err = row.Scan(
			&s.Id,
			&s.Code,
			&s.Category,
			pq.Array(&s.Parameter),
			pq.Array(&s.Img),
			&s.Title,
			&s.Description,
			&s.Price,
			&s.Price1,
			&s.Price2,
			&s.Quantity,
			&s.Sold,
			&s.Views,
			&s.Created,
			&s.Updated,
			&s.Enable,
		)
	}
	if rows!=nil {
		err = rows.Scan(
			&s.Id,
			&s.Code,
			&s.Category,
			pq.Array(&s.Parameter),
			pq.Array(&s.Img),
			&s.Title,
			&s.Description,
			&s.Price,
			&s.Price1,
			&s.Price2,
			&s.Quantity,
			&s.Sold,
			&s.Views,
			&s.Created,
			&s.Updated,
			&s.Enable,
		)
	}
	if err == sql.ErrNoRows {
		//panic(err)
		return 'n'
	}
	if err != nil {
		panic(err)
		return 'e'
	}

	return 'o'
}


