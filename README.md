# oo - Opaque Objects
[![Build Status](https://travis-ci.org/cmars/oo.svg?branch=master)](https://travis-ci.org/cmars/oo)

This repository is a [gb](http://getgb.io) project build for
[oostore](https://github.com/cmars/oostore) and
[ooclient](https://github.com/cmars/ooclient).

# Prerequisites

Install [gb](http://getgb.io).

# Build

`gb build`

There, now wasn't that easy? :)

# Tooling

[Cloudflare CFSSL](https://github.com/cloudflare/cfssl) will be used to manage
X.509 certificates for deploying with TLS. This tool is also vendored. To build
a starting point PKI for this:

```
gb build -R tools
make -C pki
```

# License

Source code under `tools/src` and `vendor/src` are copyright their respective
authors. Refer to the various "LICENSE", "LICENCE" and "COPYING" files
contained therein for the specific license terms.

Everything else is Copyright 2015 Casey Marshall and licensed under the Apache
License, Version 2.0 (the "License"); you may not use this file except in
compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
