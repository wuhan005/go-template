// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package response

import (
	"github.com/wuhan005/go-template/internal/db"
)

type User struct {
	UID      string `json:"uid"`
	Email    string `json:"email"`
	NickName string `json:"nickName"`
}

func ConvertUser(u *db.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		UID:      u.UID,
		Email:    u.Email,
		NickName: u.NickName,
	}
}

func ConvertUsers(users []*db.User) []*User {
	if users == nil {
		return nil
	}
	converted := make([]*User, len(users))
	for i, user := range users {
		user := user
		converted[i] = ConvertUser(user)
	}
	return converted
}

type ListUser struct {
	Data  []*User `json:"data"`
	Total int64   `json:"total"`
}
