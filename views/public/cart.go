package public

import (
	"Malina/entity"
	"strconv"
	"fmt"
	"Malina/config"
)

func Order()string {
	out = `
<div class="container">
	<div class="main">
        	<div class="reservation_top">
            		<div class=" contact_right">
				<h3>Cart Details</h3>
				<div class="table-responsive" id="listing">

					<form>
						<div class="form-group">
							<label for="customer_name">Name<span class=error"">*</span></label>
							<input type="text" class="form-control" id="customer_name" placeholder="Your name">
					  	</div>
				  		<div class="form-group">
							<label for="email">Email address<span class=error"">*</span></label>
							<input type="email" class="form-control" id="email" placeholder="Email">
					  	</div>
				  		<div class="form-group">
							<label for="phone">Phone<span class=error"">*</span></label>
							<input type="phone" class="form-control" id="phone" placeholder="Phone">
					  	</div>
				  		<div class="form-group">
							<label for="city">City</label>
							<input type="text" class="form-control" id="city" placeholder="City">
					  	</div>
					  	<div class="form-group">
							<label for="address">Address</label>
							<input type="text" class="form-control" id="address" placeholder="Address">
					  	</div>
					  	<div class="form-group">
							<label for="comment">Comment</label>
							<input type="comment" class="form-control" id="comment" placeholder="Comment">
					  	</div>

				  <button type="submit" class="btn btn-default">Submit</button>
				</form>

				</div>
			</div>
		</div>
	</div>

	`
	return out
}

func Cart(CartPublicList []*entity.CartPublic)string {

	out = `
<div class="container">
	<div class="main">
        	<div class="reservation_top">
            		<div class=" contact_right">
				<h3>Cart Details</h3>
				<div class="table-responsive" id="listing">

					` +ListOfCartProducts(CartPublicList)+`

				</div>
			</div>
		</div>
	</div>

	`
	return out
}

func ListOfCartProducts(CartPublicList []*entity.CartPublic) string {
	var out, idStr, imgSrc string
	var total float64
	if len(CartPublicList) < 1{
		return `<h1>Cart is empty</h1>`
	}
	out = `

		<table class="table table-bordered">
			<thead>
				<tr>
					<th style="width: 5%">Delete</th>
					<th style="width: 5%">Logo</th>
					<th>Product</th>
					<th style="width: 5%">Price</th>
					<th style="width: 5%">Quantity</th>
					<th style="width: 5%">Total</th>
				</tr>
			</thead>
		<tbody>
	`
	for i:=0;i<len(CartPublicList);i++ {
		idStr = strconv.FormatInt(CartPublicList[i].Product_id,10)
		imgSrc = app.No_image()
		if len(CartPublicList[i].Product_img)>0{imgSrc = CartPublicList[i].Product_img[0]}
		out += `
			<tr>
				<th scope="row">
					<button data-item_id="` + idStr + `" class="btn btn-danger del_item"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
				</th>
				<td><img src="` + imgSrc + `" width='60' height='50' /></td>
				<td>` + CartPublicList[i].Product_title + `</td>
				<td>` + fmt.Sprintf("%.2f",CartPublicList[i].Product_price) + `</td>
				<td><input data-item_id="` + idStr + `" class="inlist_cart_quantity" type="number" value="` + strconv.FormatInt(CartPublicList[i].Product_quantity,10) +`"/></td>
				<td class="cart_subtotal">` + fmt.Sprintf("%.2f",CartPublicList[i].Subtotal) + `</td>
			</tr>
		`
		total += CartPublicList[i].Subtotal
	}
	return out + `
			<tr>
				<td colspan="5">
					<button class="btn btn-default" id="cart_update_button">Update Cart</button>
				</td>
				<td id="cart_total">
					` + fmt.Sprintf("%.2f",total) +`
				</td>
			</tr>
		</tbody>
	</table>

	<div class="col-md-6">
		<div class="form-group">
			<label for="coupon">Coupon</label>
			<input type="text" class="form-control" id="coupon" placeholder="Coupon Code">
		</div>
		<button type="submit" class="btn btn-default">Apply Coupon</button>
	</div>
	<div class="col-md-6">
		<div class="table-responsive">
			<table class="table table-bordered">
				<thead>
					<tr>
						<th >Subtotal</th>
						<th ></th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<th scope="row">Shipping</th>
						<td id="cart_total_shipping">
							` + fmt.Sprintf("%.2f",total) +`
						</td>
					</tr>
					<tr>
						<th scope="row">Tax</th>
						<td id="cart_tax">` + fmt.Sprintf("%.3f",app.TAX()) + `</td>
					</tr>
					<tr>
						<th scope="row">Total</th>
						<td id="cart_total_shipping_with_tax">
							` + fmt.Sprintf("%.2f",total+total*app.TAX()) +`
						</td>
					</tr>
					<tr>
						<td colspan="2">
							<a href="public/cart/order/" class="btn btn-default">Proceed to Checkout</a>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>`
}
