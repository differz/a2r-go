// Package data provides interfaces for data functions.
package data

import (
	"strings"
	"sync"
)

// Data is the common interface implemented by all data functions.
type Data interface {

	// Log writes any transactions
	//	Log()
	// Reset resets the Data to its initial state.
	//	Reset()
	// Size returns the number of bytes Sum will return.
	//	Size() int
}

// In is the common interface implemented by all Input data functions.
type In interface {
	Data
	Get()
	//	Sum32() uint32
}

// Out is the common interface implemented by all Output data functions.
type Out interface {
	Data
	//	Sum32() uint32
}

/////////////////////////////////////////////////////////////////////////////////

const (
	iNCOMING = "dongle-incoming"
)

// ProjectTable contain all phone & project numbers
type ProjectTable struct {
	Num       string
	CID       string
	ProjectID int64
}

var (
	pt     *[]ProjectTable
	ptOnce sync.Once
)

// NewProjectTable return pointer to global project table
func NewProjectTable() *[]ProjectTable {
	ptOnce.Do(func() {
		pt = &[]ProjectTable{}
	})
	return pt
}

// QueryRow is struct for take this -> rows.Scan(qRow.ValuePtrs...)
type QueryRow struct {
	Columns   []string
	Types     []string
	Values    []interface{}
	ValuePtrs []interface{}
}

// Init fill pointers and return string for select statement
func (q QueryRow) Init() string {
	count := len(q.Columns)
	for i := 0; i < count; i++ {
		q.ValuePtrs[i] = &q.Values[i]
	}
	return strings.Join(q.Columns, ", ")
}

// GetValue return interface of data by field name
func (q QueryRow) GetValue(name string) interface{} {
	count := len(q.Columns)
	exist := len(q.Types)
	for i := 0; i < count; i++ {
		if q.Columns[i] == name {
			if i < exist && q.Types[i] != "" {
				return string(q.Values[i].([]byte))
			}
			return q.Values[i]
		}
	}
	return nil
}

// GetString return string by field name
func (q QueryRow) GetString(name string) string {
	var ret string
	value := q.GetValue(name)
	if value != nil {
		ret = value.(string)
	}
	return ret
}

// GetInt64 return int64 by field name
func (q QueryRow) GetInt64(name string) int64 {
	var ret int64
	value := q.GetValue(name)
	if value != nil {
		ret = value.(int64)
	}
	return ret
}

func (q QueryRow) isIncoming() bool {
	value := q.GetValue("dcontext")
	return value == iNCOMING
}

func (q QueryRow) isOutgoing() bool {
	return !q.isIncoming()
}

// ClientNumber return phone number of client
func (q QueryRow) ClientNumber() string {
	flow := ""
	if q.isIncoming() {
		flow = "src"
	} else {
		flow = "dst"
	}

	value := q.GetString(flow)

	return strings.Replace(value, "+38", "", 1)
}
