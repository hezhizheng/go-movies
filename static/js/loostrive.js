jQuery(document).ready(function() {
	jQuery('.post_hover a').attr("target","_blank");
});
//导航
jQuery(document).ready(function() {
	jQuery(".topnav ul li").hover(function() {
		jQuery(this).children("ul").show();
		jQuery(this).addClass("li01")
	},
	function() {
		jQuery(this).children("ul").hide();
		jQuery(this).removeClass("li01")
	})
});
jQuery(document).ready(function(){
	jQuery("#menu-search").click(function(){
		jQuery(".menu-search").toggleClass("current_page_item");
	});
});
//tabs
jQuery(document).ready(function() {
	jQuery('#tabnav li').click(function() {
		jQuery(this).addClass("selected").siblings().removeClass();
		jQuery("#tab-content > ul").eq(jQuery('#tabnav li').index(this)).fadeIn(800).siblings().hide()
	})
});
//图片hover
jQuery(function() {
	jQuery('.thumbnail img').hover(function() {
		jQuery(this).fadeTo("fast", 0.5)
	},
	function() {
		jQuery(this).fadeTo("fast", 1)
	})
});
jQuery(document).ready(function() {
	jQuery("a[rel='external'],a[rel='external nofollow']").click(function() {
		window.open(this.href);
		return false
	})
});
jQuery(document).ready(function() {
	jQuery('.icon1,.icon2,.icon3,.icon4,.icon5,.icon6').wrapInner('<span class="hover"></span>').css('textIndent', '0').each(function() {
		jQuery('span.hover').css('opacity', 0).hover(function() {
			jQuery(this).stop().fadeTo(350, 1)
		},
		function() {
			jQuery(this).stop().fadeTo(350, 0)
		})
	})
});
/*搜索*/
jQuery(document).ready(function(){

   $("#searchform .search-s").each(function(){
     var thisVal=$(this).val();
     //判断文本框的值是否为空，有值的情况就隐藏提示语，没有值就显示
     if(thisVal!=""){
       $(this).siblings("span").hide();
      }else{
       $(this).siblings("span").show();
      }
     //聚焦型输入框验证 
     $(this).focus(function(){
       $(this).siblings("span").hide();
      }).blur(function(){
        var val=$(this).val();
        if(val!=""){
         $(this).siblings("span").hide();
        }else{
         $(this).siblings("span").show();
        } 
      });
    })
 })
//////////////gototop//////////////////
function b(){
	h = $(window).height();
	t = $(document).scrollTop();
	if(t > h){
		$('#gotop').show();
	}else{
		$('#gotop').hide();
	}
}
$(document).ready(function(e) {
	b();
	$('#gotop').click(function(){
		$(document).scrollTop(0);	
	})
});

$(window).scroll(function(e){
	b();		
})


////////////////////////////////
$(function() {
	var sWidth = $("#focus").width(); //获取焦点图的宽度（显示面积）
	var len = $("#focus ul li").length; //获取焦点图个数
	var index = 0;
	var picTimer;
	
	//以下代码添加数字按钮和按钮后的半透明条，还有上一页、下一页两个按钮
	var btn = "<div class='button'>";
	for(var i=0; i < len; i++) {
		btn += "<span></span>";
	}
	btn += "</div><div class='preNext pre'></div><div class='preNext next'></div>";
	$("#focus").append(btn);
	$("#focus .btnBg").css("opacity",0.5);

	//为小按钮添加鼠标滑入事件，以显示相应的内容
	$("#focus .button span").css("opacity",0.4).mouseenter(function() {
		index = $("#focus .button span").index(this);
		showPics(index);
	}).eq(0).trigger("mouseenter");

	//上一页、下一页按钮透明度处理
	$("#focus .preNext").css("opacity",0.2).hover(function() {
		$(this).stop(true,false).animate({"opacity":"0.5"},300);
	},function() {
		$(this).stop(true,false).animate({"opacity":"0.2"},300);
	});

	//上一页按钮
	$("#focus .pre").click(function() {
		index -= 1;
		if(index == -1) {index = len - 1;}
		showPics(index);
	});

	//下一页按钮
	$("#focus .next").click(function() {
		index += 1;
		if(index == len) {index = 0;}
		showPics(index);
	});

	//本例为左右滚动，即所有li元素都是在同一排向左浮动，所以这里需要计算出外围ul元素的宽度
	$("#focus ul").css("width",sWidth * (len));
	
	//鼠标滑上焦点图时停止自动播放，滑出时开始自动播放
	$("#focus").hover(function() {
		clearInterval(picTimer);
	},function() {
		picTimer = setInterval(function() {
			showPics(index);
			index++;
			if(index == len) {index = 0;}
		},4000); //此4000代表自动播放的间隔，单位：毫秒
	}).trigger("mouseleave");
	
	//显示图片函数，根据接收的index值显示相应的内容
	function showPics(index) { //普通切换
		var nowLeft = -index*sWidth; //根据index值计算ul元素的left值
		$("#focus ul").stop(true,false).animate({"left":nowLeft},300); //通过animate()调整ul元素滚动到计算出的position
		//$("#focus .btn span").removeClass("on").eq(index).addClass("on"); //为当前的按钮切换到选中的效果
		$("#focus .btn span").stop(true,false).animate({"opacity":"0.4"},300).eq(index).stop(true,false).animate({"opacity":"1"},300); //为当前的按钮切换到选中的效果
	}
});
//////////////////////////////////////////////////////////////
(function($){

    $.fn.extend({ 

        hoverZoom: function(settings) {
 
            var defaults = {
                overlay: true,
                overlayColor: '#000',
                overlayOpacity: 0.5,
                zoom: 25,
                speed: 300
            };
             
            var settings = $.extend(defaults, settings);
         
            return this.each(function() {
            
                var s = settings;
                var hz = $(this);
                var image = $('img', hz);

                image.load(function() {
                    
                    if(s.overlay === true) {
                        $(this).parent().append('<div class="zoomOverlay" />');
                        $(this).parent().find('.zoomOverlay').css({
                            opacity:0, 
                            display: 'block', 
                            backgroundColor: s.overlayColor
                        }); 
                    }
                
                    var width = $(this).width();
                    var height = $(this).height();
                
                    $(this).fadeIn(1000, function() {
                        $(this).parent().css('background-image', 'none');
                        hz.hover(function() {
                            $('img', this).stop().animate({
                                height: height + s.zoom,
                                marginLeft: -(s.zoom),
                                marginTop: -(s.zoom)
                            }, s.speed);
                            if(s.overlay === true) {
                                $(this).parent().find('.zoomOverlay').stop().animate({
                                    opacity: s.overlayOpacity
                                }, s.speed);
                            }
                        }, function() {
                            $('img', this).stop().animate({
                                height: height,
                                marginLeft: 0,
                                marginTop: 0
                            }, s.speed);
                            if(s.overlay === true) {
                                $(this).parent().find('.zoomOverlay').stop().animate({
                                    opacity: 0
                                }, s.speed);
                            }
                        });
                    });
                });    
            });
        }
    });
})(jQuery);

 $(function() {$('.zoom').hoverZoom({zoom:0});});

///////////////////文章图片自适应/////////////////////////////
$(function(){
  $("#post_content img").addClass("img-responsive");
});


////////////////边栏滚动//////////////////////
SidebarFollow = function() {

	this.config = {
		element: null, // 处理的节点
		distanceToTop: 0 // 节点上边到页面顶部的距离
	};

	this.cache = {
		originalToTop: 0, // 原本到页面顶部的距离
		prevElement: null, // 上一个节点
		parentToTop: 0, // 父节点的上边到顶部距离
		placeholder: jQuery('<div>') // 占位节点
	}
};

SidebarFollow.prototype = {

	init: function(config) {
		this.config = config || this.config;
		var _self = this;
		var element = jQuery(_self.config.element);

		// 如果没有找到节点, 不进行处理
		if(element.length <= 0) {
			return;
		}

		// 获取上一个节点
		var prevElement = element.prev();
		while(prevElement.is(':hidden')) {
			prevElement = prevElement.prev();
			if(prevElement.length <= 0) {
				break;
			}
		}
		_self.cache.prevElement = prevElement;

		// 计算父节点的上边到顶部距离
		var parent = element.parent();
		var parentToTop = parent.offset().top;
		var parentBorderTop = parent.css('border-top');
		var parentPaddingTop = parent.css('padding-top');
		_self.cache.parentToTop = parentToTop + parentBorderTop + parentPaddingTop;

		// 滚动屏幕
		jQuery(window).scroll(function() {
			_self._scrollScreen({element:element, _self:_self});
		});

		// 改变屏幕尺寸
		jQuery(window).resize(function() {
			_self._scrollScreen({element:element, _self:_self});
		});
	},

	/**
	 * 修改节点位置
	 */
	_scrollScreen: function(args) {
		var _self = args._self;
		var element = args.element;
		var prevElement = _self.cache.prevElement;

		// 获得到顶部的距离
		var toTop = _self.config.distanceToTop;

		// 如果 body 有 top 属性, 消除这些位移
		var bodyToTop = parseInt(jQuery('body').css('top'), 10);
		if(!isNaN(bodyToTop)) {
			toTop += bodyToTop;
		}

		// 获得到顶部的绝对距离
		var elementToTop = element.offset().top - toTop;

		// 如果存在上一个节点, 获得到上一个节点的距离; 否则计算到父节点顶部的距离
		var referenceToTop = 0;
		if(prevElement && prevElement.length === 1) {
			referenceToTop = prevElement.offset().top + prevElement.outerHeight();
		} else {
			referenceToTop = _self.cache.parentToTop - toTop;
		}

		// 当节点进入跟随区域, 跟随滚动
		if(jQuery(document).scrollTop() > elementToTop) {
			// 添加占位节点
			var elementHeight = element.outerHeight();
			_self.cache.placeholder.css('height', elementHeight).insertBefore(element);
			// 记录原位置
			_self.cache.originalToTop = elementToTop;
			// 修改样式
			element.css({
				top: toTop + 'px',
				position: 'fixed'
			});

		// 否则回到原位
		} else if(_self.cache.originalToTop > elementToTop || referenceToTop > elementToTop) {
			// 删除占位节点
			_self.cache.placeholder.remove();
			// 修改样式
			element.css({
				position: 'static'
			});
		}
	}
};
/* <![CDATA[ */
(new SidebarFollow()).init({
	element: jQuery('#sidebar-follow'),
	distanceToTop: 50
});
/* ]]> */

//////////////////评论///////////////////////////////////////////
	$('.comt-addsmilies').click(function(){
		$('.comt-smilies').toggle();
	})
	
	$('.comt-smilies a').click(function(){
		$(this).parent().hide();
	})
    function grin(tag) {
    	var myField;
    	tag = ' ' + tag + ' ';
        if (document.getElementById('comment') && document.getElementById('comment').type == 'textarea') {
    		myField = document.getElementById('comment');
    	} else {
    		return false;
    	}
    	if (document.selection) {
    		myField.focus();
    		sel = document.selection.createRange();
    		sel.text = tag;
    		myField.focus();
    	}
    	else if (myField.selectionStart || myField.selectionStart == '0') {
    		var startPos = myField.selectionStart;
    		var endPos = myField.selectionEnd;
    		var cursorPos = endPos;
    		myField.value = myField.value.substring(0, startPos)
    					  + tag
    					  + myField.value.substring(endPos, myField.value.length);
    		cursorPos += tag.length;
    		myField.focus();
    		myField.selectionStart = cursorPos;
    		myField.selectionEnd = cursorPos;
    	}
    	else {
    		myField.value += tag;
    		myField.focus();
    	}
    }
///////内容字数限制
jQuery.fn.myWords=function(options){
	//初始化
	// alert("a");
	var defaults={
		obj_opts:"textarea",
		obj_Maxnum:400,
		obj_Lnum:".comt-num"
	}
	var opts=$.extend(defaults,options);
	return this.each(function(){
        // 找到相应对象
		var _this=$(this).find(opts.obj_opts);
		var num=parseInt(opts.obj_Maxnum/2);
		var _obj_Lnum=$(this).find(opts.obj_Lnum);
		$(_obj_Lnum).find("em").text(num);
		if(_this.val()!=""){
			//如果文本框的值不为空，防止刷新浏览器之后 对文本框里面文字个数判断失误
			var len= _this.val().replace(/[^\x00-\xff]/g, "**").length;//将两个字母转换为一个汉字
			var _num=num-parseInt(len/2);//parseInt这个方法 就是len/2转换为整数
			html="还能输入"+"<em>"+_num+"</em>"+"字";
			$(_obj_Lnum).html(html);
		}
		_this.focus(function(){
			var html;
			$(this).keyup(function(){
				//键盘输入
				var lend= $(this).val().replace(/[^\x00-\xff]/g, "**").length;//将两个字母转换为一个汉字
				// alert(len);
				var _num=num-parseInt(lend/2);//parseInt这个方法 就是len/2转换为整数
				html="还能输入"+"<em>"+_num+"</em>"+"字";
				$(_obj_Lnum).html(html);
				if(lend>opts.obj_Maxnum){
					html="已经超出"+"<em>"+(-_num)+"</em>"+"字";
					$(_obj_Lnum).html(html);
					$(_obj_Lnum).find("em").css("color","#C30");									
				}
				else{					
					// 移除css样式
					$(obj_Lnum).find("em").removeAttr("style");
				}
			});
		});
		
	});
}
 $(function(){
   //插件
   $(".comt-box").myWords({   //输入框字数
        obj_opts:"textarea",
        obj_Maxnum:480,//要是只能输入140个字  那这里就是280
        obj_Lnum:".comt-num"
    });
})
