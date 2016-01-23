package dbschema

type User struct {
	Id        int    `xorm:"not null pk autoincr INT(10)"`
	Uname     string `xorm:"not null VARCHAR(30)"`
	Passwd    string `xorm:"not null CHAR(64)"`
	Salt      string `xorm:"not null CHAR(64)"`
	Email     string `xorm:"not null default '' VARCHAR(100)"`
	Mobile    string `xorm:"not null default '' VARCHAR(15)"`
	LoginTime int    `xorm:"not null default 0 INT(10)"`
	LoginIp   string `xorm:"not null default '' VARCHAR(40)"`
	Created   int    `xorm:"not null default 0 created INT(10)"`
	Updated   int    `xorm:"not null default 0 updated INT(10)"`
	Active    string `xorm:"not null default 'Y' ENUM('Y','N')"`
	Avatar    string `xorm:"not null default '' VARCHAR(200)"`
}
