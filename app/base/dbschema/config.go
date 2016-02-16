package dbschema

type Config struct {
	Id      int    `xorm:"not null pk autoincr INT(10)"`
	Key     string `xorm:"not null VARCHAR(60)"`
	Val     string `xorm:"not null VARCHAR(200)"`
	Updated int    `xorm:"not null default 0 updated INT(10)"`
}
