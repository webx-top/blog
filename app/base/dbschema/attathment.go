package dbschema

type Attathment struct {
	Id        int    `xorm:"not null pk autoincr INT(10)"`
	Name      string `xorm:"not null VARCHAR(100)"`
	Path      string `xorm:"not null VARCHAR(255)"`
	Extension string `xorm:"not null VARCHAR(5)"`
	Type      string `xorm:"not null default 'image' ENUM('media','other','image')"`
	Size      int64  `xorm:"not null BIGINT(20)"`
	Uid       int    `xorm:"not null INT(10)"`
	Deleted   int    `xorm:"not null default 0 INT(10)"`
	Created   int    `xorm:"not null created INT(10)"`
	Audited   int    `xorm:"not null default 0 INT(10)"`
	RcId      int    `xorm:"not null default 0 INT(10)"`
	RcType    string `xorm:"not null default '' VARCHAR(30)"`
	Tags      string `xorm:"not null default '' VARCHAR(255)"`
}
