package model

import (
	"database/sql"
	"Malina/entity"
)

const FieldsApp  =`
		Site_name,
		Admin_email,
		Admin_email_pass,

		Crypt_salt,
		Session_expiration,
		Cookie_expiration,

		Cookie_name,
		Per_page,
		Radius,

		Language_list,
		Assets_path,
		Assets_backend_path,

		Assets_public_path,
		 Assets_cabinet_path,
		 Upload_path,

		  No_image,
		 Image_path,
		 Image_max_width,

		 Image_max_height,
		 Image_resize_width,
		 Image_resize_height,

		 Image_max_size,
		 Image_upload_limit,
		 Login_attempts,

		 Login_attempts_period,
		 Recaptcha_secret,
		 Recaptcha_public,

		 Use_recaptcha_register,
		 Use_recaptcha_login,
		 Use_recaptcha_forgot,

		 Use_recaptcha_change_password,
		 Email_activation,
		 Email_activation_expire,

		 Email_app_details,
		 Uri_auth_deny,
		 Uri_auth_banned,

		 Uri_permission,
		 Uri_permission_ajax,
		 Uri_permission_add,

		 Uri_permission_edit,
		 Uri_permission_del,
		 Uri_permission_get,

		 Uri_position,
		 Uri_position_ajax,
		 Uri_position_add,

		 Uri_position_edit,
		 Uri_position_del,
		 Uri_position_get,

		 Uri_app,
		 Uri_app_ajax,
		 Uri_app_add,

		 Uri_app_edit,
		 Uri_app_del,
		 Uri_app_get,

		 Uri_cabinet,
		 Uri_cabinet_add,
		 Uri_cabinet_edit,

		 Uri_cabinet_del,
		 Uri_cabinet_get,
		 Uri_auth_login,

		 Uri_auth_register,
		 Uri_auth_activation,
		 Uri_auth_forgot,

		 Uri_auth_change_password,
		 Uri_auth_logout,
		 Uri_auth_delete_app,

		 PasswordMinLength,
		 PasswordMaxLength 
`
var AppModel *appModel = new(appModel)
type appModel struct {}
func (this *appModel)Get(id int64)*entity.App{
	var querySql = `SELECT `+FieldsApp+` FROM app
	WHERE app_id = $1 LIMIT 1`
	row := Crud.GetRow(querySql, []interface{}{id})
	return entity.ScanRowApp(row)
}

func (this *appModel)Del(id int64)int64{
	var querySql = `
	DELETE FROM app WHERE app_id = $1`
	return Crud.Delete(querySql, []interface{}{id})
}

func (this *appModel)Add(app *entity.App)int64{
	var querySql = `
	INSERT INTO app ()
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,$61,$62,$63,$64,$65,$66,$67) RETURNING app_id`
	return Crud.Insert(querySql, app.Exec())
}

func (this *appModel)Edit(app *entity.App)int64{
	var querySql = `
	UPDATE app
	SET (`+FieldsApp+`)
	=
	($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,$61,$62,$63,$64,$65,$66,$67,$68)
	WHERE app_id = $1`
	return Crud.Update(querySql, app.ExecWithId())
}

func (this *appModel) GetList(search,page,perPage string,order_by string) []*entity.App {
	var where = this.getSearchSql(search)
	if order_by == "" {
		order_by = "app_id ASC"
	}
	var querySql = `
	SELECT ` + FieldsApp + `
	FROM app
	` + where + `
	ORDER BY ` + order_by + `
	LIMIT ` + perPage + ` OFFSET ` + page

	rows := Crud.GetRows(querySql, []interface{}{})
	defer rows.Close()

	var appList = []*entity.App{}

	for rows.Next() {
		appList = append(appList, entity.ScanRowsApp(rows))
	}
	return appList
}

func (this *appModel) CountItems(search string) int64{
	var querySql = `SELECT count(*) FROM app ` + this.getSearchSql(search)
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

func (this *appModel)getSearchSql(search string)string{
	if search != "" {
		return `WHERE Site_name  LIKE '%`+search+`%'
	OR Admin_email  LIKE '%`+search+`%'
	OR Admin_email_pass string
	OR Crypt_salt LIKE '%`+search+`%'
	OR Cookie_name LIKE '%`+search+`%'
	OR Assets_path LIKE '%`+search+`%'
	OR Assets_backend_path LIKE '%`+search+`%'
	OR Assets_public_path LIKE '%`+search+`%'
	OR Assets_cabinet_path LIKE '%`+search+`%'
	OR Upload_path  LIKE '%`+search+`%'
	OR No_image LIKE '%`+search+`%'
	OR Image_path LIKE '%`+search+`%'
	OR Recaptcha_secret LIKE '%`+search+`%'
	OR Recaptcha_public LIKE '%`+search+`%'
	OR Uri_auth_deny LIKE '%`+search+`%'
	OR Uri_auth_banned LIKE '%`+search+`%'
	OR Uri_permission LIKE '%`+search+`%'
	OR Uri_permission_ajax LIKE '%`+search+`%'
	OR Uri_permission_add LIKE '%`+search+`%'
	OR Uri_permission_edit LIKE '%`+search+`%'
	OR Uri_permission_del LIKE '%`+search+`%'
	OR Uri_permission_get LIKE '%`+search+`%'
	OR Uri_position LIKE '%`+search+`%'
	OR Uri_position_ajax LIKE '%`+search+`%'
	OR Uri_position_add LIKE '%`+search+`%'
	OR Uri_position_edit LIKE '%`+search+`%'
	OR Uri_position_del LIKE '%`+search+`%'
	OR Uri_position_get LIKE '%`+search+`%'
	OR Uri_account LIKE '%`+search+`%'
	OR Uri_account_ajax LIKE '%`+search+`%'
	OR Uri_account_add LIKE '%`+search+`%'
	OR Uri_account_edit LIKE '%`+search+`%'
	OR Uri_account_del LIKE '%`+search+`%'
	OR Uri_account_get LIKE '%`+search+`%'
	OR Uri_cabinet LIKE '%`+search+`%'
	OR Uri_cabinet_add LIKE '%`+search+`%'
	OR Uri_cabinet_edit LIKE '%`+search+`%'
	OR Uri_cabinet_del LIKE '%`+search+`%'
	OR Uri_cabinet_get LIKE '%`+search+`%'
	OR Uri_auth_login LIKE '%`+search+`%'
	OR Uri_auth_register LIKE '%`+search+`%'
	OR Uri_auth_activation LIKE '%`+search+`%'
	OR Uri_auth_forgot LIKE '%`+search+`%'
	OR Uri_auth_change_password  LIKE '%`+search+`%'
	OR Uri_auth_logout LIKE '%`+search+`%'
	OR Uri_auth_delete_account LIKE '%`+search+`%' `
	}
	return ""
}
