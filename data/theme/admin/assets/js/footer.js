$(function(){
  	webx.initPage();
	if ($(window).width()>1025 && $('div.contents').length && $('.contents_wrapper').length && $('aside.sidebar').length) {
		$('div.contents').css('width',$('.contents_wrapper').width()-$('aside.sidebar').width());
	}
	$('table[data-table-url]').each(function(){
		webx.table(this);
	});
	webx.autoAjax();
});