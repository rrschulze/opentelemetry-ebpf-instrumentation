// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !s390x

package fastelf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFastElf_Data(t *testing.T) {
	ctx, err := NewElfContextFromData(libbsd_so_0_12_2)
	require.NoError(t, err)

	testFastElf(t, ctx)
}

func TestFastElf_File(t *testing.T) {
	filePath, cleanup, err := writeTempFile(libbsd_so_0_12_2)

	require.NoError(t, err)

	defer cleanup()

	ctx, err := NewElfContextFromFile(filePath)
	require.NoError(t, err)

	testFastElf(t, ctx)
}

func TestFastElf_FileNoSections(t *testing.T) {
	filePath, cleanup, err := writeTempFile(minimal_elf)

	require.NoError(t, err)

	defer cleanup()

	ctx, err := NewElfContextFromFile(filePath)
	require.NoError(t, err)

	require.Empty(t, ctx.Sections)
	require.Len(t, ctx.Segments, 4)
	require.False(t, ctx.HasSymbol("setprogname"))
	require.False(t, ctx.HasSection(".gnu_debuglink"))

	require.NoError(t, ctx.Close())
}
