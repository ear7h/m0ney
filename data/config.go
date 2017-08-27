package data

import "os"

var mysqlCreds map[string]string

func init() {

	if os.Getenv("EAR7H_ENV") == "prod" {
		mysqlCreds = map[string]string{
			"user":     "root",
			"password": "",
			"host":     "db",
			"port":     "3306",
			"database": "stocks",
		}
	} else {
		mysqlCreds = map[string]string{
			"user":     "root",
			"password": "",
			"host":     "",
			"port":     "3306",
			"database": "stocks",
		}
	}
}
