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


import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/ilyaran/Malina/app"
)

type State struct {
	Id        int64    `json:"id"`
	Title     string   `json:"title"`
	Code      string   `json:"code"`
	Img       []string `json:"img"`
	Sort      int64    `json:"sort"`
	Flag      string   `json:"flag"`
	PhoneCode string   `json:"phone_code"`
	Enable		bool	`json:"enable"`
}
func(s *State)SetId(id int64){s.Id=id}
func(s *State)GetId()int64{return s.Id}
func(s *State)GetTitle()string{return s.Title}
func(s *State)GetCode()string{return s.Code}
func(s *State)GetImg()[]string{return s.Img}
func(s *State)GetSort()int64{return s.Sort}
func(s *State)GetPhone()string{return s.PhoneCode }
func(s *State)SetFlag(){
	if s.Img!=nil && len(s.Img) > 0 {
		s.Flag = s.Img[0]
	}else {
		s.Flag = app.Url_no_image
	}

}

func (s *State) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &State{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *State) Item(row *sql.Row) Scanable {
	item := State{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}
func(s *State)Scanning (row *sql.Row,rows *sql.Rows)byte{
	var err error
	if row!=nil{
		err = row.Scan(

			&s.Id,
			&s.Title,
			&s.Code,

			pq.Array(&s.Img),
			&s.Sort,

			&s.Enable,
		)
	}
	if rows!=nil{
		err = rows.Scan(
			&s.Id,
			&s.Title,
			&s.Code,

			pq.Array(&s.Img),
			&s.Sort,

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





