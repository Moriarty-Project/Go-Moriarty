package moriarty

import (
	"bufio"
	"bytes"
	"io"
	"sort"
)

// search function to search through an IO buffer for the following keyword, returns true if the keyword is found.
func SearchBufferFor(data *bufio.Reader, key []byte) bool {
	readData := make([]byte, len(key))
	for {
		// scan ahead until the first character of the key
		if _, err := data.ReadBytes(key[0]); err != nil {
			return false
		}
		data.UnreadByte()
		if _, err := data.Read(readData); err != nil {
			// we hit the end before it could begin
			return false
		}
		if bytes.Equal(key, readData) {
			// we found the key!
			return true
		}
		// reset the read data
		readData = make([]byte, len(key))
	}
}

func SearchBufferForAll(data io.Reader, keys [][]byte) bool {
	if len(keys) == 0 {
		return false
	}
	// find the key size ranges
	longest := len(keys[0])
	for _, key := range keys {
		if len(key) > longest {
			longest = len(key)
		}
	}
	// then we'll make our buffer, and scan over all the chars one at a time until we find a key.
	buffer := make([]byte, longest)
	_, err := data.Read(buffer)
	newByte := make([]byte, 1) //the buffer filler byte.
	for err == nil {
		// check the buffer for accuracy.
		for i := len(keys) - 1; i >= 0; i-- {
			// check if this key is in the buffer.
			key := keys[i]
			if bytes.Equal(key, buffer[:len(key)]) {
				// remove this key
				keys = append(keys[i+1:], keys[:i]...)
			}
		}
		// check if we're done
		if len(keys) == 0 {
			return true
		}
		// move the buffer forward one.
		_, err = data.Read(newByte)
		buffer = append(buffer[1:], newByte...)
	}
	// if all the keys were removed, then we found all of them. If any remain, then we missed something/
	return len(keys) == 0
}

func SearchBufferForAny(data io.Reader, keys [][]byte) bool {
	if len(keys) == 0 {
		return false
	}
	// we need to sort it to put them in decreasing sized order.
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	// then we'll make our buffer, and scan over all the chars one at a time until we find a key.
	buffer := make([]byte, len(keys[0]))
	data.Read(buffer) //we ignore the first error. So if the largest key is the same size as the whole file, we can still find it.
	var err error = nil
	newByte := make([]byte, 1) //the buffer filler byte.
	for err == nil {
		// check the buffer for accuracy.
		for _, key := range keys {
			// check if this key is in the buffer.
			if bytes.Equal(key, buffer[:len(key)]) {
				// we found a match!
				return true
			}
		}
		// move the buffer forward one.
		_, err = data.Read(newByte)
		buffer = append(buffer[1:], newByte...)
	}
	// if we hit an error, we need to shorten the buffer past all the keys.
	for len(buffer) > len(keys[len(keys)-1]) {
		for i := len(keys) - 1; i >= 0; i-- {
			key := keys[i]
			if len(key) > len(buffer) {
				// remove this key.
				keys = keys[i+1:]
			} else if bytes.Equal(key, buffer[:len(key)]) {
				// we found a match!
				return true
			}
		}
		buffer = buffer[1:]
	}
	// if we've gone through the whole buffer without finding anything, then we return false.
	return false
}
