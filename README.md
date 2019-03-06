# putio-go-aria2
Script to add putio links to aria2.
The URIs are added as a json opject.

The aria2 server has to be run on localhost with port 6800.
Currently, the server should not have a sercret.

Furthermore, the OAuth token from put.io has to be added by hand.

## Futur Work
1. add O-Auth secret handling by config file
1. add config parameters like aria2 address, port and secret
1. add list of downloaded content
1. add check if content was downloaded by aria2 (Status check)
1. add remove function for downloaded content

