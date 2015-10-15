# oo pki utilities

The files in this directory are useful for creating server certificates for TLS
endpoints.

# Prerequisites

## cfssl

`apt-get install libltdl-dev`, required by cfssl.
`gb -R ../tools` to build cfssl from vendored source.

## jq

`apt-get install jq`, used to extract PEM certificates from JSON.

# Roll your own CA

`$ make`

Will generate a self-signed CA and issue some example certificates including
one for localhost.

