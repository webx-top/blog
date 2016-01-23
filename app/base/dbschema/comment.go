package dbschema

type Comment struct {
	Id           int64  `xorm:"BIGINT(20)"`
	Content      string `xorm:"not null TEXT"`
	Quote        string `xorm:"TEXT"`
	Etype        string `xorm:"not null default 'html' CHAR(10)"`
	RootId       int64  `xorm:"not null default 0 BIGINT(20)"`
	RId          int64  `xorm:"not null default 0 BIGINT(20)"`
	RType        string `xorm:"ENUM('reply','append')"`
	RelatedTimes int    `xorm:"not null default 0 INT(10)"`
	RootTimes    int    `xorm:"not null default 0 INT(10)"`
	Uid          int64  `xorm:"not null default 0 BIGINT(20)"`
	Uname        string `xorm:"not null VARCHAR(30)"`
	Up           int64  `xorm:"not null default 0 BIGINT(20)"`
	Down         int64  `xorm:"not null default 0 BIGINT(20)"`
	Created      int    `xorm:"not null default 0 created INT(10)"`
	Updated      int    `xorm:"not null default 0 updated INT(10)"`
	Status       int    `xorm:"not null default 0 INT(1)"`
	RcId         int64  `xorm:"not null default 0 BIGINT(20)"`
	RcType       string `xorm:"not null default 'article' CHAR(30)"`
	ForUname     string `xorm:"VARCHAR(30)"`
	ForUid       int64  `xorm:"not null default 0 BIGINT(20)"`
}
