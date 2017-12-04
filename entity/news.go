/**
 *
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ package entity
/*

import (
	"time"
	"database/sql"
	"github.com/lib/pq"
)

type News struct {
	Id int64			`json:"id"`
	Account int64		`json:"account"`
	Category int64		`json:"category"`
	Img []string		`json:"img"`
	Comment []int64		`json:"comment"`
	Title string		`json:"title"`
	Description string	`json:"description"`
	Like int64			`json:"like"`
	Views int64			`json:"views"`
	Created time.Time	`json:"created"`
	Updated time.Time	`json:"updated"`
	Enable bool			`json:"enable"`
}
func (s *News) Append(list []Scanable,rows *sql.Rows){
	item := &News{}
	if item.Scanning(nil,rows) == 'o'{
		list = append(list, item)
	}
}
func (s *News) GetId() int64 {return s.Id}
func (s *News) SetId(id int64) {s.Id = id}

func (s *News) GetAccount() int64 {return s.Account}
func (s *News) GetCategory() int64 {return s.Category}
func (s *News) GetComment() []int64 {return s.Comment}
func (s *News) GetImg() []string {return s.Img}
func (s *News) GetTitle() string {return s.Title}
func (s *News) GetDescription() string {return s.Description}
func (s *News) GetLike() int64 {return s.Like}
func (s *News) GetViews() int64 {return s.Views}
func (s *News) GetCreated() time.Time {return s.Created}
func (s *News) GetUpdated() time.Time {return s.Updated}
func (s *News) GetEnable() bool {return s.Enable}

func (s *News) Scanning(row *sql.Row,rows *sql.Rows)byte {
	var err error
	if row!=nil {
		err = row.Scan(
			&s.Id,
			&s.Account,
			&s.Category,
			pq.Array(&s.Img),
			pq.Array(&s.Comment),
			&s.Title,
			&s.Description,
			&s.Like,
			&s.Views,
			&s.Created,
			&s.Updated,
			&s.Enable,
		)
	}
	if rows!=nil {
		err = rows.Scan(
			&s.Id,
			&s.Account,
			&s.Category,
			pq.Array(&s.Img),
			pq.Array(&s.Comment),
			&s.Title,
			&s.Description,
			&s.Like,
			&s.Views,
			&s.Created,
			&s.Updated,
			&s.Enable,
		)
	}
	if err == sql.ErrNoRows {
		panic(err)
		return 'n'
	}
	if err != nil {
		panic(err)
		return 'e'
	}

	return 'o'
}


*/
