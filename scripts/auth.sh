#!/bin/env bash

curl http://{AGENT}/${1}/keys?fingerprint=`hostname`
