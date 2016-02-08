var formater={};
formater.date=function(data,type,row){
	if (data==0) return 'N/A';
	return webx.date('Y-m-d H:i:s',data);
};
formater.onOff=function(data,type,row){
	if (data==0) return '<i class="fa fa-square" style="color:#eee"></i>';
	return '<i class="fa fa-check-square" style="color:#ba1016"></i>';
};
