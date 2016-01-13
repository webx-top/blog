$(function(){
	webx.doCalls();
	if ($(window).width()>1025 && $('div.contents').length && $('.contents_wrapper').length && $('aside.sidebar').length) {
		$('div.contents').css('width',$('.contents_wrapper').width()-$('aside.sidebar').width());
	}
});