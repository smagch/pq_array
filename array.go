package pq_array

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type IntArray []int

// Scan implements the Scanner interface.
func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok || len(b) < 2 || b[0] != '{' || b[len(b)-1] != '}' {
		return fmt.Errorf("Invalid value: %v", value)
	}
	// empty array
	if len(b) == 2 {
		*a = make(IntArray, 0)
		return nil
	}
	nums := strings.Split(string(b[1:len(b)-1]), ",")
	*a = make(IntArray, len(nums))
	for i, s := range nums {
		var err error
		if (*a)[i], err = strconv.Atoi(s); err != nil {
			return err
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
// output format is like '{1,2,3,4}'
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, num := range a {
		if i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(strconv.Itoa(num))
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
