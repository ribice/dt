// Package dt implements types for normal time, a time-zone-independent
// representation of time that follows the rules of the proleptic
// Gregorian calendar with exactly 24-hour days, 60-minute hours, and 60-second
// minutes.
//
// Because they lack location information, these types do not represent unique
// moments or intervals of time. Use time.Time for that purpose.
package dt
