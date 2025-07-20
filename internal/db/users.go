// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"

	"github.com/wuhan005/go-template/internal/dbutil"
)

var _ UsersStore = (*users)(nil)

// Users is the default instance of the UsersStore.
var Users UsersStore

// UsersStore is the persistent interface for users.
type UsersStore interface {
	// Authenticate checks the user's email and password, returning the user if valid.
	// If the credentials are invalid, it returns ErrBadCredentials.
	Authenticate(ctx context.Context, email, password string) (*User, error)
	// Create creates a new user with the given options.
	Create(ctx context.Context, options CreateUserOptions) (*User, error)
	// List retrieves a list of users based on the provided options.
	List(ctx context.Context, options ListUsersOptions) ([]*User, int64, error)
	// GetByID retrieves a user by their ID.
	GetByID(ctx context.Context, id string) (*User, error)
	// GetByUID retrieves a user by their UID.
	GetByUID(ctx context.Context, uid string) (*User, error)
	// Update updates the user with the given ID using the provided options.
	Update(ctx context.Context, id uint, options UpdateUserOptions) error
	// Delete removes a user by its ID
	Delete(ctx context.Context, id uint) error
}

// NewUsersStore returns a UsersStore instance with the given database connection.
func NewUsersStore(db *gorm.DB) UsersStore {
	return &users{db}
}

type User struct {
	dbutil.Model
	Email    string
	Password string
	Salt     string
	NickName string
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if err := u.Model.BeforeCreate(db); err != nil {
		return errors.Wrap(err, "before create model")
	}

	u.Salt = randstr.String(10)
	u.EncodePassword()
	return nil
}

// EncodePassword hashes the user's password using PBKDF2 with SHA-256.
func (u *User) EncodePassword() {
	newPassword := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPassword)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(password string) bool {
	newUser := &User{Password: password, Salt: u.Salt}
	newUser.EncodePassword()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

type users struct {
	*gorm.DB
}

var ErrBadCredentials = errors.New("invalid email or password")

func (db *users) Authenticate(ctx context.Context, email, password string) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Model(&User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, ErrBadCredentials
	}

	if !user.ValidatePassword(password) {
		return nil, ErrBadCredentials
	}
	return &user, nil
}

type CreateUserOptions struct {
	Email    string
	Password string
	NickName string
}

func (db *users) Create(ctx context.Context, options CreateUserOptions) (*User, error) {
	newUser := &User{
		Email:    options.Email,
		Password: options.Password,
		NickName: options.NickName,
	}
	if err := db.WithContext(ctx).Create(&newUser).Error; err != nil {
		return nil, errors.Wrap(err, "create user")
	}
	return newUser, nil
}

type ListUsersOptions struct {
	dbutil.Pagination
}

func (db *users) List(ctx context.Context, options ListUsersOptions) ([]*User, int64, error) {
	query := db.WithContext(ctx).Model(&User{})

	// TODO: where clause for filtering

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count")
	}

	limit, offset := options.LimitOffset()

	var users []*User
	if err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, errors.Wrap(err, "find")
	}
	return users, count, nil
}

var ErrUserNotFound = errors.New("user does not exist")

func (db *users) getBy(ctx context.Context, where string, args ...interface{}) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Where(where, args...).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "get")
	}
	return &user, nil
}

func (db *users) GetByID(ctx context.Context, id string) (*User, error) {
	return db.getBy(ctx, "id = ?", id)
}

func (db *users) GetByUID(ctx context.Context, uid string) (*User, error) {
	return db.getBy(ctx, "uid = ?", uid)
}

type UpdateUserOptions struct {
	NickName string
}

func (db *users) Update(ctx context.Context, id uint, options UpdateUserOptions) error {
	return db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"nick_name": options.NickName,
		}).Error
}

func (db *users) Delete(ctx context.Context, id uint) error {
	return db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}
