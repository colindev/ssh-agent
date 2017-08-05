#!/bin/env bash

curl http://{AGENT}/users/${1}/keys?fingerprint=`hostname`
