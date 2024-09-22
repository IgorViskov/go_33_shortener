package storage

type User struct {
	ID   uint64    `gorm:"primary_key;auto_increment"`
	URLs []*Record `gorm:"many2many:user_urls"`
}
