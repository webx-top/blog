$(function(){
$("#formPost input[name='Display']").click(function(){
	if($(this).prop('checked')&&$(this).val()=='PWD'){
		$('#inputPasswd').show();
		$('#formPost input[name=Passwd]').prop('disabled',false);
		return;
	}
	if ($('#inputPasswd').is(':visible')) {
		$('#inputPasswd').hide();
		$('#formPost input[name=Passwd]').prop('disabled',true);
	};
});
$("#formPost input[name='Display']:checked").trigger('click');
webx.includes('editor/editor.js');
$("#formPost input[name='Etype']").click(function(){
	if ($(this).prop('checked')) {
		var etype=$(this).val();
		var texta=$("#formPost textarea[name='Content']");
		texta.data("editor",etype);
		return switchEditor(texta);
	};
});
$("#formPost input[name='Etype']:checked").trigger('click');
if (errors) {
	for(var i in errors){
		$("#formPost [name='"+i+"']").parent().append('<div class="field_notice error_tips" rel="'+i+'">'+errors[i]+'</div>');
	}
	if($("#formPost div.error_tips").length){
		var ipt=$("#formPost [name='"+i+"']:first");
		if (!ipt.is(":visible")) {
			webx.scrollTo(ipt.parent());
		}else{
			webx.scrollTo(ipt);
		}
	}
};
});