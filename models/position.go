package model

import (
	"Malina/entity"
	"database/sql"
)

var FieldsPosition = `
	position_id,
	position_parent,
	coalesce((SELECT a.position_title FROM position a WHERE a.position_id = position.position_parent),'user'),
	position_title,
	position_sort,
	position_enable`

var PositionModel = new(positionModel)
type positionModel struct {}

func (this *positionModel)Get(id int64)*entity.Position{
	var querySql = `SELECT `+FieldsPosition+` FROM position
	WHERE position_id = $1 LIMIT 1`
	row := Crud.GetRow(querySql, []interface{}{id})
	return entity.PositionScanRow(row)
}

func (this *positionModel)Del(id int64)int64{
	var querySql = `
	WITH t AS (
		UPDATE account SET account_position = 1 WHERE account_position = $1
	)
	DELETE FROM position WHERE position_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *positionModel)Add(position *entity.Position)int64{
	var querySql = `
	INSERT INTO position (
		position_parent,
		position_title,
		position_sort,
		position_enable
	)
	VALUES ($1,$2,$3,$4) RETURNING position_id`
	return Crud.Insert(querySql, position.Exec())
}

func (this *positionModel)Edit(position *entity.Position)int64{

	var querySql = `
	UPDATE position SET (
		position_parent,
		position_title,
		position_sort,
		position_enable
	) = ($2,$3,$4,$5)
	WHERE position_id = $1`
	return Crud.Update(querySql, position.ExecWithId())
}

func (this *positionModel) GetList(search,page,perPage string,order_by string) []*entity.Position {
	var where = this.getSearchSql(search)
	var limit = ""
	if page != "" && perPage != ""{
		limit = `LIMIT ` + perPage + ` OFFSET ` + page
	}

	if order_by == "" {order_by = "position_sort ASC"}
	var querySql = `
	SELECT ` + FieldsPosition + ` FROM position
	` + where + `
	ORDER BY ` + order_by + ` ` + limit

	rows := Crud.GetRows(querySql, []interface{}{})
	defer rows.Close()

	var positionList = []*entity.Position{}

	for rows.Next() {
		positionList = append(positionList, entity.PositionScanRows(rows))
	}
	return positionList
}

func (this *positionModel) CountItems(search string) int64{
	var querySql = `SELECT count(*) FROM position ` + this.getSearchSql(search)
	var all int64
	row := Crud.GetRow(querySql,[]interface{}{})
	err := row.Scan(&all)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
		return -1
	}
	return all
}

func (this *positionModel)getSearchSql(search string)string{
	if search != "" {
		return `WHERE position_title LIKE '%`+search+`%'`
	}
	return ""
}

