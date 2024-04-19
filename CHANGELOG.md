# Changelog

## 1.2.2

Bump go version to 1.22.

## 1.2.1

Bump go version to 1.20.
Proper handling of OPTION requests and HTTP 405 responses.
Default delay is set to 15 seconds and is now fixed to default correctly.

## 1.2.0

New configuration flags to disable metrics, so disk space can be saved.
New metric added `icon_room_connected`.

## 1.1.1

Stripping executables, removes all debug information to reduce executable size.

## 1.1.0

Bump go version to 1.18
Remove limit of 8kb for response data
Update thermostat ignore filter with to check both Live and Enabled thermostats

## 1.0.1

Fail faster in case of server is down

## 1.0.0

Initial release
