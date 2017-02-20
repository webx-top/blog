package dbschema

type Ocontent struct {
	Id      int64  `xorm:"BIGINT(30)"`
	RcId    int64  `xorm:"not null BIGINT(20)"`
	RcType  string `xorm:"not null default 'post' VARCHAR(30)"`
	Content string `xorm:"not null TEXT"`
	Etype   string `xorm:"not null default 'markdown' ENUM('markdown')"`
}
