$(function(){
  	webx.initPage();
	if ($(window).width()>1025 && $('div.contents').length && $('.contents_wrapper').length && $('aside.sidebar').length) {
		$('div.contents').css('width',$('.contents_wrapper').width()-$('aside.sidebar').width());
	}
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
});