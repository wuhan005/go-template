// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Transactor is an interface that defines a method for executing a transaction.
type Transactor interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}

// IsUniqueViolation checks if the given error is a unique constraint violation.
func IsUniqueViolation(err error, constraint string) bool {
	if err != nil {
		return strings.Contains(err.Error(), fmt.Sprintf("unique constraint %q", constraint))
	}
	return false
}
