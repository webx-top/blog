$("#formPost input[name='display']").click(function(){
	if($(this).prop('checked')&&$(this).val()=='PWD'){
		$('#inputPasswd').show();
		$('#formPost input[name=passwd]').prop('disabled',false);
		return;
	}
	if ($('#inputPasswd').is(':visible')) {
		$('#inputPasswd').hide();
		$('#formPost input[name=passwd]').prop('disabled',true);
	};
});
$("#formPost input[name='display']:checked").trigger('click');
webx.includes('editor/editor.js');
$("#formPost input[name='etype']").click(function(){
	if ($(this).prop('checked')) {
		var etype=$(this).val();
		var texta=$("#formPost textarea[name='content']");
		texta.data("editor",etype);
		return switchEditor(texta);
	};
});
$("#formPost input[name='etype']:checked").trigger('click');