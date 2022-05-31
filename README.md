# Scada-gobr

Open Source, web-based, multi-platform solution for building your own SCADA   
(Supervisory Control and Data Acquisition) system inspired by [ScadaLTS](https://github.com/SCADA-LTS/Scada-LTS).
Code released under [the GPL license](https://github.com/SCADA-LTS/Scada-LTS/blob/develop/LICENSE).

## Tech stack

* Go lang
* Postgres with timescale DB for time series
* Svelte for web client

## Roadmap

* Custom dashboard view
* Realtime Dashboard
* User permission
* Custom web client changes, like images and text
* Control runtime manager in the web
* More data sources
  * Modbus TCP/IP
  * Modbus Serial
  * DNP3
  * IEC 101
  * OPC DA 2.0
  * ASCII Serial and File readers 
  * MongoDb
  * Elasticsearch
  * Big query
  * CSV
  * MQTT
  * AMQP
  * Kafka
  * Google cloud IOT CORE
  * gRPC
* Open telemetry
* Internationalization 
* Prometheus metrics
* Load data sources runtime across a cluster of s-gobr with consensus and then create a k8s operator
* Control with script

## Building the application

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

## Dev mode

We need this because the production mode all the assets is bundled inside the go binary, so if you can have real time web client change you may need to fallow the steps below.

```shell
# Setup api
go run cmd/api/api.go
# Setup web client (with another terminal)
cd scadagobr-client
# Instal the dependencies
npm i
npm run dev
```