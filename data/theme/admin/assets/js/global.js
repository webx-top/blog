(function() {
	window.webx = {
		lang: 'zh-cn',
		staticUrl: '',
		siteUrl: '',
		route: '',
		data: {},
		pageJs: null,
		libs: {layer:['Dialog/layer/min.js']},
		msgs: {
			err:null,
			suc:null,
			code:null //-2:no permission; -1:no auth; 0:failure; 1:success
		},
		calls: [],
		include: function(file, location) {
			if (location == null) location = "head";
			if (location == "head" && typeof(webx.data["include"]) == "undefined") {
				var jsAfter = $("#js-lazyload-begin"),
					cssAfter = $("#css-lazyload-begin");
				webx.data.include = {
					before: {},
					after: {}
				};
				if (jsAfter.length) {
					webx.data.include.after.script = jsAfter;
				} else {
					var jsBefore = $("#js-lazyload-end");
					if (jsBefore.length) webx.data.include.before.script = jsBefore;
				}
				if (cssAfter.length) {
					webx.data.include.after.link = cssAfter;
				} else {
					var cssBefore = $("#css-lazyload-end");
					if (cssBefore.length) webx.data.include.before.link = cssBefore;
				}
			}
			var files = typeof(file) == "string" ? [file] : file;
			for (var i = 0; i < files.length; i++) {
				var name = files[i].replace(/^\s|\s$/g, ""),
					att = name.split('.');
				var ext = att[att.length - 1].toLowerCase(),
					isCSS = ext == "css";
				var tag = isCSS ? "link" : "script";
				var attr = isCSS ? ' type="text/css" rel="stylesheet"' : ' type="text/javascript"';
				attr += ' charset="utf-8" ';
				var link = (isCSS ? "href" : "src") + "='" + name + "'";
				if ($(tag + "[" + link + "]").length == 0) {
					var ej = $("<" + tag + attr + link + "></" + tag + ">");
					if (location == "head") {
						if (typeof(webx.data.include.after[tag]) != 'undefined') {
							webx.data.include.after[tag].after(ej);
						} else if (typeof(webx.data.include.before[tag]) != 'undefined') {
							webx.data.include.before[tag].before(ej);
						} else {
							$(location).append(ej);
						}
					} else {
						$(location).append(ej);
					}
				}
			}
		},
		defined: function(vType, key, callback) {
			if (vType != 'undefined' || key == null) {
				if (key != null) return callback();
				return;
			}
			if (typeof(key) == 'string' && typeof(webx.libs[key]) != 'undefined') key = webx.libs[key];
			webx.includes(key);
			if (callback != null) return callback();
		},
		includes: function(js){
			if (!js) return;
			switch (typeof(js)) {
			case 'string':
				webx.include(webx.staticUrl + 'js/' + js);
				return;
			case 'boolean':
				webx.include(webx.staticUrl + 'js/pages' + webx.route.replace('*', '').split(':')[0] + '.js');
				return;
			default:
				if (typeof(js.length) == 'undefined') return;
				for (var i = 0; i < js.length; i++) {
					js[i] = webx.staticUrl + 'js/' + js[i];
				}
				webx.include(js);
			}
		},
		scrollTo: function(element, time) {
			if (!time) time = 1000;
			$('html,body').animate({
				scrollTop: $(element).offset().top
			}, time);
		},
		addCalls: function(func) {
			webx.calls.push(func);
		},
		doCalls: function() {
			for (var i = 0, len = webx.calls.length; i < len; i++) {
				webx.calls[i]();
			}
			webx.calls = [];
		},
		/* 解析模板 */
		parseTmpl: function(template, data) {
			return template.replace(/\{%([\w\.]*)%\}/g, function(str, key) {
				var keys = key.split("."),
					v = data[keys.shift()];
				for (var i = 0, l = keys.length; i < l; i++) v = v[keys[i]];
				return typeof(v) !== "undefined" && v !== null ? v : "";
			});
		},
		/* 调用译文 */
		t: function(k, obj) {
			var lang;
			if (typeof(Lang) == "undefined" || typeof(Lang[k]) == "undefined") {
				lang = k;
			} else {
				lang = Lang[k];
			}
			if (obj != null) return webx.parseTmpl(lang, obj);
			return lang;
		},
		/* 插入数据到光标位置 */
		insertAtCursor: function(myField, myValue) { /* IE support */
			if (document.selection) {
				myField.focus();
				sel = document.selection.createRange();
				sel.text = myValue;
				sel.select();
			} /* MOZILLA/NETSCAPE support */
			else if (myField.selectionStart || myField.selectionStart == '0') {
				var startPos = myField.selectionStart;
				var endPos = myField.selectionEnd; /* save scrollTop before insert */
				var restoreTop = myField.scrollTop;
				myField.value = myField.value.substring(0, startPos) + myValue + myField.value.substring(endPos, myField.value.length);
				if (restoreTop > 0) myField.scrollTop = restoreTop;
				myField.focus();
				myField.selectionStart = startPos + myValue.length;
				myField.selectionEnd = startPos + myValue.length;
			} else {
				myField.value += myValue;
				myField.focus();
			}
		},
		/* 复选框全选 */
		checkedAll: function(checkbox, target) {
			if (target == null) target = 'input[type=checkbox]';
			$(target).not(':disabled').prop('checked', $(checkbox).prop('checked'));
		},
		/* 确认关闭窗口 */
		confirmClose: function(msg) {
			if (msg == null) {
				if ($('body[onbeforeunload]').length) $('body[onbeforeunload]').removeAttr('onbeforeunload');
				return;
			}
			if ($('body').attr('onbeforeunload')) return;
			if (!msg) msg = webx.t('您填写的数据没有提交，如果离开本页面这些数据将会丢失。\n确定丢弃这些内容吗？');
			$('body').attr('onbeforeunload', "return '" + msg + "';");
		},
		/* 回车键事件 */
		enterKeyEvent: function(ele, callback) {
			$(ele).unbind('keydown');
			$(ele).keydown(function(event) {
				if (event.keyCode == 13) {
					callback.call($(this));
					return false;
				}
			});
			return false;
		},
		/* 左右键翻页 */
		turningPage: function(prevPage, nextPage, isElement) {
			$(document).keyup(function(event) {
				if (event.keyCode == 37) {
					if (!isElement) {
						if (!prevPage) {
							alert(webx.t('没有了。这已经是第一页了。'));
							return;
						}
						window.location = prevPage;
					} else {
						if ($(prevPage).length < 1) {
							alert(webx.t('没有了。这已经是第一页了。'));
							return;
						}
						$(prevPage).click();
					}
				} else if (event.keyCode == 39) {
					if (!isElement) {
						if (!nextPage) {
							alert(webx.t('没有了。这已经是最末页了。'));
							return;
						}
						window.location = nextPage;
					} else {
						if ($(nextPage).length < 1) {
							alert(webx.t('没有了。这已经是最末页了。'));
							return;
						}
						$(nextPage).click();
					}
				}
			});
			$(':text,textarea').keyup(function(event) {
				event.stopPropagation();
			});
		},
		unbindKeyEvent: function() {
			$(document).unbind('keyup');
			$(':text,textarea').unbind('keyup');
		},
		/* 禁止复制 */
		disabledCopy: function(el) {
			var fn = function() {return false;};
			$(el).attr('unselectable', 'on').css({
				'-moz-user-select': '-moz-none',
				'-moz-user-select': 'none',
				'-o-user-select': 'none',
				'-khtml-user-select': 'none',
				'-webkit-user-select': 'none',
				'-ms-user-select': 'none',
				'user-select': 'none'
			}).bind('selectstart', fn).bind('contextmenu', fn).bind('dragstart', fn).bind('selectstart', fn).bind('beforecopy', fn);
		},
		/* 级联选择(使用前请确保第一个下拉框已有选中项)
		使用方法：nestedSelect(["country_id","province_id","city_id"]) */
		nestedSelect: function(ids, initVal, attrName, timeout) {
			if (typeof(ids) == 'object') {
				var obj = ids;
				if (typeof(obj.initVal) != 'undefined') initVal = obj.initVal;
				if (typeof(obj.attrName) != 'undefined') attrName = obj.attrName;
				if (typeof(obj.timeout) != 'undefined') timeout = obj.timeout;
				if (typeof(obj.ids) != 'undefined') ids = obj.ids;
				obj = null;
			}
			var id = ids[0],
				id2 = ids[1];
			if (initVal == null) initVal = '';
			if (attrName == null) attrName = 'rel';
			if (timeout == null) timeout = 5000;
			var attr = $('#' + id2).attr(attrName);
			if (!attr) return false;
			if ($('#' + id).val() == initVal) return false;
			if ($('#' + id2 + ' option:last').val() != initVal) return false;
			$('#' + id).trigger('change');
			var i = 0;
			var ptimer = window.setInterval(function() {
				i++;
				if ($('#' + id2 + ' option:last').val() != initVal || i * 200 > timeout) {
					window.clearInterval(ptimer);
					var sel = $('#' + id2 + ' option[value="' + attr + '"]');
					if (sel.length <= 0) return;
					sel.prop('selected', true);
					ids.shift();
					if (ids.length > 1) webx.nestedSelect(ids, initVal, attrName, timeout);
				}
			}, 200);
			return true;
		},
		initPage: function(js) {
			if(js==null)js=webx.pageJs;
			webx.doCalls();
			webx.includes(js);
			webx.showMsgs(true);
		},
		showMsgs:function(once){
			if(once==null)once=false;
			if(webx.msgs.err&&webx.msgs.suc){
				webx.dialog().msg('<div>'+webx.msgs.err+'</div><div>'+webx.msgs.suc+'</div>',{offset:'10px',shift:6,icon:0,time:10000});
				if (once) webx.resetMsgs();
			}else if(webx.msgs.err){
				webx.dialog().msg(webx.msgs.err,{offset:'10px',shift:6,icon:5,time:8000});
				if (once) webx.resetMsgs();
			}else if(webx.msgs.suc){
				webx.dialog().msg(webx.msgs.suc,{offset:'10px',icon:6,time:5000});
				if (once) webx.resetMsgs();
			}
		},
		resetMsgs:function(){
			webx.msgs={err:null,suc:null,code:null};
		},
		asMsgs:function(obj){
			webx.msgs.code=obj.Status;
			if (obj.Status==1) {
				webx.msgs.suc=obj.Message;
			}else{
				webx.msgs.err=obj.Message;
			}
		},
		setMsgs:function(code,msg){
			webx.msgs.code=code;
			if (code==1) {
				webx.msgs.suc=msg;
			}else{
				webx.msgs.err=msg;
			}
		},
		dialog: function() {
			var type=typeof(layer);
			if (type=='undefined') {
				window.LAYER_PATH=webx.staticUrl+'js/Dialog/layer/';
				webx.defined(type, 'layer');
				layer.config({
    				extend: ['extend/layer.ext.js','skin/moon/style.css'],
    				skin: 'layer-ext-moon'
				});
			}
			return layer;
        }
	};
})();
function T(k, obj) {
	return webx.t(k, obj);
}
function D() {
	return webx.dialog();
}