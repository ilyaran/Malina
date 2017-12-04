package entity

import (
	"time"
	"database/sql"
	"github.com/lib/pq"
)

type Permission struct {
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

func (s *Permission) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &Permission{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *Permission) Item(row *sql.Row) Scanable {
	item := Permission{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}


func (s *Permission) GetId() int64 {return s.Id}
func (s *Permission) SetId(id int64) {s.Id = id}

func (s *Permission) GetCode() string {return s.Code}
func (s *Permission) GetCategory() int64 {return s.Category}
func (s *Permission) GetParameter() []int64 {return s.Parameter}
func (s *Permission) GetImg() []string {return s.Img}
func (s *Permission) GetTitle() string {return s.Title}
func (s *Permission) GetDescription() string {return s.Description}
func (s *Permission) GetPrice() float64 {return s.Price}
func (s *Permission) GetPrice1() float64 {return s.Price1}
func (s *Permission) GetQuantity() float64 {return s.Quantity}
func (s *Permission) GetSold() int64 {return s.Sold}
func (s *Permission) GetViews() int64 {return s.Views}
func (s *Permission) GetCreated() time.Time {return s.Created}
func (s *Permission) GetUpdated() time.Time {return s.Updated}
func (s *Permission) GetEnable() bool {return s.Enable}



func (s *Permission) Scanning(row *sql.Row,rows *sql.Rows)byte {
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


