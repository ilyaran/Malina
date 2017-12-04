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
	"github.com/ilyaran/Malina/berry"
	"html/template"
	"github.com/ilyaran/Malina/app"
)



type Welcome struct {

}


func (s *Welcome)Index(malina *berry.Malina,w http.ResponseWriter){
	malina.Content = template.HTML(`

<div class="shoes-grid">
			<a href="single.html">
				<div class="wrap-in">
					<div class="wmuSlider example1 slide-grid">
				   		<div class="wmuSliderWrapper">
							<article style="position: absolute; width: 100%; opacity: 0;">
								<div class="banner-matter">
									<div class="col-md-5 banner-bag">
										<img class="img-responsive " src="`+app.Url_assets_public+`images/bag.jpg" alt=" " />
									</div>
									<div class="col-md-7 banner-off">
											<h2>FLAT 50% 0FF</h2>
											<label>FOR ALL PURCHASE <b>VALUE</b></label>
											<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et </p>
											<span class="on-get">GET NOW</span>
									</div>

									<div class="clearfix"> </div>
								</div>

							</article>
							<article style="position: absolute; width: 100%; opacity: 0;">
								<div class="banner-matter">
									<div class="col-md-5 banner-bag">
										<img class="img-responsive " src="`+app.Url_assets_public+`images/bag1.jpg" alt=" " />
									</div>
									<div class="col-md-7 banner-off">
										<h2>FLAT 50% 0FF</h2>
										<label>FOR ALL PURCHASE <b>VALUE</b></label>
										<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et </p>
										<span class="on-get">GET NOW</span>
									</div>

									<div class="clearfix"> </div>
								</div>
							</article>
							<article style="position: absolute; width: 100%; opacity: 0;">
								<div class="banner-matter">
								<div class="col-md-5 banner-bag">
									<img class="img-responsive " src="`+app.Url_assets_public+`images/bag.jpg" alt=" " />
									</div>
									<div class="col-md-7 banner-off">
										<h2>FLAT 50% 0FF</h2>
										<label>FOR ALL PURCHASE <b>VALUE</b></label>
										<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et </p>
										<span class="on-get">GET NOW</span>
									</div>

									<div class="clearfix"> </div>
								</div>

							</article>

						</div>

						<ul class="wmuSliderPagination">
							<li><a href="#" class="">0</a></li>
							<li><a href="#" class="">1</a></li>
							<li><a href="#" class="">2</a></li>
						</ul>
						<script src="`+app.Url_assets_public+`js/jquery.wmuSlider.js"></script>
						<script>
								$('.example1').wmuSlider();
						</script>
	            	</div>
	          	</div>
			</a>
	   		      <!---->
			<div class="shoes-grid-left">
				<a href="single.html">
	   		     	<div class="col-md-6 con-sed-grid">

	   		     		<div class=" elit-grid">

		   		     		<h4>consectetur  elit</h4>
		   		     		<label>FOR ALL PURCHASE VALUE</label>
							<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, </p>
							<span class="on-get">GET NOW</span>
						</div>
						<img class="img-responsive shoe-left" src="`+app.Url_assets_public+`images/sh.jpg" alt=" " />

						<div class="clearfix"> </div>

	   		     	</div>
				</a>
				<a href="single.html">
	   		     	<div class="col-md-6 con-sed-grid sed-left-top">
	   		     		<div class=" elit-grid">
		   		     		<h4>consectetur  elit</h4>
		   		     		<label>FOR ALL PURCHASE VALUE</label>
							<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, </p>
							<span class="on-get">GET NOW</span>
						</div>
						<img class="img-responsive shoe-left" src="`+app.Url_assets_public+`images/wa.jpg" alt=" " />

						<div class="clearfix"> </div>
	   		     	</div>
				</a>
			</div>
	   		<div class="products">
				<h5 class="latest-product">LATEST PRODUCTS</h5>
				<a class="view-all" href="product.html">VIEW ALL<span> </span></a>
	   		</div>
			<div class="product-left">
				<div class="col-md-4 chain-grid">
	   		    	<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/ch.jpg" alt=" " /></a>
	   		     	<span class="star"> </span>
	   		     	<div class="grid-chain-bottom">
	   		     		<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     		<div class="star-price">
	   		     			<div class="dolor-grid">
		   		     			<span class="actual">300$</span>
		   		     			<span class="reducedfrom">400$</span>
		   		     			<span class="rating">
									<input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									<label for="rating-input-1-5" class="rating-star1"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									<label for="rating-input-1-4" class="rating-star1"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									<label for="rating-input-1-3" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									        <label for="rating-input-1-2" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									        <label for="rating-input-1-1" class="rating-star"> </label>
							    </span>
	   		     			</div>
	   		     			<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     			<div class="clearfix"> </div>
						</div>
					</div>
				</div>
	   		    <div class="col-md-4 chain-grid">
					<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/ba.jpg" alt=" " /></a>
	   		     	<span class="star"> </span>
	   		     	<div class="grid-chain-bottom">
	   		     		<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     		<div class="star-price">
	   		     			<div class="dolor-grid">
		   		     			<span class="actual">300$</span>
		   		     			<span class="reducedfrom">400$</span>
		   		     			<span class="rating">
									<input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									<label for="rating-input-1-5" class="rating-star1"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									        <label for="rating-input-1-4" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									        <label for="rating-input-1-3" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									        <label for="rating-input-1-2" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									        <label for="rating-input-1-1" class="rating-star"> </label>
								</span>
							</div>
	   		     			<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     			<div class="clearfix"> </div>
						</div>
					</div>
				</div>
	   		    <div class="col-md-4 chain-grid grid-top-chain">
					<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/bo.jpg" alt=" " /></a>
	   		     	<span class="star"> </span>
					<div class="grid-chain-bottom">
	   		     		<h6><a href="single.html">Lorem ipsum dolor</a></h6>
						<div class="star-price">
	   		     			<div class="dolor-grid">
		   		     			<span class="actual">300$</span>
		   		     			<span class="reducedfrom">400$</span>
		   		     			<span class="rating">
									<input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									<label for="rating-input-1-5" class="rating-star1"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									<label for="rating-input-1-4" class="rating-star1"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									<label for="rating-input-1-3" class="rating-star"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									<label for="rating-input-1-2" class="rating-star"> </label>
									<input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									<label for="rating-input-1-1" class="rating-star"> </label>
								</span>
							</div>
	   		     			<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     			<div class="clearfix"> </div>
						</div>
					</div>
				</div>
	   		    <div class="clearfix"> </div>
			</div>
	   		<div class="products">
				<h5 class="latest-product">LATEST PRODUCTS</h5>
				<a class="view-all" href="product.html">VIEW ALL<span> </span></a>
			</div>
	   		<div class="product-left">
	   		     	<div class="col-md-4 chain-grid">
	   		     		<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/bott.jpg" alt=" " /></a>
	   		     		<span class="star"> </span>
	   		     		<div class="grid-chain-bottom">
	   		     			<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     			<div class="star-price">
	   		     				<div class="dolor-grid">
		   		     				<span class="actual">300$</span>
		   		     				<span class="reducedfrom">400$</span>
		   		     				  <span class="rating">
									        <input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									        <label for="rating-input-1-5" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									        <label for="rating-input-1-4" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									        <label for="rating-input-1-3" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									        <label for="rating-input-1-2" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									        <label for="rating-input-1-1" class="rating-star"> </label>
							    	   </span>
	   		     				</div>
	   		     				<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     				<div class="clearfix"> </div>
							</div>
	   		     		</div>
	   		     	</div>
	   		     	<div class="col-md-4 chain-grid">
	   		     		<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/bottle.jpg" alt=" " /></a>
	   		     		<span class="star"> </span>
	   		     		<div class="grid-chain-bottom">
	   		     			<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     			<div class="star-price">
	   		     				<div class="dolor-grid">
		   		     				<span class="actual">300$</span>
		   		     				<span class="reducedfrom">400$</span>
		   		     				  <span class="rating">
									        <input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									        <label for="rating-input-1-5" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									        <label for="rating-input-1-4" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									        <label for="rating-input-1-3" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									        <label for="rating-input-1-2" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									        <label for="rating-input-1-1" class="rating-star"> </label>
							    	   </span>
	   		     				</div>
	   		     				<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     				<div class="clearfix"> </div>
							</div>
	   		     		</div>
	   		     	</div>
	   		     	<div class="col-md-4 chain-grid grid-top-chain">
	   		     		<a href="single.html"><img class="img-responsive chain" src="`+app.Url_assets_public+`images/baa.jpg" alt=" " /></a>
	   		     		<span class="star"> </span>
	   		     		<div class="grid-chain-bottom">
	   		     			<h6><a href="single.html">Lorem ipsum dolor</a></h6>
	   		     			<div class="star-price">
	   		     				<div class="dolor-grid">
		   		     				<span class="actual">300$</span>
		   		     				<span class="reducedfrom">400$</span>
		   		     				  <span class="rating">
									        <input type="radio" class="rating-input" id="rating-input-1-5" name="rating-input-1">
									        <label for="rating-input-1-5" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-4" name="rating-input-1">
									        <label for="rating-input-1-4" class="rating-star1"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-3" name="rating-input-1">
									        <label for="rating-input-1-3" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-2" name="rating-input-1">
									        <label for="rating-input-1-2" class="rating-star"> </label>
									        <input type="radio" class="rating-input" id="rating-input-1-1" name="rating-input-1">
									        <label for="rating-input-1-1" class="rating-star"> </label>
							    	   </span>
	   		     				</div>
	   		     				<a class="now-get get-cart" href="#">ADD TO CART</a>
	   		     				<div class="clearfix"> </div>
							</div>
	   		     		</div>
	   		     	</div>
	   		     	 <div class="clearfix"> </div>
	   		     </div>
	   		<div class="clearfix"> </div>
		</div>


	`)



	PublicLayoutView.response(malina,w)
}
