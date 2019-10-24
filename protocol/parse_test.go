package protocol

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/AddOffsetsToTxnRequest.json")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	got, err := Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	want := Message{
		ApiKey: 25,
		Type:   "request",
		Name:   "AddOffsetsToTxnRequest",
		ValidVersions: ValidVersions{
			To: 1,
		},
		FlexibleVersions: "none",
		Fields: []Field{
			{
				Name:     "TransactionalId",
				Type:     "string",
				Versions: Versions{UpToCurrent: true},
				About:    "The transactional id corresponding to the transaction.",
			},
			{
				Name:     "ProducerId",
				Type:     "int64",
				Versions: Versions{UpToCurrent: true},
				About:    "Current producer id in use by the transactional id.",
			},
			{
				Name:     "ProducerEpoch",
				Type:     "int16",
				Versions: Versions{UpToCurrent: true},
				About:    "Current epoch associated with the producer id.",
			},
			{
				Name:     "GroupId",
				Type:     "string",
				Versions: Versions{UpToCurrent: true},
				About:    "The unique group identifier.",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}
