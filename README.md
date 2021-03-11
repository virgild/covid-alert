# covid-alert

This Go program scrapes the St. Joseph's Unity Health Toronto page for Covid-19 figures,
and saves the results to a local SQLite3 database file. It also alerts you via SMS when 
the figures have changed.

# Build process

## Build locally
```
make build/covid-alert
```

## Build for Raspberry Pi 3

### Other requirements
* Docker

```
# Build the Docker image capable of building the program binary for Raspberry Pi 3:
make covid-alert-builder

# Build the program binary for Raspberry Pi 3:
make build-arm
```
