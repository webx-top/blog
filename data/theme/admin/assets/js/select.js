function select(elem,pagerows,wrapper){
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
					hideChildren($(this),wrapper);
					return;
				}
				rselect(url,$(this),sync,ignore,pagerows,wrapper);
			});
			if (rel) e.trigger('change');
		});
	},'json');
}
function hideChildren(el,wrapper){
	var elem;
	var e=el;
	if(wrapper&&el.parent(wrapper).length>0){
		e=el.parent(wrapper);
		elem=wrapper+':visible';
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
function rselect(url,e,sync,ignore,pagerows,wrapper){
	var pid=e.val();
	$.get(url,{format:'json',pid:pid,ignore:ignore,pagerows:pagerows},function(r){
		webx.ajaxr(r,function(res,done){
            if (res.Data.data.length<1) {
                if(sync)$(sync).val(e.val());
                hideChildren(e,wrapper);
                return;
            }
			var hasWrap=wrapper&&e.parent(wrapper).length>0;
			var ne=hasWrap?e.parent(wrapper).next(wrapper):e.next('select');
			if (ne.length<1) {
				var tag='<select></select>';
				if(hasWrap){
					var tagName=e.parent(wrapper).get(0).tagName.toLowerCase();
					tag='<'+tagName+' class="'+e.parent(wrapper).attr('class')+'">'+tag+'</'+tagName+'>';
				}
				ne=$(tag);
				if(hasWrap){
					e.parent(wrapper).after(ne);
				}else{
					e.after(ne);
				}
				
			}else if(!ne.is(':visible')){
                ne.show();
            }
			var selectObj=hasWrap?ne.find('select'):ne;
            if (!ne.data('event')) {
                ne.data('event','1');
				selectObj.change(function(){
                    var v=$(this).val();
					if (v==''||v=='0') {
                        if(sync)$(sync).val(e.val());
                        hideChildren($(this),wrapper);
						return;
					}
                    if(sync)$(sync).val(v);
					rselect(url,$(this),sync,ignore,pagerows,wrapper);
				});
            }
			var opts='';
			var rel=hasWrap?ne.find('select').attr('rel'):ne.attr('rel');
			for(var i=0;i<res.Data.data.length;i++){
				var v=res.Data.data[i];
				var s=rel==v.Id?' selected="selected"':'';
				opts+='<option value="'+v.Id+'"'+s+'>'+v.Name+'</option>';
			}
			opts='<option value="">- '+webx.t('请选择')+' -</option>'+opts;
			selectObj.html(opts);
			
			if (rel) selectObj.trigger('change');
		});
	},'json');
}