// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package goexec // import "go.opentelemetry.io/obi/pkg/internal/goexec"

import (
	"golang.org/x/arch/s390x/s390xasm"
)

// FindReturnOffsets returns the byte offsets within data of every function
// return instruction relative to baseOffset.
//
// On s390x, function returns are encoded as BCR 15,%r14 (unconditional branch
// to the link register R14).  s390x instructions are 2, 4, or 6 bytes long;
// the Len field returned by s390xasm.Decode gives the exact size so we advance
// correctly over variable-length encodings.
func FindReturnOffsets(baseOffset uint64, data []byte) ([]uint64, error) {
	var returnOffsets []uint64
	index := 0
	for index < len(data) {
		inst, err := s390xasm.Decode(data[index:])
		if err != nil || inst.Len == 0 {
			// Unknown or truncated instruction; skip one byte to resync.
			index++
			continue
		}

		if isBCRReturn(inst) {
			returnOffsets = append(returnOffsets, baseOffset+uint64(index))
		}

		index += inst.Len
	}
	return returnOffsets, nil
}

// isBCRReturn reports whether inst is an unconditional branch to R14 (the
// s390x link register), i.e. BCR 15,%r14 (mnemonic: BR %r14).
func isBCRReturn(inst s390xasm.Inst) bool {
	if inst.Op != s390xasm.BCR {
		return false
	}
	mask, ok1 := inst.Args[0].(s390xasm.Mask)
	reg, ok2 := inst.Args[1].(s390xasm.Reg)
	return ok1 && ok2 && mask == 15 && reg == s390xasm.R14
}
