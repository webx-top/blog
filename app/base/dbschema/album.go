package dbschema

type Album struct {
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Title        string `xorm:"not null VARCHAR(180)"`
	Description  string `xorm:"not null default '' VARCHAR(200)"`
	Content      string `xorm:"not null TEXT"`
	Created      int    `xorm:"not null created INT(10)"`
	Updated      int    `xorm:"not null default 0 updated INT(10)"`
	Views        int    `xorm:"not null default 0 INT(10)"`
	Comments     int    `xorm:"not null default 0 INT(10)"`
	Likes        int    `xorm:"not null default 0 INT(10)"`
	Display      string `xorm:"not null default 'ALL' ENUM('FRIEND','PWD','ALL','SELF')"`
	Deleted      int    `xorm:"not null default 0 INT(10)"`
	AllowComment string `xorm:"not null default 'Y' ENUM('Y','N')"`
	Tags         string `xorm:"not null default '' VARCHAR(255)"`
	Catid        int    `xorm:"not null default 0 INT(10)"`
}
