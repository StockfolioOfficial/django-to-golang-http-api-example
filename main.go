package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/article"
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/author"
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/core"
	"github.com/stockfolioofficial/django-to-golang-rest-api-example/supporter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbConnect(article.Migrate, author.Migrate)
	httpServerStart(article.URLPatterns...)
}

func dbConnect(wire ...supporter.DBWire) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	di := make([]func(*gorm.DB), len(wire))
	models := make([]interface{}, len(wire))

	for i, w := range wire {
		di[i], models[i] = w()
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	for _, d := range di {
		d(db)
	}
}

type echoValidator struct {
	validator *validator.Validate
}

func (ev *echoValidator) Validate(i interface{}) error {
	var wrapper struct {
		Value interface{} `validator:"dive"`
	}
	wrapper.Value = i
	return ev.validator.Struct(&wrapper)
}

func httpServerStart(routes ...supporter.Route) {
	e := echo.New()
	for _, route := range routes {
		route(e)
	}

	e.Validator = &echoValidator{core.DefaultValidator}
	e.Start(":8000")
}

