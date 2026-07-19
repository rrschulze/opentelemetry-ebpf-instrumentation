// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package meta // import "go.opentelemetry.io/obi/pkg/appolly/meta"

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"

	attr "go.opentelemetry.io/obi/pkg/export/attributes/names"
)

// goarchToOTelArch maps Go GOARCH values to the OTel semconv host.arch strings.
// Only the GOARCH values that OBI can realistically run on are listed;
// unknown architectures are omitted so the attribute is simply absent.
var goarchToOTelArch = map[string]string{
	"amd64": string(semconv.HostArchAMD64.Value.Emit()),
	"arm":   string(semconv.HostArchARM32.Value.Emit()),
	"arm64": string(semconv.HostArchARM64.Value.Emit()),
	"ppc64": string(semconv.HostArchPPC64.Value.Emit()),
	"s390x": string(semconv.HostArchS390x.Value.Emit()),
	"386":   string(semconv.HostArchX86.Value.Emit()),
}

func linuxLocalFetcher(_ context.Context) (NodeMeta, error) {
	mid, err := fetchMachineID()
	if err != nil {
		// If we can't read host ID, we don't retry as it is mostly
		// (1) this linux distribution does not have the files where we are supposing
		// (2) there is some unrecoverable disk error
		// (3) we lack permissions
		// Then in this case, we only log a debug message
		slog.Debug("can't get local machine ID",
			"component", "meta.linuxLocalFetcher",
			"error", err)
	}
	nm := NodeMeta{HostID: mid}
	if archVal, ok := goarchToOTelArch[runtime.GOARCH]; ok {
		nm.Metadata = append(nm.Metadata, Entry{
			Key:   attr.Name(semconv.HostArchKey),
			Value: archVal,
		})
	}
	return nm, nil
}

func fetchMachineID() (string, error) {
	if result, err := os.ReadFile("/etc/machine-id"); err == nil && len(bytes.TrimSpace(result)) > 0 {
		return string(bytes.TrimSpace(result)), nil
	}

	if result, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil && len(bytes.TrimSpace(result)) > 0 {
		return string(bytes.TrimSpace(result)), nil
	} else {
		return "", fmt.Errorf("can't read host ID: %w", err)
	}
}
