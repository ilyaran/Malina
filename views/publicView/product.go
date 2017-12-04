package publicView

import (
	"net/http"
	"github.com/ilyaran/Malina/berry"
	"html/template"
	"strconv"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/lang"
	"github.com/ilyaran/Malina/app"
	"fmt"

	"github.com/ilyaran/Malina/libraries"
)

type Product struct {

}

func (s *Product)Index(malina *berry.Malina,category *entity.Category,w http.ResponseWriter){
	var content = `
		<div class="women-product">
		<div class=" w_content">
			<div class="women">
				<a href="#"><h4>Enthecwear - <span>4449 itemms</span> </h4></a>
				<ul class="w_nav">
					<li>Per page : </li>
					<li>
						`+PublicLayoutView.Per_page_select+`
					</li> |
					<li>Sort : </li>
					<li>
						<select id="order_by">
							<option value="5">Price Low to High</option>
							<option value="4">Price High to Low</option>
						</select>
					</li>

			     <div class="clearfix"> </div>
			     </ul>
			     <div class="clearfix"> </div>
			</div>
		</div>
		<!-- grids_of_4 -->
		<div class="grid-product" id="product_listing">

		`+ s.Listing(malina,category)+`

		<div class="clearfix"> </div>
		</div>
	</div>

		`

	malina.Content=template.HTML(content)
	PublicLayoutView.response(malina,w)
}

func (s *Product)Listing(malina *berry.Malina,category *entity.Category)string{

	if malina.List == nil {
		return `<h1>` + lang.T("no items") + `</h1>`
	}

	var out, idStr, imgSrc string
	if category !=nil{
		out += `
		<ol class="breadcrumb">
			<li><a href="`+app.Url_public_product_list+`">All products</a></li>
			`+category.PublicBreadcrumbPathView+`
			<li class="active">`+category.Title+`</li>
		</ol>`
	}

	out += `
	<div class="clearfix"> </div>
	<ul class="pagination">` + malina.Paging + `</ul>
	<div class="clearfix"> </div>`

	for _,v:=range *malina.List {
		idStr = strconv.FormatInt(v.(*entity.Product).Id, 10)
		if v.(*entity.Product).GetImg() != nil {
			imgSrc = app.Url_assets_uploads + v.(*entity.Product).Img[0]
		} else {
			imgSrc = app.Url_no_image
		}
		out += `

		<div class="  product-grid">
			<div class="content_box">
				<a href="` +app.Url_public_product_item +`?id=` +idStr + `">
			   	<div class="left-grid-view grid-view-left">
			   	   	 <img src="` + imgSrc + `" class="img-responsive watch-right" alt=""/>
				   	   	<div class="mask">
	                        <div class="info">Quick View</div>
			            </div>
				   	  </a>
				   </div>
				    <h4><a href="` +app.Url_public_product_item +`?id=` +idStr + `"> ` + v.(*entity.Product).Title + `</a></h4>
				     <p>` + libraries.CategoryLib.ListPublic[v.(*entity.Product).Category].Title + `</p>
				     Price:` + fmt.Sprintf("%.2f $", v.(*entity.Product).Price) + `
			   	</div>
                 </div>
				     `
	}

	out += `
	<div class="clearfix"> </div>
	<ul class="pagination">` + malina.Paging + `</ul>`

	return out
}



func (s *Product)Single(malina *berry.Malina,product *entity.Product,w http.ResponseWriter)string{
	var images string
	if product.GetImg()!=nil{
		for k,v := range product.Img {
			if k < 1{
				images += `
			<li>
				<a href="optionallink.html">
					<img class="etalage_thumb_image" src="`+app.Url_assets_uploads+v+`" class="img-responsive" />
					<img class="etalage_source_image" src="`+app.Url_assets_uploads+v+`" class="img-responsive" title="" />
				</a>
			</li>`
			}else {
				images += `
			<li>
				<img class="etalage_thumb_image" src="`+app.Url_assets_uploads+v+`" class="img-responsive" />
				<img class="etalage_source_image" src="`+app.Url_assets_uploads+v+`" class="img-responsive" title="" />
			</li>`
			}
		}
	}else {
		images = `<li><img class="etalage_thumb_image" src="`+app.Url_no_image+`" class="img-responsive" /></li>`
	}
	var out=`
		<link rel="stylesheet" href="`+ app.Url_assets_public+`css/etalage.css" type="text/css" media="all" />
		<script src="`+ app.Url_assets_public+`js/jquery.etalage.min.js"></script>
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

		<div class=" single_top">
	    	<div class="single_grid">`
				var category *entity.Category
	    		if product.Category > 0 {
					if v,ok:=libraries.CategoryLib.MapPublic[product.Category];ok{
						category = v
					}
				}

				if category !=nil{
					out += `
				<ol class="breadcrumb">
					<li><a href="`+app.Url_public_product_list+`">All products</a></li>
						`+category.PublicBreadcrumbPathView+`
					<li class="active"><a href="`+app.Url_public_product_list+`?category_id=`+strconv.FormatInt(category.Id,10)+`">`+category.Title+`</a></li>
				</ol>`
				}
				out+=`
				<div class="grid images_3_of_2">
					<ul id="etalage">
					` + images + `
					</ul>
					<div class="clearfix"> </div>
				</div>
				<div class="desc1 span_3_of_2">


				<h4>` + product.Title + `</h4>
				<div class="cart-b">
					<div class="left-n ">$` + fmt.Sprintf("%.2f",product.Price) + `</div>
				    <a class="now-get get-cart-in" href="#">ADD TO CART</a>
				    <div class="clearfix"></div>
				</div>
				<h6>` + fmt.Sprintf("%.f", product.Price) + ` items in stock</h6>
			   	<p>` + product.Description + `</p>
			   	<div class="share">
					<h5>Share Product :</h5>
					<ul class="share_nav">
						<li><a href="#"><img src="`+ app.Url_assets_public+`images/facebook.png" title="facebook"></a></li>
						<li><a href="#"><img src="`+ app.Url_assets_public+`images/twitter.png" title="Twiiter"></a></li>
						<li><a href="#"><img src="`+ app.Url_assets_public+`images/rss.png" title="Rss"></a></li>
						<li><a href="#"><img src="`+ app.Url_assets_public+`images/gpluse.png" title="Google+"></a></li>
				    </ul>
				</div>
			</div>
          	<div class="clearfix"> </div>
        </div>
        <ul id="flexiselDemo1">
			<li><img src="`+ app.Url_assets_public+`images/pi.jpg" /><div class="grid-flex"><a href="#">Bloch</a><p>Rs 850</p></div></li>
			<li><img src="`+ app.Url_assets_public+`images/pi1.jpg" /><div class="grid-flex"><a href="#">Capzio</a><p>Rs 850</p></div></li>
			<li><img src="`+ app.Url_assets_public+`images/pi2.jpg" /><div class="grid-flex"><a href="#">Zumba</a><p>Rs 850</p></div></li>
			<li><img src="`+ app.Url_assets_public+`images/pi3.jpg" /><div class="grid-flex"><a href="#">Bloch</a><p>Rs 850</p></div></li>
			<li><img src="`+ app.Url_assets_public+`images/pi4.jpg" /><div class="grid-flex"><a href="#">Capzio</a><p>Rs 850</p></div></li>
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
		<script type="text/javascript" src="`+ app.Url_assets_public+`js/jquery.flexisel.js"></script>

        <div class="toogle">
			<h3 class="m_3">Product Details</h3>
			<p class="m_text">Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum.</p>
		</div>
    </div>
	`


	malina.Content=template.HTML(out)
	PublicLayoutView.response(malina,w)
	return out
}
























