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

import "database/sql"

type Scanable interface {
	GetId() int64
	SetId(id int64)
	Scanning(row *sql.Row,rows *sql.Rows)byte
	Appending(rows *sql.Rows,size int64)*[]Scanable
	Item(row *sql.Row) Scanable
}

type Hierarchical interface {
	Scanable
	GetParent()int64
	SetLevel(level int)
	SetChildren(item Hierarchical)
	SetDescendants(item Hierarchical)
	SetAncestors(item Hierarchical)

}