# dt
Go's missing DateTime package

## Why dt?

Go's standard library contains single date package - `time`. The type provided by it, `Time`, contains date, time and location information.

More often than not we don't location info, or we need to represent date/time only.

dt provides exactly that, a time-zone-independent representation of time that follows the rules of the proleptic Gregorian calendar with exactly 24-hour days, 60-minute hours, and 60-second minutes.

## What is provided?

dt provides three types to work with:

- Time: Contains time info in HH:mm
- Date: Contains date info: (YYYY-MM-DD)
- DateTime: Contains date and time information (YYYY-MM-DDTHH:mm)

Unlike `time.Time` these types contain an additional `Valid` field representing whether the data inside it was scanned.
This prevents situations like saving default date in database when nothing was received or responding via JSON with default date even though the date was empty.

## Why not civil package?

Google already offers something similar in [civil](https://github.com/googleapis/google-cloud-go/tree/master/civil) package.

- It's not an independent library, but a small package in a very big project which leads to its own problems.
- It doesn't implement the Scan/Value sql interfaces.
- It marshalls to zero date/time/datetime which is horrible (`time.Time` does this as well.) You can't differentiate inputed `00:00` and empty value.
- Slower development cycle

## License

dt is licensed under the Apache2 license. Check the [LICENSE](LICENSE) file for details.

## Author

[Emir Ribic](https://ribice.ba)