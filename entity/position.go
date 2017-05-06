package entity

import (
	"database/sql"
)

type Position struct {
	Id int64	 `json:"id"`
	Parent *Position `json:"parent"`
	Title string	 `json:"title"`
	Sort int64	 `json:"sort"`
	Level int	 `json:"level"`
	Enable bool	 `json:"enable"`
}
func (s *Position)GetId()int64{return s.Id}
func (s *Position)GetParent()*Position{return s.Parent}
func (s *Position)GetTitle()string{return s.Title}
func (s *Position)GetLevel()int{return s.Level}
func (s *Position)GetSort()int64{return s.Sort}
func (s *Position)GetEnable()bool{return s.Enable}

func (s *Position)SetLevel(level int){s.Level = level}

func (this *Position) ExecWithId() []interface{} {
	return []interface{}{
		this.Id,
		this.Parent.GetId(),
		this.Title,
		this.Sort,
		this.Enable,
	}
}
func (this *Position) Exec() []interface{} {
	return []interface{}{
		this.Parent.GetId(),
		this.Title,
		this.Sort,
		this.Enable,
	}
}
func PositionScanRow(row *sql.Row) *Position {
	s := &Position{Parent:&Position{}}
	err := row.Scan(&s.Id,&s.Parent.Id,&s.Parent.Title,&s.Title,&s.Sort,&s.Enable)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}
func PositionScanRows(rows *sql.Rows) *Position {
	s := &Position{Parent:&Position{}}
	err := rows.Scan(&s.Id,&s.Parent.Id,&s.Parent.Title,&s.Title,&s.Sort,&s.Enable)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}

func NewPosition(id int64, parent int64, title string,sort int64,enable bool) *Position {
	return &Position{Id:id,Parent:&Position{Id:parent},Title:title,Sort:sort,Enable:enable}
}










