package author

import (
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/supporter"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:true"`
	Name      string         `gorm:"type:varchar(200);not null"`
	UpdatedAt supporter.Time `gorm:"type:datetime(6);not null"`
	CreatedAt supporter.Time `gorm:"type:datetime(6);not null"`
}

func (Model) TableName() string {
	return "author"
}


func Migrate() (func(db *gorm.DB), interface{}) {
	return func(db *gorm.DB) {
		db.Create([]Model{
			{
				ID: 1,
				Name: "Iman Tumorang",
				UpdatedAt: supporter.Time(time.Now()),
				CreatedAt: supporter.Time(time.Now()),
			},
		})
	}, &Model{}
}