package mergesort

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type tape struct {
	fp    *os.File
	count int
}

type readRecordFunc func(file *os.File, context interface{}) ([]byte, int, error)
type writeRecordFunc func(file *os.File, buf []byte, context interface{}) (int, error)
type compareRecordsFunc func(buf1, buf2 []byte, context interface{}) int

func MergeSort(unsortedFile, sortedFile *os.File, read readRecordFunc, write writeRecordFunc, compare compareRecordsFunc, context interface{}, blockSize int) error {
	var err error
	sourceTape := make([]tape, 2)
	record := make([][]byte, 2)

	// create temporary files sourceTape[0] and sourceTape[1]
	sourceTape[0].fp, err = ioutil.TempFile("", "goms")
	if err != nil {
		return err
	}
	defer os.Remove(sourceTape[0].fp.Name())
	sourceTape[1].fp, err = ioutil.TempFile("", "goms")
	if err != nil {
		return err
	}
	defer os.Remove(sourceTape[1].fp.Name())

	// read blocks, sort them in memory, and write the alternately to tapes 0 and 1
	blockCount := 0
	destination := 0
	list := newRecordsList(blockSize, compare, context)
	for {
		var recordSize int
		record[0], recordSize, err = read(unsortedFile, context)
		if err != nil && err != io.EOF {
			// error reading, return
			return err
		}
		if recordSize != 0 {
			// not EOF, add record to in memory list
			list.add(record[0])
			blockCount++
		}
		if blockCount == blockSize || recordSize == 0 && blockCount != 0 {
			// sort the in memory list
			sort.Sort(list)
			// now write them out
			for _, rec := range list.records {
				_, err := write(sourceTape[destination].fp, rec, context)
				if err != nil {
					return err
				}
				sourceTape[destination].count++
			}
			list = newRecordsList(blockSize, compare, context)
			destination ^= 1 // toggle tape
			blockCount = 0
		}
		if recordSize == 0 {
			break // all done
		}
	}
	if sortedFile == unsortedFile {
		unsortedFile.Seek(0, os.SEEK_SET)
	}
	sourceTape[0].fp.Seek(0, os.SEEK_SET)
	sourceTape[1].fp.Seek(0, os.SEEK_SET)

	// FIXME (what?) delete the unsorted file here, if required (see instructions)

	if sourceTape[1].count == 0 {
		// handle case where memory sort is all that is required
		err = sourceTape[1].fp.Close()
		if err != nil {
			return err
		}
		sourceTape[1] = sourceTape[0]
		sourceTape[0].fp = sortedFile
		for sourceTape[1].count != 0 {
			record[0], _, err = read(sourceTape[1].fp, context)
			if err != nil {
				return err
			}
			_, err := write(sourceTape[0].fp, record[0], context)
			if err != nil {
				return err
			}
			sourceTape[1].count--
		}
	} else {
		// merge tapes, two by two, until every record is in source_tape[0]
		for sourceTape[1].count != 0 {
			destination := 0
			destinationTape := make([]tape, 2)
			if sourceTape[0].count <= blockSize {
				destinationTape[0].fp = sortedFile
			} else {
				destinationTape[0].fp, err = ioutil.TempFile("", "goms")
				if err != nil {
					return err
				}
				defer os.Remove(destinationTape[0].fp.Name())
			}
			destinationTape[1].fp, err = ioutil.TempFile("", "goms")
			if err != nil {
				return err
			}
			defer os.Remove(destinationTape[1].fp.Name())
			record[0], _, err = read(sourceTape[0].fp, context)
			if err != nil {
				return err
			}
			record[1], _, err = read(sourceTape[1].fp, context)
			if err != nil {
				return err
			}
			for sourceTape[0].count != 0 {
				count := make([]int, 2)
				count[0] = sourceTape[0].count
				if count[0] > blockSize {
					count[0] = blockSize
				}
				count[1] = sourceTape[1].count
				if count[1] > blockSize {
					count[1] = blockSize
				}
				for count[0]+count[1] != 0 {
					sel := 0
					if count[0] == 0 {
						sel = 1
					} else if count[1] == 0 {
						sel = 0
					} else if compare(record[0], record[1], context) < 0 {
						sel = 0
					} else {
						sel = 1
					}
					_, err = write(destinationTape[destination].fp, record[sel], context)
					if err != nil {
						return err
					}
					if sourceTape[sel].count > 1 {
						record[sel], _, err = read(sourceTape[sel].fp, context)
						if err != nil {
							return err
						}
					}
					sourceTape[sel].count--
					count[sel]--
					destinationTape[destination].count++
				}
				destination ^= 1
			}
			sourceTape[0].fp.Close()
			sourceTape[1].fp.Close()
			destinationTape[0].fp.Seek(0, os.SEEK_SET)
			destinationTape[1].fp.Seek(0, os.SEEK_SET)
			// fixme memcmp?
			sourceTape[0] = destinationTape[0]
			sourceTape[1] = destinationTape[1]
			blockSize <<= 1
		}
	}
	sourceTape[1].fp.Close()
	return nil
}
