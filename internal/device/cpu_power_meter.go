// SPDX-FileCopyrightText: 2025 The Kepler Authors
// SPDX-License-Identifier: Apache-2.0

package device

import "context"

// EnergyZone represents a measurable energy or power zone/domain exposed by a power meter.
// An EnergyZone typically represents a logical zone of the hardware unit, e.g. cpu core, cpu package
// dram, uncore etc.
// Reference: https://firefox-source-docs.mozilla.org/performance/power_profiling_overview.html
type EnergyZone interface {
	// Name() returns the zone name
	Name() string

	// Index() returns the index of the zone
	Index() int

	// Path() returns the path from which the energy usage value ie being read
	Path() string

	// Energy() returns energy consumed by the zone.
	Energy() Energy

	// MaxEnergy returns  the maximum value of energy usage that can be read.
	// When energy usage reaches this value, the energy value returned by Energy()
	// will wrap around and start again from zero.
	MaxEnergy() Energy
}

// CPUPowerMeter implements powerMeter
type CPUPowerMeter interface {
	powerMeter

	// Zones() returns a slice of the energy measurement zones
	Zones() ([]EnergyZone, error)
}

var _ CPUPowerMeter = (*cpuPowerMeter)(nil)

type cpuPowerMeter struct{}

func (c *cpuPowerMeter) Zones() ([]EnergyZone, error) {
	return nil, nil
}

func (c *cpuPowerMeter) Name() string {
	// TODO: set a proper name when rapl is implemented
	return "cpu"
}

func (c *cpuPowerMeter) Start(ctx context.Context) error {
	// TODO: Implement power monitoring logic
	return nil
}

func (c *cpuPowerMeter) Stop() error {
	// TODO: Implement stop logic
	return nil
}

func NewCPUPowerMeter() *cpuPowerMeter {
	return &cpuPowerMeter{}
}
