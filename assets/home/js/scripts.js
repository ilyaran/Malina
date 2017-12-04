$(document).ready(function () {
    var baseUrl = $("base").attr("href");
    var Url_assets_uploads = $("#Url_assets_uploads").val();
    var Url_no_image = $("#Url_no_image").val();
    var current_controller = $("#controller").val();
    var objects_list = {
        "product":new Product(),
        "parameter":new Parameter(),
        "category":new Category(),
        "news":new News(),
        "permission":new Permission()
    };


    function Crud (d,c,t) {
        var department=d;
        var controller=c;
        var tabname=t;


        this.req = function (uri,action,data,type,method){
            return $.ajax({
                beforeSend: function (xhr) {
                    $("#loader").show();
                    xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
                },
                method: method,
                url: uri,
                data: data,
                dataType: type
            });
        };
        this.run = function (action,descriptor,data) {

            this.req(baseUrl+department+"/"+controller+"/"+action+"/",descriptor,data,"json","POST")
                .done(function (msg) {

                    $("#loader").hide();
                    objects_list[tabname].runResult(msg,descriptor);

                }).fail(function (data, textStatus, xhr) {
                $("#loader").hide();
                setError(data.responseJSON);
                objects_list[tabname].setStatusError(data.responseJSON);
            });
        };

        this.getList = function (data) {

            if($("#per_page")) data.per_page = $("#per_page").val();
            if($("#search") && $("#search").val()!="") data.search = $("#search").val();
            if($("#order_by")) data.order_by = $("#order_by").val();

            this.req(baseUrl+department+"/"+controller+"/list/","list",data,"html","GET")
                .done(function (msg) {
                    $("#loader").hide();
                    $("#"+controller+"_listing").html(msg);
                    $("[data-toggle='toggle']")
                        .bootstrapToggle('destroy')
                        .bootstrapToggle();
                    objects_list[tabname].getListResult(msg);

                }).fail(function (data, textStatus, xhr) {
                $("#loader").hide();
                objects_list[tabname].setListStatusError(data, textStatus, xhr);

            });
        };
    }
    function setError(msg){

        $.each(msg, function (key, value) {
            if ($("#"+key).length){
                $("#"+key).after('<p class="error" style="color:red;">'+value+'<p>');
            }else {
                alert(key+" : "+value);
            }
        });

    }
    $(".form_send").on("click",function () {
        $(".error").html("");
        objects_list[current_controller].run_crud();
    });
    $(".listing_button").on("click",function () {
        objects_list[current_controller].getList(0);
    });
    $("body").on("click",".page_"+current_controller,function () {
        var current_page = $(this).children("span").attr("data-page");
        objects_list[current_controller].getList(current_page);
    }).on("click",".edit_item_"+current_controller,function () {
        var id = $(this).attr("data-item_id");
        $("#action").val("get");
        $("#item_id").val(id);
        objects_list[current_controller].run_crud();
    }).on("click",".del_item_"+current_controller,function () {
        var id = $(this).attr("data-item_id");
        $("#action").val("del");
        $("#item_id").val(id);

        new Alert_dialog("Cancel","Ok").run("Alert","Are you sure to delete item with ID:"+id+"?");
    });

    $("#send_inlist_"+current_controller).on("click", function () {
        var fields = [];
        $(".inlist_fields").each(function (index, elem) {
            fields.push($(this).val()+"");
        });
        fields.push(current_controller+"_del_item");
        console.log(fields);
        var data={};
        $.each(fields, function (index, value) {
            data[value + "_inlist"] = [];
            $("."+ value+"_inlist").each(function () {
                var item_id = $(this).attr("data-item_id");
                var item_value = "";
                if ($(this).is(':checkbox')) {
                    if ($(this).prop('checked')) {
                        item_value = true;
                    } else item_value = false;
                } else item_value = $(this).val();
                data[value + "_inlist"].push(item_id + "|" + item_value);
            });
        });
        console.log(data);
        inlist_request(data);
    });


    // alert dialog
    $("#modal_button2").on("click",function () {
        objects_list[current_controller].modal_button2_click();
    });
    function Alert_dialog (title_button1,title_button2,class_button1,class_button2){
        var modal = $("#modal_dialog");
        var button1 = $("#modal_button1");
        var button2 = $("#modal_button2");
        button1.text(title_button1);
        button2.text(title_button2);
        if(class_button1 !== null || class_button1 !== undefined)button1.addClass(class_button1);
        if(class_button2 !== null || class_button2 !== undefined)button2.addClass(class_button2);

        this.run = function(title,body){
            this._title = title;
            this._body = body;
            $("#modal_title").text(title);
            $("#modal_body").text(body);

            $('#modal_dialog').modal("show");
        }
    }

    // tree components

    // tree elements collapse component
    if($(".button_collapse")){
        $("body").on("click",".button_collapse",function () {
            var id_p = $(this).attr("data-item_id");
            var parent_p = $(this).attr("data-item_parent");
            var level_p = $(this).attr("data-item_level");

            $(".tree_row").each(function () {
                var id = $(this).attr("data-item_id");
                var parent = $(this).attr("data-item_parent");
                var level = $(this).attr("data-item_level");
                if(parent==id_p){
                    if($(this).is(":visible")){
                        $(this).hide(300);
                        var cl = function (p) {
                            $(".tree_row").each(function () {
                                var parent1 = $(this).attr("data-item_parent");
                                var id1 = $(this).attr("data-item_id");
                                if(parent1==p){
                                    $(this).hide(300);
                                    cl(id1);
                                }
                            });
                        };
                        cl(id_p);

                    }else {
                        $(this).show(300);
                    }
                }
            });

        });
    }
    // end tree elements collapse component


    function Product () {
        this.crud = new Crud("home","product","product");
        this.data={};
        this.current_page=0;
        this.run_crud = function () {
            var action  = $("#action").val();
            if (action == "add" || action == "edit") {
                if($("#image_preview")){
                    var hrefs = "";
                    $("#image_preview").children("img").each(function(){
                        hrefs += "|"+$(this).attr("src");
                    });
                    if (hrefs!=""){
                        this.data.img = hrefs.substring(1);
                    }
                }
                if($("#related_products_preview")){
                    var ids = "";
                    $("#related_products_preview").children("img").each(function(){
                        ids += ","+$(this).attr("data-product_id");
                    });
                    if (ids!=""){
                        this.data.related_products = ids.substring(1);
                    }
                }

                this.data.category_id = $("#category_id").val();
                this.data.title = $("#title").val();
                this.data.description = CKEDITOR.instances['description'].getData();
                this.data.short_description = $('#short_description').val();
                this.data.price = $('#price').val();
                this.data.price1 = $('#price1').val();
                this.data.code = $('#code').val();
                if ($('#enable').prop('checked')) {
                    this.data.enable = 1;
                }
            }
            this.crud.run(action,action,this.data);
        };
        this.setStatusError = function (data) {

        };
        this.runResult = function (result,descriptor) {
            console.log(result);
            console.log(descriptor);
            if (descriptor=="get_related_product"){
                //related_products_preview
                var image = new Image();
                if(result.item.img!=null){
                    image.setAttribute("src",Url_assets_uploads+result.item.img[0]);
                }else{
                    image.setAttribute("src",Url_no_image);
                }
                image.setAttribute("style","width:120px;");
                image.setAttribute("data-product_id",result.item.id);
                image.setAttribute("class","logo");
                $("#related_products_preview").prepend(image);
            }

            if (descriptor=="add"){
                $("#items_table-tab").tab('show');
                this.getList(this.current_page);
            }
            if (descriptor=="get_product"){
                if(result!==null) {
                    $("#category_id").val(result.item.category);
                    $("#title").val(result.item.title);
                    CKEDITOR.instances['description'].setData(result.item.description);
                    $('#enable').prop('checked', result.item.enable);
                    $("#price").val(result.item.price);
                    $("#price1").val(result.item.price1);
                    $("#code").val(result.item.code);
                }
                $("#action").val("edit");
                $("#item_id_bar").show();
            }
        };
        this.getList = function (page) {
            this.data={};
            this.data.category_id = $("#category_id").val()*1;
            this.data.price_max = $('#price_max').val()*1.00;
            this.data.price_min = $('#price_min').val()*1.00;
            this.data.page = page;
            this.current_page=page;
            this.crud.getList(this.data);
        };
        this.modal_button2_click = function () {
            $('#modal_dialog').modal('hide');
            this.run_crud();
        };


        $('body').on('click',".related_product",function () {
            var id = $(this).attr("data-product_id");
            objects_list.product.crud.run("get","get_related_product",{id:id});

            $("#item_form-tab").tab('show');
        });
        $('.parameters_list').change(function() {
            var ids = "";
            var parent = $(this).attr("data-parent_id");
            var idcurr=$(this).attr("data-parameter_id");
            var checked = $(this).prop('checked');
            $(".parameter_parent_"+parent).each(function () {
                if($(this).attr("data-parameter_id")!=idcurr){
                    if (checked) {
                        $(this).bootstrapToggle('off');
                    }
                }
            });
            $('.parameters_list').each(function () {
                if ($(this).prop('checked')) {
                    ids += ","+$(this).attr("data-parameter_id");
                }
            });
            if (ids!=""){
                objects_list.product.data.parameters = ids.substring(1);
            }
        });


    }


















    // end tree components

    // *********************************************
    // ***************** objects *******************
    // *********************************************

    function Permission () {
        this.crud = new Crud();
        this.run_crud = function () {
            var data={};
            var action  = $("#action").val();
            if (action == "add" || action == "edit") {
                data.title = $("#title").val();
                data.description = $("#description").val();
                if ($('#enable').prop('checked')) {
                    data.enable = 1;
                }
            }
            this.crud.run(action,data);
        };
        this.runResult = function (result) {
            if(result!==null) {
                $("#title").val(result.item.title);
                $("#description").val(result.item.description);
                $('#enable').prop('checked', result.item.enable);
            }
            $("#action").val("edit");
            $("#item_id_bar").show();
        };
        this.getList = function (page) {
            var data={};
            data.page = page;
            this.crud.getList(data);
        };
        this.modal_button2_click = function () {
            $('#modal_dialog').modal('hide');
            this.run_crud();
        };
    }

    function News () {
        this.crud = new Crud();
        this.run_crud = function () {
            var data={};
            var action  = $("#action").val();
            if (action == "add" || action == "edit") {
                data.category_id = $("#category_id").val();
                data.title = $("#title").val();
                data.description = CKEDITOR.instances['description'].getData();
                if ($('#enable').prop('checked')) {
                    data.enable = 1;
                }
            }
            this.crud.run(action,data);
        };
        this.runResult = function (result) {
            if(result!==null) {
                $("#category_id").val(result.item.category);
                $("#title").val(result.item.title);
                CKEDITOR.instances['description'].setData(result.item.description);
                $('#enable').prop('checked', result.item.enable);
            }
            $("#action").val("edit");
            $("#item_id_bar").show();
        };
        this.getList = function (page) {
            var data={};
            data.category_id = $("#category_id").val()*1;
            data.page = page;
            this.crud.getList(data);
        };
        this.modal_button2_click = function () {
            $('#modal_dialog').modal('hide');
            this.run_crud();
        };
    }

    function Parameter () {
        this.crud = new Crud();
        this.action = "";
        this.run_crud = function () {
            var data={};
            this.action  = $("#action").val();
            if (this.action == "add" || this.action == "edit") {
                data.parent = $("#parent").val();
                data.title = $("#title").val();
                data.sort = $("#sort").val();
                data.description = CKEDITOR.instances['description'].getData();
                if ($('#enable').prop('checked')) {
                    data.enable = 1;
                }
            }
            this.crud.run(this.action,data);
        };
        this.runResult = function (result) {
            if(this.action=="del"){
                $("#parent").html(result.select_options);
                $("#action").val("add");
                $("#item_id_bar").hide();
            }
            if(this.action=="add" || this.action=="edit"){
                $("#parent").html(result.select_options);
                $("#parent").val(result.item.parent);
            }
            if(this.action=="get"){
                $("#parent").val(result.item.parent);
                $("#title").val(result.item.title);
                $("#sort").val(result.item.sort);
                CKEDITOR.instances['description'].setData(result.item.description);
                $('#enable').prop('checked', result.item.enable);
            }
        };
        this.getList = function (page) {
            var data={};
            data.parent = $("#parent").val();
            data.page = page;
            this.crud.getList(data);
        };
        this.modal_button2_click = function () {
            $('#modal_dialog').modal('hide');
            this.run_crud();
        };

    }

    function Category () {
        this.crud = new Crud();
        this.action = "";
        this.run_crud = function () {
            var data={};
            this.action  = $("#action").val();
            if (this.action == "add" || this.action == "edit") {
                data.parent = $("#parent").val()*1;
                data.title = $("#title").val();
                data.sort = $("#sort").val();
                data.description = CKEDITOR.instances['description'].getData();
                if ($('#enable').prop('checked')) {
                    data.enable = 1;
                }
            }
            this.crud.run(this.action,data);
        };
        this.runResult = function (result) {
            if(this.action=="del"){
                $("#parent").html(result.select_options);
                $("#action").val("add");
                $("#item_id_bar").hide();
            }
            if(this.action=="add" || this.action=="edit"){
                $("#parent").html(result.select_options);
                $("#parent").val(result.item.parent);
            }
            if(this.action=="get"){
                $("#parent").val(result.item.parent);
                $("#title").val(result.item.title);
                $("#sort").val(result.item.sort);
                CKEDITOR.instances['description'].setData(result.item.description);
                $('#enable').prop('checked', result.item.enable);
            }
        };
        this.getList = function (page) {
            var data={};
            data.parent = $("#parent").val()*1;
            data.page = page;
            this.crud.getList(data);
        };
        this.modal_button2_click = function () {
            $('#modal_dialog').modal('hide');
            this.run_crud();
        };

    }




    // *********************************************
    // ***************** end objects ***************
    // *********************************************

    function inlist_request (data) {
        objects_list[current_controller].crud.req(baseUrl + "home/" + current_controller + "/inlist/","inlist",data,"json","POST")
            .done(function (msg) {
                $("#loader").hide();
                if (msg.Status == 200) {
                    var p = $("#page").val();
                    objects_list[current_controller].getList(p);
                } else {
                    objects_list[current_controller].crud.setError(msg);
                }
            }).fail(function () {
            $("#loader").hide();
            alert("server error");
        });
    }
    //******************************
    //image manipulation
    //******************************
    /*function clean_image_preview(selector) {
        $(selector).html("");
        $("#error_"+selector).html("");
    }*/
    function ImgM(image_preview_selector){
        var image_preview = $("#"+image_preview_selector);
        $("#clean_"+image_preview_selector).on("click",function (event) {
            $("#"+image_preview_selector).html("");
            $("#error_"+image_preview_selector).html("");
        });
        $("#select_image").click(function () {
            openCustomRoxy2();
        });
        function openCustomRoxy2(){
            $('#roxyCustomPanel2').dialog({modal:true, width:875,height:600});
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
    }
    new ImgM("image_preview");
    new ImgM("related_products_preview");

});

function closeCustomRoxy2(){
    $('#roxyCustomPanel2').dialog('close');
}



















