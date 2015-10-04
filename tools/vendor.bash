#!/bin/bash
git subtree add -P vendor/src/$1 https://$1.git master
