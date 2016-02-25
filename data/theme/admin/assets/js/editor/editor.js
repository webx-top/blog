webx.libs.editormdPreview=['editor/markdown/lib/marked.min.js','editor/markdown/lib/prettify.min.js','editor/markdown/lib/raphael.min.js','editor/markdown/lib/underscore.min.js','editor/markdown/css/editormd.preview.min.css','editor/markdown/editormd.min.js'];
webx.libs.editormd=['editor/markdown/css/editormd.min.css','editor/markdown/editormd.min.js'];
webx.libs.flowChart=['editor/markdown/lib/flowchart.min.js','editor/markdown/lib/jquery.flowchart.min.js'];
webx.libs.sequenceDiagram=['editor/markdown/lib/sequence-diagram.min.js'];
webx.libs.xheditor=['editor/xheditor/xheditor.min.js','editor/xheditor/xheditor_lang/'+webx.lang+'.js'];
webx.libs.codehighlight=['codeHighlight/loader/run_prettify.js?skin=sons-of-obsidian'];

/* 解析markdown为html */
function parseMarkdown2HTML(viewZoneId,markdownData,options){
	var defaults={
       markdown        : markdownData,
       //htmlDecode    : true,  // 开启HTML标签解析，为了安全性，默认不开启
       htmlDecode      : "style,script,iframe",  // you can filter tags decode
       //toc           : false,
       tocm            : true,  // Using [TOCM]
       //gfm           : false,
       //tocDropdown   : true,
       emoji           : true,
       taskList        : true,
       tex             : true,  // 默认不解析
       flowChart       : true,  // 默认不解析
       sequenceDiagram : true,  // 默认不解析
    };
	var params = $.extend({}, defaults, options || {});
	if(params.flowChart)webx.defined(typeof($.fn.flowChart),'flowChart');
	if(params.sequenceDiagram)webx.defined(typeof($.fn.sequenceDiagram),'sequenceDiagram');
	webx.defined(typeof(editormd),'editormdPreview');
	var EditormdView = editormd.markdownToHTML(viewZoneId, params);
	return EditormdView;
}
/* 初始化Markdown编辑器 */
function initEditorMarkdown(editorElement,uploadUrl,options){
	webx.defined(typeof(editormd),'editormd');
	if(uploadUrl!=null){
		if(uploadUrl.indexOf('?')>=0){
			uploadUrl+='&';
		}else{
			uploadUrl+='?';
		}
		uploadUrl+='format=json';
		uploadUrl+='&client=markdown';
	}
	var container=$(editorElement).parent();
	var containerId=container.attr('id');
	if (containerId===undefined) {
		containerId='webx-md-'+window.location.href;
		container.attr('id',containerId);
	};
	var defaults={
       width : "100%",
       height : container.height(),
       path : webx.staticUrl+"js/editor/markdown/lib/",
       markdown : $(editorElement).text(),
       codeFold : true,
       saveHTMLToTextarea : true,			// 保存HTML到Textarea
       searchReplace : true,
       watch : true,						// 关闭实时预览
       htmlDecode : "style,script,iframe",	// 开启HTML标签解析，为了安全性，默认不开启
       emoji : true,
       taskList : true,
       tocm : true,					 // Using [TOCM]
       tex : true,                   // 开启科学公式TeX语言支持，默认关闭
       flowChart : true,             // 开启流程图支持，默认关闭
       sequenceDiagram : true,       // 开启时序/序列图支持，默认关闭,
       imageUpload : true,
       imageFormats : ["jpg", "jpeg", "gif", "png", "bmp"],
       imageUploadURL : uploadUrl,
       crossDomainUpload : true,
       uploadCallbackURL : webx.staticUrl+"js/editor/markdown/plugins/image-dialog/upload_callback.htm",
       onload : function(){}
    };
	var params = $.extend({}, defaults, options || {});
	if (!uploadUrl) params.imageUpload=false;
	var EditormdView = editormd(containerId, params);
	return EditormdView;
}

/* 初始化xheditor */
function initEditorX(editorElement,uploadUrl,uploadType){
	webx.defined(typeof($.fn.xheditor),'xheditor');
	var editor;
	if(!uploadUrl){editor=$(editorElement).xheditor({});}else{
	if(uploadUrl.indexOf('?')>=0){
		uploadUrl+='&';
	}else{
		uploadUrl+='?';
	}
	uploadUrl+='format=json&client=xheditor';
	var plugins={
		Code:{c:'xhe_btnCode',t:'插入代码',h:1,e:function(){
			var that=this;
			var lang=[/*"aea","agc","apollo",basic","cbm","cl","clj",*/"css","dart",/*"el","erl",*/"erlang",/*"fs",*/"go",/*"hs",*/"html","javascript",/*"latex","lisp","ll","llvm","lsp","lua","matlab","ml","mumps","n","nemerle","pascal",*/"php",/*"proto","r","rd","rkt","s",*/"scala",/*"scm","Splus",*/"sql",/*"ss","tcl","tex","vb","vbs","vhd","vhdl","wiki","xq",*/"xquery","xml","yaml","yml"];
			var htmlCode='<div><select id="xheCodeType">';
			for(var i=0;i<lang.length;i++){
				var s=lang[i]=='go'?' selected="selected"':'';
				htmlCode+='<option value="'+lang[i]+'"'+s+'>'+lang[i]+'</option>';
			}
			htmlCode+='<option value="">其它</option></select></div><div><textarea id="xheCodeValue" wrap="soft" spellcheck="false" style="width:300px;height:100px;" /></div><div style="text-align:right;"><input type="button" id="xheSave" value="确定" /></div>';
			var jCode=$(htmlCode),jType=$('#xheCodeType',jCode),
				jValue=$('#xheCodeValue',jCode),jSave=$('#xheSave',jCode);
			jSave.click(function(){
				that.loadBookmark();
				that.pasteHTML('<pre class="prettyprint linenums lang-'+jType.val()+'">'+that.domEncode(jValue.val())+'</pre>');
				that.hidePanel();
				return false;
			});
			that.saveBookmark();
			that.showDialog(jCode);
		}},
        EndInput:{c:'xhe_btnEndInput',t:'末尾新行 (Shift+End)',s:'shift+end',e:function(){
			this.appendHTML('<p><br /></p>');/*解决光标无法移出容器的问题*/
        }}
	};
	var option={
	'skin':'default',//'shortcuts':{'ctrl+enter':submitForm},'loadCSS':'<style></style>',
	'plugins':plugins,
	'upLinkUrl':uploadUrl+'&type=file',
	'upLinkExt':"zip,rar,7z,tar,gz,txt,xls,doc,docx,ppt,pptx,et,wps,rtf,dps",
	'upImgUrl':uploadUrl+'&type=image',
	'upImgExt':"jpg,jpeg,gif,png",
	'upFlashUrl':uploadUrl+'&type=flash',
	'upFlashExt':"swf",
	'upMediaUrl':uploadUrl+'&type=media',
	'upMediaExt':"avi,wmv,wma,mp3,mp4,mpeg,mkv,rm,rmv,mid"
	};
	if(uploadType!=null && typeof(uploadType)=="object"){
		for (var i in option) {
			if(typeof(uploadType[i])!="undefined")option[i]=uploadType[i];
		}
	}
	/* IE10以下不支持HTML5中input:file域的mutiple属性，采用iframe加载swfupload实现批量选择上传 */
	if($.browser.msie && parseFloat($.browser.version) < 10.0) {
		uploadUrl='!{editorRoot}xheditor_plugins/multiupload/multiupload.html?uploadurl='+encodeURIComponent(uploadUrl);
		if (option.upLinkUrl) {
			option.upLinkUrl=uploadUrl+'&ext=Attachment('+'*.'+option.upLinkExt.replace(/,/g,';*.')+')';
			option.upLinkExt='';
		}
		if (option.upImgUrl) {
			option.upImgUrl=uploadUrl+'&ext=Image('+'*.'+option.upImgExt.replace(/,/g,';*.')+')';
			option.upImgExt='';
		}
		if (option.upFlashUrl) {
			option.upFlashUrl=uploadUrl+'&ext=Flash('+'*.'+option.upFlashExt.replace(/,/g,';*.')+')';
			option.upFlashExt='';
		}
		if (option.upMediaUrl) {
			option.upMediaUrl=uploadUrl+'&ext=Media('+'*.'+option.upMediaUrl.replace(/,/g,';*.')+')';
			option.upMediaExt='';
		}
	}
	editor=$(editorElement).xheditor(option);
	}
	return editor;
}

//例如：switchEditor($('textarea'))
function switchEditor(texta,cancelFn){
	var upurl=texta.data("upload-url");
	var etype=texta.data("editor");
	var ctype=texta.attr("data-current-editor");
	if (ctype==etype) return;
	var className=texta.data("class");
	if (className===undefined) {
		className=texta.attr("class");
		if (!className) className='';
		texta.data("class",className);
	}
	var content=texta.data("content-elem"),cElem=content;
	if (content) cElem=webx.parseTmpl(content,{type:etype});
	var obj=texta.get(0);
	switch(etype){
		case 'markdown':
		if(typeof(texta.xheditor)!='undefined'){
			/*
			var cc=webx.parseTmpl(content,{type:ctype});
			if (cc&&$(cc).length>0) {
				if (texta.val()!=$(cc).val()&&!confirm('确定要切换吗？切换编辑器将会丢失您当前所做的修改。')) {
					if (cancelFn!=null) cancelFn();
					return false;
				};
			};
			*/
			texta.xheditor(false);
		}
		if(cElem&&$(cElem).length>0){
			texta.text($(cElem).val());
			texta.val($(cElem).val());
		}
		initEditorMarkdown(obj,upurl);
		texta.attr("data-current-editor",etype);
		break;
		default:
		if(cElem&&$(cElem).length>0){
			var cc=webx.parseTmpl(content,{type:ctype});
			if (cc&&$(cc).length>0) {
				$(cc).text(texta.val());
				var ht=$('textarea[name="'+texta.parent().attr('id')+'-html-code"]');
				if (ht.length>0&&ht.val()!="") {
					$(cElem).text(ht.val());
					$(cElem).val(ht.val());
				}
			};
			texta.val($(cElem).val());
			texta.text($(cElem).val());
		}
		texta.parent().removeAttr('class');
		texta.attr('class',className).siblings().remove();
		initEditorX(obj,upurl);
		texta.attr("data-current-editor","html");
	};
	return true;
}