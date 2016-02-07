(function() {
	window.webx = {
		lang: 'zh-cn',
		staticUrl: '',
		siteUrl: '',
		appUrl: '',
		appName: '',
	    controllerName: '',
	 	actionName: '',
		data: {},
		pageJs: null,
		libs: {
			layer: ['dialog/layer/min.js'],
			noty: ['dialog/noty/min.js'],
			validate: ['validate/min.js'],
			table:['dataTables/min.js']
		},
		msgs: {
			err: null,
			suc: null,
			code: null //-2:no permission; -1:no auth; 0:failure; 1:success
		},
		calls: [],
		getLang: function(){
			if(webx.data.lang!=null)return webx.data.lang;
			var part=webx.lang.split('-');
			if(part.length>1){
				part[1]=part[1].toUpperCase();
				webx.data.lang=part.join('-');
			}else{
				webx.data.lang=lang;
			}
			return webx.data.lang;
		},
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
			$.ajaxSetup({cache: true});
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
				if ($(tag + "[" + link + "]").length > 0) continue;
				var ej = $("<" + tag + attr + link + "></" + tag + ">");
				if (location == "head") {
					if (typeof(webx.data.include.after[tag]) != 'undefined') {
						webx.data.include.after[tag].after(ej);
						continue;
					} else if (typeof(webx.data.include.before[tag]) != 'undefined') {
						webx.data.include.before[tag].before(ej);
						continue;
					}
				}
				$(location).append(ej);
			}
			$.ajaxSetup({cache: false});
		},
		defined: function(vType, key, callback) {
			if (vType != 'undefined' || key == null) {
				if (key != null && callback != null) return callback();
				return;
			}
			if (typeof(key) == 'string' && typeof(webx.libs[key]) != 'undefined') key = webx.libs[key];
			webx.includes(key);
			if (callback != null) return callback();
		},
		includes: function(js) {
			if (!js) return;
			switch (typeof(js)) {
			case 'string':
				webx.include(webx.staticUrl + 'js/' + js);
				return;
			case 'boolean':
				webx.include(webx.staticUrl + 'js/pages/' + webx.appName + '/' + webx.controllerName + '/' + webx.actionName + '.js');
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
		/** 
		 * 级联选择(使用前请确保第一个下拉框已有选中项)
		 * 使用方法：nestedSelect(["country_id","province_id","city_id"]) 
		 **/
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
			if (js == null) js = webx.pageJs;
			webx.doCalls();
			webx.includes(js);
			webx.showMsgs(true);
			webx.autoValidateForm();
		},
		autoValidateForm: function() {
			$('form[data-validate="true"]').each(function() {
				var sucFn = $(this).attr('data-validate-callback'),
					option = $(this).attr('data-validate-option');
				if (sucFn == undefined) sucFn = null;
				if (option == undefined) option = null;
				if (option) {
					option = $.parseJSON(option);
					var t = typeof(option.validate);
					if (t != 'undefined' && option.validate && t != 'function') {
						try {
							option.validate = eval(option.validate);
						} catch (e) {
							console.log(e);
						}
						if (typeof(option.validate) != 'function') {
							option.validate = function() {
								return true;
							}
						}
					}
				}
				if (sucFn) {
					try {
						sucFn = eval(sucFn);
					} catch (e) {
						console.log(e);
					}
					if (typeof(sucFn) != 'function') {
						sucFn = null;
					}
				}
				webx.validate($(this), sucFn, option);
			});
		},
		showMsgs: function(once) {
			if (once == null) once = false;
			if (webx.msgs.err && webx.msgs.suc) {
				webx.dialog().msg('<div>' + webx.msgs.err + '</div><div>' + webx.msgs.suc + '</div>', {
					offset: '10px',
					shift: 6,
					icon: 0,
					time: 10000
				});
				if (once) webx.resetMsgs();
			} else if (webx.msgs.err) {
				webx.dialog().msg(webx.msgs.err, {
					offset: '10px',
					shift: 6,
					icon: 5,
					time: 8000
				});
				if (once) webx.resetMsgs();
			} else if (webx.msgs.suc) {
				webx.dialog().msg(webx.msgs.suc, {
					offset: '10px',
					icon: 6,
					time: 5000
				});
				if (once) webx.resetMsgs();
			}
		},
		resetMsgs: function() {
			webx.msgs = {
				err: null,
				suc: null,
				code: null
			};
		},
		asMsgs: function(obj) {
			webx.msgs.code = obj.Status;
			if (obj.Status == 1) {
				webx.msgs.suc = obj.Message;
			} else {
				webx.msgs.err = obj.Message;
			}
		},
		setMsgs: function(code, msg) {
			webx.msgs.code = code;
			if (code == 1) {
				webx.msgs.suc = msg;
			} else {
				webx.msgs.err = msg;
			}
		},
		dialog: function() {
			var type = typeof(layer);
			if (type == 'undefined') {
				window.LAYER_PATH = webx.staticUrl + 'js/Dialog/layer/';
				webx.defined(type, 'layer');
				layer.config({
					extend: ['extend/layer.ext.js', 'skin/moon/style.css'],
					skin: 'layer-ext-moon'
				});
			}
			return layer;
		},
		noty: function(option, timeout, maxVisible) { //webx.noty({text:'webx'});
			if (timeout == null) timeout = 3000;
			if (maxVisible == null) maxVisible = 5;
			if (typeof(option) != 'object') option = {
				text: option
			};
			var defaults = {
				text: 'text',
				//type: warning/error/information/success/notification
				type: 'information',
				layout: 'topRight',
				theme: 'relax',
				maxVisible: maxVisible,
				closeWith: ['click'],
				//timeout: false
				timeout: timeout,
				animation: {
					open: 'animated bounceInRight',
					close: 'animated bounceOutRight',
					easing: 'swing',
					speed: 500
				},
				tmpl: '<div class="activity-item"><i class="fa fa-{%icon%}"></i><div class="activity">{%content%}</div></div>'
			};
			option = $.extend({}, defaults, option || {})
			webx.defined(typeof(noty), 'noty');
			if (option.tmpl) {
				if (typeof(option.text) != 'object') option.text = {
					content: option.text
				};
				if (!option.text.icon) {
					switch (option.type) {
					case 'success':
						option.text.icon = 'ok';
						break; //smile-o
					case 'warning':
						option.text.icon = 'warning';
						break;
					case 'information':
						option.text.icon = 'info';
						break;
					case 'error':
						option.text.icon = 'ban';
						break; //meh-o
					case 'notification':
						option.text.icon = 'bullhorn';
						break;
					}
				}
				option.text = webx.parseTmpl(option.tmpl, option.text);
			}
			return noty(option);
		},
		captcha: {
			show: function(element, app, ident) {
				if (ident == null) ident = 'captcha';
				if ($('#' + ident + 'Image').length > 0) return;
				$.get(webx.siteUrl + 'captcha/reload', {
					format: 'json',
					app: app,
					v: Math.random()
				}, function(r) {
					if (typeof(r) != 'object') return;
					var id = r.Data.Id;
					var rel = $(element).attr('rel');
					var style = '';
					if (rel != 'nostyle') style = 'border-radius:5px;border:1px solid #DDD;box-shadow:0 0 5px #EEE;';
					var captcha = $('<img id="' + ident + 'Image" src="' + webx.siteUrl + 'captcha/' + id + '.png" alt="Captcha image" title="' + webx.t("点击这里刷新验证码") + '" onclick="webx.captcha.click(this,\'' + app + '\');" onerror="webx.captcha.action(this,\'' + app + '\');" style="' + style + 'cursor:pointer" /><input type="hidden" name="captchaId" id="' + ident + 'Id" value="' + id + '" />');
					$(element).html(captcha);
				}, 'json');
			},
			click: function(element, app) {
				var spt = $(element).attr('src').split('?');
				webx.data.captchaErrorTimes = 0;
				$(element).attr('src', spt[0] + '?app=' + app + '&reload=' + Math.random());
			},
			monitor: function(element, app) {
				$(element).error(function() {
					webx.captcha.action(this, app);
				});
			},
			action: function(element, app) {
				if (webx.data.captchaErrorTimes > 1) {
					if (webx.data.captchaErrorTimes < 9) {
						alert(webx.t("验证码图片已经失效，请刷新页面重试。"));
						webx.data.captchaErrorTimes = 9;
					}
					return;
				}
				var obj = $(element);
				$.get(webx.siteUrl + 'captcha/reload', {
					format: 'json',
					app: app,
					v: Math.random()
				}, function(r) {
					if (typeof(r) != 'object') return;
					var id = r.Data.Id;
					webx.data.captchaErrorTimes++;
					obj.attr('src', webx.siteUrl + 'captcha\/' + id + '.png');
					obj.next('input[type=hidden]').val(id);
				}, 'json');
			}
		},
		validate: function(element, sucFn, options) {
			webx.defined(typeof($.fn.html5Validate), 'validate');
			var defaults = {
				novalidate: false,
				validate: function() {
					return true;
				},
				submit: true
			};
			var params = $.extend({}, defaults, options || {});
			var object;
			if (typeof(element) == 'object' && typeof(element.length) != 'undefined') {
				object = element
			} else {
				object = $(element);
			}
			return object.html5Validate(function() {
				if (sucFn != null) sucFn();
				if (params.submit && $(this).get(0).tagName.toLowerCase() == 'form') this.submit();
				return true;
			}, params);
		},
		parseAjaxSetting: function(data) {
			if (typeof(data)!='string') return data;
			var ti = data.indexOf("json:");
			if (ti === 0) {
				data = $.parseJSON(data.substring(5));
				return data;
			}
			var elem = data;
			data = {};
			ti = elem.indexOf("elem:");
			if (ti !== 0) return data;
			elem = elem.substring(5);
			if ($(elem).length < 1) return data;
			if (elem.indexOf(",") === -1) {
				switch ($(elem).get(0).tagName) {
				case "FORM":
					data = $(elem).serialize();
					break;
				case "INPUT":
					var ob = $(elem),
						tp = ob.first().attr('type');
					var me = ob.first().attr('name');
					if (tp.toUpperCase() == "CHECKBOX") {
						data[me] = [];
						ob.has(":checked").each(function() {
							data[me].push($(this).val())
						});
						break;
					}
					data[me] = $(elem).first().val();
					break;
				case "SELECT":
					var ob = $(elem),
						multi = ob.first().attr('multiple');
					var me = ob.first().attr('name');
					if (multi) {
						data[me] = [];
						ob.find("option:selected").each(function() {
							data[me].push($(this).val())
						});
						break;
					}
					data[me] = $(elem).first().val();
					break;
				default:
					var me = $(elem).first().attr('name');
					data[me] = $(elem).first().val();
				}
				return data;
			}
			var ob = $(elem);
			ob.each(function() {
				var me = $(this).attr('name');
				switch ($(this).get(0).tagName) {
				case "INPUT":
					var tp = $(this).attr('type');
					if (tp.toUpperCase() == "CHECKBOX") {
						if (typeof(data[me]) == 'undefined') {
							data[me] = [];
						}
						data[me].push($(this).val());
						break;
					}
					data[me] = $(this).val();
					break;
				case "SELECT":
					var multi = $(this).attr('multiple');
					if (multi) {
						if (typeof(data[me]) == 'undefined') {
							data[me] = [];
						}
						$(this).find('option:selected').each(function() {
							data[me].push($(this).val())
						});
						break;
					}
					data[me] = $(this).val();
					break;
				default:
					data[me] = $(this).val();
				}
			});
			return data;
		},
		autoAjax: function(element) {
			if (element == null) element = '';
			$(element + '[data-ajax-url]').click(function(e) {
				e.preventDefault();
				var that = $(this);
				var href = that.data('ajax-url'),
					type = that.data('ajax-type'),
					dataType = that.data('ajax-dataType'),
					confirmMsg = that.data('ajax-confirm'),
					data = that.data('ajax-data'),
					formId = that.data('ajax-formId');
				if (confirmMsg && !confirm(confirmMsg)) return;
				if (!href) return;
				if (!type) type = 'get';
				if (!dataType) dataType = 'json';
				if (formId) data = $('#' + formId).serializeArray();
				data = data ? webx.parseAjaxSetting(data) : {format: 'json'};
				if (typeof(data) != 'object') {
					data += '&format=json';
				} else {
					data["format"] = 'json';
				}
				$.ajax({
					url: href,
					type: type,
					data: data,
					cache: false,
					dataType: dataType,
					success: function(data, textStatus) {
						webx.ajaxr(data, function(resp, done) {
							if (that.data('ajax-reload')) {
								webx.noty(resp.Message);
								window.setTimeout(function() {
									window.location.reload();
								}, 2000);
								return;
							}
							var c = that.data('ajax-callback');
							if (c) {
								//c(resp,done);
								window.setTimeout(c, 0);
								return;
							}
							if (that.attr('type') == 'checkbox') that.prop('checked', !that.prop('checked'));
							webx.noty(resp.Message);
						});
					}
				});
			});
		},
		ajaxr: function(resp, callback, respType) {
			if (respType == null) respType = typeof(resp) == 'object' ? 'json' : '';
			if (callback == null) {
				callback = {
					'1': function() {
						webx.noty(webx.t('操作成功'));
					}
				};
			} else {
				var dataType = typeof(callback);
				switch (dataType) {
				case 'function':
					callback = {
						'1': callback
					};
					break;
				case 'string':
					callback = {
						'1': function() {
							webx.noty(callback);
						}
					};
					break;
				}
			}
			var done = null;
			switch (respType) {
			case 'json':
				if (typeof(resp) != 'object' || resp == null) return done;
				if (typeof(callback[resp.Status]) == 'function') done = callback[resp.Status](resp, done);
				if (done == null) {
					switch (resp.Status) {
					case -1:
						/*未登录*/
						webx.dialog().confirm(webx.t('登录状态已经失效，您需要重新登录。现在要前往登录界面吗？'), {
							icon: 4
						}, function(index) {
							window.location = resp.Data.Location;
						});
						break;
					case -2:
						/*无权限*/
					case 0:
						/*操作失败*/
						webx.noty({
							text: resp.Message,
							type: 'error'
						});
						break;
					}
				}
				break;
			default:
				if (typeof(callback['1']) == 'function') {
					done = callback['1'](resp, done);
				}
			}
			return done;
		},
		table:function(element,options){
			webx.defined(typeof($.fn.dataTable),'table');
			var url=$(element).data('table-url');
			var cols=$(element).data('table-cols');
			var defaults={
				"processing": true,
        		"serverSide": true,
        		"ajax": {
        			url: url,
    				data: function(d){
    					d.format='json';
    					return d;//附加提交参数
            		},
        			result: function(d){
        				d=webx.ajaxr(d,function(r, done){
        					d={
        						draw:r.Data.draw,
        						recordsTotal:r.Data.recordsTotal,
        						recordsFiltered:r.Data.recordsFiltered,
        						data:r.Data.data
        					};
        					return d;
        				});
                		return d;
            		}
    			}, 
    			"columnDefs": [],
    			"autoWidth": false,
        		"columns": [],
        		"sDom": "<'dtTop'<'dtShowPer'l><'dtFilter'f>><'dtTables't><'dtBottom'<'dtInfo'i><'dtPagination'p>>",
        		"language": {
            		'emptyTable': webx.t('没有数据'),  
                	'loadingRecords': webx.t('加载中...'),  
                	'processing': webx.t('查询中...'),  
                	'search': webx.t('检索:'),  
                	'lengthMenu': webx.t('每页 _MENU_ 行'),  
                	'zeroRecords': webx.t('没有数据'),  
                	'paginate': {  
                    	'first': webx.t('第一页'),  
                    	'last': webx.t('最后一页'),  
                    	'next': '&gt;',  
                    	'previous': '&lt;'  
                	},  
                	'info': webx.t('第 _PAGE_ 页 / 共 _PAGES_ 页'),  
                	'infoEmpty': webx.t('没有数据'),  
                	'infoFiltered': webx.t('(共有 _MAX_ 行)')
        		},
        		"sPaginationType": "full_numbers",
        		"initComplete": function(){
        			var id=$(element).attr("id")+"_wrapper";
        			$("#"+id).find(".dtShowPer select").uniform();
        			$("#"+id).find(".dtFilter input").addClass("simple_field").css({
        				"width": "auto",
        				"margin-left": "15px",
        			});
        		},
        		"drawCallback": function(settings) {
        			var id=$(element).attr("id")+"_wrapper";
        			$("#"+id).find("td .simple_form").uniform();
    			}
    		};
    		options=$.extend({},defaults,options||{});
    		var toBool=function(v){
    			if(v=='0'||v=='false'||v=='off'||v=='no'||v=='n'||!v){
					v=false;
				}else{
					v=true;
				}
				return v
    		}
    		var parser=function(k,v){
				switch(k){
					case 'field':
						k='data';
					break;
					case 'orderable':
					case 'searchable':
						v=toBool(v);
					break;
				}
				return [k,v];
    		}
			if (cols) {
				cols=cols.split(';');
				for (var i = 0; i < cols.length; i++) {
					cols[i]=$.trim(cols[i]);
					var kvs=cols[i].split(','),cd={};
					for (var j = 0; j < kvs.length; j++) {
						kvs[j]=$.trim(kvs[j]);
						var kv=kvs[j].split(':');
						if(kv.length<2){
							cd.data=kvs[j];
						}else{
							kv=parser($.trim(kv[0]),$.trim(kv[1]));
							cd[kv[0]]=kv[1];
						}
					}
					options.columns.push(cd);
				}
			} else {
				var hideCols = [];
    			$(element).find("thead > tr > th[data-col-field]").each(function(k,item){
    				var cd={};
					cd.data=$(this).data("col-field");
					var v=$(this).data("col-orderable");
					if(v!==undefined)cd.orderable=toBool(v);
					v=$(this).data("col-searchable")
					if(v!==undefined)cd.searchable=toBool(v);
					options.columns.push(cd);
					v=$(this).data("col-render");
					if(v!==undefined){
						var td = v;
						if(v.length>0&&v.substring(0,1)=="#")td=$(v).html();
						cd = {
							render:function(data,type,row){
                    			return webx.parseTmpl(td+'',{data:data,row:row,type:type});
                			},
                			targets:k
                		};
						options.columnDefs.push(cd);
					}
					v=$(this).data("col-visible");
					if(v!==undefined&&toBool(v))hideCols.push(k);
    			});
    			if (hideCols.length>0) options.columnDefs.push({visible:false,targets:hideCols});
    		}
			if(!options.ajax.url){
				options.serverSide=false;
				options.ajax=null;
			}
			return $(element).dataTable(options);
		}
	};
})();

function T(k, obj) {
	return webx.t(k, obj);
}

function D() {
	return webx.dialog();
}