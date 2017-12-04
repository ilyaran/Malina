package models



type Parameter struct {

}

func (this *Parameter)Upsert(action byte, args... interface{})int64{
	if args[0].(int64) == 0 {
		args[0]=nil
	}
	if action=='a'{
		return CrudGeneral.Insert(`
		INSERT INTO parameter
		(parameter_parent,parameter_title,parameter_description,parameter_sort,parameter_enable)
		VALUES
		($1,$2,$3,$4,$5)
		RETURNING parameter_id
		`, args...)
	}
	return CrudGeneral.Update(`
		UPDATE parameter SET
		(parameter_parent,parameter_title,parameter_description,parameter_sort,parameter_enable)
		=
		($1,$2,$3,$4,$5)
		WHERE parameter_id = $6
		`, args...)
}
