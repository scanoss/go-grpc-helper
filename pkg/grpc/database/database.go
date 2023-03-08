// SPDX-License-Identifier: MIT
/*
 * Copyright (c) 2023, SCANOSS
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

// OpenDBConnection establishes a connection specified database.
func OpenDBConnection(dsn, driver, user, passwd, host, schema, sslMode string) (*sqlx.DB, error) {
	if len(dsn) == 0 {
		dsn = fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s",
			driver,
			user,
			passwd,
			host,
			schema,
			sslMode)
	}
	zlog.S.Debug("Connecting to Database...")
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		zlog.S.Errorf("Failed to open database: %v", err)
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	return db, nil
}

// SetDBOptionsAndPing configures DB connections and attempts to ping it.
func SetDBOptionsAndPing(db *sqlx.DB) error {
	db.SetConnMaxIdleTime(30 * time.Minute) // TODO add to app config
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	err := db.Ping()
	if err != nil {
		zlog.S.Errorf("Failed to ping database: %v", err)
		return fmt.Errorf("failed to ping database: %v", err)
	}
	return nil
}

// CloseDBConnection closes specified database connection.
func CloseDBConnection(db *sqlx.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			zlog.S.Warnf("Problem closing DB: %v", err)
		}
	}
}

// CloseSQLConnection closes the specified database connection.
func CloseSQLConnection(conn *sqlx.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			zlog.S.Warnf("Warning: Problem closing database connection: %v", err)
		}
	}
}
