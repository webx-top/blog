function select(elem,pagerows){
    if(pagerows==null)pagerows=50;
	var e=$(elem);
	var url=e.data('url');
    var sync=e.data('sync');
    var ignore=e.data('ignore');
	$.get(url,{format:'json',pid:0,ignore:ignore,pagerows:pagerows},function(r){
		webx.ajaxr(r,function(res,done){
			var opts='';
			var rel=e.attr('rel');
			for(var i=0;i<res.Data.data.length;i++){
				var v=res.Data.data[i];
				var s=rel==v.Id?' selected="selected"':'';
				opts+='<option value="'+v.Id+'"'+s+'>'+v.Name+'</option>';
			}
			e.append(opts);
			e.change(function(){
                var v=$(this).val();
                if(sync)$(sync).val(e.val());
				if (v==''||v=='0') {
					hideChildren($(this));
					return;
				}
				rselect(url,$(this),sync,ignore,pagerows);
			});
			if (rel) e.trigger('change');
		});
	},'json');
}
function hideChildren(e){
	var elem;
	if(e.parent('.selector').length>0){
		e=e.parent('.selector');
		elem='.selector:visible';
	}else{
		elem='select:visible';
	}
    if (e.next(elem).length<1) {
        return;
    }
    var ss=e.siblings(elem);
    for (var index = ss.index(e.next(elem)); index < ss.length; index++) {
        ss.eq(index).hide();
    }
}
function rselect(url,e,sync,ignore,pagerows){
	var pid=e.val();
	$.get(url,{format:'json',pid:pid,ignore:ignore,pagerows:pagerows},function(r){
		webx.ajaxr(r,function(res,done){
            if (res.Data.data.length<1) {
                if(sync)$(sync).val(e.val());
                hideChildren(e);
                return;
            }
			var ne=e.next('select');
			if (ne.length<1) {
				
				var tag='<select></select>';
				if(e.parent('.selector').length>0){
					var tagName=e.parent('.selector').get(0).tagName.toLowerCase()
					tag='<'+tagName+' class="'+e.parent('.selector').attr('class')+'">'+tag+'</'+tagName+'>';
				}
				ne=$(tag);
				if(e.parent('.selector').length>0){
					e.parent('.selector').after(ne);
				}else{
					e.after(ne);
				}
				
			}else if(!ne.is(':visible')){
                ne.show();
            }
            if (!ne.data('event')) {
                ne.data('event','1');
				ne.change(function(){
                    var v=$(this).val();
					if (v==''||v=='0') {
                        if(sync)$(sync).val(e.val());
                        hideChildren($(this));
						return;
					}
                    if(sync)$(sync).val(v);
					rselect(url,$(this),sync,ignore,pagerows);
				});
            }
			var opts='',rel=ne.attr('rel');
			for(var i=0;i<res.Data.data.length;i++){
				var v=res.Data.data[i];
				var s=rel==v.Id?' selected="selected"':'';
				opts+='<option value="'+v.Id+'"'+s+'>'+v.Name+'</option>';
			}
			opts='<option value="">- '+webx.t('请选择')+' -</option>'+opts;
			if(e.parent('.selector').length>0){
				ne.find('select').html(opts);
			}else{
				ne.html(opts);
			}
			
			if (rel) ne.trigger('change');
		});
	},'json');
}