package publicView

import (
	"github.com/ilyaran/Malina/berry"
	"net/http"
	"html/template"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/lang"
)

type Auth struct {

}


func (s *Auth)Index(malina *berry.Malina,message string,w http.ResponseWriter){
	content :=`
	<div class="account_grid">
		<div class=" login-right">
			<div class="col-md-12">
				<div class="create-account col">
					<div class="co">
						<h3>`+message+`</h3>
					</div>
				</div>
			</div>
		</div>
	</div>`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func (s *Auth)Login(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	content := `
	<div class="account_grid">
		<div class=" login-right">
			<h3>REGISTERED CUSTOMERS</h3>
			<p>If you have an account with us, please log in.</p>
			<form method="POST" >
				<div>
					<span>Email Address<label>*</label></span>
					<input value="`+r.FormValue("email")+`" name="email" type="email" required/>
					<span class="error">`;if v,ok:= malina.Result["email"]; ok {content += v.(string)};content += `</span>
				</div>
				<div>
					<span>Password<label>*</label></span>
					<input type="text" name="password" required/>
					<span class="error">`;if v,ok:= malina.Result["password"]; ok {content += v.(string)};content += `</span>
				</div>
				<a class="forgot" href="#">Forgot Your Password?</a>
				<input type="submit" value="Login">
			</form>
		</div>
		<div class=" login-left">
			<h3>NEW CUSTOMERS</h3>
			<p>By creating an account with our store, you will be able to move through the checkout process faster, store multiple shipping addresses, view and track your orders in your account and more.</p>
			<a class="acount-btn" href="auth/register">Create an Account</a>
		</div>
		<div class="clearfix"> </div>
	</div>

	`

	malina.Content = template.HTML(content)

	PublicLayoutView.response(malina,w)
}

func (s *Auth)Register(malina *berry.Malina,w http.ResponseWriter, r *http.Request){
	content := `
<div class="register">
		  	  <form method="POST" >
				 <div class="  register-top-grid">
					<h3>PERSONAL INFORMATION</h3>
					<div class="mation">
						<span>First Name<label>*</label></span>
						<input value="`+r.FormValue("first_name")+`" name="first_name" type="text" required/>
						<span class="error">`;if v,ok:= malina.Result["first_name"]; ok {content += v.(string)};content += `</span>
						<span>Last Name<label>*</label></span>
						<input value="`+r.FormValue("last_name")+`" name="last_name" type="text" required/>
						<span class="error">`;if v,ok:= malina.Result["last_name"]; ok {content += v.(string)};content += `</span>

						 <span>Email Address<label>*</label></span>
						 <input value="`+r.FormValue("email")+`" name="email" type="email" required/>
						<span class="error">`;if v,ok:= malina.Result["email"]; ok {content += v.(string)};content += `</span>


					</div>
					 <div class="clearfix"> </div>
					   <a class="news-letter" href="#">
						 <label class="checkbox"><input type="checkbox" name="checkbox" checked=""><i> </i>Sign Up</label>
					   </a>
					 </div>
				     <div class="  register-bottom-grid">
						    <h3>LOGIN INFORMATION</h3>
							<div class="mation">
								<span>Password<label>*</label></span>
								<input type="text" name="password" required/>
						<span class="error">`;if v,ok:= malina.Result["password"]; ok {content += v.(string)};content += `</span>
								<span>Confirm Password<label>*</label></span>
								<input type="text" name="confirm_password" required/>
						<span class="error">`;if v,ok:= malina.Result["confirm_password"]; ok {content += v.(string)};content += `</span>
							</div>
					 </div><div class="clearfix"> </div>
						<input type="submit" value="submit">
					   <div class="clearfix"> </div>
				</form>
				<div class="clearfix"> </div>
				<div class="register-but">
				   <form method="POST" >

				   </form>
				</div>
		   </div>

	`


	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func(s *Auth) RegistrationResult(malina *berry.Malina,email string,w http.ResponseWriter){

	content :=`

		<div class="col-md-12">
					<div class="create-account col">
						<div class="co">
							<h3>Завершите регистрацию</h3>
							<h6>Перейдите в свой почтовый ящик и нажмите <br> кнопку «Завершить регистрацию» в письме от `+app.Site_name+`</h6>
							<h6><b>`+email+`</b></h6>
						</div>
						<div class="col-descr">
							<h6>1. Если вы не увидели наше письмо, проверьте папку нежелательной почты (Spam).</h6>
							<h6>2. Убедитесь, что вы правильно ввели e-mail: <b>`+email+`</b>. Если вы допустили ошибку, создайте аккаунт заново или обратитесь в <a href="`+app.Url_support+`">Центр поддержки</a></h6>
						</div>
					</div>
		</div><!-- /.row -->

	`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func(s *Auth) ActivationSuccess(malina *berry.Malina,w http.ResponseWriter){

	content :=`
<div class="account_grid">
		<div class=" login-right">
		<div class="col-md-12">
				<div class="row justify-content-md-center">
						<div class="create-account">
							<h3>Поздравляем! Ваш аккаунт создан</h3>
							<h6>Вы можете перейти в свой <a href="`+app.Base_url+`cabinet/profile">Личный профиль</a> или <br>перейти на <a href="/">Главную страницу</a></h6>
						</div>
						<div class="col-md-9 main-shop">
						<h3>Магазин <a href="#" class="view-all">Смотреть все</a></h3>
							<div class="row">
								<div class="col-sm-3">
									<div class="shop-view">
										<div class="screen">
											<span>CS:GO</span>
											<div class="name">USP-S | Неонуар</div>
											<img src="`+ app.Url_assets_public+`theme/img/weapon/1.png" alt="">
										</div>
										<div class="descript">
											<div class="balance"><b></b> 5000</div>
											<a class="buy-icon" href="#">
												<span>В корзину</span>
											</a>
										</div>
									</div><!-- /.shop-view -->
								</div>
								<div class="col-sm-3">
									<div class="shop-view">
										<div class="screen">
											<span>CS:GO</span>
											<div class="name">MP9 - Бульдозер</div>
											<img src="`+ app.Url_assets_public+`theme/img/weapon/2.png" alt="">
										</div>
										<div class="descript">
											<div class="balance"><b></b> 5000</div>
											<a class="buy-icon" href="#">
												<span>В корзину</span>
											</a>
										</div>
									</div><!-- /.shop-view -->
								</div>
								<div class="col-sm-3">
									<div class="shop-view">
										<div class="screen">
											<span>CS:GO</span>
											<div class="name">M4A4 - Вой</div>
											<img src="`+ app.Url_assets_public+`theme/img/weapon/3.png" alt="">
										</div>
										<div class="descript">
											<div class="balance"><b></b> 5000</div>
											<a class="buy-icon" href="#">
												<span>В корзину</span>
											</a>
										</div>
									</div><!-- /.shop-view -->
								</div>
								<div class="col-sm-3">
									<div class="shop-view">
										<div class="screen">
											<span>CS:GO</span>
											<div class="name">AWP - Медуза</div>
											<img src="`+ app.Url_assets_public+`theme/img/weapon/4.png" alt="">
										</div>
										<div class="descript">
											<div class="balance"><b></b> 5000</div>
											<a class="buy-icon" href="#">
												<span>В корзину</span>
											</a>
										</div>
									</div><!-- /.shop-view -->
								</div>
							</div>
						</div><!-- /.main-shop -->
				</div><!-- /.row -->

		</div>	</div>	</div>

	`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func(s *Auth) ForgotForm(malina *berry.Malina,message string,w http.ResponseWriter, r *http.Request){

	content :=`
<div class="account_grid">
		<div class=" login-right">
	<div class="col-md-12">
		<div class="create-account col">
			<div class="co">`
	if message == `` {
		content += `
				<h3>Забыли пароль</h3>
				<h6>Введите email.</h6>
				<form action="` + app.Base_url + `auth/forgot/" method="POST">
					<h6><b><input value="` + r.FormValue("email") + `" name="email" type="email"  placeholder="email"/></b></h6>
					<h3 class="error">` +message + `</h3>
					<h6><b><input type="submit" value="Отправить"/></b></h6>
				</form>`
	}else{
		content += `<h3>`+message+`</h3>`
	}
	content += `
			</div>
		</div>
	</div>
		</div>
	</div>

	<!-- /.row -->`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func(s *Auth) ForgotResult(malina *berry.Malina,message string,w http.ResponseWriter){

	content :=`

		<div class="col-md-12">
					<div class="create-account col">
						<div class="co">
							<h3>Новый пароль отправлен.</h3>
							<div class="col-descr">
								<h6>1. Перейдите в свой почтовый ящик <br> и вы найдете новый сгенерированный пароль в письме от `+app.Site_name+`</h6>
								<h6>2. Если вы не увидели наше письмо, проверьте папку нежелательной почты (Spam).</h6>
								<h6>3. Убедитесь, что вы правильно ввели e-mail: Если вы допустили ошибку, <a href="`+app.Base_url+`auth/forgot/">попробуйте</a> заново или обратитесь в <a href="`+app.Url_support+`">Центр поддержки</a></h6>
							</div>
						</div>
					</div>
		</div><!-- /.row -->

	`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func(s *Auth) ChangePassword(malina *berry.Malina,message string,w http.ResponseWriter, r *http.Request){

	content :=`

		<div class="col-md-12">
				<div class="create-account col">
					<div class="co">`
	if message != `` {
		content += `<h3>`+message+`</h3>`
	}else {
		content += `
						<h3>` + lang.T("reset password") + `</h3>
						<div class="col-descr">
							<form method="post" action="` + app.Url_auth_change_password + `">

								<h3 class="error">`;
		if v, ok := malina.Result["no_rows"]; ok {
			content += v.(string)
		};
		content += `</h3>
								<h3 class="error">`;
		if v, ok := malina.Result["error"]; ok {
			content += v.(string)
		};
		content += `</h3>

		                        <div class="form-group">
		                            <input class="form-control" value="` + r.FormValue("old_password") + `" placeholder="` + lang.T("old password") + `" name="old_password" type="text" autofocus>
		                            <h3 class="error">`;
		if v, ok := malina.Result["old_password"]; ok {
			content += v.(string)
		};
		content += `</h3>
		                        </div>

		                        <div class="form-group">
		                            <input class="form-control" value="` + r.FormValue("new_password") + `" placeholder="` + lang.T("new password") + `" name="new_password" type="text" autofocus>
		                            <span class="error">`;
		if v, ok := malina.Result["new_password"]; ok {
			content += v.(string)
		};
		content += `</span>
		                        </div>

		                        <div class="form-group">
		                            <input class="form-control" value="` + r.FormValue("confirm_new_password") + `" placeholder="` + lang.T("confirm new password") + `" name="confirm_new_password" type="text" autofocus>
		                            <h3 class="error">`;
		if v, ok := malina.Result["confirm_new_password"]; ok {
			content += v.(string)
		};
		content += `</h3>
		                        </div>

		                        <!-- Change this to a button or input when using this as a form -->
		                        <button type="submit" class="btn btn-lg btn-success btn-block">` + lang.T("send") + `</button>

							</form>
						</div>`
	}
	content +=`</div>
		</div>

	`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}

