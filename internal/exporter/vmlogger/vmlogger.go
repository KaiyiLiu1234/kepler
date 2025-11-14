// SPDX-FileCopyrightText: 2025 The Kepler Authors
// SPDX-License-Identifier: Apache-2.0

package vmlogger

// VM Logger only supports node cpu usage record. Need to include a feature for recording virtual machine power field

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sustainable-computing-io/kepler/internal/monitor"
	"github.com/sustainable-computing-io/kepler/internal/service"
)

type VMLogger struct {
	logger     *slog.Logger
	monitor    monitor.PowerDataProvider
	file       *os.File
	writer     *csv.Writer
	duration   time.Duration
	outputPath string
	vmID       string
	startTime  time.Time
}

type options struct {
	logger     *slog.Logger
	duration   time.Duration
	outputPath string
	vmID       string
}

type OptionFn func(*options)

func defaultOptions() options {
	return options{
		logger:     slog.Default(),
		duration:   5 * time.Minute, // Run for 5 minutes by default
		outputPath: "/tmp/kepler_vm_energy.csv",
		vmID:       "fedora40",
	}
}

func WithLogger(logger *slog.Logger) OptionFn {
	return func(o *options) {
		o.logger = logger
	}
}

func WithDuration(duration time.Duration) OptionFn {
	return func(o *options) {
		o.duration = duration
	}
}

func WithOutputPath(path string) OptionFn {
	return func(o *options) {
		o.outputPath = path
	}
}

func WithVMID(vmID string) OptionFn {
	return func(o *options) {
		o.vmID = vmID
	}
}

// NewVMLogger create new VM logger service
func NewVMLogger(pm monitor.PowerDataProvider, opts ...OptionFn) *VMLogger {
	options := defaultOptions()
	for _, apply := range opts {
		apply(&options)
	}
	return &VMLogger{
		logger:     options.logger.With("service", "vmlogger"),
		monitor:    pm,
		duration:   options.duration,
		outputPath: options.outputPath,
		vmID:       options.vmID,
	}
}

func (c *VMLogger) Name() string {
	return "vmlogger"
}

func (c *VMLogger) Init() error {
	c.logger.Info("Initializing VM logger",
		"output", c.outputPath,
		"duration", c.duration,
		"vmID", c.vmID,
	)

	// Create CSV file (will truncate)
	file, err := os.Create(c.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file at %s: %w", c.outputPath, err)
	}
	c.file = file
	c.writer = csv.NewWriter(file)

	// Write CSV headers
	header := []string{"timestamp", "vm_core_power_w", "vm_core_energy_j"}
	if err := c.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}
	c.writer.Flush()

	return nil
}

func (c *VMLogger) Run(ctx context.Context) error {
	c.logger.Info("Starting VM logger collection",
		"duration", c.duration,
		"note", "using monitor's configured interval for data collection",
	)

	c.startTime = time.Now()
	dataCh := c.monitor.DataChannel()

	durationTimer := time.NewTimer(c.duration)
	defer durationTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("VM logger stopped due to context cancellation")
			return nil

		case <-durationTimer.C:
			c.logger.Info("VM logger duration completed",
				"duration", c.duration,
				"time_elapsed", time.Since(c.startTime),
			)
			return nil

		case <-dataCh:
			elapsed := time.Since(c.startTime)

			if elapsed >= c.duration {
				c.logger.Info("VM logger duration completed",
					"duration", c.duration,
					"time_elapsed", elapsed,
				)
				return nil
			}

			if err := c.recordSnapshot(); err != nil {
				c.logger.Error("Failed to record snapshot", "error", err)
				// Continue on Error
			}
			// Continue after recording Snapshot
		}
	}
}

func (c *VMLogger) recordSnapshot() error {
	snapshot, err := c.monitor.Snapshot()
	vmID := c.vmID
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	// Extract Resource metrics from snapshot
	timestamp := snapshot.Timestamp.Format(time.RFC3339Nano)

	// Extract VM Energy metrics from snapshot
	var vmName string
	var vmCorePowerW float64
	var vmCoreEnergyJ float64

	// print out virtual machines
	for key := range snapshot.VirtualMachines {
		c.logger.Info("recording virtual machines", "key", key)
	}

	vm := snapshot.VirtualMachines[vmID]
	for zone, usage := range vm.Zones {
		if zone.Name() == "core" {
			vmCorePowerW = usage.Power.Watts()
			vmCoreEnergyJ = usage.EnergyDelta.Joules()
			break // we only need core for experiment
		}
	}

	c.logger.Info("Recording VM snapshot",
		"vm_name", vmName,
		"vm_id", vmID,
		"snapshot_timestamp", timestamp,
		"node_timestamp", snapshot.Node.Timestamp.Format(time.RFC3339Nano),
		"vm_core_power_w", vmCorePowerW,
		"vm_core_energy_j", vmCoreEnergyJ,
	)

	// Write to CSV
	record := []string{
		timestamp,
		fmt.Sprintf("%.6f", vmCorePowerW),
		fmt.Sprintf("%.6f", vmCoreEnergyJ),
	}
	if err := c.writer.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV record: %w", err)
	}
	c.writer.Flush()

	if err := c.writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	c.logger.Debug("Recorded CPU usage",
		"timestamp", timestamp,
		"vm_core_power_w", vmCorePowerW,
		"vm_core_energy_j", vmCoreEnergyJ,
		"time_elapsed", time.Since(c.startTime),
	)

	return nil
}

func (c *VMLogger) Shutdown() error {
	c.logger.Info("Shutting down VM logger",
		"total_duration", time.Since(c.startTime),
	)

	if c.writer != nil {
		c.writer.Flush()
	}

	if c.file != nil {
		if err := c.file.Close(); err != nil {
			return fmt.Errorf("failed to close CSV file: %w", err)
		}
	}

	c.logger.Info("VM logger fully shutdown")
	return nil
}

var _ service.Service = (*VMLogger)(nil)
var _ service.Initializer = (*VMLogger)(nil)
var _ service.Runner = (*VMLogger)(nil)
var _ service.Shutdowner = (*VMLogger)(nil)
