# putio-go-aria2
[![Build Status](https://travis-ci.org/azak-azkaran/agent.svg?branch=master)](https://travis-ci.org/azak-azkaran/agent)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_putio-go-aria2&metric=alert_status)](https://sonarcloud.io/dashboard?id=azak-azkaran_putio-go-aria2)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=azak-azkaran_putio-go-aria2&metric=coverage)](https://sonarcloud.io/dashboard?id=azak-azkaran_putio-go-aria2)

Script to add putio links to aria2.
The URIs are added as a json object.

The aria2 server has to be run on localhost with port 6800.
Currently, the server should not have a secret.

Furthermore, the OAuth token from put.io has to be added by hand.

## Future Work
1. add config parameters like aria2 address, port and secret
1. add check if content was downloaded by aria2 (Status check)

