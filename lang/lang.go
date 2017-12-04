package lang

import (
	"strconv"
	"fmt"
	"github.com/ilyaran/Malina/app"
)

var Lang string = "ru"

func T(key string) string {

	if Lang!=""{
		if v,ok:=Dict[Lang][key];ok{
			return v
		}
	}else {
		if v,ok:=Dict["en"][key];ok{
			return v
		}
	}



	return key
}

var Dict = map[string]map[string]string{
	"ru":{
		`sign in failed your password is an incorrect`:`Вход не удался, пароль и email не совпадают`,
		`invalid`:`не корректный`,
		`validation required`	:`Поле %s обязательно требуется.`,
		`validation min length`	: `Поле %s должно содержать не менее %d символов.`,
		`validation max length`	: `Поле %s должно содержать не более %d символов.`,
		`auth exceed max attempts`:"Attempts exceed max attempts in config. No more than "+ strconv.FormatInt(app.Login_attempts,10)+ " times per 24 hours",
		`auth_confirm_pass_valid`:"Поле повторить пароль и поле пароль должны совпадать.",

		"exists allready":"уже существует",
		`auth_logout_first`:"Вы должны сперва выйти. Вы уже залогинены.",
		`server error`: `ошибка сервера`,
		`wrong_recaptcha`:"Каптча не правильная. Может быть вы робот.",
		`not found`:`не найдено`,
		`auth_email_not_exist` : "Email не существует",
		"required":"обязательное поле для заполнения",
		`auth_login_allready`:"Вы уже залогинены.",

		`auth first`:"Вы должны сперва авторизоваться, прежде чем пытаться сбросить пароль.",
		"old password":"Старый пароль.",
		"new password":"Новый пароль.",
		"confirm new password":"Подтвердить новый пароль.",
		"reset password" : "Новый пароль",
		"send":"Отправить",
		"password success changed":"Пароль успешно изменен.",
		"incorrect password":"Не корректный пароль.",
		`auth_activation_incorrect_code`:"Код активации неправильный или истек срок годности кода. Пожалуйста проверьте email снова.",




		`validation text field`	:  `The %s field must contain %s`,

		`validation integer`	:  `The %s field must contain an integer.`,
		`validation email`	:"The Email field must contain all valid email addresses.",
		`validation phone`	:"The PhoneCode field must contain all valid phone numbers.",


		"no post data":"No Post Data",
		"not auth":"You must login",
		"no post request":"It is not a Post Method Request",

		`upload no base64`:               `Unable to find a base64 format.`,
		`upload error base64 decoding`:   `Unable to decode a base64 string.`,
		`upload destination error`:       `A problem was encountered while attempting to move the uploaded file to the final destination.`,
		`upload invalid filetype`:        `The filetype you are attempting to upload is not allowed.`,
		`upload_invalid_dimensions`:        `The file dimensions you are attempting to upload is not allowed.`,
		`upload_unable_to_write_file`:        `Unable to write file.`,
		`upload_file_exceeds_limit`:      `The uploaded file exceeds the maximum allowed size in ` + fmt.Sprintf("%d", app.Image_max_size) + ` bytes.`, //checked

		`upload_one_image_title`:       `Images (gif,jpg,png) pictures, size no greater than: ` + fmt.Sprintf("%d", app.Image_max_size) + ` bytes each, max dimensions: ` + fmt.Sprintf("%d x %d", app.Image_max_width, app.Image_max_height) + ` pixels.`,
		`upload no image`:            `The file you submitted is not an image.`,



		`auth_login_email_or_password_incorrect`:"Email or password was incorrect.",


		//login attempts exceed max attempts in config


		`auth_nick_not_exist`:"Nick does not exist",
		`auth_phone_not_exist`:"PhoneCode does not exist",
		`auth_nick_or_email_not_exist`:"Nick or email address not exist.",
		`auth nick or email or phone not exist`:"Nick or email address or phone number not exist.",
		`auth_nick_or_email_required`:"Nick or email is required.",
		`auth_email_required`:"Email is required.",//checked

		`auth_nick_or_email_valid`:"The Nick or Email field must contain valid nick name or email.",

		`auth_reg_disabled`:"Registration has been disabled.",
		`auth email or nick or phone and password are not exist`:"Email or Nick or PhoneCode and Password are not exist.",


		`auth_email_exist_allready`:"Email allready exists choose other email",
		`auth_nick_exist_allready`:"Nick allready exists choose other nick",
		"auth phone exist allready":"PhoneCode number allready exists choose other phone number",
		`auth email or nick or phone exist allready`:"Email or nick or phone allready exists choose other email or nick or phone",
		"auth email or phone exist allready":"Email or phone allready exists choose other email or phone",

		`auth_email_sent_allready`:"There is already activation code waiting to be activated",

		`auth_activation_success`:"Your account have been successfully activated.",//checked
		`auth_activation_allready_activated`:"Your account had been allready activated.",
		`auth_i_agree`:"You must agree before register.",
		`auth_success_reg_check`:"You have successfully registered. Check your email address to activate your account.",
		`auth_success_reg`:"You have successfully registered.",



		`auth_activation_notyet`:"Your account hasn't been activated yet. Please check your email.",
		`auth_activation_sent_allready`:"Your request to change password is already sent. Please check your email.",
		`auth_not_activated`:`Your account hasn't been activated yet. Please check your email.`,

		//******************** Email subject
		`auth_account_subject`:         `данные аккаунта`,
		`auth_activate_subject`:        `активация`,
		`auth_forgot_password_subject`: `запрос нового пароля`,

		`auth_activate_content`:`
%v !

Для активации Вашего аккаунта, вы должны пройти по следующей ссылке:
%v

Пожалуйста активируйте аккаунт в течение %d часов, иначе исчетет отведенный срок для активации.

Ваш имейл и пароль:

Email: %v
Password: %v

Желаем всего наилучшего,
 команда %v`,
		//******************************
		`auth_forgot_password_content`:`
%v,

Вы запросили изменение вашего пароля, потому что вы забыли пароль.

Email: %v
Password: %v

Если у вас возникли проблемы, тогда свяжитесь по %v.

Желаем всего наилучшего,
 команда %v`,

		//******************************
		`auth_account_content`: `
%v,

Спасибо за регистрацию. Ваш аккаунт был успешно создан.

Вы можете залогиниться по имейлу и паролю:

Email: %s
Password: %v

Залогинившись Вы можете попасть в личный кабинет по ссылке:
%v

Мы надеемся на что Вы проведете хорошее время с нами.

Желаем всего наилучшего,
 команда %v`,
	},


}
