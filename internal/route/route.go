package route

import (
	"github.com/flamego/flamego"
	"gorm.io/gorm"

	"github.com/wuhan005/go-template/internal/context"
)

func New(db *gorm.DB) *flamego.Flame {
	f := flamego.Classic()

	f.Use(context.Contexter(db))

	f.Group("/api", func() {

	})

	f.Get("/healthz")

	return f
}
