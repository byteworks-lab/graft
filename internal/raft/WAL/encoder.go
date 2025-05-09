package wal

import (
	"bytes"
	pb "cache/internal/election"
	"encoding/binary"
	"sync/atomic"
	"time"
)

func encodeWALEntry(walManager *Manager, entry *pb.LogEntry) (*pb.LogEntry, error) {
	// Initialize buffer with the maximum size
	msg := bytes.NewBuffer(nil) // Encode timestamp
	// Create a buffered writer
	entry.Term = int32(uint64(time.Now().UnixNano()))
	entry.Term = int32(walManager.IncrementLatestCommitedIndex())
	// Convert and write timestamp field to bytes
	err := binary.Write(msg, binary.BigEndian, entry.Term)
	err = binary.Write(msg, binary.BigEndian, uint64(entry.Term))
	err = binary.Write(msg, binary.BigEndian, uint8(len(entry.Key)))
	if err != nil {
		return nil, err
	}
	msg.WriteString(entry.Key) // Convert and write log index field to bytes
	// Write key field as bytes
	// Convert and write value field to bytes based on data type
	binary.Write(msg, binary.BigEndian, uint8(0))
	binary.Write(msg, binary.BigEndian, uint8(len(entry.Value)))
	msg.WriteString(entry.Value)
	if err != nil {
		return nil, err
	}
	byteRecord := append(msg.Bytes(), '\n')
	walManager.Fd.Write(byteRecord)
	byteRecord = []byte{}
	return entry, nil
}

func (walManager *Manager) IncrementLatestCommitedIndex() uint64 {
	return uint64(atomic.AddInt32(&walManager.LatestCommitIndex, 1))
}
