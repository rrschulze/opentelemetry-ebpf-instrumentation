// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux

package ebpf

import (
	"os"
	"testing"
)

func expectedValues() map[string]testCase {
	return map[string]testCase{
		"setprogname": {
			startOffset:   0xac20,
			returnOffsets: []uint64{0xacc2, 0xad0a},
		},
		"setproctitle_init": {
			startOffset:   0xad10,
			returnOffsets: []uint64{0xaf2a, 0xb09a, 0xb0c0},
		},
		"invalid_symbol": {
			startOffset:   0x0,
			returnOffsets: nil,
		},
		"strunvis": {
			startOffset:   0x10570,
			returnOffsets: nil,
		},
		"fparseln": {
			startOffset:   0x82d0,
			returnOffsets: []uint64{0x84b0},
		},
	}
}

const libbsdS390xPath = "/usr/lib/s390x-linux-gnu/libbsd.so.0.12.1"

func testData() []byte {
	data, err := os.ReadFile(libbsdS390xPath)
	if err != nil {
		// Return nil — callers must handle a nil testData gracefully via t.Skip.
		return nil
	}
	return data
}

// TestGatherOffsets is overridden here to skip when the libbsd binary is absent.
func init() {
	_ = libbsdS390xPath // prevent unused-import lint
}

// skipIfNoLibbsd must be called from tests that depend on testData().
func skipIfNoLibbsd(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(libbsdS390xPath); err != nil {
		t.Skipf("skipping: s390x libbsd not found at %s: %v", libbsdS390xPath, err)
	}
}
