# Scada-gobr

Open Source, web-based, multi-platform solution for building your own SCADA   
(Supervisory Control and Data Acquisition) system inspired by [ScadaLTS](https://github.com/SCADA-LTS/Scada-LTS).
Code released under [the GPL license](https://github.com/SCADA-LTS/Scada-LTS/blob/develop/LICENSE).

## What is the Goal?

The goal of this project is to create a web-based SCADA system that is easy to use, easy to install and
easy to maintain.

## What is the SCADA?

SCADA (Supervisory Control and Data Acquisition) is a control system architecture that uses computers,
networks, and graphical user interfaces to monitor and control a process. A SCADA system collects data
from remote field devices, such as sensors and switches, via a communication network and presents the
information on operator workstations for monitoring and control.

### Tech stack

* Go lang
* Postgres with timescale DB for time series
* React for web client

### Roadmap

* [ ] Custom dashboard view
* [ ] Realtime Dashboard
* [ ] Realtime data visualization
    * [ ] Histogram
    * [ ] Pie chart
    * [ ] Donut chart
    * [ ] Table
    * [ ] Map
    * [ ] Gauge
    * [ ] Tree view
    * [ ] Image
    * [ ] Text
    * [ ] Heatmap
    * [ ] Custom (user can add custom visualization)
* [ ] User permission
    * [ ] View permission
    * [ ] No Auth permission
* [ ] Custom web client changes, like images and text
* [ ] Control runtime manager in the web
* [ ] Export data
    * [ ] Export to CSV
    * [ ] Export to JSON
    * [ ] Export to XML
* [ ] More data sources
    * [ ] Modbus TCP/IP
    * [ ] Modbus Serial
    * [ ] DNP3
    * [ ] IEC 101
    * [ ] OPC DA 2.0
    * [ ] ASCII Serial and File readers
    * [ ] MongoDb
    * [ ] Elasticsearch
    * [ ] Big query
    * [ ] CSV
    * [ ] MQTT
    * [ ] AMQP
    * [ ] Kafka
    * [ ] gRPC
    * [x] REST
    * [x] Request Http
    * [x] Sql (Postgres, Mysql, Sqlite, Sqlserver)
* [x] Open telemetry
    * [ ] Tracing
    * [ ] Metrics
    * [ ] Logs
    * [ ] Custom exporter
* [ ] Internationalization
* [ ] Prometheus metrics
* [ ] Metrics dashboard
* [ ] Schedule
    * [ ] Cron like
    * [ ] Persistent
    * [ ] Queue
    * [ ] Run scripts
* [ ] Swagger
* [ ] Load data sources runtime across a cluster of s-gobr with consensus and then create a k8s operator
    * [ ] Consensus Raft
* [ ] Control with script
    * [ ] Pipeline data
    * [ ] Listen and publish to internal events
    * [ ] Listen and publish to external events like mqtt, amqp, kafka, grpc, rest, http
    * [ ] Listen and publish to internal data points
    * [ ] Notifications
        * [ ] Email
        * [ ] SMS
        * [ ] Webhook
* [ ] Support for multiple database as main and time series
    * [ ] Sqlite
    * [x] Postgres

### Building the application

To run the tasks below you need, **docker**, **npm** and **go** instaled

```shell
# create api binary
make api
# setup postgres
docker-compose up -d postgres
# Execute aplication
# The default port is 11139
./api
```

### Dev mode

We need this because the production mode all the assets is bundled inside the go binary, so if you can have real time
web client change you may need to fallow the steps below.

```shell
# Setup api
go run cmd/api/api.go
# Setup web client (with another terminal)
cd scadagobr-client
# Instal the dependencies
npm i
npm run dev
```

## Contribute

1. Fork the repository
2. Create a pull request
