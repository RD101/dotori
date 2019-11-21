#!/bin/sh
go run assets/asset_generate.go
go install
sudo dotori -http :80

