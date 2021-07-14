# Telegraf Input Plug-in for Piaware

This project is an application that pulls data from a
[Piaware](https://flightaware.com/adsb/piaware/)
ADS-B 
receiver and outputs it in a format suitable for ingestion by
[Telegraf](https://www.influxdata.com/time-series-platform/telegraf/),
where it can then be routed to storage of your choice.


## Key Components

[Telegraf](https://www.influxdata.com/time-series-platform/telegraf/),
is an open source server agent to help you collect metrics from servers,
sensors and other systems.  [Telegraf](https://docs.influxdata.com/telegraf/)
lets you route this data to more than 50 different services including InfluxDB,
Graphite, Loki, MQTT, Kafka, New Relic, OpenTSDB, Prometheus, SQL servers of
various sorts, Websockets, and many more.  Telegraf runs on your local
network and is not tied to any particular internet or cloud service. 

[Piaware](https://flightaware.com/adsb/piaware/)
is an [ADS-B](https://www.faa.gov/nextgen/programs/adsb/)
ground station that you can make from a Raspbery Pi and a cheap
USB radio stick like the RTL-SDR dongle.  The USB receivers range
from $20 to $50 depending on who makes them.  Flightaware also
offers their own antennas, filters, and antennas.  A full description
of the hardware requirements can be found
[here](https://flightaware.com/adsb/piaware/build).


## How it works

Aircraft observations from Piaware look like this:
```json
{
    "hex":"a18284",
    "flight":"DAL1095 ",
    "alt_baro":33125,
    "alt_geom":34775,
    "gs":422.1,
    "ias":285,
    "tas":474,
    "mach":0.800,
    "track":269.0,
    "track_rate":0.00,
    "roll":0.5,
    "mag_heading":276.9,
    "baro_rate":1664,
    "geom_rate":1664,
    "squawk":"7162",
    "emergency":"none",
    "category":"A5",
    "nav_qnh":1013.6,
    "nav_altitude_mcp":36000,
    "nav_heading":277.7,
    "lat":40.173138,
    "lon":-76.133381,
    "nic":8,
    "rc":186,
    "seen_pos":1.4,
    "version":2,
    "nic_baro":1,
    "nac_p":9,
    "nac_v":1,
    "sil":3,
    "sil_type":"perhour",
    "gva":2,
    "sda":2,
    "mlat":[],
    "tisb":[],
    "messages":532,
    "seen":0.2,
    "rssi":-30.6 
}
```

The report's _key_ is the 
[ICAO 24-bit address](https://en.wikipedia.org/wiki/Aviation_transponder_interrogation_modes)
that shows up in the data as the field 'hex'.  This ID is assigned when the
Mode S transponder is registered to the aircraft.  That assignment is semi-
permanent as there are a only 2<sup>24</sup> possible codes and they do sometimes
get recycled.  Transponders also occasionally have to be replaced.

The data includes a field "flight" that is sometimes the flight number,
sometimes the aircraft tail number, and sometimes not set at all.  It often
is padded with whitespace.

This executable polls data from the /data/aircraft.json data feed of your Piaware
system.  This feed is a file, usually somewhere in /tmp, that is updated once
a second by Piaware.  It is accessible via the Piaware web server via a URL like this:

* http://piware-server/dump1090/data/aircraft.json
* http://piaware-server/tar1090/data/aircraft.json

depending on the version you are running.  There are a limited number of ICAO
addresses so they are occasionally re-used.  The FAA periodically releases the 
[U.S. registration database](http://registry.faa.gov/database/ReleasableAircraft.zip)
that includes the current ICAO assignments for aircraft that fly in the
U.S. along with the registration information such as the tail number,
aircraft type identifier, and owner.

This data still isn't enough to generate the context that explains
the report, so 
[this project](https://github.com/wiedehopf/tar1090-db) and
[This one](https://github.com/wiedehopf/readsb) pull together
detailed registration data from multiple sources, including those
outside of the jurisdiction of the FAA (US).

This project pulls all of this data together to generate more complete
records that are useful for storage and analysis.  The aircraft info and
registration databases are pulled from the internet and cached locally,
and updated in the background from time to time as needed.

The aircraft type code is interpreted, as best possible, into an aircraft
description.  However, the aircraft type designators in the ICAO database
use the same code for many (similar) different models, so accuracy is
definitely not guaranteed.

## Usage

The main purpose of this is to be executed as an
(execd)[https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd]
external input plug-in to telegraf.  To do this, add to your
existing telegraf.conf:

```toml
[[inputs.execd]]
  command = ["telegraf_piaware", "-url", "http://your-server/tar1090/data/aircraft.json"]
```
 
 replacing the url with the correct one for your Piaware instance.
 
 
## License

This project is licensed under the
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)
and is generally free for you to use how you want as long as you agree
to the terms of the license.


## Credits

The main guts of parsing the piaware data and merging it with registration and
aircraft type data was done by [Ed Welch](https://github.com/slim-bean/) as
part of his [adsb-loki project located here](github.com/slim-bean/adsb-loki).
His project is a standalone executable that writes piaware data to
 [Grafana Loki](https://grafana.com/oss/loki/).  He generously made some
 changes (on his day off) to the organization of his project so I could use
 it as a library for mine.
 

