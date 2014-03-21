package mergesort

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readNewlineString(file *os.File, context interface{}) (interface{}, error) {
	cur, err := file.Seek(0, os.SEEK_CUR)
	r := bufio.NewReader(file)
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	file.Seek(cur+int64(len(line)), os.SEEK_SET)
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
