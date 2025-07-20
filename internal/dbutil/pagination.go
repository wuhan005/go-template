// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

// DefaultPageSize is the default number of items per page for pagination.
var DefaultPageSize = 20

// Pagination represents pagination parameters for database queries.
type Pagination struct {
	Page     int
	PageSize int
}

// Normalize ensures that the pagination parameters are valid.
func (p Pagination) Normalize() Pagination {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	return p
}

// LimitOffset returns the LIMIT and OFFSET parameters for SQL queries based on the pagination parameters.
func (p Pagination) LimitOffset() (limit, offset int) {
	return LimitOffset(p.Page, p.PageSize)
}

// LimitOffset returns LIMIT and OFFSET parameter for SQL.
// The first page is page 0.
func LimitOffset(page, pageSize int) (limit, offset int) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return pageSize, (page - 1) * pageSize
}
