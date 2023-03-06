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
)

func TestServerConfigLoadFile(t *testing.T) {
	_, err := LoadFile("")
	if err == nil {
		t.Errorf("Did not get expected error when loading a file")
	}
	filename := "./test/does-not-exist.txt"
	_, err = LoadFile(filename)
	if err == nil {
		t.Errorf("Did not get expected error when loading a file: %v", filename)
	}
	filename = "./tests/allow_list.txt"
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
