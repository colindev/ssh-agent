#!/bin/env bash

project=`curl "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google"`
name=`curl "http://metadata.google.internal/computeMetadata/v1/instance/name" -H "Metadata-Flavor: Google"`
curl http://{AGENT}/${1}/keys?fingerprint=${name}.${project}
