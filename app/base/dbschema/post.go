package dbschema

type Post struct {
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Title        string `xorm:"not null VARCHAR(180)" valid:"required"`
	Description  string `xorm:"not null VARCHAR(200)"`
	Content      string `xorm:"not null TEXT" valid:"required"`
	Etype        string `xorm:"not null default 'html' ENUM('html','markdown')"`
	Created      int    `xorm:"not null default 0 created INT(10)"`
	Updated      int    `xorm:"not null default 0 updated INT(10)"`
	Display      string `xorm:"not null default 'ALL' ENUM('ALL','SELF','FRIEND','PWD')"`
	Uid          int    `xorm:"not null default 0 INT(10)"`
	Uname        string `xorm:"not null default '' VARCHAR(30)"`
	Passwd       string `xorm:"not null default '' VARCHAR(64)"`
	Views        int    `xorm:"not null default 0 INT(10)"`
	Comments     int    `xorm:"not null default 0 INT(10)"`
	Likes        int    `xorm:"not null default 0 INT(10)"`
	Deleted      int    `xorm:"not null default 0 INT(10)"`
	Year         int    `xorm:"not null INT(5)"`
	Month        int    `xorm:"not null TINYINT(1)"`
	AllowComment string `xorm:"not null default 'Y' ENUM('Y','N')"`
	Tags         string `xorm:"not null default '' VARCHAR(255)"`
	Catid        int    `xorm:"not null default 0 INT(10)"`
}
