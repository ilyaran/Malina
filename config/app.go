package app

const (
	HOST = "localhost:5432"
	USER = "postgres"
	PASSWORD = "postgres"
	NAME = "malina"
)
var base_url = "http://localhost:3001/"
type  App struct {

	Id int64
	Site_name string

	Admin_email string
	Admin_email_pass string

	Crypt_salt string

	Session_expiration int64
	Cart_cookie_expiration int
	Sms_activation_expiration int64
	Cookie_expiration int
	Cookie_name string

	Per_page int64
	Radius int64

	Tax float64
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
	Login_attempts int64
	Login_attempts_period int
	Recaptcha_secret string
	Recaptcha_public string
	Use_recaptcha_register bool

	Use_recaptcha_login bool
	Use_recaptcha_forgot bool
	Use_recaptcha_change_password bool
	Image_limit_per_item int
	Email_activation bool

	Email_activation_expire int64
	Email_account_details bool

	Uri_auth_deny string
	Uri_auth_banned string

	Uri_public_cart string
	Uri_public_product_list              string
	Uri_public_product_list_ajax         string
	Uri_public_product_get          string
	
	Uri_cart              string
	Uri_cart_ajax         string
	Uri_cart_add          string
	Uri_cart_edit         string
	Uri_cart_del          string
	Uri_cart_get          string
	
	Uri_category              string
	Uri_category_ajax         string
	Uri_category_add          string
	Uri_category_edit         string
	Uri_category_del          string
	Uri_category_get          string
	
	Uri_product              string
	Uri_product_ajax         string
	Uri_product_add          string
	Uri_product_edit         string
	Uri_product_del          string
	Uri_product_get          string
	
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
	Uri_position_del         string
	Uri_position_get         string

	Uri_home_account string
	Uri_account_ajax string
	Uri_account_add  string
	Uri_account_edit string
	Uri_account_del  string
	Uri_account_get  string

	Uri_cabinet              string
	Uri_cabinet_add          string
	Uri_cabinet_edit         string
	Uri_cabinet_del          string
	Uri_cabinet_get          string

	Uri_auth_login           string
	Uri_auth_register        string
	Uri_auth_activation      string
	Uri_auth_forgot          string
	Uri_auth_change_password string
	Uri_auth_logout          string
	Uri_auth_delete_account  string


	PasswordMinLength        int
	PasswordMaxLength        int

	Pattern_title string
	Pattern_phone string
	Pattern_email string

	//Cart
	ADD_PRODUCT int8
	REMOVE_PRODUCT int8
	UPDATE_PRODUCTS_QUANTITIES int8
	SAVE_PRODUCT_FOR_LATER int8
	MOVE_PRODUCT_TO_CART int8
}

var app = &App{
	Tax:0.003,
	//Cart
	ADD_PRODUCT : 1,
	REMOVE_PRODUCT : 2,
	UPDATE_PRODUCTS_QUANTITIES : 3,
	SAVE_PRODUCT_FOR_LATER : 4,
	MOVE_PRODUCT_TO_CART : 5,

	Site_name : "Malina",

	Admin_email : "johnxiaran@gmail.com",
	Admin_email_pass : "",

	Crypt_salt : "tyr2ty4u1rt6xbO7po9c3oa5h",

	Session_expiration : int64(3600),
	Cart_cookie_expiration: 3600,
	Sms_activation_expiration : int64(600),
	Cookie_expiration : 3600,
	Cookie_name : "malina",

	Per_page : int64(15),
	Radius : int64(8),

	Language_list : []string{"en", "ru", "zh-CN"},
	Assets_path : base_url + "assets/",
	Assets_backend_path : base_url +"assets/home/",
	Assets_public_path : base_url + "assets/public/",
	Assets_cabinet_path : base_url + "assets/cabinet/",

	Upload_path : "assets/uploads/",
	No_image : base_url +"assets/img/noimg.jpg",
	Image_path : base_url +"assets/img/",
	Image_max_width : 2000,
	Image_max_height : 2000,
	Image_resize_width : int64(400),
	Image_resize_height : int64(400),

	Image_max_size : int64(1512000),
	Image_upload_limit : 1,
	Image_limit_per_item : 12,

	//auth
	Login_attempts : int64(12),
	Login_attempts_period : 3600,
	Recaptcha_secret : "6Le_zhITAAAAAL2T7WoeckNYcWjeNYj0fg2-3eWe",
	Recaptcha_public  : "6Le_zhITAAAAAFb3_plfOi_dlxq_BQ5nXppR61tA",
	Use_recaptcha_register : false,

	Use_recaptcha_login : false,
	Use_recaptcha_forgot : true,
	Use_recaptcha_change_password : true,

	Email_activation : true,

	Email_activation_expire : int64(3600),
	Email_account_details : true,

	Uri_auth_deny : "",
	Uri_auth_banned : "",

	Uri_public_cart              : "public/cart/crud/",
	Uri_public_product_list              : "public/product/list/",
	Uri_public_product_list_ajax         : "public/product/list_ajax/",
	Uri_public_product_get		     : "public/product/get/",

	Uri_cart              : "home/cart/",
	Uri_cart_ajax         : "home/cart/ajax_list/",
	Uri_cart_add          : "home/cart/add/",
	Uri_cart_edit         : "home/cart/edit/",
	Uri_cart_del          : "home/cart/del/",
	Uri_cart_get          : "home/cart/get/",
	
	Uri_category              : "home/category/",
	Uri_category_ajax         : "home/category/ajax_list/",
	Uri_category_add          : "home/category/add/",
	Uri_category_edit         : "home/category/edit/",
	Uri_category_del          : "home/category/del/",
	Uri_category_get          : "home/category/get/",

	Uri_product              : "home/product/",
	Uri_product_ajax         : "home/product/ajax_list/",
	Uri_product_add          : "home/product/add/",
	Uri_product_edit         : "home/product/edit/",
	Uri_product_del          : "home/product/del/",
	Uri_product_get          : "home/product/get/",

	Uri_permission : "home/permission/",
	Uri_permission_ajax : "home/permission/ajax_list/",
	Uri_permission_add : "home/permission/add/",
	Uri_permission_edit : "home/permission/edit/",
	Uri_permission_del : "home/permission/del/",
	Uri_permission_get : "home/permission/get/",

	Uri_position : "home/position/",
	Uri_position_ajax : "home/position/ajax_list/",
	Uri_position_add : "home/position/add/",
	Uri_position_edit : "home/position/edit/",
	Uri_position_del : "home/position/del/",
	Uri_position_get : "home/position/get/",

	Uri_home_account:  "home/account/",
	Uri_account_ajax : "home/account/ajax_list/",
	Uri_account_add :  "home/account/add/",
	Uri_account_edit : "home/account/edit/",
	Uri_account_del :  "home/account/del/",
	Uri_account_get :  "home/account/get/",

	Uri_cabinet : "cabinet/",
	Uri_cabinet_add : "cabinet/add/",
	Uri_cabinet_edit : "cabinet/edit/",
	Uri_cabinet_del : "cabinet/del/",
	Uri_cabinet_get : "cabinet/get/",

	Uri_auth_login : base_url+"auth/login/",
	Uri_auth_register : base_url+"auth/register/",
	Uri_auth_activation : base_url+"auth/activation/",
	Uri_auth_forgot : base_url+"auth/forgot/",
	Uri_auth_change_password : base_url+"auth/change_password/",
	Uri_auth_logout : base_url+"auth/logout/",
	Uri_auth_delete_account : base_url+"auth/delete_account/",

	PasswordMinLength : 6,
	PasswordMaxLength : 64,

	Pattern_title : `^[\w_-\s]{1,255}$`,
	Pattern_phone : `^[\+]?[0-9]{8,16}$`,
	Pattern_email:`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`,
}

func DB_HOST()string{return HOST}
func DB_USER()string{return USER
}
func DB_PASSWORD()string{return PASSWORD
}
func DB_NAME()string{return NAME
}


func Language_list()[]string{return app.Language_list}
func TAX()float64{return app.Tax}
func ADD_PRODUCT()int8{return app.ADD_PRODUCT}
func REMOVE_PRODUCT()int8{return app.REMOVE_PRODUCT}
func UPDATE_PRODUCTS_QUANTITIES()int8{return app.UPDATE_PRODUCTS_QUANTITIES}
func SAVE_PRODUCT_FOR_LATER()int8{return app.SAVE_PRODUCT_FOR_LATER}
func MOVE_PRODUCT_TO_CART()int8{return app.MOVE_PRODUCT_TO_CART}
func Image_limit_per_item()int{return app.Image_limit_per_item}

func Session_expiration() int64 {return app.Session_expiration}
func Cart_cookie_expiration() int{return app.Cart_cookie_expiration}
func Cookie_expiration() int {return app.Cookie_expiration}
func Cookie_name() string {return app.Cookie_name}
func Crypt_salt() string {return app.Crypt_salt}
func Base_url() string {return base_url}
func Site_name() string {return app.Site_name}
func Admin_email() string {return app.Admin_email}
func Admin_email_pass() string {return app.Admin_email_pass}
func Radius() int64 {return app.Radius}
func Per_page() int64 {return app.Per_page}
func Image_max_size() int64 {return app.Image_max_size}
func Image_upload_limit() int {return app.Image_upload_limit}
func Image_max_width() int {return app.Image_max_width}
func Image_max_height() int {return app.Image_max_height}

func No_image() string {return app.No_image}
func Image_path() string {return app.Image_path}

func Use_recaptcha_register() bool {return app.Use_recaptcha_register}
func Use_recaptcha_login() bool {return app.Use_recaptcha_login}
func Use_recaptcha_forgot() bool {return app.Use_recaptcha_forgot}
func Use_recaptcha_change_password() bool {return app.Use_recaptcha_change_password}

func Login_attempts() int64 {return app.Login_attempts}
func Login_attempts_period() int {return app.Login_attempts_period}
func Recaptcha_secret() string {return app.Recaptcha_secret}
func Recaptcha_public() string {return app.Recaptcha_public}
func Email_activation()bool {return app.Email_activation}

func Email_activation_expire() int64 {return app.Email_activation_expire}
func Email_account_details() bool {return app.Email_account_details}

func Uri_auth_login() string {return app.Uri_auth_login}
func Uri_auth_register() string {return app.Uri_auth_register}
func Uri_auth_activation() string {return app.Uri_auth_activation}
func Uri_auth_forgot() string {return app.Uri_auth_forgot}
func Uri_auth_change_password() string {return app.Uri_auth_change_password}
func Uri_auth_logout() string {return app.Uri_auth_logout}
func Uri_auth_delete_account() string {return app.Uri_auth_delete_account}

func Uri_public_cart() string {return app.Uri_public_cart}
func Uri_public_product_list() string {return app.Uri_public_product_list}
func Uri_public_product_list_ajax() string {return app.Uri_public_product_list_ajax}
func Uri_public_product_get() string {return app.Uri_public_product_get}

func Uri_cart() string {return app.Uri_cart}
func Uri_cart_ajax() string {return app.Uri_cart_ajax}
func Uri_cart_add() string {return app.Uri_cart_add}
func Uri_cart_edit() string {return app.Uri_cart_edit}
func Uri_cart_del() string {return app.Uri_cart_del}
func Uri_cart_get() string {return app.Uri_cart_get}

func Uri_category() string {return app.Uri_category}
func Uri_category_ajax() string {return app.Uri_category_ajax}
func Uri_category_add() string {return app.Uri_category_add}
func Uri_category_edit() string {return app.Uri_category_edit}
func Uri_category_del() string {return app.Uri_category_del}
func Uri_category_get() string {return app.Uri_category_get}

func Uri_product() string {return app.Uri_product}
func Uri_product_ajax() string {return app.Uri_product_ajax}
func Uri_product_add() string {return app.Uri_product_add}
func Uri_product_edit() string {return app.Uri_product_edit}
func Uri_product_del() string {return app.Uri_product_del}
func Uri_product_get() string {return app.Uri_product_get}

func Uri_permission() string {return app.Uri_permission}
func Uri_permission_ajax() string {return app.Uri_permission_ajax}
func Uri_permission_add() string {return app.Uri_permission_add}
func Uri_permission_edit() string {return app.Uri_permission_edit}
func Uri_permission_del() string {return app.Uri_permission_del}
func Uri_permission_get() string {return app.Uri_permission_get}

func Uri_position() string {return app.Uri_position}
func Uri_position_ajax() string {return app.Uri_position_ajax}
func Uri_position_add() string {return app.Uri_position_add}
func Uri_position_edit() string {return app.Uri_position_edit}
func Uri_position_del() string {return app.Uri_position_del}
func Uri_position_get() string {return app.Uri_position_get}

func Uri_account() string {return app.Uri_home_account }
func Uri_account_ajax() string {return app.Uri_account_ajax}
func Uri_account_add() string {return app.Uri_account_add}
func Uri_account_edit() string {return app.Uri_account_edit}
func Uri_account_del() string {return app.Uri_account_del}
func Uri_account_get() string {return app.Uri_account_get}

func Uri_cabinet() string {return app.Uri_cabinet}
func Uri_cabinet_add() string {return app.Uri_cabinet_add}
func Uri_cabinet_edit() string {return app.Uri_cabinet_edit}
func Uri_cabinet_del() string {return app.Uri_cabinet_del}
func Uri_cabinet_get() string {return app.Uri_cabinet_get}

func PasswordMinLength() int {return app.PasswordMinLength}
func PasswordMaxLength() int {return app.PasswordMaxLength}

func Upload_path() string {return app.Upload_path}
func Assets_path() string {return app.Assets_path}
func Assets_backend_path() string {return app.Assets_backend_path}
func Assets_public_path() string {return app.Assets_public_path}

func Pattern_title()string { return app.Pattern_title }
func Pattern_phone()string { return app.Pattern_phone }
func Pattern_email()string { return app.Pattern_email }

/*func SetTable()  {
	s := reflect.ValueOf(app).Elem()
	typeOfT := s.Type()
	var def,ty string
	var t = `CREATE TABLE app (app_id BIGSERIAL NOT NULL PRIMARY KEY UNIQUE,`
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if "string" == fmt.Sprintf("%v",f.Type()){
			def = fmt.Sprintf("'%v'",f.Interface())
			ty="VARCHAR(512)"
		}else if "bool" == fmt.Sprintf("%v",f.Type()){
			def = fmt.Sprintf("%v",f.Interface())
			ty="BOOLEAN"
		}else {
			def = fmt.Sprintf("%v",f.Interface())
			ty="INTEGER"
		}
		t+=fmt.Sprintf(`%v %v DEFAULT %v NOT NULL,
		`,typeOfT.Field(i).Name,ty,def)
	}
	fmt.Println(t+`);`)
}*/

















