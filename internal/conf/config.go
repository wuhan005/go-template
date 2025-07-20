// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var App struct {
	IpHeader string `envconfig:"IP_HEADER"`
}

var Postgres struct {
	DSN string `envconfig:"POSTGRES_DSN"`
}

var Redis struct {
	Address  string `envconfig:"REDIS_ADDRESS"`
	Username string `envconfig:"REDIS_USERNAME"`
	Password string `envconfig:"REDIS_PASSWORD"`
	Database int    `envconfig:"REDIS_DATABASE"`
}

var Tracing struct {
	Endpoint    string `envconfig:"TRACING_ENDPOINT"`
	Token       string `envconfig:"TRACING_TOKEN"`
	ServiceName string `envconfig:"TRACING_SERVICE_NAME"`
	HostName    string `envconfig:"HOSTNAME"`
}

// Init initializes the configuration by reading environment variables.
func Init() error {
	if err := envconfig.Process("", &App); err != nil {
		return errors.Wrap(err, "parse app")
	}
	if err := envconfig.Process("", &Postgres); err != nil {
		return errors.Wrap(err, "parse postgres")
	}
	if err := envconfig.Process("", &Redis); err != nil {
		return errors.Wrap(err, "parse redis")
	}
	if err := envconfig.Process("", &Tracing); err != nil {
		return errors.Wrap(err, "parse tracing")
	}

	return nil
}
