# dt
Go's missing DateTime package

[![Build Status](https://travis-ci.org/ribice/dt.svg?branch=master)](https://travis-ci.org/ribice/dt)
[![codecov](https://codecov.io/gh/ribice/dt/branch/master/graph/badge.svg)](https://codecov.io/gh/ribice/dt)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/dt)](https://goreportcard.com/report/github.com/ribice/dt)

## Why dt?

Go's standard library contains a single date package - `time`. The type provided by it, `Time`, contains date, time and location information.

More often than not we don't need location info, or we need to represent date/time only.

dt provides exactly that, a time-zone-independent representation of time that follows the rules of the proleptic Gregorian calendar with exactly 24-hour days, 60-minute hours, and 60-second minutes.

## What is provided?

dt provides three types to work with:

- Time: Contains time info: HH:mm
- Date: Contains date info: YYYY-MM-DD
- DateTime: Contains date and time information: YYYY-MM-DDTHH:mm

Unlike `time.Time` these types contain an additional `Valid` field representing whether the data inside it was scanned/marshaled. This prevents situations like saving default date in a database when nothing was received or responding via JSON with default date even though the date was empty.

Types provided in dt represent sql types `time`, `date` and `timestamp`.

## Why not civil package?

Google already offers something similar in [civil](https://github.com/googleapis/google-cloud-go/tree/master/civil) package.

- It's not an independent library, but a small package in a very big project which leads to its problems.
- It doesn't implement the Scan/Value SQL interfaces.
- It marshalls to zero date/time/datetime (`time.Time` does this as well.) You can't differentiate inputted zero date/time/datetime and empty value.
- Slower development cycle

## License

dt is licensed under the Apache2 license. Check the [LICENSE](LICENSE) file for details.

## Author

[Emir Ribic](https://ribice.ba)
