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
	"testing"

	"golang.org/x/net/context"

	_ "github.com/lib/pq"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	_ "modernc.org/sqlite"
)

type Persons struct {
	FirstName string
	LastName  string
}

func TestQuerySQLite(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	db, err := OpenDBConnection(":memory:", "sqlite", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer CloseDBConnection(db)
	db.MustExec("CREATE TABLE person (firstname text, lastname text)")
	db.MustExec("INSERT INTO person (firstname, lastname) VALUES ('harry', 'potter')")
	ctx := context.Background()

	q1 := NewDBSelectContext(zlog.S, db, nil, true)
	var results1 []Persons
	err = q1.SelectContext(ctx, &results1, "SELECT * FROM person where firstname = $1", "harry")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Results1: %v\n", results1)

	conn, err := db.Connx(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	q2 := NewDBSelectContext(zlog.S, db, conn, true)
	var results2 []Persons
	err = q2.SelectContext(ctx, &results2, "SELECT * FROM person where firstname = $1", "harry")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Results2: %v\n", results2)
}
