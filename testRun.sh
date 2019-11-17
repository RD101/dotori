#!/bin/sh
go run assets/asset_generate.go
go install
dotori -http :80