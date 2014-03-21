package mergesort

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readNewlineString(file *os.File, context interface{}) ([]byte, int, error) {
	curPos, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, 0, err
	}
	r := bufio.NewReader(file)
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, 0, err
	}
	file.Seek(curPos+int64(len(line)), os.SEEK_SET)
	return line, len(line), nil
}

func writeNewlineString(file *os.File, buf []byte, context interface{}) (int, error) {

	_, err := file.WriteString(string(buf))
	if err != nil {
		return 0, err
	}
	return len(buf) + 1, nil
}

func compareStrings(buf1, buf2 []byte, context interface{}) int {
	str1 := string(buf1)
	str2 := string(buf2)

	if str1 < str2 {
		return -1
	} else if str1 > str2 {
		return 1
	}
	return 0
}

func TestSmallString(t *testing.T) {
	unsortedFile, err := os.Open("unsorted_strings.txt")
	if err != nil {
		t.Fatal(err)
	}
	sortedFile, err := os.OpenFile("sorted_strings.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer sortedFile.Close()
	defer os.Remove(sortedFile.Name())
	err = MergeSort(unsortedFile, sortedFile, readNewlineString, writeNewlineString, compareStrings, nil, 2)
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

func compareNumbers(buf1, buf2 []byte, context interface{}) int {
	str1 := string(buf1)
	str1 = strings.Trim(str1, "\n")
	num1, _ := strconv.ParseInt(str1, 10, 64)

	str2 := string(buf2)
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
	unsortedFile, err := os.Open("unsorted_numbers.txt")
	if err != nil {
		t.Fatal(err)
	}
	sortedFile, err := os.OpenFile("sorted_numbers.txt", os.O_RDWR|os.O_CREATE, 0666)
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
