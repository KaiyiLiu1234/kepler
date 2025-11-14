// SPDX-FileCopyrightText: 2025 The Kepler Authors
// SPDX-License-Identifier: Apache-2.0

package csvlogger

// CSV Logger only supports node cpu usage record. Need to include a feature for recording virtual machine power field

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

type CSVLogger struct {
	logger     *slog.Logger
	monitor    monitor.PowerDataProvider
	file       *os.File
	writer     *csv.Writer
	duration   time.Duration
	outputPath string
	startTime  time.Time
}

type options struct {
	logger     *slog.Logger
	duration   time.Duration
	outputPath string
}

type OptionFn func(*options)

func defaultOptions() options {
	return options{
		logger:     slog.Default(),
		duration:   5 * time.Minute, // Run for 5 minutes by default
		outputPath: "/tmp/kepler_cpu_usage.csv",
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

// NewCSVLogger create new CSV logger service
func NewCSVLogger(pm monitor.PowerDataProvider, opts ...OptionFn) *CSVLogger {
	options := defaultOptions()
	for _, apply := range opts {
		apply(&options)
	}
	return &CSVLogger{
		logger:     options.logger.With("service", "csvlogger"),
		monitor:    pm,
		duration:   options.duration,
		outputPath: options.outputPath,
	}
}

func (c *CSVLogger) Name() string {
	return "csvlogger"
}

func (c *CSVLogger) Init() error {
	c.logger.Info("Initializing CSV logger",
		"output", c.outputPath,
		"duration", c.duration,
	)

	// Create CSV file (will truncate)
	file, err := os.Create(c.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file at %s: %w", c.outputPath, err)
	}
	c.file = file
	c.writer = csv.NewWriter(file)

	// Write CSV headers
	header := []string{"timestamp", "cpu_usage_ratio"} // include cpu_delta_time?
	if err := c.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}
	c.writer.Flush()

	return nil
}

func (c *CSVLogger) Run(ctx context.Context) error {
	c.logger.Info("Starting CSV logger collection",
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
			c.logger.Info("CSV logger stopped due to context cancellation")
			return nil

		case <-durationTimer.C:
			c.logger.Info("CSV logger duration completed",
				"duration", c.duration,
				"time_elapsed", time.Since(c.startTime),
			)
			return nil

		case <-dataCh:
			elapsed := time.Since(c.startTime)

			if elapsed >= c.duration {
				c.logger.Info("CSV logger duration completed",
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

func (c *CSVLogger) recordSnapshot() error {
	snapshot, err := c.monitor.Snapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	// Extract Resource metrics from snapshot
	timestamp := snapshot.Timestamp.Format(time.RFC3339Nano)
	cpuUsageRatio := snapshot.Node.UsageRatio

	// Write to CSV
	record := []string{
		timestamp,
		fmt.Sprintf("%.6f", cpuUsageRatio),
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
		"cpu_usage_ratio", cpuUsageRatio,
		"time_elapsed", time.Since(c.startTime),
	)

	return nil
}

func (c *CSVLogger) Shutdown() error {
	c.logger.Info("Shutting down CSV logger",
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

	c.logger.Info("CSV logger fully shutdown")
	return nil
}

var _ service.Service = (*CSVLogger)(nil)
var _ service.Initializer = (*CSVLogger)(nil)
var _ service.Runner = (*CSVLogger)(nil)
var _ service.Shutdowner = (*CSVLogger)(nil)
