package dbschema

type Category struct {
	Id       int    `xorm:"not null pk autoincr INT(10)"`
	Pid      int    `xorm:"not null default 0 INT(10)"`
	Name     string `xorm:"not null VARCHAR(30)"`
	Haschild string `xorm:"not null default 'N' ENUM('N','Y')"`
	Updated  int    `xorm:"not null default 0 updated INT(10)"`
	RcType   string `xorm:"not null default 'post' VARCHAR(30)"`
}
