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
 */ package publicView


import (
	"net/http"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/berry"
	"github.com/ilyaran/Malina/libraries"
	"html/template"
	"strconv"
	"github.com/ilyaran/Malina/lang"
)


func PublicLayoutViewInit(){
	PublicLayoutView =&PublicLayout{}
	PublicLayoutView.SetLayout()
}

var PublicLayoutView *PublicLayout
type PublicLayout struct {
	LayoutMain string
	Per_page_select string
}
func(s *PublicLayout) response(malina *berry.Malina, w http.ResponseWriter){
	var navAuthBar string
	if malina.CurrentAccount!=nil{
		navAuthBar = `
	<div class="header-bottom-right">
		<div class="account">
			<a href="cabinet/profile/"><span> </span>`
		if malina.CurrentAccount.Phone!=""{navAuthBar +=malina.CurrentAccount.Phone}else
		if malina.CurrentAccount.Nick!=""{navAuthBar +=malina.CurrentAccount.Nick}else
		if malina.CurrentAccount.FirstName!=""{navAuthBar +=malina.CurrentAccount.FirstName}else{
			navAuthBar +=malina.CurrentAccount.Email
		}
		navAuthBar +=`</a>
		</div>`
	}else {
		navAuthBar = `
	<div class="header-bottom-right">
		<ul class="login">
			<li><a href="auth/login"><span> </span>LOGIN</a></li> |
			<li ><a href="auth/register">SIGNUP</a></li>
		</ul>`
	}
	navAuthBar +=`
		<div class="cart"><a href="cabinet/cart/"><span> </span>CART</a></div>
		<div class="clearfix"> </div>
	</div>`

	malina.NavAuth= template.HTML(navAuthBar)

	t := template.New(malina.ControllerName)
	t, _ = t.Parse(PublicLayoutView.LayoutMain)
	t.Execute(w, malina)
}
func (s *PublicLayout) SetPer_page_select()string{
	var out string
	var step = (app.Per_page_max -1-app.Per_page)/app.Per_page_select_options_widget_num
	for i:=int64(app.Per_page); i<(app.Per_page_max -step); i+= step {
		out += `<option value="`+strconv.FormatInt(i,10)+`">`+strconv.FormatInt(i,10)+`</option>`
	}
	out += `<option value="`+strconv.FormatInt(app.Per_page_max-1,10)+`">`+strconv.FormatInt(app.Per_page_max-1,10)+`</option>`

	s.Per_page_select = lang.T("per page") + `<select name="per_page" id="per_page" >`+out+`</select>`
	return lang.T("per page") + `<select name="per_page" id="per_page" >`+out+`</select>`
}

func (s *PublicLayout) navbar(malina *berry.Malina) string {
	var out string
	out+=`


	`
	return out
}


func (s *PublicLayout)SetLayout(){
	s.SetPer_page_select()
	s.LayoutMain = `
<!--A Design by W3layouts
Author: W3layout
Author URL: http://w3layouts.com
License: Creative Commons Attribution 3.0 Unported
License URL: http://creativecommons.org/licenses/by/3.0/
-->
<!DOCTYPE html>
<html>
<head>
	<title>Big shope A Ecommerce Category Flat Bootstarp Resposive Website Template | Product :: w3layouts</title>
	<base href="`+ app.Base_url +`"/>
	<link href="`+app.Url_assets_public+`css/bootstrap.css" rel="stylesheet" type="text/css" media="all" />
	<!--theme-style-->
	<link href="`+app.Url_assets_public+`css/style.css?fff=dfdf" rel="stylesheet" type="text/css" media="all" />
	<!--//theme-style-->
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<script type="application/x-javascript"> addEventListener("load", function() { setTimeout(hideURLbar, 0); }, false); function hideURLbar(){ window.scrollTo(0,1); } </script>
	<!--fonts-->
	<link href='http://fonts.googleapis.com/css?family=Open+Sans:400,300,600,700,800' rel='stylesheet' type='text/css'>
	<!--//fonts-->
	<script src="`+app.Url_assets_public+`js/jquery.min.js"></script>
	<link rel="stylesheet" type="text/css" href="`+ app.Url_assets_public+`css/jquery-ui1.css">

	<!--script-->
</head>
<body>

	<input type="hidden" value="6000" id="MAX_PRICE"/>


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
					<div class="clearfix"> </div>
				</div>

				<div class="top-header-right down-top">
					<a href="`+app.Url_public_product_list+`">Products</a>
					<div class="clearfix"> </div>
				</div>
				<div class="top-header-right">

					<div class="down-top">
						  <select class="in-drop">
							  <option value="English" class="in-of">English</option>
							  <option value="Japanese" class="in-of">Japanese</option>
							  <option value="French" class="in-of">French</option>
							  <option value="German" class="in-of">German</option>
							</select>
					</div>
					<div class="down-top top-down">
						  <select class="in-drop">

						  <option value="Dollar" class="in-of">Dollar</option>
						  <option value="Yen" class="in-of">Yen</option>
						  <option value="Euro" class="in-of">Euro</option>
							</select>
					</div>

					<div class="clearfix"> </div>

				</div>

				<div class="clearfix"> </div>
			</div>
		</div>
		<div class="bottom-header">
			<div class="container">
				<div class="header-bottom-left">
					<div class="logo">
						<a href="`+app.Base_url+`"><img src="`+app.Url_assets_public+`images/logo.png" alt=" " /></a>
					</div>
					<div class="search">
						<input id="search" type="text" value="" onfocus="this.value = '';" onblur="if (this.value == '') {this.value = '';}" />
						<button id="get_list_button" type="submit" class="btn btn-default">SEARCH</button>

					</div>
					<div class="clearfix"> </div>
				</div>

				{{ .NavAuth }}

				<div class="clearfix"> </div>
			</div>
		</div>
	</div>
	<!---->
	<!-- start content -->
	<div class="container">



			{{ .Content }}



	<div class="sub-cate">


			<div class="range">
                    <h3 class="sear-head">Price range</h3>
                    <ul class="dropdown-menu6">
                        <li>
                            <div id="slider-range"></div>
						<div id="amount"></div>
                            <input type="text" id="amount1" style="border: 0; color: #ffffff; font-weight: normal;" />
                        </li>
                    </ul>
                    <!---->
                    <script type="text/javascript" src="`+ app.Url_assets_public+`js/jquery-ui.js"></script>
                    <script type='text/javascript'>//<![CDATA[
                    $(window).load(function(){
                        $( "#slider-range" ).slider({
                            range: true,
                            min: 0,
                            max: $("#MAX_PRICE").val()*1,
                            values: [ 0, $("#MAX_PRICE").val()*1 ],
                            slide: function( event, ui ) {  $( "#amount" ).html( "$" + ui.values[ 0 ] + " - $" + ui.values[ 1 ] );
                            }
                        });
                        $( "#amount" ).html("$ " + $("#slider-range").slider("values",0)+" - $" + $("#slider-range").slider("values",1));

                    	function getList(current_page) {
							var arg_data = {};
							arg_data.page = current_page;
							arg_data.per_page = $("#per_page").val();
							if ($("#search").val() != "") arg_data.search = $("#search").val();
							arg_data.order_by = $("#order_by").val();
							arg_data.price_min = $( "#slider-range" ).slider( "values", 0 );
							arg_data.price_max = $( "#slider-range" ).slider( "values", 1 );

							$.get( "`+ app.Url_public_product_list+`", arg_data, function( data ) {
								$("#product_listing").html(data);
							},"html");
						}
						$("#get_list_button").on("click",function () {
							getList(0);
						});
						$("body").on("click",".page_product",function () {
							var current_page = $(this).children("span").attr("data-page");
							getList(current_page);
						});


                    }); //]]>
                    </script>

                </div>



				<div class=" top-nav rsidebar span_1_of_left">
					<h3 class="cate">CATEGORIES</h3>
		  			<ul class="menu">
                        `+libraries.CategoryLib.SetPublicListView()+`
                    </ul>
                </div>

                        {{ .CategoryParameters }}

                <!--<div class=" chain-grid menu-chain">
	   		     	<div class="grid-chain-bottom chain-watch">
		   		    	<span class="actual">Processors gg</span><br/>
                        <span class="actual">&nbsp;4-Core</span><br/>
		   		     	<span class="actual">&nbsp;&nbsp;Threads</span><br/>
                        &nbsp;&nbsp;&nbsp;<input type="checkbox"><span class="">   8 thread</span><br/>
                        &nbsp;&nbsp;&nbsp;<input type="checkbox"><span class="">   16 thread</span><br/>
					</div>
	   		    </div>
				<div class=" chain-grid menu-chain">
	   		     	<div class="grid-chain-bottom chain-watch">
		   		    	<span class="actual">Processors gg</span><br/>
                        <span class="actual">&nbsp;4-Core</span><br/>
		   		     	<span class="actual">&nbsp;&nbsp;Threads</span><br/>
                        &nbsp;&nbsp;&nbsp;<input type="checkbox"><span class="">   8 thread</span><br/>
                        &nbsp;&nbsp;&nbsp;<input type="checkbox"><span class="">   16 thread</span><br/>
					</div>
	   		    </div>-->
                <!--initiate accordion-->
				<script type="text/javascript">
					$(function() {
						var menu_ul = $('.menu > li > ul');
						var menu_a  = $('.menu > li > span > a');
						var menu_button  = $('.menu li > span > button');
						menu_ul.hide();
						menu_button.click(function(e) {
							if($(this).parent().next("ul").is(':visible')){
								$(this).html("V");
								//$(this).parent().next("ul").css('border','4px solid #f3c500 !important;');
								$(this).parent().next("ul").css('background-color','#fff');
							}else {
								$(this).html(">");
							}
							$(this).parent().next("ul").toggle(300);
						});
					});
				</script>

				<div class=" chain-grid menu-chain">
	   		    	<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/wat.jpg" alt=" " /></a>
	   		     	<div class="grid-chain-bottom chain-watch">
		   		    	<span class="actual dolor-left-grid">300$</span>
		   		     	<span class="reducedfrom">500$</span>
		   		     	<h6>Lorem ipsum dolor</h6>
	   		     	</div>
	   		    </div>
	   		    <a class="view-all all-product" href="`+app.Url_public_product_list+`">VIEW ALL PRODUCTS<span> </span></a>
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
</html>
	`

}


