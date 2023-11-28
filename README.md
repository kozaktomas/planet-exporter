# Planet Exporter

Devastating production outages are usually caused by unexpected planetary alignments. This tool will help you to avoid
them. It exports current positions of all planets in our solar system. From our experience, it's usually related to the
certain distance between two planets or between a planets and the Sun. The exporter exports the distance between all
planets and the Sun. Also, the Moon is sometimes involved - so we included distance to the Earth as well. Now you can
create/amend alerts based on those metrics to avoid unexpected outages.

### Usage

#### Run using Docker

```
docker run -p 9042:9042 docker pull ghcr.io/kozaktomas/planet-exporter:main
curl localhost:9042/metrics
```

#### Alerts

todo

### Metrics example

```
# HELP distance_between_objects Current distance between objects in the space in kilometers
# TYPE distance_between_objects gauge
distance_between_objects{from="earth",to="jupiter"} 6.13527435e+08
distance_between_objects{from="earth",to="mars"} 3.76274473e+08
distance_between_objects{from="earth",to="moon"} 384925.0368589994
distance_between_objects{from="earth",to="neptune"} 4.435408063e+09
distance_between_objects{from="earth",to="saturn"} 1.467320168e+09
distance_between_objects{from="earth",to="uranus"} 2.799938039e+09
# ....
```

### Development

Make sure you download data files. You can use `make download` to download all files to data folder. Then you have to
set `VSOP87` environment variable to point to the folder with datafiles: `export VSOP87=$(pwd)/data`. Now you are good
to go.
