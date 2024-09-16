Copyright 2024 Kirk Rader

# Manual Integration Tests for Home Automation Interfaces

This console application sends the results of various "smart" device
interactions to `stdout`. It supports querying device configuration data from
either or both of a hard-coded collection of Philips Hue "bridges" and a
PowerView "hub." It does so by using the software interfaces implemented by the
[../automation/hue/README.md](../automation/hue/README.md) and
[../automation/powerview/README.md](../automation/powerview/README.md)
packages.

In addition, it demonstrates the use of
[../automation/trigger/](../automation/trigger/) to send asynchronous events at
specific times each day based on the sun's position on that date. Those events
can be used to trigger automations such as turning on lights at sunset and
turing them off again at sunrise.

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

```
$ ./automation_integration -help
Usage of /tmp/go-build2119991573/b001/exe/automation_integration:
  -bedtime int
        desired bedtime (0-23) (default 22)
  -help
        display usage and exit
```

When launched, it writes the output of various device API invocations and
asynchronous events to a file named `output.txt` until it is terminted by
pressing the return key. It also logs the operations of various goroutines to
`stdout`, while writing error messages to `stderr`.
