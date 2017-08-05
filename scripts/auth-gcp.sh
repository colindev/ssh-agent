#!/bin/env bash

project=`curl "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google" 2>/dev/null`
name=`curl "http://metadata.google.internal/computeMetadata/v1/instance/name" -H "Metadata-Flavor: Google" 2>/dev/null`
curl http://{AGENT}/users/${1}/keys?fingerprint=${name}.${project}
