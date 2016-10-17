function resizeWith(){
	if ($('div.contents').length && $('.contents_wrapper').length && $('aside.sidebar').length) {
		$('div.contents').css('width',$('.contents_wrapper').width()-$('aside.sidebar').width());
	}
}
$(function(){
	webx.initPage();
	resizeWith();
	$(window).resize(resizeWith);
	$('table[data-table-url]').each(function(){
		webx.table(this);
	});
	webx.autoAjax();
/*
	$(document).on('click', 'aside.sidebar ul.tab_nav li a', function(event) {
  		var container = $('#contents_wrapper');
  		$.pjax.click(event, container);
  		$('aside.sidebar ul.tab_nav li.active_tab').removeClass('active_tab');
  		$(event.target).parent('li').addClass('active_tab');
	}).on('pjax:send',function(){
  		window.loading=webx.dialog().load();
	}).on('pjax:start',function(){
  		webx.dialog().close(window.loading);
	});
*/
	if (typeof(webx.errFor)!='undefined' && webx.errFor!='' && $("form [name='"+webx.errFor+"']").length>0) {
  		var ipt=$("form [name='"+webx.errFor+"']:first");
		var msg=webx.msgs.err?webx.msgs.err:webx.msgs.suc;
  		ipt.parent().append('<div class="field_notice error_tips" rel="'+i+'">'+msg+'</div>');
    	if(!ipt.is(":visible")){
      		webx.scrollTo(ipt.parent());
    	}else{
      		webx.scrollTo(ipt);
    	}
	}
});