// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package appconst

var (
	// BuildCommit records the current Git Commit Hash, set at compile time using -ldflags
	BuildCommit string
	// BuildDate records the current Git Commit date, set at compile time using -ldflags
	BuildDate string
)
