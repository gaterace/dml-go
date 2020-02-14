// Copyright 2019 Demian Harvill
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package dml provides golang support for github.com/gaterace/xdml DmlExtension struct.
package dml

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"encoding/hex"

	"github.com/google/uuid"
	sdec "github.com/shopspring/decimal"
)

var validDecimal = regexp.MustCompile("^(\\+-)?\\d+(.\\d+)$")
var validGuid = regexp.MustCompile("^[0-9a-fA-F]{32}$")
var invalidDecimalError = errors.New("not a valid decimal string")
var invalidGuidError = errors.New("not a valid guid")

// Create a dml.DateTime structure from a string.
// String can be in format
//    YYYY-MM-DD HH:MM:SS
// or
//    YYYY-MM-DD
func DateTimeFromString(s string) *DateTime {

	var millis int64

	if (len(s) == 10) || (len(s) == 19) {
		year, err := strconv.Atoi(s[0:4])
		if err != nil {
			return nil
		}

		month, err := strconv.Atoi(s[5:7])
		if err != nil {
			return nil
		}

		day, err := strconv.Atoi(s[8:10])
		if err != nil {
			return nil
		}

		hour := 12
		min := 0
		sec := 0

		if len(s) == 19 {
			hour, err = strconv.Atoi(s[11:13])
			if err != nil {
				return nil
			}

			min, err = strconv.Atoi(s[14:16])
			if err != nil {
				return nil
			}

			sec, err = strconv.Atoi(s[17:19])
			if err != nil {
				return nil
			}
		}

		t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)

		millis = t.Unix() * 1000
	}

	result := &DateTime{Milliseconds: millis}

	return result
}

// Create a dml.Datetime structure from a time.Time value
func DateTimeFromTime(t time.Time) *DateTime {
	var millis int64

	millis = t.Unix() * 1000
	result := &DateTime{Milliseconds: millis}

	return result
}

// Create a time.Time value from a dml.DateTime instance.
func (m *DateTime) TimeFromDateTime() time.Time {
	millis := m.Milliseconds
	secs := millis / 1000

	t := time.Unix(secs, 0)

	return t
}

func (m *DateTime) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	t := m.TimeFromDateTime()
	tDisplay := t.Format(time.UnixDate)

	buffer.WriteString(fmt.Sprintf("\"timestamp\": \"%s\",", tDisplay))

	buffer.WriteString(fmt.Sprintf("\"milliseconds\": %d", m.Milliseconds))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// Get the string representation of a dml.Decimal instance.
func (m *Decimal) StringFromDecimal() string {
	return m.Plaintext
}

// Create a dml.Decimal struct from a string.
func DecimalFromString(s string) (*Decimal, error) {
	// TODO: validate string
	if validDecimal.MatchString(s) {
		result := &Decimal{Plaintext: s}
		return result, nil
	} else {

	}
	return nil, invalidDecimalError
}

// Convert dml.Decimal instance to shopspring decimal.Decimal
func (m *Decimal) ConvertDecimal() (sdec.Decimal, error) {
	result, err := sdec.NewFromString(m.Plaintext)
	return result, err
}

// Convert shopspring decimal.Decimal to dml.Decimal
func ConvertDecimal(d sdec.Decimal) *Decimal {
	result := &Decimal{Plaintext: d.String()}
	return result
}

// Convert dml.Guid to Uuid.
func (m *Guid) ConvertUuid() (uuid.UUID, error) {
	result, err := uuid.FromBytes(m.Guid)
	return result, err
}

// Convert Uuid to dml.Guid
func ConvertUuid(id uuid.UUID) *Guid {
	b := id[:]
	result := &Guid{Guid: b}
	return result
}

// Create  dml.Guid from byte slice
func GuidFromBytes(b []byte) (*Guid, error) {
	if len(b) != 16 {
		return nil, invalidGuidError
	}

	result := &Guid{Guid: b}
	return result, nil
}

// Create dml.Guid from string with hex representation
func GuidFromString(s string) (*Guid, error) {
	if !validGuid.MatchString(s) {
		return nil, invalidGuidError
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	result := &Guid{Guid: b}
	return result, nil
}

// Create a new dml.Guid
func NewGuid() *Guid {
	nid := uuid.New()
	result := ConvertUuid(nid)
	return result
}

func (m *Guid) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	hexstr := hex.EncodeToString(m.Guid)
	buffer.WriteString(fmt.Sprintf("\"guid\": \"%s\"", hexstr))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
