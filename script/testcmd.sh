#!/bin/bash

service nginx restart
go test "$@" ./...
