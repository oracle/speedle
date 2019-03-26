//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package suid

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// SUID is the type of Speedle Unique ID
type SUID []byte

var sequenceID uint32
var hostname = getHostname()

var hashAlg = sha1.New()

// New returns a new SUID instance
func New() SUID {
	// Get current timestamp
	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb[:], uint64(time.Now().UnixNano()))
	pidb := make([]byte, 4)
	binary.BigEndian.PutUint32(pidb[:], uint32(os.Getpid()))
	seqb := make([]byte, 4)
	binary.BigEndian.PutUint32(seqb[:], atomic.AddUint32(&sequenceID, 1))

	hashAlg.Write(tb[:])
	hashAlg.Write(hostname)
	hashAlg.Write(pidb[:])
	hashAlg.Write(seqb[:])

	ret := make([]byte, 13)
	copy(ret[:], hashAlg.Sum(nil))
	return ret
}

func (s SUID) String() string {
	// Omit last char
	return strings.ToLower(base32.StdEncoding.EncodeToString(s)[:20])
}

func getHostname() []byte {
	hostName, err := os.Hostname()
	if err != nil {
		// This error should not happen
		panic(err)
	}
	return []byte(hostName)
}

// ParseSUID extracts the first SUID from the src string
func ParseSUID(src string) string {
	if len(src) == 0 {
		return ""
	}
	re := regexp.MustCompile(`([A-Za-z0-9]{20})`)
	match := re.FindStringSubmatch(src)
	if nil != match {
		return match[0]
	}

	log.Errorf("Failed to parse SUID from %s", src)
	return src
}
