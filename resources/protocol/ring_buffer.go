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

package protocol

import (
	"context"
	"sync/atomic"
)

// RingBuffer implements a ring buffer that supports a single reader and
// single writer.
type RingBuffer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	data      []byte        // data provides a temporary internal buffer
	start     int32         // start position to read from
	next      int32         // next position to write to
	lock      chan struct{} // lock enables signaling between read and write goroutines
	lockState int32         // lockState will either be 0 or 1
	size      int32         // size of internal buffer
}

// NewRingBuffer returns a new ringBuffer suitable for a single reader and a single writer
func NewRingBuffer(size int) *RingBuffer {
	ctx, cancel := context.WithCancel(context.Background())
	return &RingBuffer{
		ctx:    ctx,
		cancel: cancel,
		data:   make([]byte, size),
		lock:   make(chan struct{}),
		size:   int32(size),
	}
}

func (r *RingBuffer) Close() error {
	r.cancel()
	return nil
}

// WriteN writes the first n bytes from data to the buffer.  If buffer is full, WriteN
// blocks until ReadN clears space
func (r *RingBuffer) WriteN(data []byte, n int) {
	var (
		pos   = atomic.LoadInt32(&r.next)  // position next byte should be written
		start = atomic.LoadInt32(&r.start) // start of data
		wrote int                          // bytes written since last lock
	)

	for i := 0; i < n; i++ {
		r.data[pos] = data[i]
		pos++
		if pos == r.size {
			pos = 0
		}

		// whenever pos catches up to the starting position, wait for
		// additional bytes to be written before continuing
		if pos == start {
			// important to store pos prior to reading from lockMux
			atomic.StoreInt32(&r.next, pos)

			// release the read goroutine if its locked on us
			if atomic.CompareAndSwapInt32(&r.lockState, 1, 0) {
				wrote = i
				select {
				case <-r.ctx.Done():
					return
				case <-r.lock:
					// ok
				}
			}

			// lock this goroutine and wait for the read goroutine to unlock us
			if atomic.CompareAndSwapInt32(&r.lockState, 0, 1) {
				select {
				case <-r.ctx.Done():
					return
				case r.lock <- struct{}{}:
					// ok
				}
				start = atomic.LoadInt32(&r.start)
			}
		}
	}

	atomic.StoreInt32(&r.next, pos)

	// if at least one byte has been written since the last time r.lock
	// was cleared, check to see if it needs to be cleared again
	if n > wrote+1 {
		if atomic.CompareAndSwapInt32(&r.lockState, 1, 0) {
			select {
			case <-r.ctx.Done():
				return
			case <-r.lock:
				// ok
			}
		}
	}
}

// ReadN reads n bytes in the byte array provided.  If insufficient bytes
// are available, blocks until WriteN refills the buffer.
func (r *RingBuffer) ReadN(data []byte, n int) {
	var (
		pos  = atomic.LoadInt32(&r.start)
		next = atomic.LoadInt32(&r.next)
		read = 0
	)

	for i := 0; i < n; i++ {
		if pos == next {
			if read > 0 {
				atomic.StoreInt32(&r.start, pos)
				if atomic.CompareAndSwapInt32(&r.lockState, 1, 0) {
					select {
					case <-r.ctx.Done():
						return
					case <-r.lock:
						// ok
					}
					read = 0
				}
			}
			if atomic.CompareAndSwapInt32(&r.lockState, 0, 1) {
				select {
				case <-r.ctx.Done():
					return
				case r.lock <- struct{}{}:
					// ok
				}

				next = atomic.LoadInt32(&r.next)
			}
		}

		data[i] = r.data[pos]
		pos++
		if pos == r.size {
			pos = 0
		}
		read++
	}

	atomic.StoreInt32(&r.start, pos)
	if read > 0 {
		if atomic.CompareAndSwapInt32(&r.lockState, 1, 0) {
			select {
			case <-r.ctx.Done():
				return
			case <-r.lock:
				// ok
			}
		}
	}
}