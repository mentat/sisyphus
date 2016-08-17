# Sisyphus

Constantly rolling stuff up a hill.

## Overview

A microservice that efficiently downloads a file from CloudFiles in a certain
DC, only to re-upload it to another DC.  It does this extremely efficiently
by parallelizing the download and upload.

It relies on gocloudfiles: https://github.com/mentat/gocloudfiles

## download

    cd $GOPTH/src
    go get github.com/mentat/gocloudfiles
    go get github.com/mentat/sisyphus

## Testing

    export TEST_USERNAME="blah"
    export TEST_KEY="blah"

    go test

## Building

    go build

## Install

    go install
    $GOPATH/bin/sisyphus

## Usage

    ./sisyphus
