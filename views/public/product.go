package public

import (
	"Malina/entity"
	"fmt"
	"strconv"
	"Malina/config"
	"Malina/language"
)

func ProductItem(product *entity.Product)string {
	var images string
	if len(product.Get_img())>0{
		for k,v := range product.Get_img(){
			if k < 1{
				images += `
			<li>
				<a href="optionallink.html">
					<img class="etalage_thumb_image" src="`+v+`" class="img-responsive" />
					<img class="etalage_source_image" src="`+v+`" class="img-responsive" title="" />
				</a>
			</li>`
			}else {
				images += `
			<li>
				<img class="etalage_thumb_image" src="`+v+`" class="img-responsive" />
				<img class="etalage_source_image" src="`+v+`" class="img-responsive" title="" />
			</li>`
			}
		}
	}else {
		images = `<li><img class="etalage_thumb_image" src="`+app.No_image()+`" class="img-responsive" /></li>`
	}
	out = `
<script src="`+app.Assets_public_path()+`js/jquery.etalage.min.js"></script>
<script>
	jQuery(document).ready(function($){
		$('#etalage').etalage({
			thumb_image_width: 300,
			thumb_image_height: 400,
			source_image_width: 900,
			source_image_height: 1200,
			show_hint: true,
			click_callback: function(image_anchor, instance_id){
			alert('Callback example:\nYou clicked on an image with the anchor: "'+image_anchor+'"\n(in Etalage instance: "'+instance_id+'")');
		}
	});

});
</script>

<div class="container">
	<div class=" single_top">
		<div class="single_grid">
			<div class="grid images_3_of_2">
				<ul id="etalage">
					` + images + `
				</ul>
				<div class="clearfix"> </div>
			</div>
			<div class="desc1 span_3_of_2">
				<h4>`+product.Get_title()+`</h4>
				<div class="cart-b">
					<div class="left-n ">$` + fmt.Sprintf("%.2f",product.Get_price()) + `</div>
					<a class="now-get get-cart-in add-to-cart" href="#">`+lang.T("add to cart")+`</a>
					<div class="clearfix"></div>
				</div>
				<h6>100 items in stock</h6>
			   	<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.</p>
			   	<div class="share">
					<h5>Share Product :</h5>
					<ul class="share_nav">
						<li><a href="#"><img src="`+app.Assets_public_path()+`images/facebook.png" title="facebook"></a></li>
						<li><a href="#"><img src="`+app.Assets_public_path()+`images/twitter.png" title="Twiiter"></a></li>
						<li><a href="#"><img src="`+app.Assets_public_path()+`images/rss.png" title="Rss"></a></li>
						<li><a href="#"><img src="`+app.Assets_public_path()+`images/gpluse.png" title="Google+"></a></li>
				    	</ul>
				</div>
			</div>
          		<div class="clearfix"> </div>
          	</div>
          	<ul id="flexiselDemo1">
			<li><img src="`+app.Assets_public_path()+`images/pi.jpg" /><div class="grid-flex"><a href="#">Bloch</a><p>Rs 850</p></div></li>
			<li><img src="`+app.Assets_public_path()+`images/pi1.jpg" /><div class="grid-flex"><a href="#">Capzio</a><p>Rs 850</p></div></li>
			<li><img src="`+app.Assets_public_path()+`images/pi2.jpg" /><div class="grid-flex"><a href="#">Zumba</a><p>Rs 850</p></div></li>
			<li><img src="`+app.Assets_public_path()+`images/pi3.jpg" /><div class="grid-flex"><a href="#">Bloch</a><p>Rs 850</p></div></li>
			<li><img src="`+app.Assets_public_path()+`images/pi4.jpg" /><div class="grid-flex"><a href="#">Capzio</a><p>Rs 850</p></div></li>
		</ul>
	    	<script type="text/javascript">
		 $(window).load(function() {
			$("#flexiselDemo1").flexisel({
				visibleItems: 5,
				animationSpeed: 1000,
				autoPlay: true,
				autoPlaySpeed: 3000,
				pauseOnHover: true,
				enableResponsiveBreakpoints: true,
		    	responsiveBreakpoints: {
		    		portrait: {
		    			changePoint:480,
		    			visibleItems: 1
		    		},
		    		landscape: {
		    			changePoint:640,
		    			visibleItems: 2
		    		},
		    		tablet: {
		    			changePoint:768,
		    			visibleItems: 3
		    		}
		    	}
		    });

		});
		</script>
		<script type="text/javascript" src="`+app.Assets_public_path()+`js/jquery.flexisel.js"></script>

          	<div class="toogle">
			<h3 class="m_3">Product Details</h3>
			<p class="m_text">`+product.Get_description()+`</p>
		</div>
        </div>
       <!---->
	`
	return out
}

func Product(productList []*entity.Product, paging string)string {

	out = `
<!-- start content -->
<div class="container">
	<div class="women-product">
		<div class="w_content">
			<div class="women">
				<a href="#"><h4>Enthecwear - <span>4449 itemms</span> </h4></a>
				<ul class="w_nav">
					<li>Sort : </li>
					<li><a class="active" href="#">popular</a></li> |
					<li><a href="#">new </a></li> |
					<li><a href="#">discount</a></li> |
					<li><a href="#">price: Low High </a></li>
					<div class="clearfix"> </div>
			     	</ul>
			     	<div class="clearfix"> </div>
			</div>
		</div>
		<!-- grids_of_4 -->
		<div class="grid-product">

			`+ProductListing(productList,paging)+`

			<div class="clearfix"> </div>
		</div>
	</div>`
	return out
}

func ProductListing(productList []*entity.Product, paging string)string{
	var out,idStr string
	if productList != nil && len(productList)>0 {
		for _, i := range productList {
			idStr = strconv.FormatInt(i.Get_id(),10)
			var imgSrc string = app.No_image()
			if len(i.Get_img())>0{imgSrc = i.Get_img()[0]}
			out += `
			<div class="product-grid">
				<div class="content_box">
					<a href="single.html">
						<div class="left-grid-view grid-view-left">
			   				<img style="height:285px;" src="` + imgSrc + `" class="img-responsive watch-right" alt=""/>
				   			<div class="mask">
	                        				<div class="info">Quick View</div>
			            			</div>
						</a>
				</div>
				<h4><a href="`+app.Uri_public_product_get()+`?id=`+idStr+`">` + i.Get_title() +`</a></h4>
				<p>It is a long established fact that a reader</p>
				` + fmt.Sprintf("%.2f",i.Get_price()) + `
				</div>
			</div>`
		}
		out += `<nav aria-label="...">
				<ul class="pagination">
					`+paging+`
				</ul>
			</nav>`
	}else {
		out += `<h1>` + lang.T("no items") + `</h1>`
	}
	return out
}

