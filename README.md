# intelite-mqtt
[![Build Status](https://drone.esc.pp.ua/api/badges/alexcom/intelite-mqtt/status.svg)](https://drone.esc.pp.ua/alexcom/intelite-mqtt)

## What it is

Service that allows to control Maxus Intelite lamps via IR blasters with Tasmota IR firmware on board. Consumes MQTT
messages of certain own format and produces messages that Tasmota IR can interpret and emit. Should work with lamp
models SMT005, SMT006. Base implementation written by [AlexNk](https://github.com/AlexNk) for Android and (allegedly) tested against SMT006. In my turn I
ported it to Go and still testing against SMT005.

## Deployment

1. Building and running binary
2. Building and running in Docker

## Links

[Base implementation](https://github.com/AlexNk/intelite_smt006_remote)
