package dbschema

type Tag struct {
	Id      int    `xorm:"not null pk autoincr INT(10)"`
	Name    string `xorm:"not null VARCHAR(30)"`
	Uid     int    `xorm:"not null INT(10)"`
	Created int    `xorm:"not null created INT(10)"`
	Times   int    `xorm:"not null INT(10)"`
	RcType  string `xorm:"not null default 'post' VARCHAR(30)"`
}
