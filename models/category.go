package models

type Category struct {

}

func (this *Category)Upsert(action byte,imgUrlsString string, args... interface{})int64{
	if imgUrlsString!=""{
		imgUrlsString = `ARRAY[`+imgUrlsString+`]::VARCHAR(255)[]`
	}else {
		imgUrlsString = "NULL"
	}
	if args[0].(int64) == 0 {
		args[0]=nil
	}
	if action=='a'{
		return CrudGeneral.Insert(`
		INSERT INTO category
		(category_parent,category_title,category_description,category_sort,category_enable,category_img)
		VALUES
		($1,$2,$3,$4,$5,`+imgUrlsString+`)
		RETURNING category_id
		`, args...)
	}
	return CrudGeneral.Update(`
		UPDATE category SET
		(category_parent,category_title,category_description,category_sort,category_enable,category_img)
		=
		($1,$2,$3,$4,$5,`+imgUrlsString+`)
		WHERE category_id = $6
		`, args...)
}

func (this *Category)UpdateAfterDeleted(args... interface{})int64{
	if args[1].(int64) == 0 {
		args[1]=nil
	}
	return CrudGeneral.Update(`
		UPDATE category SET
		category_parent = $2
		WHERE category_parent = $1
		`, args...)
}















