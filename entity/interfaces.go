package entity

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