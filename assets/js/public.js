/**
 * Public JavaScripts.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on    JQuery
 * @email      	il.aranov@gmail.com
 * @link
    * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Aranzhanovich Toxanbayev)
 */
$(document).ready(function () {
    var BaseUrl = $("base").attr("href");
    var postData = {};
    postData.cart_action = 3;
    request(BaseUrl+"public/cart/crud/","cart_json_list",postData,"json","POST");

    function request(uri,action,postData,type,method) {
        $(".error").html("");
        $.ajax({
            beforeSend: function (xhr) {
                //$("#loader").show();
                xhr.setRequestHeader('X-Requested-With', 'xmlhttprequest');
            },
            method: method,
            url: uri,
            data: postData,
            dataType: type
        }).done(function (msg) {
            //$("#loader").hide();
            console.log(msg);
            if (action=="cart_update"){
                $("#listing").html(msg);
            }
            if (action=="add_to_cart"){
                if(msg.Status > 0 ){

                }else {
                    if(msg.Result) {
                        setCartContent(msg.Result);
                    }
                }
            }
            if (action=="cart_json_list"){
                if(msg.Status > 0 ){

                }else {
                    if(msg.Result) {
                        setCartContent(msg.Result);
                    }
                }
            }
            if (action=="delete_from_cart"){
                if(msg.Status > 0 ){

                }else {
                    if(msg.Result) {
                        setCartContent(msg.Result);
                        getAjaxList();
                    }
                }
            }
        }).fail(function () {
            //$("#loader").hide();
            alert("server error");
        });
    }
    function getAjaxList() {
        request(BaseUrl+"public/cart/ajax_list/","cart_update",postData,"html","POST");
    }

    function setCartContent(cart_deatails_array) {
        var total = 0.0;
        var quantity = 0;
        $.each(cart_deatails_array, function (key, value) {
            total += value.subtotal;
            quantity += value.product_quantity;
        });
        $("#cart_content").html(quantity + " total: $"+total);
    }

    $("#listing").on("click","#cart_update_button",function () {
        postData = {};
        var fields = ["cart_quantity"];
        $.each(fields, function (index, value) {
            postData[value+"_inlist"]=[];
            $(".inlist_"+value).each(function(){
                var item_id = $(this).attr("data-item_id");
                var item_value = "";
                if ($(this).is(':checkbox')){
                    if ($(this).prop('checked')) {
                        item_value = true;
                    }else item_value = false;
                }else item_value = $(this).val();
                postData[value+"_inlist"].push(item_id+"|"+item_value);
            });
        });
        getAjaxList();
    });
    $("#listing").on("click", ".del_item", function () {
        postData = {};
        postData.cart_action = 2;
        postData.product_id = $(this).attr("data-item_id");
        request(BaseUrl+"public/cart/crud/","delete_from_cart",postData,"json","POST");
    });

    $('.add-to-cart').click(function () {
        postData = {};
        postData.cart_action = 1;
        postData.product_id = $(this).attr("data-item_id");
        request(BaseUrl+"public/cart/crud/","add_to_cart",postData,"json","POST");
    });

    function flyToElement(flyer, flyingTo) {
        var $func = $(this);
        var divider = 3;
        var flyerClone = $(flyer).clone();
        $(flyerClone).css({position: 'absolute', top: $(flyer).offset().top + "px", left: $(flyer).offset().left + "px", opacity: 1, 'z-index': 1000});
        $('body').append($(flyerClone));
        var gotoX = $(flyingTo).offset().left + ($(flyingTo).width() / 2) - ($(flyer).width()/divider)/2;
        var gotoY = $(flyingTo).offset().top + ($(flyingTo).height() / 2) - ($(flyer).height()/divider)/2;

        $(flyerClone).animate({
                opacity: 0.4,
                left: gotoX,
                top: gotoY,
                width: $(flyer).width()/divider,
                height: $(flyer).height()/divider
            }, 700,
            function () {
                $(flyingTo).fadeOut('fast', function () {
                    $(flyingTo).fadeIn('fast', function () {
                        $(flyerClone).fadeOut('fast', function () {
                            $(flyerClone).remove();
                        });
                    });
                });
            });
    }
});
