@echo off

set db_host=127.0.0.1
set db_username=root
set db_password=root
set db_table=blog
set db_charset=utf8
set db_savename=backup
set db_savepath=.

set time_hh=%time:~0,2%
if /i %time_hh% LSS 10 (set time_hh=0%time:~1,1%)
set filename=%date:~,4%%date:~5,2%%date:~8,2%_%time_hh%%time:~3,2%%time:~6,2%

mkdir %db_savepath%\%filename%

mysqldump -h%db_host% -u%db_username% -p%db_password% %db_table% --default-character-set=%db_charset% --extended-insert --single-transaction > %db_savepath%\%filename%\%db_savename%.sql

7zr.exe a %db_savepath%\%filename%\%db_savename%.sql.zip %db_savepath%\%filename%\*.sql -r
del %db_savepath%\%filename%\*.sql

echo set db_host=%db_host% > %db_savepath%\%filename%\restore.bat
echo set db_username=%db_username% >> %db_savepath%\%filename%\restore.bat
echo set db_password=%db_password% >> %db_savepath%\%filename%\restore.bat
echo set db_table=%db_table% >> %db_savepath%\%filename%\restore.bat
echo set db_savename=%db_savename% >> %db_savepath%\%filename%\restore.bat

echo mysql -h%%db_host%% -u%%db_username%% -p%%db_password%% %%db_table%% ^< %%db_savename%%.sql >> %db_savepath%\%filename%\restore.bat
echo pause >> %db_savepath%\%filename%\restore.bat

copy %db_savepath%\%filename%\%db_savename%.sql.zip ..\..\%db_savename%.sql.zip /y

rd /s /Q %db_savepath%\%filename%
pause