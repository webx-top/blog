$(function(){
    
    webx.includes('select.js');
	select('#select-rcategeries',50,'.prettySelector');
    
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
});