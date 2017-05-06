package public

import (
	"Malina/config"
	"Malina/libraries"
	"Malina/language"
	"fmt"
	"time"
)
var out string
func Header(cartPublicListJSON []byte)string  {
	out = `<!DOCTYPE html>
<html>
<head>
	<title>Big shope A Ecommerce Category Flat Bootstarp Resposive Website Template | Home :: w3layouts</title>
	<base href="`+ app.Base_url() +`">
	<link href="`+app.Assets_public_path()+`css/bootstrap.css" rel="stylesheet" type="text/css" media="all" />
	<!--theme-style-->
	<link href="`+app.Assets_public_path()+`css/style.css" rel="stylesheet" type="text/css" media="all" />
	<!--//theme-style-->
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<script type="application/x-javascript">
		addEventListener("load", function() {
			setTimeout(hideURLbar, 0); }, false);
			function hideURLbar(){ window.scrollTo(0,1);
		}
		var CART_DETAILS = ` + string(cartPublicListJSON) + `;
	</script>
	<!--fonts-->
	<link href='http://fonts.googleapis.com/css?family=Open+Sans:400,300,600,700,800' rel='stylesheet' type='text/css'>
	<!--//fonts-->
	<script src="` + app.Assets_backend_path() + `js/jquery.min.js"></script>
	<script src="` + app.Assets_path() + `js/public.js?`+fmt.Sprintf("%v",time.Now())+`"></script>
	<!--script-->
</head>
<body>
	<div id="loader" class="loader" style="background-color: black;position: absolute;top:30%; left: 50%; ;z-index: 1500;display: none;">
		<img src="`+app.Assets_path()+`img/ajax_loader_blue_256.gif" />
	</div>

	<!--header-->
	<div class="header">
		<div class="top-header">
			<div class="container">
				<div class="top-header-left">
					<ul class="support">
						<li><a href="#"><label> </label></a></li>
						<li><a href="#">24x7 live<span class="live"> support</span></a></li>
					</ul>
					<ul class="support">
						<li class="van"><a href="#"><label> </label></a></li>
						<li><a href="#">Free shipping <span class="live">on order over 500</span></a></li>
					</ul>
					<ul class="support">
						<li><a href="#"></a></li>
						<li><a href="public/product/list/"><span class="live">View all products</span></a></li>
					</ul>
					<div class="clearfix"> </div>
				</div>
				<div class="top-header-right">

					<div class="down-top">
						<select class="in-drop">
							` + lang.LangSelectOptions + `
						</select>
					 </div>
					<div class="down-top top-down">
						  <select class="in-drop">

						  <option value="Dollar" class="in-of">Dollar</option>
						  <option value="Yen" class="in-of">Yen</option>
						  <option value="Euro" class="in-of">Euro</option>
							</select>
					 </div>
					 <!---->
					<div class="clearfix"> </div>
				</div>
				<div class="clearfix"> </div>
			</div>
		</div>
		<div class="bottom-header">
			<div class="container">
				<div class="header-bottom-left">
					<div class="logo">
						<a href="/"><img src="`+ app.Assets_public_path() +`images/logo.png" alt=" " /></a>
					</div>
					<div class="search">
						<input type="text" value="" onfocus="this.value = '';" onblur="if (this.value == '') {this.value = '';}" >
						<input type="submit"  value="SEARCH">
					</div>
					<div class="clearfix"> </div>
				</div>
				<div class="header-bottom-right">`;
					if library.SESSION.GetSessionObj().GetAccount_id() > 0 {
						out+=`
						<div class="account"><a href="public/cart/details/"><span> </span>`+library.SESSION.GetSessionObj().GetEmail()+`&nbsp;`+library.SESSION.GetSessionObj().GetNick()+`&nbsp;`+library.SESSION.GetSessionObj().GetPhone()+`</a></div>
						<ul class="login">
							<li><a href="` + app.Base_url()+app.Uri_auth_logout() + `"><span> </span>LOGOUT</a></li>
						</ul>`
					}else {
						out+=`
						<ul class="login">
							<li><a href="` + app.Base_url()+app.Uri_auth_login() + `"><span> </span>` + lang.T("login") + `</a></li> |
							<li ><a href="` + app.Base_url()+app.Uri_auth_register() + `">` + lang.T("sign up") + `</a></li>
						</ul>`
					}; out+=`
						<div class="cart">
							<a href="public/cart/details/">
								<span> </span>` + lang.T("cart") + `&nbsp;(<a id="cart_content">0</a>)
							</a>
						</div>
						<div class="clearfix"> </div>
				</div>
				<div class="clearfix"> </div>
			</div>
		</div>
	</div>`
	return out
}


func Footer()string  {
	return `
<div class="sub-cate">
	<div class=" top-nav rsidebar span_1_of_left">
		<h3><a href="public/product/list/">CATEGORIES</a></h3>
		<ul class="menu">
			`+library.CATEGORY.PublicListView+`
		</ul>
	</div>
	<!--initiate accordion-->
		<script type="text/javascript">
			$(function() {
			    var menu_ul = $('.menu li > ul'),
			    menu_a = $('.menu li > a');
			    menu_ul.hide();
			    menu_a.click(function(e) {
			        e.preventDefault();
			        if(!$(this).hasClass('active')) {
			            menu_a.removeClass('active');
			            //menu_ul.filter(':visible').slideUp('normal');
			            $(this).addClass('active').next().stop(true,true).slideDown('normal');
			        } else {
			            $(this).removeClass('active');
			            $(this).next().stop(true,true).slideUp('normal');
			        }
			    });
			});
		</script>
					<div class=" chain-grid menu-chain">
	   		     		<a href="single.html"><img class="img-responsive chain" src="`+ app.Assets_public_path() +`images/wat.jpg" alt=" " /></a>
	   		     		<div class="grid-chain-bottom chain-watch">
		   		     		<span class="actual dolor-left-grid">300$</span>
		   		     		<span class="reducedfrom">500$</span>
		   		     		<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     		</div>
	   		     	</div>
	   		     	 <a class="view-all all-product" href="product.html">VIEW ALL PRODUCTS<span> </span></a>
			</div>
	   		    <div class="clearfix"> </div>
		</div>
	<!---->
	<div class="footer">
		<div class="footer-top">
			<div class="container">
				<div class="latter">
					<h6>NEWS-LETTER</h6>
					<div class="sub-left-right">
						<form>
							<input type="text" value="Enter email here"onfocus="this.value = '';" onblur="if (this.value == '') {this.value = 'Enter email here';}" />
							<input type="submit" value="SUBSCRIBE" />
						</form>
					</div>
					<div class="clearfix"> </div>
				</div>
				<div class="latter-right">
					<p>FOLLOW US</p>
					<ul class="face-in-to">
						<li><a href="#"><span> </span></a></li>
						<li><a href="#"><span class="facebook-in"> </span></a></li>
						<div class="clearfix"> </div>
					</ul>
					<div class="clearfix"> </div>
				</div>
				<div class="clearfix"> </div>
			</div>
		</div>
		<div class="footer-bottom">
			<div class="container">
				<div class="footer-bottom-cate">
					<h6>CATEGORIES</h6>
					<ul>
						<li><a href="#">Curabitur sapien</a></li>
						<li><a href="#">Dignissim purus</a></li>
						<li><a href="#">Tempus pretium</a></li>
						<li ><a href="#">Dignissim neque</a></li>
						<li ><a href="#">Ornared id aliquet</a></li>
						<li><a href="#">Ultrices id du</a></li>
						<li><a href="#">Commodo sit</a></li>
						<li ><a href="#">Urna ac tortor sc</a></li>
						<li><a href="#">Ornared id aliquet</a></li>
						<li><a href="#">Urna ac tortor sc</a></li>
						<li ><a href="#">Eget nisi laoreet</a></li>
						<li ><a href="#">Faciisis ornare</a></li>
					</ul>
				</div>
				<div class="footer-bottom-cate bottom-grid-cat">
					<h6>FEATURE PROJECTS</h6>
					<ul>
						<li><a href="#">Curabitur sapien</a></li>
						<li><a href="#">Dignissim purus</a></li>
						<li><a href="#">Tempus pretium</a></li>
						<li ><a href="#">Dignissim neque</a></li>
						<li ><a href="#">Ornared id aliquet</a></li>
						<li><a href="#">Ultrices id du</a></li>
						<li><a href="#">Commodo sit</a></li>
					</ul>
				</div>
				<div class="footer-bottom-cate">
					<h6>TOP BRANDS</h6>
					<ul>
						<li><a href="#">Curabitur sapien</a></li>
						<li><a href="#">Dignissim purus</a></li>
						<li><a href="#">Tempus pretium</a></li>
						<li ><a href="#">Dignissim neque</a></li>
						<li ><a href="#">Ornared id aliquet</a></li>
						<li><a href="#">Ultrices id du</a></li>
						<li><a href="#">Commodo sit</a></li>
						<li ><a href="#">Urna ac tortor sc</a></li>
						<li><a href="#">Ornared id aliquet</a></li>
						<li><a href="#">Urna ac tortor sc</a></li>
						<li ><a href="#">Eget nisi laoreet</a></li>
						<li ><a href="#">Faciisis ornare</a></li>
					</ul>
				</div>
				<div class="footer-bottom-cate cate-bottom">
					<h6>OUR ADDERSS</h6>
					<ul>
						<li>Aliquam metus  dui. </li>
						<li>orci, ornareidquet</li>
						<li> ut,DUI.</li>
						<li >nisi, dignissim</li>
						<li >gravida at.</li>
						<li class="phone">PH : 6985792466</li>
						<li class="temp"> <p class="footer-class">Design by <a href="http://w3layouts.com/" target="_blank">W3layouts</a> </p></li>
					</ul>
				</div>
				<div class="clearfix"> </div>
			</div>
		</div>
	</div>
</body>
</html>`

}


func Error(msg string)string{
	return `
	<html>
		<body>
			<h1>` + msg + `</h1>
		</body>
	</html>`
}

/*

<li><a href="#">Curabitur sapien<img class="arrow-img" src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/> </a>
					<ul>
						<li><a href="product.html">Cute Kittens </a></li>
						<li class="subitem2"><a href="product.html">Strange Stuff </a></li>
						<li class="item11"><a href="#">Automatic Fails<img class="arrow-img" src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/> </a>
							<ul>
								<li class="subitem1"><a href="product.html">Cute Kittens </a></li>
								<li class="item111"><a href="#">11Automatic Fails<img class="arrow-img" src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/> </a>
									<ul>
										<li class="subitem1"><a href="product.html">Cute Kittens </a></li>
										<li class="subitem2"><a href="product.html">Strange Stuff </a></li>
										<li class="subitem3"><a href="product.html">Automatic Fails </a></li>
									</ul>
								</li>
								<li class="subitem3"><a href="product.html">Automatic Fails </a></li>
							</ul>
						</li>
					</ul>
				</li>
		<li class="item2"><a href="#">Dignissim purus <img class="arrow-img " src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/></a>
			<ul class="cute">
				<li class="subitem1"><a href="product.html">Cute Kittens </a></li>
				<li class="subitem2"><a href="product.html">Strange Stuff </a></li>
				<li class="subitem3"><a href="product.html">Automatic Fails </a></li>
			</ul>
		</li>
		<li class="item3"><a href="#">Ultrices id du<img class="arrow-img img-arrow" src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/> </a>
			<ul class="cute">
				<li class="subitem1"><a href="product.html">Cute Kittens </a></li>
				<li class="subitem2"><a href="product.html">Strange Stuff </a></li>
				<li class="subitem3"><a href="product.html">Automatic Fails</a></li>
			</ul>
		</li>
		<li class="item4"><a href="#">Cras iacaus rhone <img class="arrow-img img-left-arrow" src="`+ app.Assets_public_path() +`images/arrow1.png" alt=""/></a>
			<ul class="cute">
				<li class="subitem1"><a href="product.html">Cute Kittens </a></li>
				<li class="subitem2"><a href="product.html">Strange Stuff </a></li>
				<li class="subitem3"><a href="product.html">Automatic Fails </a></li>
			</ul>
		</li>
				<li>
			<ul class="kid-menu">
				<li><a href="product.html">Tempus pretium</a></li>
				<li ><a href="product.html">Dignissim neque</a></li>
				<li ><a href="product.html">Ornared id aliquet</a></li>
			</ul>
		</li>
		<ul class="kid-menu ">
				<li><a href="product.html">Commodo sit</a></li>
				<li ><a href="product.html">Urna ac tortor sc</a></li>
				<li><a href="product.html">Ornared id aliquet</a></li>
				<li><a href="product.html">Urna ac tortor sc</a></li>
				<li ><a href="product.html">Eget nisi laoreet</a></li>
				<li><a href="product.html">Faciisis ornare</a></li>
				<li class="menu-kid-left"><a href="contact.html">Contact us</a></li>
			</ul>

 */