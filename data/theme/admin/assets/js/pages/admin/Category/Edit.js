function select(elem){
	var e=$(elem);
	var url=e.data('url');
	$.get(url,{format:'json',pid:0},function(r){
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
				if (e.val()==''||e.val()=='0') {
					while(e.next('select').length>0)e.next('select').remove();
					return;
				};
				rselect(url,e);
			});
			if (rel) e.trigger('change');
		});
	},'json');
}
function rselect(url,e){
	var pid=e.val(),id=e.attr('id');
	$.get(url,{format:'json',pid:pid},function(r){
		webx.ajaxr(r,function(res,done){
			var opts='',rel=e.attr('rel');
			for(var i=0;i<res.Data.data.length;i++){
				var v=res.Data.data[i];
				var s=rel==v.Id?' selected="selected"':'';
				opts+='<option value="'+v.Id+'"'+s+'>'+v.Name+'</option>';
			}
			if (!opts) return;
			var ne=e.next('select[parent="'+id+'"]');
			if (ne.length<1) {
				var id2='select-'+Math.random();
				ne=$('<select id="'+id2+'" parent="'+id+'"></select>');
				e.after(ne);
				ne.change(function(){
					if (ne.val()==''||ne.val()=='0') {
						while(ne.next('select').length>0)ne.next('select').remove();
						return;
					};
					rselect(url,ne);
				});
			}
			ne.html(opts);
			if (rel) ne.trigger('change');
		});
	},'json');
}
$(function(){
	select('#select-rcategeries');
	webx.nestedSelect(["select-rcategeries","province_id","city_id"])
});