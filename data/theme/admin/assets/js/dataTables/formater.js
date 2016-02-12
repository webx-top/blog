var formater={};
formater.date=function(data,type,row){
	if (data==0) return 'N/A';
	return webx.date('Y-m-d H:i:s',data);
};
formater.onOff=function(data,type,row){
	if (data==0) return '<i class="fa fa-square" style="color:#eee"></i>';
	return '<i class="fa fa-check-square" style="color:#ba1016"></i>';
};
var eventer={};
eventer.onInit=function(element){
    var id=$(element).attr("id")+"_wrapper";
    $("#"+id).find(".dtShowPer select").uniform();
    $("#"+id).find(".dtFilter input").addClass("simple_field").css({
        "width": "auto","margin-left": "15px",
    });
};
eventer.onDraw=function(element){
	var id=$(element).attr("id")+"_wrapper";
	$("#"+id).find("td .simple_form").uniform();
};