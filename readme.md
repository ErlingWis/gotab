# gotab
[![Docker](https://github.com/ErlingWis/gotab/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/ErlingWis/gotab/actions/workflows/docker-publish.yml)

minimal object store

## build

`nerdctl build . -t gotab:dev`

## usage

`nerdctl run -p 3000:8080 -v /tmp/gotab:/gotab-disk --rm gotab:dev --disk /gotab-disk --verbosity 4`