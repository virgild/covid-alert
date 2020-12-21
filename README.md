# covid-alert

This program scrapes the St. Joseph's Unity Health Toronto page for Covid-19 figures,
and saves the results to a local SQLite3 database file. It also alerts you via SMS when 
the figures have changed.

# Build process

Build locally
```
make build/covid-alert
```

Build for Raspberry Pi 3
```
# Build the Docker image capable of building the Raspberry Pi 3 image:
make covid-alert-builder

# Build the program binary for Raspberry Pi 3:
make build/covid-alert-arm
```
