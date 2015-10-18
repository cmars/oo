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
one for localhost. This is useful for development purposes.

# Generating keys, certs, CSRs

1. Create a
   [CSRJSON](https://github.com/cloudflare/cfssl#generating-certificate-signing-request-and-private-key)
   file named <hostname>.csr.json in this directory.
2. `make <hostname>.pem` to generate a private key and issue a certificate with
   the self-signed CA.
3. `make <hostname>-csr.pem` to generate a private key and CSR.

# Role of TLS & X.509 in the oo cryptosystem

The security of opaque object content does not rely on the integrity of TLS,
X.509 or even the oostore service, because objects are encrypted end-to-end.

Trust in a third-party service that gates access to opaque objects is
bootstrapped by requesting public keys from the service over HTTPS.

