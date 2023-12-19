# Planet Exporter

Correlate production issues with planet constellations using Prometheus. Do not let the five whys stop you from finding
the root cause and go right to the cause!

### Usage

#### Run using Docker

```
docker run -p 9042:9042 ghcr.io/kozaktomas/planet-exporter:main
curl localhost:9042/metrics
```

#### Alerts

todo

### Metrics example

```
# HELP distance_between_objects Current distance between objects in the space in meters
# TYPE distance_between_objects gauge
distance_between_objects{from="earth",to="jupiter"} 6.47142851255e+11
distance_between_objects{from="earth",to="mars"} 3.68829062259e+11
distance_between_objects{from="earth",to="moon"} 3.7115189771399796e+08
distance_between_objects{from="earth",to="neptune"} 4.490103341003e+12
distance_between_objects{from="earth",to="saturn"} 1.518308348414e+12
# ....
```

### Development

Make sure you download data files. You can use `make download` to download all files to data folder. Then you have to
set `VSOP87` environment variable to point to the folder with datafiles: `export VSOP87=$(pwd)/data`. Now you are good
to go.
