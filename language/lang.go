package lang

import (
	"fmt"
	"github.com/ilyaran/Malina/config"
	"strconv"
)

func T(key string)string  {
	if v,ok:= Dict[key];ok{
		return v
	}
	return key
}
var LangSelectOptions string
var FullLanguageList = map[string]string{
	"en" : "English",
	"ru" : "Russian",
	"zh-CN" : "Simplified	Chinese",
	"zh-TW" : "Traditional	Chinese",
	"hi" : "Hindi",
	"ko" : "Korean",
	"yi" : "Yiddish",
	"am" : "Armenian",
	"es" : "Spanish",
	"fr" : "French",
	"nl" : "Dutch",
	"cs" : "Czech",
	"no" : "Norwegian",
	"sv" : "Swedish",
	"el" : "Greek",

}
var Dict = map[string]string{

	`not find`:`not find any item`,

	`validation text field`	:  `The %s field must contain %s`,

	`validation integer`	:  `The %s field must contain an integer.`,
	`validation email`	:"The Email field must contain all valid email addresses.",
	`validation phone`	:"The Phone field must contain all valid phone numbers.",
	`validation required`	:`The %s field is required.`,
	`validation min length`	: `The %s field must be at least %d characters in length.`,
	`validation max length`	: `The %s field cannot exceed %d characters in length.`,

	`server error`: `Server Error`,
	"no post data":"No Post Data",
	"not auth":"You must login",
	"no post request":"It is not a Post Method Request",

	`upload no base64`:               `Unable to find a base64 format.`,
	`upload error base64 decoding`:   `Unable to decode a base64 string.`,
	`upload destination error`:       `A problem was encountered while attempting to move the uploaded file to the final destination.`,
	`upload invalid filetype`:        `The filetype you are attempting to upload is not allowed.`,
	`upload_invalid_dimensions`:        `The file dimensions you are attempting to upload is not allowed.`,
	`upload_unable_to_write_file`:        `Unable to write file.`,
	`upload_file_exceeds_limit`:      `The uploaded file exceeds the maximum allowed size in ` + fmt.Sprintf("%d", app.Image_max_size()) + ` bytes.`,  //checked

	`upload_one_image_title`:       `Images (gif,jpg,png) pictures, size no greater than: ` + fmt.Sprintf("%d", app.Image_max_size()) + ` bytes each, max dimensions: ` + fmt.Sprintf("%d x %d", app.Image_max_width(),app.Image_max_height()) + ` pixels.`,
	`upload no image`:            `The file you submitted is not an image.`,

	`wrong_recaptcha`:"The Captcha field is wrong. May be you are a robot.",

	`auth_login_email_or_password_incorrect`:"Email or password was incorrect.",
	`auth_login_allready`:"You had been allready sign in.",

	//login attempts exceed max attempts in config
	`auth exceed max attempts`:"Attempts exceed max attempts in config. No more than "+ strconv.FormatInt(app.Login_attempts(),10)+ " times per 24 hours",

	`auth_email_not_exist`:"Email does not exist",
	`auth_nick_not_exist`:"Nick does not exist",
	`auth_phone_not_exist`:"Phone does not exist",
	`auth_nick_or_email_not_exist`:"Nick or email address not exist.",
	`auth nick or email or phone not exist`:"Nick or email address or phone number not exist.",
	`auth_nick_or_email_required`:"Nick or email is required.",
	`auth_email_required`:"Email is required.",//checked

	`auth_nick_or_email_valid`:"The Nick or Email field must contain valid nick name or email.",

	`auth_reg_disabled`:"Registration has been disabled.",
	`auth email or nick or phone and password are not exist`:"Email or Nick or Phone and Password are not exist.",
	`auth_logout_first`:"You have to logout first, before registering.",
	`auth_reset_first`:"You have to login first, before reeset password.",
	`auth_email_exist_allready`:"Email allready exists choose other email",
	`auth_nick_exist_allready`:"Nick allready exists choose other nick",
	"auth phone exist allready":"Phone number allready exists choose other phone number",
	`auth email or nick or phone exist allready`:"Email or nick or phone allready exists choose other email or nick or phone",
	"auth email or phone exist allready":"Email or phone allready exists choose other email or phone",

	`auth_email_sent_allready`:"There is already activation code waiting to be activated",
	`auth_confirm_pass_valid`:"The Confirm Password field must be same with password field.",
	`auth_activation_success`:"Your account have been successfully activated.",//checked
	`auth_activation_allready_activated`:"Your account had been allready activated.",
	`auth_i_agree`:"You must agree before register.",
	`auth_success_reg_check`:"You have successfully registered. Check your email address to activate your account.",
	`auth_success_reg`:"You have successfully registered.",

	`auth_activation_incorrect_code`:"The activation code you entered was incorrect. Please check your email again.",

	`auth_activation_notyet`:"Your account hasn't been activated yet. Please check your email.",
	`auth_activation_sent_allready`:"Your request to change password is already sent. Please check your email.",
	`auth_not_activated`:`Your account hasn't been activated yet. Please check your email.`,

	//******************** Email subject
	`auth_account_subject`:         ` account details`,
	`auth_activate_subject`:        ` activation`,
	`auth_forgot_password_subject`: `New password request`,

	`auth_activate_content`:`
Welcome to %v,

To activate your account, you must follow the activation link below:
%s

Please activate your account within %d hours, otherwise your registration will become invalid and you will have to register again.

You can use either you username or email address to login.
Your login details are as follows:

Email: %s
Password: %s

We hope that you enjoy your stay with us :)

Regards,
The %s Team	`,
	//******************************
	`auth_forgot_password_content`:`
%v,

You have requested your password to be changed, because you forgot the password.

Email: %s
Nick : %s
Password: %v

If you have any more problems with gaining access to your account please contact %v.

Regards,
The %v Team`,

	//******************************
	`auth_account_content`:`
Welcome to %v,

Thank you for registering. Your account was successfully created.

You can login with either your username or email address:

Email: %s
Nick : %s
Password: %v

You can try logging in now by going to %v

We hope that you enjoy your stay with us.

Regards,
The %v Team`,
}
