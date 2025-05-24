package wal

type Command int

// Declare constants using iota
const (
	SET Command = iota
	GET
	DELETE
)

type WALLogEntry struct {
	Comm      Command `json:"Command"`
	Key       string  `json:"Key"`
	Value     any     `json:"Value"`
	Timestamp uint64  `json:"Timestamp"`
	LogIndex  uint64  `json:"LogIndex"`
}

func NewWALLogEntry(comm Command, key string, value any, commitedLogIndex uint64, clientTerm uint64) WALLogEntry {
	return WALLogEntry{
		Comm:     comm,
		Key:      key,
		Value:    value,
		LogIndex: commitedLogIndex,
	}
}
