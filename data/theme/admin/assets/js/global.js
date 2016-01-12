window.webx={calls:[],lang:'zh-cn',libs:{

},staticUrl:'',siteUrl:'',cachedData:{}};
webx.include=function(file,location){
	if(location==null)location="head";
	if(location=="head" && typeof(webx.cachedData["include"])=="undefined"){
	var jsAfter=$("#js-lazyload-begin"),cssAfter=$("#css-lazyload-begin");
	webx.cachedData.include={before:{},after:{}};
	if(jsAfter.length){
		webx.cachedData.include.after.script=jsAfter;
	}else{
		var jsBefore=$("#js-lazyload-end");
		if(jsBefore.length)webx.cachedData.include.before.script=jsBefore;
	}
	if(cssAfter.length){
		webx.cachedData.include.after.link=cssAfter;
	}else{
		var cssBefore=$("#css-lazyload-end");
		if(cssBefore.length)webx.cachedData.include.before.link=cssBefore;
	}
	}
	var files = typeof(file)=="string" ? [file] : file;
	for (var i = 0; i < files.length; i++) {
		var name = files[i].replace(/^\s|\s$/g, ""),att = name.split('.');
		var ext = att[att.length - 1].toLowerCase(),isCSS = ext == "css";
		var tag = isCSS ? "link" : "script";
		var attr = isCSS ? ' type="text/css" rel="stylesheet"' : ' type="text/javascript"';
		attr += ' charset="utf-8" ';
		var link = (isCSS ? "href" : "src") + "='" + $.IncPath + name + "'";
		if ($(tag + "[" + link + "]").length == 0) {
			var ej = $("<" + tag + attr + link + "></" + tag + ">");
			if(location=="head"){
				if(typeof(webx.cachedData.include.after[tag])!='undefined'){
					webx.cachedData.include.after[tag].after(ej);
				}else if(typeof(webx.cachedData.include.before[tag])!='undefined'){
					webx.cachedData.include.before[tag].before(ej);
				}else{
					$(location).append(ej);/*prepend*/
				}
			}else{
				$(location).append(ej);
			}
		}
	}
};
webx.defined=function(vType,key,callback){
	if(vType!='undefined'||key==null){
		if(key!=null)return callback();
		return;
	}
	if(typeof(key)=='string' && typeof(webx.libs[key])!='undefined') key=webx.libs[key];
	webx.include(key);
	if(callback!=null)return callback();
};
webx.scrollTo=function(element,time){
	if(!time) time = 1000;
	$('html,body').animate({scrollTop:$(element).offset().top},time);
};
function addCalls(func){
	webx.calls.push(func);
}
function doCalls(){
	for(var i=0,len=webx.calls.length;i<len;i++){
		webx.calls[i]();
	}
	webx.calls=[];
}
/* 调用译文 */
function tr(k, obj){
	var lang;
	if (typeof(Lang) == "undefined" || typeof(Lang[k]) == "undefined") {
		lang = k;
	} else {
		lang = Lang[k];
	}
	if (obj != null) return parseTmpl(lang, obj);
	return lang;
}
function Tr(k,obj){
	return tr(k,obj);
}
/* 解析模板 */
function parseTmpl(template, data) {
  return template.replace(/\{%([\w\.]*)%\}/g, function(str, key) {
	var keys = key.split("."), v = data[keys.shift()];
	for (var i = 0, l = keys.length; i < l; i++) v = v[keys[i]];
	return typeof(v)!== "undefined" && v !== null ? v : "";
  });
}
/* 插入数据到光标位置 */
function insertAtCursor(myField, myValue) {
	 /* IE support */
	 if (document.selection) {
		 myField.focus();
		 sel = document.selection.createRange();
		 sel.text = myValue;
		 sel.select();
	 }
	 /* MOZILLA/NETSCAPE support */
	 else if (myField.selectionStart || myField.selectionStart == '0') {
		 var startPos = myField.selectionStart;
		 var endPos = myField.selectionEnd;
		 /* save scrollTop before insert */
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
}
/* 复选框全选 */
function checkedAll(checkbox,target){
	if(target==null)target='input[type=checkbox]';
	$(target).not(':disabled').prop('checked', $(checkbox).prop('checked'));
}
/* 确认关闭窗口 */
function confirmCloseWindow(msg){
	if(msg==null){
		if($('body[onbeforeunload]').length)$('body[onbeforeunload]').removeAttr('onbeforeunload');
		return;
	}
	if($('body').attr('onbeforeunload'))return;
	if(!msg)msg=Tr('您填写的数据没有提交，如果离开本页面这些数据将会丢失。\n确定丢弃这些内容吗？');
	$('body').attr('onbeforeunload',"return '"+msg+"';");
}
/* 回车键事件 */
function enterKeyEvent(ele,callback){
	$(ele).unbind('keydown');
	$(ele).keydown(function(event){
		if(event.keyCode==13){callback.call($(this));return false;}
	});
	return false;
}
/* 左右键翻页jQuery版 */
function turningPage(prevPage,nextPage,isElement){
	$(document).keyup(function(event){
		if(event.keyCode==37){
			if(!isElement){
				if(!prevPage){
					alert(Tr('没有了。这已经是第一页了。'));
					return;
				}
				window.location=prevPage;
			}else{
				if($(prevPage).length<1){
					alert(Tr('没有了。这已经是第一页了。'));
					return;
				}
				$(prevPage).click();
			}
		}else if(event.keyCode==39){
			if(!isElement){
				if(!nextPage){
					alert(Tr('没有了。这已经是最末页了。'));
					return;
				}
				window.location=nextPage;
			}else{
				if($(nextPage).length<1){
					alert(Tr('没有了。这已经是最末页了。'));
					return;
				}
				$(nextPage).click();
			}
		}
	});
	$(':text,textarea').keyup(function(event){
		event.stopPropagation();
	});
}
function unbindKeyEvent(){
	$(document).unbind('keyup');
	$(':text,textarea').unbind('keyup');
}
/* 禁止复制 */
function disabledCopy(el){
	var fn=function(){return false;};
	$(el).attr('unselectable','on').css({
		'-moz-user-select':'-moz-none',
		'-moz-user-select':'none',
		'-o-user-select':'none',
		'-khtml-user-select':'none',
		'-webkit-user-select':'none',
		'-ms-user-select':'none',
		'user-select':'none'
	}).bind('selectstart',fn).bind('contextmenu',fn)
	.bind('dragstart',fn).bind('selectstart',fn).bind('beforecopy',fn);
}
/* 级联选择(使用前请确保第一个下拉框已有选中项)
使用方法：nestedSelect(["country_id","province_id","city_id"]) */
function nestedSelect(ids, initVal, attrName, timeout){
	if(typeof(ids)=='object'){
		var obj=ids;
		if(typeof(obj.initVal)!='undefined') initVal=obj.initVal;
		if(typeof(obj.attrName)!='undefined') attrName=obj.attrName;
		if(typeof(obj.timeout)!='undefined') timeout=obj.timeout;
		if(typeof(obj.ids)!='undefined') ids=obj.ids;
		obj=null;
	}
	var id=ids[0],id2=ids[1];
	if(initVal==null)initVal='';
	if(attrName==null)attrName='rel';
	if(timeout==null)timeout=5000;
	var attr=$('#'+id2).attr(attrName);
	if(!attr) return false;
	if($('#'+id).val()==initVal) return false;
	if($('#'+id2+' option:last').val()!=initVal) return false;
	$('#'+id).trigger('change');
	var i=0;
	var ptimer=window.setInterval(function(){
		i++;
		if($('#'+id2+' option:last').val()!=initVal || i*200>timeout){
			window.clearInterval(ptimer);
			var sel=$('#'+id2+' option[value="'+attr+'"]');
			if(sel.length<=0)return;
			sel.prop('selected',true);
			ids.shift();
			if(ids.length>1)nestedSelected(ids,initVal,attrName,timeout);
		}
	},200);
	return true;
}