// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

// CreateUser is used for creating a new user.
type CreateUser struct {
	// Email is the user's email address.
	Email string `json:"email" valid:"required;email" label:"电子邮箱"`
	// Password is the user's password.
	Password string `json:"password" valid:"required" label:"密码"`
	// NickName is the user's nickname.
	NickName string `json:"nickName" valid:"required" label:"昵称"`
}

// UpdateUser is used for updating user information.
type UpdateUser struct {
	// NickName is the user's nickname.
	NickName string `json:"nickName" valid:"required" label:"昵称"`
}
