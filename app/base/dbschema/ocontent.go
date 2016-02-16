package dbschema

type Ocontent struct {
	Id      int    `xorm:"not null pk autoincr INT(10)"`
	RcId    int    `xorm:"not null INT(10)"`
	RcType  string `xorm:"not null default 'post' VARCHAR(30)"`
	Content string `xorm:"not null TEXT"`
	Etype   string `xorm:"not null default 'markdown' ENUM('markdown')"`
}
