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
	"context"
	"regexp"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var sqlRegex = regexp.MustCompile(`\$\d+`) // regex to check for SQL parameters

type DBQueryContext struct {
	conn  *sqlx.Conn
	s     *zap.SugaredLogger
	trace bool
}

// NewDBSelectContext creates a new instance of the DBQueryContext service.
func NewDBSelectContext(s *zap.SugaredLogger, conn *sqlx.Conn, trace bool) *DBQueryContext {
	return &DBQueryContext{s: s, conn: conn, trace: trace}
}

// SelectContext logs the give query before executing it and the result afterward, if tracing is enabled?
func (q *DBQueryContext) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if q.trace {
		q.SQLQueryTrace(query, args...)
	}
	err := q.conn.SelectContext(ctx, dest, query, args...)
	if err == nil && q.trace {
		q.SQLResultsTrace(dest)
	}
	return err
}

// SQLQueryTrace logs the given SQL query if debug is enabled.
func (q *DBQueryContext) SQLQueryTrace(query string, args ...interface{}) {
	q.s.Debugf("SQL Query: "+sqlRegex.ReplaceAllString(query, "%v"), args...)
}

// SQLResultsTrace logs the given SQL result if debug is enabled.
func (q *DBQueryContext) SQLResultsTrace(results interface{}) {
	q.s.Debugf("SQL Results: %#v", results)
}
