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

// Package files provides some utilities for checking and loading files into memory
package files

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

// checkFile validates if the given file exists or not.
func checkFile(filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("no file specified to check")
	}
	fileDetails, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file doest no exist")
		}
		return err
	}
	if fileDetails.IsDir() {
		return fmt.Errorf("is a directory and not a file")
	}
	if fileDetails.Size() == 0 {
		return fmt.Errorf("specified file is empty")
	}
	return nil
}

// CheckTLS tests if TLS should be enabled or not.
func CheckTLS(certFile, keyFile string) (bool, error) {
	var startTLS = false
	if len(certFile) > 0 && len(keyFile) > 0 {
		err := checkFile(certFile)
		if err != nil {
			zlog.S.Errorf("Cert file is not accessible: %v", keyFile)
			return false, err
		}
		err = checkFile(keyFile)
		if err != nil {
			zlog.S.Errorf("Key file is not accessible: %v", keyFile)
			return false, err
		}
		startTLS = true
	}
	return startTLS, nil
}

// LoadFiltering loads the IP filtering options if available.
func LoadFiltering(allowListFile, denyListFile string) ([]string, []string, error) {
	// load the 'allow' list details
	var allowedIPs []string
	var deniedIPs []string
	var err error
	if len(allowListFile) > 0 {
		allowedIPs, err = loadListFile(allowListFile, "allow")
		if err != nil {
			return nil, nil, err
		}
	}
	// load the 'deny' list details
	if len(denyListFile) > 0 {
		deniedIPs, err = loadListFile(denyListFile, "deny")
		if err != nil {
			return nil, nil, err
		}
	}
	return allowedIPs, deniedIPs, nil
}

// loadListFile loads the given file, parses it and returns its contents in a string array.
func loadListFile(listFile, name string) ([]string, error) {
	var listFileContents []string
	if len(listFile) > 0 {
		err := checkFile(listFile)
		if err != nil {
			zlog.S.Errorf("%s List file is not accessible: %v", name, listFile)
			return nil, err
		}
		listFileContents, err = LoadFile(listFile)
		if err != nil {
			return nil, err
		}
	}
	return listFileContents, nil
}

// LoadFile loads the specified file and returns its contents in a string array.
func LoadFile(filename string) ([]string, error) {
	if len(filename) == 0 {
		return nil, fmt.Errorf("no file supplied to load")
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v - %v", filename, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var list []string
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			list = append(list, line)
		}
	}
	return list, nil
}
