#!/bin/bash -xe
git subtree push -P src/github.com/cmars/ooclient git@github.com:cmars/ooclient.git master
git subtree push -P src/github.com/cmars/oostore git@github.com:cmars/oostore.git master
git subtree push -P src/github.com/cmars/quorum git@github.com:cmars/quorum.git master
git subtree push -P charms/layers/pgsql git@github.com:cmars/juju-relation-pgsql.git master

