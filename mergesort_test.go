//  Copyright (c) 2014 Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package mergesort

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readNewlineString(file *os.File, context interface{}) (interface{}, error) {
	// get current pos
	cur, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(file)
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// seek past the string we just read
	_, err = file.Seek(cur+int64(len(line)), os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return string(line[:len(line)-1]), nil
}

func writeNewlineString(file *os.File, record interface{}, context interface{}) error {
	str := record.(string)
	_, err := file.WriteString(str + "\n")
	if err != nil {
		return err
	}
	return nil
}

func compareStrings(rec1, rec2 interface{}, context interface{}) int {
	str1 := rec1.(string)
	str2 := rec2.(string)

	if str1 < str2 {
		return -1
	} else if str1 > str2 {
		return 1
	}
	return 0
}

func TestSmallString(t *testing.T) {
	unsortedFile, err := os.Open("test/unsorted_strings.txt")
	if err != nil {
		t.Fatal(err)
	}
	sortedFile, err := os.OpenFile("test/sorted_strings.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer sortedFile.Close()
	defer os.Remove(sortedFile.Name())
	err = MergeSort(unsortedFile, sortedFile, readNewlineString, writeNewlineString, compareStrings, nil, 10)
	if err != nil {
		t.Error(err)
	}

	sortedReader := bufio.NewReader(sortedFile)
	lastLine := ""
	line, err := sortedReader.ReadBytes('\n')
	for err == nil {
		if string(line) < lastLine {
			t.Errorf("out of order %s before %s", lastLine, line)
		}
		lastLine = string(line)
		line, err = sortedReader.ReadBytes('\n')
	}
}

func compareNumbers(rec1, rec2 interface{}, context interface{}) int {
	str1 := rec1.(string)
	str1 = strings.Trim(str1, "\n")
	num1, _ := strconv.ParseInt(str1, 10, 64)

	str2 := rec2.(string)
	str2 = strings.Trim(str2, "\n")
	num2, _ := strconv.ParseInt(str2, 10, 64)

	if num1 < num2 {
		return -1
	} else if num1 > num2 {
		return 1
	}
	return 0
}

func TestMediumNumbers(t *testing.T) {
	unsortedFile, err := os.Open("test/unsorted_numbers.txt")
	if err != nil {
		t.Fatal(err)
	}
	sortedFile, err := os.OpenFile("test/sorted_numbers.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer sortedFile.Close()
	defer os.Remove(sortedFile.Name())
	err = MergeSort(unsortedFile, sortedFile, readNewlineString, writeNewlineString, compareNumbers, nil, 64)
	if err != nil {
		t.Error(err)
	}

	sortedReader := bufio.NewReader(sortedFile)
	var lastNum int64 = 0
	line, err := sortedReader.ReadBytes('\n')
	for err == nil {
		strx := strings.Trim(string(line), "\n")
		numx, _ := strconv.ParseInt(strx, 10, 64)
		if numx < lastNum {
			t.Errorf("out of order %d before %d", lastNum, numx)
		}
		lastNum = numx
		line, err = sortedReader.ReadBytes('\n')
	}
}

func TestSameFile(t *testing.T) {
	// create a file with 100 random numbers
	unsortedFile, err := os.OpenFile("test/new_unsorted_file.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer unsortedFile.Close()
	defer os.Remove(unsortedFile.Name())
	for i := 0; i < 100; i++ {
		rnd := rand.Intn(100)
		_, err = unsortedFile.WriteString(strconv.Itoa(rnd) + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	unsortedFile.Seek(0, os.SEEK_SET)

	err = MergeSort(unsortedFile, unsortedFile, readNewlineString, writeNewlineString, compareNumbers, nil, 64)
	if err != nil {
		t.Error(err)
	}

	// now verify its sorted correctly
	sortedReader := bufio.NewReader(unsortedFile)
	var lastNum int64 = 0
	line, err := sortedReader.ReadBytes('\n')
	for err == nil {
		strx := strings.Trim(string(line), "\n")
		numx, _ := strconv.ParseInt(strx, 10, 64)
		if numx < lastNum {
			t.Errorf("out of order %d before %d", lastNum, numx)
		}
		lastNum = numx
		line, err = sortedReader.ReadBytes('\n')
	}
}
