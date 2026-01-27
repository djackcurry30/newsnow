package model

type Cache struct {
	ID      string    `gorm:"primaryKey;column:id"`
	Updated int64     `gorm:"column:updated;not null"`
	Data    string    `gorm:"column:data;not null"`
}

func (Cache) TableName() string {
	return "cache"
}

type User struct {
	ID      string    `gorm:"primaryKey;column:id"`
	Email   string    `gorm:"column:email;not null"`
	Data    string    `gorm:"column:data;not null;default:''"`
	Type    string    `gorm:"column:type;not null;default:'github'"`
	Created int64     `gorm:"column:created;not null"`
	Updated int64     `gorm:"column:updated;not null"`
}

func (User) TableName() string {
	return `"user"`
}
