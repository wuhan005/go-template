// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Model is a base model for GORM that includes common fields such as ID, UID, CreatedAt, UpdatedAt, and DeletedAt.
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UID       string         `gorm:"uniqueIndex" json:"uid"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook that sets the UID field to a new unique ID if it is not already set.
func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.UID == "" {
		m.UID = xid.New().String()
	}
	return nil
}
