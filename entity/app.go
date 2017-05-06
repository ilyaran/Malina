package entity

import "database/sql"

type  App struct {

	Id int64
	Site_name string

	Admin_email string
	Admin_email_pass string

	Crypt_salt string

	Session_expiration int64
	Cookie_expiration int
	Cookie_name string

	Per_page int64
	Radius int64

	Language_list []string
	Assets_path string
	Assets_backend_path string
	Assets_public_path string
	Assets_cabinet_path string

	Upload_path  string
	No_image string
	Image_path string
	Image_max_width int
	Image_max_height int
	Image_resize_width int64
	Image_resize_height int64

	Image_max_size int64
	Image_upload_limit int

	//auth
	Login_attempts int
	Login_attempts_period int
	Recaptcha_secret string
	Recaptcha_public string
	Use_recaptcha_register bool

	Use_recaptcha_login bool
	Use_recaptcha_forgot bool
	Use_recaptcha_change_password bool

	Email_activation bool

	Email_activation_expire int64
	Email_account_details bool

	Uri_auth_deny string
	Uri_auth_banned string

	Uri_permission string
	Uri_permission_ajax string
	Uri_permission_add string
	Uri_permission_edit string
	Uri_permission_del string
	Uri_permission_get string

	Uri_position string
	Uri_position_ajax string
	Uri_position_add string
	Uri_position_edit string
	Uri_position_del string
	Uri_position_get string

	Uri_account string
	Uri_account_ajax string
	Uri_account_add string
	Uri_account_edit string
	Uri_account_del string
	Uri_account_get string

	Uri_cabinet string
	Uri_cabinet_add string
	Uri_cabinet_edit string
	Uri_cabinet_del string
	Uri_cabinet_get string

	Uri_auth_login string
	Uri_auth_register string
	Uri_auth_activation string
	Uri_auth_forgot string
	Uri_auth_change_password  string
	Uri_auth_logout string
	Uri_auth_delete_account string

	PasswordMinLength int
	PasswordMaxLength int
}

func ScanRowApp(row *sql.Row)*App{
	s := &App{}
	err := row.Scan(
		&s.Id,
		&s.Site_name,
		&s.Admin_email,
		&s.Admin_email_pass,
		&s.Crypt_salt,
		&s.Session_expiration,
		&s.Cookie_expiration,
		&s.Cookie_name,
		&s.Per_page,
		&s.Radius,
		&s.Language_list,
		&s.Assets_path,
		&s.Assets_backend_path,
		&s.Assets_public_path,
		&s.Assets_cabinet_path,
		&s.Upload_path,
		&s.No_image,
		&s.Image_path,
		&s.Image_max_width,
		&s.Image_max_height,
		&s.Image_resize_width,
		&s.Image_resize_height,
		&s.Image_max_size,
		&s.Image_upload_limit,
		&s.Login_attempts,
		&s.Login_attempts_period,
		&s.Recaptcha_secret,
		&s.Recaptcha_public,
		&s.Use_recaptcha_register,
		&s.Use_recaptcha_login,
		&s.Use_recaptcha_forgot,
		&s.Use_recaptcha_change_password,
		&s.Email_activation,
		&s.Email_activation_expire,
		&s.Email_account_details,
		&s.Uri_auth_deny,
		&s.Uri_auth_banned,
		&s.Uri_permission,
		&s.Uri_permission_ajax,
		&s.Uri_permission_add,
		&s.Uri_permission_edit,
		&s.Uri_permission_del,
		&s.Uri_permission_get,
		&s.Uri_position,
		&s.Uri_position_ajax,
		&s.Uri_position_add,
		&s.Uri_position_edit,
		&s.Uri_position_del,
		&s.Uri_position_get,
		&s.Uri_account,
		&s.Uri_account_ajax,
		&s.Uri_account_add,
		&s.Uri_account_edit,
		&s.Uri_account_del,
		&s.Uri_account_get,
		&s.Uri_cabinet,
		&s.Uri_cabinet_add,
		&s.Uri_cabinet_edit,
		&s.Uri_cabinet_del,
		&s.Uri_cabinet_get,
		&s.Uri_auth_login,
		&s.Uri_auth_register,
		&s.Uri_auth_activation,
		&s.Uri_auth_forgot,
		&s.Uri_auth_change_password,
		&s.Uri_auth_logout,
		&s.Uri_auth_delete_account,
		&s.PasswordMinLength,
		&s.PasswordMaxLength,
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
func ScanRowsApp(rows *sql.Rows)*App{
	s := &App{}
	err := rows.Scan(
		&s.Id,
		&s.Site_name,
		&s.Admin_email,
		&s.Admin_email_pass,
		&s.Crypt_salt,
		&s.Session_expiration,
		&s.Cookie_expiration,
		&s.Cookie_name,
		&s.Per_page,
		&s.Radius,
		&s.Language_list,
		&s.Assets_path,
		&s.Assets_backend_path,
		&s.Assets_public_path,
		&s.Assets_cabinet_path,
		&s.Upload_path,
		&s.No_image,
		&s.Image_path,
		&s.Image_max_width,
		&s.Image_max_height,
		&s.Image_resize_width,
		&s.Image_resize_height,
		&s.Image_max_size,
		&s.Image_upload_limit,
		&s.Login_attempts,
		&s.Login_attempts_period,
		&s.Recaptcha_secret,
		&s.Recaptcha_public,
		&s.Use_recaptcha_register,
		&s.Use_recaptcha_login,
		&s.Use_recaptcha_forgot,
		&s.Use_recaptcha_change_password,
		&s.Email_activation,
		&s.Email_activation_expire,
		&s.Email_account_details,
		&s.Uri_auth_deny,
		&s.Uri_auth_banned,
		&s.Uri_permission,
		&s.Uri_permission_ajax,
		&s.Uri_permission_add,
		&s.Uri_permission_edit,
		&s.Uri_permission_del,
		&s.Uri_permission_get,
		&s.Uri_position,
		&s.Uri_position_ajax,
		&s.Uri_position_add,
		&s.Uri_position_edit,
		&s.Uri_position_del,
		&s.Uri_position_get,
		&s.Uri_account,
		&s.Uri_account_ajax,
		&s.Uri_account_add,
		&s.Uri_account_edit,
		&s.Uri_account_del,
		&s.Uri_account_get,
		&s.Uri_cabinet,
		&s.Uri_cabinet_add,
		&s.Uri_cabinet_edit,
		&s.Uri_cabinet_del,
		&s.Uri_cabinet_get,
		&s.Uri_auth_login,
		&s.Uri_auth_register,
		&s.Uri_auth_activation,
		&s.Uri_auth_forgot,
		&s.Uri_auth_change_password,
		&s.Uri_auth_logout,
		&s.Uri_auth_delete_account,
		&s.PasswordMinLength,
		&s.PasswordMaxLength,
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

func (s *App) ExecWithId() []interface{} {
	return []interface{}{
		s.Id,
		s.Site_name,
		s.Admin_email,
		s.Admin_email_pass ,
		s.Crypt_salt,
		s.Session_expiration ,
		s.Cookie_expiration,
		s.Cookie_name ,
		s.Per_page ,
		s.Radius ,
		s.Language_list ,
		s.Assets_path ,
		s.Assets_backend_path ,
		s.Assets_public_path ,
		s.Assets_cabinet_path ,
		s.Upload_path ,
		s.No_image,
		s.Image_path,
		s.Image_max_width,
		s.Image_max_height,
		s.Image_resize_width,
		s.Image_max_size,
		s.Image_upload_limit,
		s.Login_attempts,
		s.Login_attempts_period,
		s.Recaptcha_secret,
		s.Recaptcha_public,
		s.Use_recaptcha_register,
		s.Use_recaptcha_login,
		s.Use_recaptcha_forgot,
		s.Use_recaptcha_change_password,
		s.Email_activation,
		s.Email_activation_expire,
		s.Email_account_details,
		s.Uri_auth_deny,
		s.Uri_auth_banned,
		s.Uri_permission,
		s.Uri_permission_ajax,
		s.Uri_permission_add,
		s.Uri_permission_edit,
		s.Uri_permission_del,
		s.Uri_permission_get,
		s.Uri_position,
		s.Uri_position_ajax,
		s.Uri_position_add,
		s.Uri_position_edit,
		s.Uri_position_del,
		s.Uri_position_get,
		s.Uri_account,
		s.Uri_account_ajax,
		s.Uri_account_add,
		s.Uri_account_edit,
		s.Uri_account_del,
		s.Uri_account_get,
		s.Uri_cabinet,
		s.Uri_cabinet_add,
		s.Uri_cabinet_edit,
		s.Uri_cabinet_del,
		s.Uri_cabinet_get,
		s.Uri_auth_login,
		s.Uri_auth_register,
		s.Uri_auth_activation,
		s.Uri_auth_forgot,
		s.Uri_auth_change_password ,
		s.Uri_auth_logout,
		s.Uri_auth_delete_account,
		s.PasswordMinLength,
		s.PasswordMaxLength,
	}
}
func (s *App) Exec() []interface{} {
	return []interface{}{
		s.Site_name,
		s.Admin_email,
		s.Admin_email_pass ,
		s.Crypt_salt,
		s.Session_expiration ,
		s.Cookie_expiration,
		s.Cookie_name ,
		s.Per_page ,
		s.Radius ,
		s.Language_list ,
		s.Assets_path ,
		s.Assets_backend_path ,
		s.Assets_public_path ,
		s.Assets_cabinet_path ,
		s.Upload_path ,
		s.No_image,
		s.Image_path,
		s.Image_max_width,
		s.Image_max_height,
		s.Image_resize_width,
		s.Image_max_size,
		s.Image_upload_limit,
		s.Login_attempts,
		s.Login_attempts_period,
		s.Recaptcha_secret,
		s.Recaptcha_public,
		s.Use_recaptcha_register,
		s.Use_recaptcha_login,
		s.Use_recaptcha_forgot,
		s.Use_recaptcha_change_password,
		s.Email_activation,
		s.Email_activation_expire,
		s.Email_account_details,
		s.Uri_auth_deny,
		s.Uri_auth_banned,
		s.Uri_permission,
		s.Uri_permission_ajax,
		s.Uri_permission_add,
		s.Uri_permission_edit,
		s.Uri_permission_del,
		s.Uri_permission_get,
		s.Uri_position,
		s.Uri_position_ajax,
		s.Uri_position_add,
		s.Uri_position_edit,
		s.Uri_position_del,
		s.Uri_position_get,
		s.Uri_account,
		s.Uri_account_ajax,
		s.Uri_account_add,
		s.Uri_account_edit,
		s.Uri_account_del,
		s.Uri_account_get,
		s.Uri_cabinet,
		s.Uri_cabinet_add,
		s.Uri_cabinet_edit,
		s.Uri_cabinet_del,
		s.Uri_cabinet_get,
		s.Uri_auth_login,
		s.Uri_auth_register,
		s.Uri_auth_activation,
		s.Uri_auth_forgot,
		s.Uri_auth_change_password ,
		s.Uri_auth_logout,
		s.Uri_auth_delete_account,
		s.PasswordMinLength,
		s.PasswordMaxLength,
	}
}