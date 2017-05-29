package model

import (
	"github.com/ilyaran/Malina/entity"
	"database/sql"
)

var FieldsPosition = `
	position_id,
	position_parent,
	position_title,
	position_sort,
	position_enable,
	(
		select '{'||string_agg('"'||cast(p.id AS TEXT)||'"' || ':' || cast(row_to_json(p) AS TEXT),',') || '}'
		from (
		       select permission_id id, permission_data
		       from permission WHERE permission_position = position_id
		     ) as p
	) AS permissions
	`

var PositionModel = new(positionModel)
type positionModel struct {
	Where string
	Query string
	All int64
	Exec []interface{}
	//Id int64
}

func ( this *positionModel ) Get ( id int64, TreeMap map[int64]*entity.Position ) *entity.Position {
	this.Query = `SELECT `+FieldsPosition+` FROM position WHERE id = $1 LIMIT 1 `
	row := Crud.GetRow(this.Query, []interface{}{id})
	return entity.PositionScan(row,nil,TreeMap)
}
func (this *positionModel)Del() int64 {
	this.Query = `DELETE FROM position WHERE position_id = $1`
	return Crud.Delete(this.Query, this.Exec)
}
func (this *positionModel)Add(positionParentId int64) int64 {
	this.Query = `
	INSERT INTO position `
	if positionParentId > 0 {
		this.Query += `
		(position_parent,
		position_title,
		position_sort,
		position_enable)
		VALUES ($1,$2,$3,$4)
		RETURNING position_id`
	}else {
		this.Query += `
		(position_title,
		position_sort,
		position_enable)
		VALUES ($1,$2,$3)
		RETURNING position_id`
	}
	//fmt.Println(this.Query)
	return Crud.Insert(this.Query, this.Exec)
}

func (this *positionModel)Edit(positionParentId int64) int64 {
	this.Query = `
	UPDATE position SET `
	if positionParentId > 0 {
		this.Query += `
		(position_parent,
		position_title,
		position_sort,
		position_enable)
		= ($1,$2,$3,$4)
		WHERE position_id = $5`
	}else {
		this.Query = `
		(position_title,
		position_sort,
		position_enable)
		= ($1,$2,$3)
		WHERE position_id = $4`
	}
	return Crud.Update(this.Query, this.Exec)
}

func (this *positionModel) GetList(search,page,perPage string,order_by string,TreeMap map[int64]*entity.Position) []*entity.Position {

	this.CountItems(search)

	var limit = ""
	if page != "" && perPage != ""{
		limit = ` LIMIT ` + perPage + ` OFFSET ` + page
	}

	if order_by == "" {order_by = "position_sort ASC"}
	this.Query = `
	SELECT ` + FieldsPosition + ` FROM position
	` + this.Where + `
	ORDER BY ` + order_by + ` ` + limit

	rows := Crud.GetRows(this.Query, []interface{}{})

	var positionList = []*entity.Position{}

	for rows.Next() {
		positionList = append(positionList, entity.PositionScan(nil,rows,TreeMap))
	}

	return positionList
}

func (this *positionModel) CountItems(search string) {
	if search != `` {
		this.Where = ` WHERE position_title LIKE '%` + search + `%'`
	}

	this.Query = `SELECT count(*) FROM position ` + this.Where

	row := Crud.GetRow(this.Query,[]interface{}{})
	err := row.Scan(&this.All)
	if err == sql.ErrNoRows {
		this.All = 0
	}
	if err != nil {
		panic(err)
		this.All = -1
	}
}

