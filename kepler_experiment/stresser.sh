#!/bin/bash

# Check if stress-ng is installed
if ! command -v stress-ng &>/dev/null; then
	echo "stress-ng is not installed. Please install it (e.g., 'sudo apt install stress-ng' on Ubuntu)."
	exit 1
fi

# Get the number of CPU cores
NUM_CORES=$(nproc)
echo "Detected $NUM_CORES CPU cores."

# Function to run stress-ng with a specific load percentage
run_stress() {
	local load=$1
	local duration=20

	# If load is 0, skip running stress-ng
	if [ "$load" -eq 0 ]; then
		echo "No stress (0% load) for $duration seconds..."
		sleep $duration
		return
	fi

	# Run stress-ng with the specified load
	echo "Applying $load% CPU load for $duration seconds..."
	timeout --foreground $duration stress-ng --cpu "$NUM_CORES" --cpu-load "$load" --quiet &
	wait $! # Wait for the timeout command to complete
}

# Main loop to create the stress curve indefinitely
while true; do
	# Ascending: 0% to 100% in 10% increments
	for load in 0 10 20 30 40 50 60 70 80 90 100; do
		run_stress $load
	done

	# Descending: 100% to 0% in 10% increments
	for load in 90 80 70 60 50 40 30 20 10 0; do
		run_stress $load
	done
done
