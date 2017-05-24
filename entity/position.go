package entity

import (
	"database/sql"
)

import (
	"fmt"
	"encoding/json"
)

type Position struct {
	Id                    int64			`json:"Id"`
	Parent                *Position			`json:"parent"`
	Title                 string			`json:"title"`
	Sort                  int64			`json:"sort"`
	Level                 int			`json:"level"`
	Enable                bool			`json:"enable"`
	PermissionsMap        map[int64]*Permission

	DescendantObjectsList []*Position
	DescendantIdsMap      map[int64]bool
	Path                  string
	DescendantIdsString   string
}

func (s *Position)GetId() int64 {return s.Id}
func (s *Position)GetPermissions() map[int64]*Permission {return s.PermissionsMap}
func (s *Position)GetParent() *Position {return s.Parent}
func (s *Position)GetTitle() string {return s.Title}
func (s *Position)GetLevel() int {return s.Level}
func (s *Position)GetSort() int64 {return s.Sort}
func (s *Position)GetEnable() bool {return s.Enable}

func (s *Position)SetLevel(level int) {s.Level = level}
func (this *Position) Get_descendants() []*Position {return this.DescendantObjectsList}
func (this *Position) GetDescendantIdsString() string {return this.DescendantIdsString}
func (this *Position) GetDescendantIdsMap() map[int64]bool{return this.DescendantIdsMap}

func (this *Position) Set_descendants(treeList []*Position) {
	this.DescendantObjectsList = []*Position{}
	this.DescendantIdsString = ""
	this.DescendantIdsMap = map[int64]bool{}
	var anon func(int64)
	anon = func(parent int64) {
		for _, v := range treeList {
			if v.Parent.Id == parent {
				this.DescendantObjectsList = append(this.DescendantObjectsList, v)
				this.DescendantIdsString += fmt.Sprintf(",%v", v.Id)
				this.DescendantIdsMap[v.Id]=false
				anon(v.Id)
			}
		}
	}
	anon(this.Id)
	if this.DescendantIdsString!=``{
		this.DescendantIdsString = this.DescendantIdsString[1:]
	}
}

func (this *Position) Exec() []interface{} {
	return []interface{}{
		this.Parent.Id,
		this.Title,
		this.Sort,
		this.Enable,
		this.Id,
	}
}

func PositionScan(row *sql.Row, rows *sql.Rows, TreeMap map[int64]*Position) *Position {
	s := &Position{}
	var parent sql.NullInt64
	var permissions []byte
	var err error
	if row != nil {
		err = row.Scan(&s.Id, &parent, &s.Title, &s.Sort, &s.Enable, &permissions)
	}
	if row == nil && rows != nil {
		err = rows.Scan(&s.Id, &parent, &s.Title, &s.Sort, &s.Enable, &permissions)
	}
	if v,ok := TreeMap[parent.Int64];ok{
		s.Parent = v
	}else {
		s.Parent = &Position{}
	}
	if len(permissions) > 0 {

		err = json.Unmarshal(permissions, &s.PermissionsMap)
		if err != nil {
			panic(err)
			return nil
		}
		for _,v := range s.PermissionsMap {
			v.Position = s
		}
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

func NewPosition(id int64, parent int64, title string, sort int64, enable bool) *Position {
	return &Position{Id:id, Parent:&Position{Id:parent}, Title:title, Sort:sort, Enable:enable}
}


