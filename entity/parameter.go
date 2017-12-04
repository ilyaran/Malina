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
	"time"
	"database/sql"
)

type Parameter struct {
	Id int64							`json:"id"`
	Parent int64						`json:"parent"`
	Title string						`json:"title"`
	Description string					`json:"description"`
	Sort int64							`json:"sort"`
	Value string						`json:"value"`
	Created time.Time					`json:"created"`
	Updated time.Time					`json:"updated"`
	Enable bool							`json:"enable"`

	Level int							`json:"level"`

	Children            []*Parameter 	`json:"-"`
	Ancestors           []*Parameter	`json:"-"`
	Descendants         []*Parameter	`json:"-"`
}

func (s *Parameter) GetId() int64 {return s.Id}
func (s *Parameter) GetParent() int64 {return s.Parent}
func (s *Parameter) SetId(id int64) {s.Id = id}
func (s *Parameter) GetLevel() int {return s.Level}
func (s *Parameter) SetLevel(level int) {s.Level = level}
func (s *Parameter) GetTitle() string {return s.Title}
func (s *Parameter) GetDescription() string {return s.Description}
func (s *Parameter) GetSort() int64 {return s.Sort}
func (s *Parameter) GetValue() string {return s.Value}
func (s *Parameter) GetCreated() time.Time {return s.Created}
func (s *Parameter) GetUpdated() time.Time {return s.Updated}
func (s *Parameter) GetEnable() bool {return s.Enable}


func (s *Parameter) Scanning(row *sql.Row,rows *sql.Rows)byte {
	var err error
	if row!=nil {
		err = row.Scan(
			&s.Id,
			&s.Parent,
			&s.Title,
			&s.Description,
			&s.Sort,
			&s.Value,
			&s.Created,
			&s.Updated,
			&s.Enable,
		)
	}
	if rows!=nil {
		err = rows.Scan(
			&s.Id,
			&s.Parent,
			&s.Title,
			&s.Description,
			&s.Sort,
			&s.Value,
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
func (s *Parameter) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &Parameter{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *Parameter) Item(row *sql.Row) Scanable {
	item := Parameter{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}
func (s *Parameter) SetAncestorsDescendants(parent int64,getDescendants bool)int{
	if getDescendants {
		if s.Parent == parent {
			s.Descendants = append(s.Descendants, s)
			//s.closure(s.Id, getDescendants)
			return 1
		}
	}else {
		if s.Id == parent {
			s.Ancestors = append([]*Parameter{s}, s.Ancestors...)
			//s.closure(s.Id, getDescendants)
			return 2
		}
	}
	return 3
}

func (this *Parameter) SetChildren(item Hierarchical) {
	if this.Children == nil{
		this.Children = []*Parameter{}
	}
	this.Children = append(this.Children, item.(*Parameter))
}

func (this *Parameter) SetAncestors(item Hierarchical) {
	if this.Ancestors == nil{
		this.Ancestors = []*Parameter{}
	}
	this.Ancestors = append([]*Parameter{item.(*Parameter)}, this.Ancestors...)
}

func (this *Parameter) SetDescendants(item Hierarchical) {
	if this.Descendants == nil{
		this.Descendants = []*Parameter{}
	}
	this.Descendants = append(this.Descendants, item.(*Parameter))
}

