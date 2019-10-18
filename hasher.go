// Copyright 2014 Spencer Kimball.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package cbfilter

import (
	"fmt"
	"hash"

	"github.com/spaolacci/murmur3"
)

// Hasher utilizes a 64 bit variant of MurmurHash3 to support
// the use of the Kirsch and MitzenMacher method to determine k
// hashed values using this formula:
//
// H[i](x) = hash[0:4] + i*hash[4:8]
//
// http://www.eecs.harvard.edu/~kirsch/pubs/bbbf/rsa.pdf
type hasher struct {
	mmh3   hash.Hash64
	hashed bool   // true if we've hashed a key
	h1     uint32 // first 4 bytes of key hash
	h2     uint32 // last 4 bytes of key hash
}

// newHasher allocates and return a new Hasher.
func newHasher() *hasher {
	return &hasher{mmh3: murmur3.New64()}
}

// HashKey writes the given key string to the hasher.
func (h *hasher) hashKey(key []byte) {
	h.mmh3.Reset() // clear current hash state
	if _, err := h.mmh3.Write(key); err != nil {
		panic(fmt.Sprintf("unable to write string to hasher: %v", key))
	}
	h.hashed = true
	sum := h.mmh3.Sum64()
	h.h1 = uint32(sum & 0xffffffff)
	h.h2 = uint32((sum >> 32) & 0xffffffff)
}

// GetHash returns the hash value at the given offset.
func (h *hasher) getHash(i uint32) uint32 {
	if !h.hashed {
		panic("hasher must be initialized first with a call to hashKey(key)")
	}
	return h.h1 + i*h.h2
}
