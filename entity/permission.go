package entity


import "database/sql"

type Permission struct {
	Id int64		`json:"id"`
	Data string		`json:"data"`
	Position *Position	`json:"position"`
}

func (s *Permission)SetData(data string){s.Data = data}
func (s *Permission)SetPosition(position *Position){s.Position = position}

func (s *Permission)GetId()int64{return s.Id}
func (s *Permission)GetData()string{return s.Data}
func (s *Permission)GetPosition()*Position{return s.Position}

func PermissionScan(row *sql.Row, rows *sql.Rows) *Permission {
	s := &Permission{}
	var position_id int64
	var positionTitle string
	var err error
	if row != nil {err = row.Scan(&s.Id, &position_id,&positionTitle,&s. Data)}
	if rows != nil {err = rows.Scan(&s.Id, &position_id,&positionTitle,&s.Data)}
	if position_id > 0 && positionTitle != ``{
		s.Position = &Position{Id : position_id,Title:positionTitle}
	}
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	return s
}

func NewPermission(data string, position int64) *Permission {
	return &Permission{Data:data,Position:&Position{Id:position}}
}