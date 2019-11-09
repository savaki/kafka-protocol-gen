// Code generated by kafka-protocol-gen. DO NOT EDIT.
//
// Copyright 2019 Matt Ho
//
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

package sizeof

import "encoding/binary"

const (
	Bool        int32 = 1     // Bool bytes
	Int8        int32 = 1     // Int8 bytes
	Int16       int32 = 2     // Int16 size
	Int32       int32 = 4     // Int32 bytes
	Int64       int32 = 8     // Int64 bytes
	ArrayLength       = Int32 // Arraylength bytes e.g. Int32
)

// Bytes returns size of []byte
func Bytes(data []byte) int32 {
	return ArrayLength + int32(len(data)) // int32 length + length of bytes
}

// Int32Array returns size of []int32
func Int32Array(ii []int32) int32 {
	return ArrayLength + int32(len(ii))*Int32 // int32 length + length of array * int32 length
}

// Int64Array returns size of []int64
func Int64Array(ii []int64) int32 {
	return ArrayLength + int32(len(ii))*Int64 // int32 length + length of array * int64 length
}

// String returns size of string
func String(s string) int32 {
	return Int16 + int32(len(s))
}

// String returns size of string array
func StringArray(ss []string) int32 {
	var sz int32
	sz += ArrayLength // int32 length of array
	for _, s := range ss {
		sz += String(s)
	}
	return sz
}

// VarBytes returns the length of a var int
func VarBytes(data []byte) int32 {
	length := len(data)
	return VarInt(int64(length)) + int32(length)
}

// VarInt returns the length of a var int
func VarInt(i int64) int32 {
	var buf [16]byte
	length := binary.PutVarint(buf[:], i)
	return int32(length)
}

// VarString returns the length of a var string
func VarString(s string) int32 {
	length := len(s)
	return VarInt(int64(length)) + int32(length)
}
