package mergesort

type record []byte

type recordsList struct {
	records []record
	compare compareRecordsFunc
	context interface{}
}

func newRecordsList(expectedSize int, compare compareRecordsFunc, context interface{}) *recordsList {
	return &recordsList{
		records: make([]record, 0, expectedSize),
		compare: compare,
		context: context,
	}
}

func (rl *recordsList) add(r record) {
	rl.records = append(rl.records, r)
}

func (rl *recordsList) Len() int      { return len(rl.records) }
func (rl *recordsList) Swap(i, j int) { rl.records[i], rl.records[j] = rl.records[j], rl.records[i] }
func (rl *recordsList) Less(i, j int) bool {
	return rl.compare(rl.records[i], rl.records[j], rl.context) < 0
}
