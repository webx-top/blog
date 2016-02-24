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
if (typeof(errors)=='object'&&errors) {
  for(var i in errors){
    $("form [name='"+i+"']").parent().append('<div class="field_notice error_tips" rel="'+i+'">'+errors[i]+'</div>');
  }
  if($("form div.error_tips").length){
    if (errorFor=='') errorFor=i;
    var ipt=$("form [name='"+errorFor+"']:first");
    if (ipt.length>0) {
    if (!ipt.is(":visible")) {
      webx.scrollTo(ipt.parent());
    }else{
      webx.scrollTo(ipt);
    }
  };
  }
};
});