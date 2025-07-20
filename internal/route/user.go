// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/go-template/internal/context"
	"github.com/wuhan005/go-template/internal/db"
	"github.com/wuhan005/go-template/internal/dbutil"
	"github.com/wuhan005/go-template/internal/form"
	"github.com/wuhan005/go-template/internal/response"
)

// UserHandler is a struct that handles user-related routes.
type UserHandler struct{}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// List
// @Summary List users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.ListUser
// @Failure 500 "Internal server error" string
// @Router /users [get]
func (*UserHandler) List(ctx context.Context) error {
	users, total, err := db.Users.List(ctx.Request().Context(), db.ListUsersOptions{
		Pagination: dbutil.Pagination{
			Page:     ctx.QueryInt("page", 1),
			PageSize: ctx.QueryInt("pageSize", dbutil.DefaultPageSize),
		},
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to list users")
		return ctx.ServerError()
	}

	responseUsers := response.ConvertUsers(users)
	return ctx.Success(response.ListUser{
		Data:  responseUsers,
		Total: total,
	})
}

// Create
// @Summary Create a user
// @Accept json
// @Produce json
// @Param form body form.CreateUser true "User creation form"
// @Success 200 {object} response.User
// @Failure 500 "Internal server error" string
// @Router /users [post]
func (*UserHandler) Create(ctx context.Context, f form.CreateUser) error {
	user, err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Email:    f.Email,
		Password: f.Password,
		NickName: f.NickName,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create user")
		return ctx.ServerError()
	}

	responseUser := response.ConvertUser(user)
	return ctx.Success(responseUser)
}

func (*UserHandler) Userer(ctx context.Context) error {
	userUID := ctx.Param("user_uid")
	user, err := db.Users.GetByUID(ctx.Request().Context(), userUID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return ctx.Error(http.StatusNotFound, "User does not exist")
		}

		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get user")
		return ctx.ServerError()
	}

	ctx.Map(user)
	return nil
}

// Get
// @Summary Get user details
// @Produce json
// @Param user_uid path string true "User UID"
// @Success 200 {object} response.User
// @Failure 404 "User does not exist" string
// @Failure 500 "Internal server error" string
// @Router /users/{user_uid} [get]
func (*UserHandler) Get(ctx context.Context, user *db.User) error {
	responseUser := response.ConvertUser(user)
	return ctx.Success(responseUser)
}

// Update
// @Summary Update user details
// @Accept json
// @Produce json
// @Param user_uid path string true "User UID"
// @Param form body form.UpdateUser true "User update form"
// @Success 200 "User updated successfully" string
// @Failure 404 "User does not exist" string
// @Failure 500 "Internal server error" string
// @Router /users/{user_uid} [put]
func (*UserHandler) Update(ctx context.Context, user *db.User, f form.UpdateUser) error {
	if err := db.Users.Update(ctx.Request().Context(), user.ID, db.UpdateUserOptions{
		NickName: f.NickName,
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update user")
		return ctx.ServerError()
	}

	return ctx.Success("User updated successfully")
}

// Delete
// @Summary Delete a user
// @Produce json
// @Param user_uid path string true "User UID"
// @Success 200 "User deleted successfully" string
// @Failure 404 "User does not exist" string
// @Failure 500 "Internal server error" string
// @Router /users/{user_uid} [delete]
func (*UserHandler) Delete(ctx context.Context, user *db.User) error {
	if err := db.Users.Delete(ctx.Request().Context(), user.ID); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to delete user")
		return ctx.ServerError()
	}
	return ctx.Success("User deleted successfully")
}
