package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

var reComment = regexp.MustCompile(`(?m)^\s*//.*$`)

// Parse message payload
func Parse(r io.Reader) (Message, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Message{}, fmt.Errorf("unable to parse message: %w", err)
	}
	data = reComment.ReplaceAll(data, nil)

	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return Message{}, fmt.Errorf("unable to parse message: %w", err)
	}

	return message, nil
}
