package dbschema

type Link struct {
	Id       int    `xorm:"not null pk autoincr INT(10)"`
	Name     string `xorm:"not null VARCHAR(30)"`
	Url      string `xorm:"not null VARCHAR(200)"`
	Logo     string `xorm:"not null VARCHAR(200)"`
	Show     string `xorm:"not null default 'N' ENUM('N','Y')"`
	Verified int    `xorm:"not null default 0 INT(10)"`
	Created  int    `xorm:"not null default 0 created INT(10)"`
	Updated  int    `xorm:"not null default 0 updated INT(10)"`
	Catid    int    `xorm:"not null default 0 INT(10)"`
	Sort     int    `xorm:"not null default 0 INT(10)"`
}
