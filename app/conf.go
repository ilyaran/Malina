/**
 *
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ package app

import (
	"errors"
	"strconv"
	"io/ioutil"
	"regexp"
	"strings"
)



var PlatformOs string
var	Per_page_max int64 = 201
var Base_cloud_url = ""

var Url_assets = ""
var Url_assets_home = ""
var Url_assets_public = ""
var Url_assets_uploads = ""

var Path_assets = ""
var Path_assets_uploads = ""

var Url_no_image = ""
var Base_url =  ``

var Url_cabinet string
var Url_public_product_list string
var Url_public_product_item string
var Url_public_news_list string
var Url_public_category_list string
var Url_public_parameter_list string

var Url_home_product_list string
var Url_home_parameter_list string
var Url_home_category_list string
var Url_home_role_list string
var Url_home_news_list string
var Url_home_account_list string
var Url_home_cart_list string
var Url_home_order_list string
var Url_home_permission_list string
var Url_home_activation_list string


var Url_auth_change_password string


var Url_support string
var Url_auth_deny                   = ""
var Url_auth_banned                 = ""

var Url_auth_activation = ""
var Url_cabinet_profile = ""

func SetConfig(os string)error{
	if os == "windows" {
		Port_addr=":3002"
		Base_url="http://localhost"+Port_addr+"/"
		Root_path = "./"
	}else if os == "linux" {
		Base_url="http://"+Domain_name+Port_addr+"/"
		Root_path = ""
	}else {
		return errors.New("Not supported  OS/Platform")
	}

	if Base_cloud_url != "" {
		Url_assets = Base_cloud_url + "/"
	}else {
		Url_assets = Base_url + "assets/"
	}

	Url_no_image = Url_assets + "img/noimg.jpg"
	Url_assets_home = Url_assets + "home/"
	Url_assets_public = Url_assets + "public/"
	Url_assets_uploads = Url_assets + "uploads/"

	Path_assets = "assets"
	Path_assets_uploads = Path_assets + "/uploads/"

	Url_cabinet =Base_url+"cabinet/"

	Url_public_product_list = Base_url+"public/product/list/"
	Url_public_news_list = Base_url+"public/news/list/"
	Url_public_category_list = Base_url+"public/category/list/"
	Url_public_parameter_list = Base_url+"public/parameter/list/"
	Url_public_product_item = Base_url+"public/product/item/"

	Url_home_product_list = Base_url+"home/product/list/"
	Url_home_parameter_list = Base_url+"home/parameter/list/"
	Url_home_category_list = Base_url+"home/category/list/"
	Url_home_role_list = Base_url+"home/role/list/"
	Url_home_news_list = Base_url+"home/news/list/"
	Url_home_account_list = Base_url+"home/account/list/"
	Url_home_cart_list = Base_url+"home/cart/list/"
	Url_home_order_list = Base_url+"home/order/list/"
	Url_home_permission_list = Base_url+"home/permission/list/"
	Url_home_activation_list = Base_url+"home/activation/list/"


	Url_auth_change_password = Base_url+"auth/change_password/"

	Url_support = Base_url+"public/support/"

	Url_auth_deny                   = ""
	Url_auth_banned                 = ""

	Url_auth_activation =  Base_url+"auth/activation/"
	Url_cabinet_profile =  Base_url+"cabinet/profile/"

	AccountSelectSqlFieldsList = []string{
		/*0*/ "COALESCE(account.account_id,0)",
		/*1*/ "COALESCE(account_role,0)",
		/*2*/ "COALESCE(account_email,'')",
		/*3*/ "COALESCE(account_phone,'')",

		/*4*/ "COALESCE(account_nick,'')",
		/*5*/ "COALESCE(account_first_name,'')",
		/*6*/ "COALESCE(account_last_name,'')",

		/*7*/ "COALESCE(account_skype,'')",

		/*8*/ "COALESCE(account_state,0)",
		/*9*/ `CASE WHEN account_img = '' OR account_img IS NULL
				THEN '`+ Url_assets_public+`theme/img/no-avatar.png'
            	ELSE '`+Base_url+Path_assets_uploads+`' || account_img
            END`,
		/*10*/ "COALESCE((SELECT SUM(balance_in) - SUM(balance_out) FROM balance WHERE balance_account = account.account_id),0)",
		/*11*/ "COALESCE(account_ban,FALSE)",

		/*12*/ "COALESCE(account_reason,'')",
		/*13*/ "COALESCE(account_created,now())",
		/*14*/ "COALESCE(account_updated,now())",
		/*15*/ "COALESCE(account_last_logged,now())",

		/*16*/ "COALESCE(account_last_ip,'')",
		/*17*/ "COALESCE(account_token,'')",
		/*18*/ "COALESCE(account_password,'')",
		/*19*/ "COALESCE(account_provider,'')",

	}




	SetFilemanagerConf(os)


	return nil
}

func SetFilemanagerConf(os string)  {
	var pathToConfTemplate string
	var pathToConf string
	if os == "windows" {
		pathToConfTemplate = Path_assets+"/filemanager/conf_template.json"
		pathToConf = Path_assets+"/filemanager/conf.json"
	}
	if os == "linux" {
		pathToConfTemplate = Root_path + "/" + Path_assets+"/filemanager/conf_template.json"
		pathToConf = Root_path + "/" + Path_assets+"/filemanager/conf.json"
	}
	b, err := ioutil.ReadFile(pathToConfTemplate) // just pass the file name
	if err != nil {
		panic(err)
	}
	str := string(b) // convert content to a 'string'
	if os == "windows" {
		str = regexp.MustCompile(`\{\{Path_assets_uploads\}\}`).ReplaceAllString(str, strings.Trim(Path_assets_uploads,"/"))
	}
	if os == "linux" {
		str = regexp.MustCompile(`\{\{Path_assets_uploads\}\}`).ReplaceAllString(str, Root_path+"/"+strings.Trim(Path_assets_uploads,"/"))
	}
	str = regexp.MustCompile(`\{\{Base_url\}\}`).ReplaceAllString(str, strings.Trim(Base_url,"/"))

	d1 := []byte(str)
	err = ioutil.WriteFile(pathToConf, d1, 0644)
	if err != nil{
		panic(err)
	}

}



var Domain_name =  `example.com`
var Site_name =  `Malina`

var DB_host =  `localhost:5432`
var DB_user =  `postgres`
var DB_password =  `postgres`
var DB_name =  `malina`


var Root_path =  ``
var Port_addr =  ``

var Per_page int64 = 3
var Per_page_select_options_widget_num int64 = 5
var Per_page_public_select_widget_html = ""


var Crypt_salt =  `9Woi6SPj5UjFAh2sd3fHBFU1msdI5fjGcuhwam2dasMZ`
var GenereatedPasswordLen = 6
var GenereatedPincodeLen = 4

var Cookie_name =  `malina`
var Cookie_expiration  = 3600

var Admin_email  = "johnxiaran@gmail.com"
var Admin_email_pass = "MuromSvetXia116692"

var Admin_role_id int64 = 0

var Pattern_phone = `^[\+]?[0-9]{8,16}$`
var Pattern_email = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

var LangList = []string{"en","ru","zh-CN","he","es"}

var Session_expiration int64 = 3600
var Session_expirationStr string = strconv.FormatInt(Session_expiration,10)


var Short_description_len = 2048

var AccountSelectSqlFieldsList []string
var PasswordMinLength = 6
var PasswordMaxLength = 64

var Title_max_len=255
var Short_description_max_len=255


//auth
var Login_attempts                  = int64(12)
var Login_attempts_period           = 3600
var Recaptcha_secret                = "6Le_zhITAAAAAL2T7WoeckNYcWjeNYj0fg2-3eWe"
var Recaptcha_public                = "6Le_zhITAAAAAFb3_plfOi_dlxq_BQ5nXppR61tA"
var Use_recaptcha_register          = false

var Use_recaptcha_login             = false
var Use_recaptcha_forgot            = false
var Use_recaptcha_change_password   = false

var Email_activation                = true

var Email_activation_expire         = int64(360000)
var Email_account_details           = true


var Image_max_width =    2000
var Image_max_height =    2000
var Image_resize_width =  int64(400)
var Image_resize_height = int64(400)

var Image_max_size = int64(1234567)
var Image_upload_limit = 1
var Image_limit_per_item = 12