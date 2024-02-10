# auto-lookup

Extensible license plate and VIN number lookup service for the auto industry.

It can do vehicle lookups via license plate numbers or VIN numbers using simple
HTTP Get requests. The code base is designed to be extensible, ie. new services
can be added without changing the lookup algorithm.

## Usage ##

Syntax: `auto-lookup[.exe] dk|se regno|vin`

Examples:

- `auto-lookup dk BX71743`
- `auto-lookup se ONR701`

You can also use the word "test" in place of the license plate number in order
to test the service. It cannot be guaranteed to work with the plate numbers
used for testing, but this _will_ verify that the API connectivity is working.

A configuration file, `conf.yml`, containing the configuration of each service,
including an API token, must be located next to the executable.

## Implementation ##

The services themselves need to be implemented as a struct that satisfies
the `AutoService` interface, and the main program will call methods on each one
in turn in order to determine if a given service supports license plate lookups
for a given country.

## Changelog

_v1.2 (2024-02-10)_ - Use Go Modules. Update Makefile to build for Linux.

_v1.1 (2018-08-18)_ - Minor refactoring and code cleanup. Removed HTTP method.

_v1.0 (2018-08-16)_ - Initial version. Supports Danish and Swedish license plates.
