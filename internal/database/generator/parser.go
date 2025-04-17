package main

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/generator/postgres"
)

func parsePostgresConnectionString(connStr string) (*postgres.DBConnection, error) {
	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	user := ""
	password := ""
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}

	host := u.Hostname()
	portStr := u.Port()
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	dbName := strings.TrimPrefix(u.Path, "/")

	components := &postgres.DBConnection{
		User:       user,
		Password:   password,
		Host:       host,
		Port:       port,
		DBName:     dbName,
		SchemaName: "public",
		SslMode:    "disable",
	}

	queryParams := u.Query()
	for key, value := range queryParams {
		if key == "schema" {
			components.SchemaName = value[0]
		} else if key == "sslmode" {
			components.SslMode = value[0]
		}
	}

	return components, nil
}
