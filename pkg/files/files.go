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
	"bufio"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

// checkFile validates if the given file exists or not.
func checkFile(filename string) (bool, error) {
	if len(filename) == 0 {
		return false, fmt.Errorf("no file specified to check")
	}
	fileDetails, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("file doest no exist")
		}
		return false, err
	}
	if fileDetails.IsDir() {
		return false, fmt.Errorf("is a directory and not a file")
	}
	if fileDetails.Size() == 0 {
		return false, fmt.Errorf("specified file is empty")
	}
	return true, nil
}

// CheckTLS tests if TLS should be enabled or not.
func CheckTLS(certFile, keyFile string, s *zap.SugaredLogger) (bool, error) {
	var startTLS = false
	if len(certFile) > 0 && len(keyFile) > 0 {
		cf, err := checkFile(certFile)
		if err != nil || !cf {
			if err != nil {
				return false, err
			} else {
				if s != nil {
					s.Errorf("Cert file is not accessible: %v", certFile)
				}
				return false, fmt.Errorf("cert file not accesible: %v", certFile)
			}
		}
		kf, err := checkFile(keyFile)
		if err != nil || !kf {
			if s != nil {
				s.Errorf("Key file is not accessible: %v", keyFile)
			}
			if err != nil {
				return false, err
			} else {
				return false, fmt.Errorf("key file not accesible: %v", keyFile)
			}
		}
		startTLS = true
	}
	return startTLS, nil
}

// LoadFiltering loads the IP filtering options if available.
func LoadFiltering(allowListFile, denyListFile string, s *zap.SugaredLogger) ([]string, []string, error) {
	allowedIPs, err := loadListFile(allowListFile, "allow", s)
	if err != nil {
		return nil, nil, err
	}
	deniedIPs, err := loadListFile(denyListFile, "deny", s)
	if err != nil {
		return nil, nil, err
	}
	//var allowedIPs []string
	//if len(allowListFile) > 0 {
	//	cf, err := checkFile(allowListFile)
	//	if err != nil || !cf {
	//		if s != nil {
	//			s.Errorf("Allow List file is not accessible: %v", allowListFile)
	//		}
	//		if err != nil {
	//			return nil, nil, err
	//		} else {
	//			return nil, nil, fmt.Errorf("allow list file not accesible: %v", allowListFile)
	//		}
	//	}
	//	allowedIPs, err = LoadFile(allowListFile)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//}
	//var deniedIPs []string
	//if len(config.Filtering.DenyListFile) > 0 {
	//	cf, err := checkFile(config.Filtering.DenyListFile)
	//	if err != nil || !cf {
	//		zlog.S.Errorf("Deny List file is not accessible: %v", config.Filtering.DenyListFile)
	//		if err != nil {
	//			return nil, nil, err
	//		} else {
	//			return nil, nil, fmt.Errorf("deny list file not accesible: %v", config.Filtering.DenyListFile)
	//		}
	//	}
	//	deniedIPs, err = LoadFile(config.Filtering.DenyListFile)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//}
	return allowedIPs, deniedIPs, nil
}

// loadListFile loads the given file, parses it and returns its contents in a string array.
func loadListFile(listFile, name string, s *zap.SugaredLogger) ([]string, error) {
	var listFileContents []string
	if len(listFile) > 0 {
		cf, err := checkFile(listFile)
		if err != nil || !cf {
			if s != nil {
				s.Errorf("%s List file is not accessible: %v", name, listFile)
			}
			if err != nil {
				return nil, err
			} else {
				return nil, fmt.Errorf("%s list file not accesible: %v", name, listFile)
			}
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
