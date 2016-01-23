db_host="localhost"
db_username="root"
db_password="root"
db_table="blog"
db_charset="utf8"
db_savename="backup"
db_savepath="."
filename=`date -u +%Y%m%d%H%M%S`
mkdir $db_savepath/$filename
mysqldump -h$db_host -u$db_username -p$db_password $db_table --default-character-set=$db_charset --extended-insert --single-transaction > $db_savepath/$filename/$db_savename.sql
zip -r $db_savepath/$filename/$db_savename.sql.zip $db_savepath/$filename/
rm -rf $db_savepath/$filename/*.sql
echo "db_host=\"$db_host\"" > $db_savepath/$filename/restore.sh
echo "db_username=\"$db_username\"" >> $db_savepath/$filename/restore.sh
echo "db_password=\"$db_password\"" >> $db_savepath/$filename/restore.sh
echo "db_table=\"$db_table\"" >> $db_savepath/$filename/restore.sh
echo "db_savename=\"$db_savename\"" >> $db_savepath/$filename/restore.sh
echo "mysql -h\$db_host -u\$db_username -p\$db_password \$db_table < \$db_savename.sql" >> $db_savepath/$filename/restore.sh
echo "echo 按任意键继续" >> $db_savepath/$filename/restore.sh
echo "read -n 1" >> $db_savepath/$filename/restore.sh
yes | cp $db_savepath/$filename/$db_savename.sql.zip $GOPATH/src/coscms/$db_savename.sql.zip
rm -rf $db_savepath/$filename
#echo 按任意键继续
#read -n 1
