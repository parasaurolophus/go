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

## Usage

```
$ go build ; ./automation
Usage of ./automation:
  -help
    	display usage and exit
  -hue
    	invoke Hue API
  -pv
    	invoke PowerView API
```

At least one of `-hue` and `-pv` are required. Invoking the program with no
arguments is equivalent to passing it `-help`. If both `-pv` and `-hue` are
supplied, the PowerView test is performed first, then the Hue test.

For `-pv`, the program uses `powerver.New` to create an interface to 