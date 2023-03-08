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

package files

import (
	"fmt"
	"testing"

	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCheckFile(t *testing.T) {
	err := checkFile("")
	if err == nil {
		t.Errorf("[1] expected to get an error")
	}
	err = checkFile("../../tests/empty-file.txt")
	if err == nil {
		t.Errorf("[2] expected to get an error")
	}
	err = checkFile("../../tests")
	if err == nil {
		t.Errorf("[3] expected to get an error")
	}
	err = checkFile("../../tests/does-not-exist.txt")
	if err == nil {
		t.Errorf("[4] expected to get an error")
	}
	err = checkFile("../../tests/server.crt")
	if err != nil {
		t.Errorf("[5] unexpected error: %v", err)
	}
}

func TestCheckTLS(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	startTLS, err := CheckTLS("", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.False(t, startTLS, "[1] should be false")

	startTLS, err = CheckTLS("../../tests/empty-file.txt", "../../tests/empty-file.txt")
	if err == nil {
		t.Errorf("[2] should've caused an error")
	}
	assert.False(t, startTLS, "[2] should be false")

	startTLS, err = CheckTLS("../../tests/server.crt", "../../tests/empty-file.txt")
	if err == nil {
		t.Errorf("[3] should've caused an error")
	}
	assert.False(t, startTLS, "[3] should be false")

	startTLS, err = CheckTLS("../../tests/server.crt", "../../tests/server.key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.True(t, startTLS, "[4] should be true")
}

func TestServerConfigLoadFile(t *testing.T) {
	_, err := LoadFile("")
	if err == nil {
		t.Errorf("Did not get expected error when loading a file")
	}
	filename := "../../tests/does-not-exist.txt"
	_, err = LoadFile(filename)
	if err == nil {
		t.Errorf("Did not get expected error when loading a file: %v", filename)
	}
	filename = "../../tests/allow_list.txt"
	res, err := LoadFile(filename)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when loading a data file", err)
	}
	if len(res) == 0 {
		t.Errorf("No data decoded from data file: %v", filename)
	} else {
		fmt.Printf("Data File details: %+v\n", res)
	}
}

func TestLoadFiltering(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	allowedIPs, deniedIPs, err := LoadFiltering("", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.True(t, allowedIPs == nil && deniedIPs == nil, "[1] should be nil")

	allowedIPs, deniedIPs, err = LoadFiltering("../../tests/does-not-exist.txt", "../../tests/does-not-exist.txt")
	if err == nil {
		t.Errorf("[2] should've caused an error")
	}
	assert.True(t, allowedIPs == nil && deniedIPs == nil, "[2] should be nil")

	allowedIPs, deniedIPs, err = LoadFiltering("../../tests/allow_list.txt", "../../tests/does-not-exist.txt")
	if err == nil {
		t.Errorf("[3] should've caused an error")
	}
	assert.True(t, allowedIPs == nil && deniedIPs == nil, "[3] should be nil")

	allowedIPs, deniedIPs, err = LoadFiltering("../../tests/does-not-exist.txt", "../../tests/deny_list.txt")
	if err == nil {
		t.Errorf("[4] should've caused an error")
	}
	assert.True(t, allowedIPs == nil && deniedIPs == nil, "[3] should be nil")

	allowedIPs, deniedIPs, err = LoadFiltering("../../tests/allow_list.txt", "../../tests/deny_list.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("Allowed IP: %v, Denied: %v", allowedIPs, deniedIPs)
	assert.True(t, allowedIPs != nil && deniedIPs != nil, "[5] should not be nil")
}
