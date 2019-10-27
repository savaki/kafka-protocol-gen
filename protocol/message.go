package protocol

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reValidVersions = regexp.MustCompile(`"(\d+)(-(\d+))?"`)
	reVersions      = regexp.MustCompile(`"(\d+)(-(\d+)|\+)?"`)
)

// Field represents a single field (or struct) with the kafka message
type Field struct {
	Name     string          `json:"name,omitempty"`     // Name of field
	Default  json.RawMessage `json:"default,omitempty"`  // Default value for field
	Type     string          `json:"type,omitempty"`     // Type of field
	Versions Versions        `json:"versions,omitempty"` // Versions field is compatible with
	About    string          `json:"about"`              // About
	Fields   []Field         `json:"fields,omitempty"`   // Fields for embedded type
}

// Message definition for kafka protocol as defined here,
// https://github.com/apache/kafka/tree/trunk/clients/src/main/resources/common/message
type Message struct {
	ApiKey           int           `json:"apiKey"`                  // ApiKey
	Type             string        `json:"type"`                    // Type of message; request or response
	Name             string        `json:"name"`                    // Name of message
	ValidVersions    ValidVersions `json:"validVersions"`           // ValidVersions contains set of valid message Versions
	FlexibleVersions string        `json:"flexibleVersions"`        // FlexibleVersions ... todo figure out what this means
	Fields           []Field       `json:"fields,omitempty"`        // Fields contained within Message
	CommonStructs    []Field       `json:"commonStructs,omitempty"` // CommonStructs used by the message
}

// ValidVersions parses the valid versions string into its semantic values
type ValidVersions struct {
	From int
	To   int
}

// UnmarshalJSON implements json.Unmarshaler
func (v *ValidVersions) UnmarshalJSON(data []byte) error {
	match := reValidVersions.FindSubmatch(data)
	if len(match) == 0 {
		return fmt.Errorf("unable to parse ValidVersions, %v", string(data))
	}

	from, err := strconv.Atoi(string(match[1]))
	if err != nil {
		return fmt.Errorf("unable to parse ValidVersions.  invalid from, %v", string(match[1]))
	}

	to := from
	if toStr := string(match[3]); len(toStr) > 0 {
		v, err := strconv.Atoi(toStr)
		if err != nil {
			return fmt.Errorf("unable to parse ValidVersions.  invalid to, %v", toStr)
		}
		to = v
	}

	*v = ValidVersions{
		From: from,
		To:   to,
	}

	return nil
}

// Versions provides a semantic interpretation of the versions field
type Versions struct {
	From        int  // From version
	To          int  // To contains max supported version; UpToCurrent takes precedence over To
	UpToCurrent bool // Versions up to current are supported
}

func (v Versions) IsValid(version int) bool {
	return version >= v.From && (version <= v.To || v.UpToCurrent)
}

func (v Versions) IsValidVersions(versions ValidVersions) bool {
	for version := versions.From; version <= versions.To; version++ {
		if v.IsValid(version) {
			return true
		}
	}
	return false
}

func (v Versions) String() string {
	if v.UpToCurrent {
		return strconv.Itoa(v.From) + "+"
	}
	return strconv.Itoa(v.From) + "-" + strconv.Itoa(v.To)
}

// UnmarshalJSON implements json.Unmarshaler
func (v *Versions) UnmarshalJSON(data []byte) error {
	match := reVersions.FindSubmatch(data)
	if len(match) == 0 {
		return fmt.Errorf("unable to parse ValidVersions, %v", string(data))
	}

	from, err := strconv.Atoi(string(match[1]))
	if err != nil {
		return fmt.Errorf("unable to parse ValidVersions.  invalid from, %v", string(match[1]))
	}

	var to int
	if toStr := string(match[3]); len(toStr) > 0 {
		v, err := strconv.Atoi(toStr)
		if err != nil {
			return fmt.Errorf("unable to parse ValidVersions.  invalid to, %v", toStr)
		}
		to = v
	}

	upToCurrent := string(match[2]) == "+"

	*v = Versions{
		From:        from,
		To:          to,
		UpToCurrent: upToCurrent,
	}

	return nil
}
