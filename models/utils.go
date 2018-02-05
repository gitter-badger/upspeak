package models

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Postgres' JSONB type. It's a byte array of already encoded JSON (like json.RawMessage)
type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not string")
	}
	*j = append((*j)[0:0], s...)

	return nil
}

// MarshalJSON returns *m as the JSON encoding of m.
func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSONB) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSONB) Equals(j1 JSONB) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

// GeneratePasswordHash generates a hash for a given password string
func GeneratePasswordHash(passwd string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return string(hashByte), err
	}

	return string(hashByte), nil
}

// VerifyPasswordHash verifies a hash against its password
func VerifyPasswordHash(passwd string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	if err != nil {
		return err
	}

	return nil
}

// NullString adds JSON support for nullable string. It extends sql.NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (n *NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON  for NullString
func (n *NullString) UnmarshalJSON(data []byte) error {
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		n.Valid = true
		n.String = *str
	} else {
		n.Valid = false
	}

	return nil
}

// NullTime adds JSON support for a nullable time. It extends pq.NullTime
type NullTime struct {
	pq.NullTime
}

// MarshalJSON for NullTime
func (n *NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", n.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for NullTime
func (n *NullTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Time = t
	n.Valid = true
	return nil
}

// NullInt64 adds JSON support for a nullable int64. It extends sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}

// UnmarshalJSON for NullInt64
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = (err == nil)
	return err
}
