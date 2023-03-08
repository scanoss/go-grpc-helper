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
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"golang.org/x/net/context"
)

func TestOpenDBConnectionSqLite(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()

	db, err := OpenDBConnection(":memory:", "sqlite3", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	CloseDBConnection(db)
}

func TestOpenDBConnectionFail(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	db, err := OpenDBConnection("", "does-not-exist", "nouser", "nopwd", "does-not-exit", "noschema", "")
	if err == nil {
		CloseDBConnection(db)
		t.Errorf("Expected to get an error")
	}
}

func TestSetDBOptionsAndPing(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	db, err := OpenDBConnection(":memory:", "sqlite3", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = SetDBOptionsAndPing(db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	CloseDBConnection(db)
	err = SetDBOptionsAndPing(db)
	if err == nil {
		t.Errorf("Expected to get an error")
	}
}

func TestCloseSQLConnection(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	db, err := OpenDBConnection(":memory:", "sqlite3", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer CloseDBConnection(db)
	conn, err := db.Connx(context.Background()) // Get a connection from the pool
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	CloseSQLConnection(conn) // should pass first time
	CloseSQLConnection(conn) // should fail second time
}
