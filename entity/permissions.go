package entity

import "database/sql"

type Permission struct {
	Id int64		`json:"id"`
	Data string		`json:"data"`
	Position *Position	`json:"position"`
}
func (s *Permission)GetId()int64{return s.Id}
func (s *Permission)GetData()string{return s.Data}
func (s *Permission)GetPosition()*Position{return s.Position}

func (s *Permission)SetId(id int64){s.Id=id}
func (s *Permission)SetData(data string){s.Data=data}
func (s *Permission)SetPosition(Position *Position){s.Position=Position}

func (this *Permission) ExecWithId() []interface{} {
	return []interface{}{
		this.Id,
		this.Data,
		this.Position.GetId(),
	}
}
func (this *Permission) Exec() []interface{} {
	return []interface{}{
		this.Data,
		this.Position.GetId(),
	}
}
func PermissionScanRow(row *sql.Row) *Permission {
	s := &Permission{Position:&Position{}}
	err := row.Scan(
		&s.Id,
		&s.Position.Id,
		&s.Position.Title,
		&s.Data,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}
func PermissionScanRows(rows *sql.Rows) *Permission {
	s := &Permission{Position:&Position{}}
	err := rows.Scan(
		&s.Id,
		&s.Position.Id,
		&s.Position.Title,
		&s.Data,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}

func NewPermission(id int64, data string, position int64) *Permission {
	return &Permission{Id:id,Data:data,Position:&Position{Id:position}}
}