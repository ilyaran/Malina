/**
 * Home JavaScripts.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
    * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */

$(document).ready(function () {
    var baseUrl = $("base").attr("href");
    var dbtable = $("#dbtable").val();
    var postData = {};
    var current_page = 0;
    $("#action").change(function () {
        if ($(this).val() == "add") {
            $("#item_id_bar").hide();
            $(".noneed").html("*");
        } else{
            $("#item_id_bar").show();
            $(".noneed").html("");
        }
    });
    $("body").click(function () {
        $("#success_bar").hide(300);
    });
    // get Ajax list
    $("#get_list").click(function () {
        current_page = 0;
        getAjaxList();
    });
    // get Ajax list by page
    $("#listing").on("click", ".paging", function () {
        current_page = $(this).children("span").attr("data-page");
        getAjaxList();
    });

    function getAjaxList() {
        postData = {};
        postData = getPostData(true);
        postData.page = current_page;
        request($("#url_" + dbtable + "_list_ajax").val(),"ajax_list",postData,"html","POST");
    }
    $("#listing").on("click", ".del_item", function () {
        var item_id = $(this).attr("data-item_id");
        var uri = $("#url_" + dbtable + "_del").val();
        postData = {};
        postData[dbtable + "_id"] = item_id;
        request(uri,"delete",postData,"json","POST");
    });
    $("#listing").on("click", ".edit_item", function () {
        $("#action").val("edit");
        $("#item_id_bar").show();
        var id = $(this).attr("data-item_id");
        $("#item_id").val(id);
        $('.nav-tabs a[href="#form"]').tab('show');
        postData = {};
        postData[dbtable + "_id"] = id;
        var uri = $("#url_" + dbtable + "_get").val();
        request(uri,"get",postData,"json","POST");
    });
    $(".submitButton").click(function () {
        postData = {};
        var action = $("#action").val();
        if (action != "add") {
            postData[dbtable + "_id"] = $("#id").val();
        }
        var uri = $("#url_" + dbtable + "_" + action).val();
        if (action == "add") {
            postData = getPostData(false);
        }
        if (action == "edit") {
            postData = getPostData(false);
            postData[dbtable + "_id"] = $("#item_id").val();
        }
        if (action == "get"){
            postData[dbtable + "_id"] = $("#item_id").val();
        }
        request(uri,action,postData,"json","POST");
    });
    function request(uri,action,postData,type,method) {
        $(".error").html("");
        $.ajax({
            beforeSend: function (xhr) {
                $("#loader").show();
                xhr.setRequestHeader('X-Requested-With', 'xmlhttprequest');
            },
            method: method,
            url: uri,
            data: postData,
            dataType: type/*,
            headers: {
                'Cookie': document.cookie
            }*/
        }).done(function (msg) {
            $("#loader").hide();
            //console.log(msg);
            if(msg.Result){
                if(msg.Result.unauth){
                    window.location.reload();
                }
                if(msg.Result.error){
                    $("#error").html(msg.Result.error);
                }
            }
            if (action == "ajax_list") {
                /*if(msg.match(/(unauth)/)){
                    window.location.reload();
                }*/
                $("#listing").html(msg);
                $("[data-toggle='toggle']").bootstrapToggle('destroy');
                $("[data-toggle='toggle']").bootstrapToggle();
            }
            if (action == "inlist") {
                getAjaxList(0);
                if(msg.Status > 0){
                    setErrors("inlist",msg);
                }else {
                    ShowSuccess("");
                    getAjaxList(0);
                }
            }
            if (action == "add") {
                if(msg.Status > 0){
                    setErrors(msg);
                }else {
                    if(msg.Result.id){
                       $("#item_id").val(msg.Result.id);
                    }
                    getAjaxList(0);
                    setSuccess(msg);
                    ShowSuccess("");

                }
            }
            if (action == "get") {
                if(msg.Status > 0 ){
                    $("#item_id_error").html(msg.Result[dbtable+"_id"]);
                }else if(msg.id){
                    setEditForm(msg);
                }
            }
            if (action == "edit") {
                if(msg.Status > 0 ){
                    setErrors(msg);
                }else {
                    getAjaxList(0);
                    setSuccess(msg);
                    ShowSuccess("");
                }
            }
            if (action == "delete") {
                getAjaxList(0);
            }
        }).fail(function () {
            $("#loader").hide();
            alert("server error");
        });
    }
    function setErrors(msg) {
        if(msg.Result) {
            $.each(msg.Result, function (key, value) {
                if(key.indexOf(dbtable+"_") >= 0){
                    if ($("#inlist_error").length) $("#inlist_error").html($("#inlist_error").html()+"<p>"+key+" : "+value+"</p>");
                }else if (key.indexOf("_id") >= 0){
                    if ($("#item_id_error").length) $("#item_id_error").html(value);
                }else {
                    if ($("#" + key + "_error").length) $("#" + key + "_error").html(value);
                }
            });
        }
    }
    function setSuccess(msg) {
        if(msg.Result) {
            $.each(msg.Result, function (key, value) {
                if ($("#" + key).length) $("#" + key).html(value);
            });
        }
    }
    function ShowSuccess(msg) {
        $('.nav-tabs a[href="#home"]').tab('show');
        var success_bar = $("#success_bar");
        if ( msg == "" ){msg = "Success!";}
        success_bar.html(msg);
        success_bar.show(300);
        //$("#home").tab('show');
    }

    function getPostData(isListing) {
        if(isListing){
            postData.per_page = $("#per_page").val();
            postData.search = $("#search").val();
            postData.order_by = $("#order_by").val();
            if($("#category")) {
                if ($("#category").val() != 0) {
                    postData.category = $("#category").val();
                }
                if($("#price_min") && $("#price_max")) {
                    postData.price_min = $("#price_min").val();
                    postData.price_max = $("#price_max").val();
                }
            }
        }else {
            if (dbtable == "account") {
                postData.nick = $("#nick").val();
                postData.email = $("#email").val();
                postData.phone = $("#phone").val();
                postData.newpass = $("#newpass").val();
                postData.position = $('#position').val();
                if ($('#ban').prop('checked')) {
                    postData.ban = 1;
                }
                postData.ban_reason = $('#ban_reason').val();

            }
            if (dbtable == "position") {
                 postData.parent = $("#parent").val();
                 postData.title = $("#title").val();
                 postData.sort = $('#sort').val();
                 if ($('#enable').prop('checked')) {
                     postData.enable = 1;
                 }
            }
            if (dbtable == "permission") {
                 postData.data = $("#data").val();
                 postData.position_id = $("#position_id").val();
            }
            if (dbtable == "category") {

                postData.parent = $("#parent").val();
                postData.title = $("#title").val();
                postData.description = CKEDITOR.instances['description'].getData();
                postData.sort = $('#sort').val();
                if ($('#enable').prop('checked')) {
                    postData.enable = 1;
                }
            }
            if (dbtable == "product") {

                postData.category_id = $("#category_id").val();
                postData.title = $("#title").val();
                postData.description = CKEDITOR.instances['description'].getData();
                postData.price = $('#price').val();
                postData.price1 = $('#price1').val();
                postData.code = $('#code').val();
                if ($('#enable').prop('checked')) {
                    postData.enable = 1;
                }
            }
            if($("#image_preview")){
                var hrefs = "";
                $("#image_preview").children("img").each(function(){
                    hrefs += "|"+$(this).attr("src");
                });
                if (hrefs!=""){
                    postData.img = hrefs.substring(1);
                }
            }
        }

        return postData;
    }
    $("#submitInlistButton").on("click", function () {
        postData = {};
        var fields = [];
        if (dbtable == "permission") {
            fields = ["permission_data","del"];
        }
        if (dbtable == "position") {
            fields = ["position_title","position_sort","position_enable","del"];
        }
        if (dbtable == "account") {
            fields = ["account_nick","account_ban","account_email","account_phone","account_newpass","del"];
        }
        if (dbtable == "category") {
            fields = ["category_sort","category_title","category_enable","del"];
        }
        if (dbtable == "product") {
            fields = ["product_price","product_price1","product_title","product_enable","product_code","del"];
        }
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
        request(baseUrl+"home/" + dbtable + "/inlist/","inlist",postData,"json","POST");
    });
    function setEditForm(msg) {
        if (dbtable == "account") {
            $("#nick").val(msg.nick);
            $("#email").val(msg.email);
            $("#phone").val(msg.phone);
            $("#position").val(msg.position+0);
            $('#ban').prop('checked', msg.ban);
            $("#ban_reason").val(msg.ban_reason);
        }
        if (dbtable == "position") {
            $("#parent").val(msg.parent.id);
            $("#title").val(msg.title);
            $('#enable').prop('checked', msg.enable);
            $("#sort").val(msg.sort);
        }
        if (dbtable == "permission") {
            $("#position_id").val(msg.position.id);
            $("#data").val(msg.data);
        }
        if (dbtable == "category") {
            $("#parent").val(msg.parent);
            $("#title").val(msg.title);
            CKEDITOR.instances['description'].setData(msg.description);
            $('#enable').prop('checked', msg.enable);
            $("#sort").val(msg.sort);
        }
        if (dbtable == "product") {
            $("#category_id").val(msg.category_id);
            $("#title").val(msg.title);
            CKEDITOR.instances['description'].setData(msg.description);
            $('#enable').prop('checked', msg.enable);
            $("#price").val(msg.price);
            $("#price1").val(msg.price1);
            $("#code").val(msg.code);
        }
        clean_image_preview();
        if(msg.img){
            $.each(msg.img, function (index, value) {
                if (value!="") {
                    var image = new Image();
                    image.setAttribute("src", value);
                    image.setAttribute("style", "width:120px;");
                    image.setAttribute("class", "logo");
                    $("#image_preview").append(image);
                }
            });
        }
    }

    //******************************
    //image manipulation
    //******************************
    $("#select_image").click(function () {
        openCustomRoxy2();
    });
    function openCustomRoxy2(){
        $('#roxyCustomPanel2').dialog({modal:true, width:875,height:600});
    }

    var image_preview = $("#image_preview");

    $("#clean_images").on("click",function (event) {
        clean_image_preview();
    });

    function clean_image_preview() {
        image_preview.html("");
        $("#img_error").html("");
        postData.img = "";
    }

    //move image to first position
    image_preview.on("dragstart", "img", function(){
        moveImg(this);
    });
    image_preview.on("click", "img", function(){
        moveImg(this);
    });

    function moveImg(o){
        var obj = o;
        image_preview.children("img").each(function(){
            image_preview.prepend(this);
        });
        var detect = false;
        image_preview.children("img").each(function(){
            if (this == obj){
                detect = true;
            }else {
                if(detect){
                    image_preview.append(obj);
                    detect = false;
                }else {
                    image_preview.append(this);
                }
            }
        });

        image_preview.children("img").each(function(){
            image_preview.prepend(this);
        });
    }

});
function closeCustomRoxy2(){
    $('#roxyCustomPanel2').dialog('close');
}

/*

 // read image to file and push into postData array
 function readImage(file) {
 var reader = new FileReader();
 reader.addEventListener("load", function () {
 var image = new Image();
 image.addEventListener("load", function () {
 postData.images.push(reader.result);
 image_preview.append(this);
 });
 image.src = useBlob ? window.URL.createObjectURL(file) : reader.result;
 image.width=128;

 // If we set the variable `useBlob` to true:
 // (Data-URLs can end up being really large
 // `src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAA...........etc`
 // Blobs are usually faster and the image src will hold a shorter blob name
 // src="blob:http%3A//example.com/2a303acf-c34c-4d0a-85d4-2136eef7d723"
 if (useBlob) {
 // Free some memory for optimal performance
 window.URL.revokeObjectURL(file);
 }

 });
 reader.readAsDataURL(file);
 }
 */