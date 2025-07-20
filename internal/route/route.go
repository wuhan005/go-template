// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	gocontext "context"
	"time"

	"github.com/MEDIGO/go-healthz"
	flamegoswagger "github.com/asjdf/flamego-swagger"
	"github.com/flamego/flamego"
	swaggerfiles "github.com/swaggo/files"
	"gorm.io/gorm"

	_ "github.com/wuhan005/go-template/docs"
	"github.com/wuhan005/go-template/internal/appconst"
	"github.com/wuhan005/go-template/internal/context"
	dbpkg "github.com/wuhan005/go-template/internal/db"
	"github.com/wuhan005/go-template/internal/form"
	"github.com/wuhan005/go-template/internal/tracing"
)

// New creates a new Flamego instance with the necessary middleware and routes.
// @Title Go Template API
// @Version 1.0
// @BasePath /api
func New(db *gorm.DB) *flamego.Flame {
	f := flamego.Classic()

	f.Use(
		tracing.Middleware("go-template"),
		context.Contexter(db),
	)

	f.Group("/api", func() {

		userHandler := NewUserHandler()
		f.Group("/users", func() {
			f.Combo("").
				Get(userHandler.List).
				Post(form.Bind(form.CreateUser{}), userHandler.Create)
			f.Combo("/{user_uid}", userHandler.Userer).
				Get(userHandler.Get).
				Put(form.Bind(form.UpdateUser{}), userHandler.Update).
				Delete(userHandler.Delete)
		})
	})

	// HACK: /swagger is 404, redirect to /swagger/index.html
	f.Any("/swagger", func(ctx context.Context) { ctx.Redirect("/swagger/index.html") })
	f.Any("/swagger/{**}", flamegoswagger.WrapHandler(swaggerfiles.Handler))

	healthz.Set("version", appconst.BuildCommit)
	healthz.Register("postgres", 10*time.Second, func() error {
		return dbpkg.Ping(gocontext.Background())
	})
	f.Get("/healthz", healthz.Handler())

	return f
}
