Copyright 2024 Kirk Rader

# Manual Integration Tests for Home Automation Interfaces

This console application sends the results of various interactions to `stdout`.
It supports querying device configuration data from either or both of a
hard-coded collection of Philips Hue "bridges" and a PowerView "hub." It does
by using the software interfaces implemented by the
[../hue/README.md](../hue/README.md) and
[../powerview/README.md](../powerview/README.md) packages.

## Environment

The program assumes the following environment variables are set:

| Variable                    | Description                                               |
|-----------------------------|-----------------------------------------------------------|
| `$GROUND_FLOOR_HUE_ADDRESS` | IP address or host name for the first of two Hue bridges  |
| `$GROUND_FLOOR_HUE_KEY`     | API security key for the first of two Hue bridges         |
| `$BASEMENT_HUE_ADDRESS`     | IP address or host name for the second of two Hue bridges |
| `$BASEMENT_HUE_KEY`         | API security key for the second of two Hue bridges        |
| `$POWERVIEW_ADDRESS`        | IP address or host name for a PowerView hub               |
| `$LATITUDE`                 | Latitude for sun position calculations                    |
| `$LONGITUDE`                | Longitude for sun position calculations                   |

## Usage

This application ignores any command-line arguments. When launched, it writes
the output of various device API invocations and asynchronous events to a file
named `output.txt` until it is terminted by pressing the return key. It also
logs the operations of various goroutines to `stdout`, while writing error
messages to `stderr`.
