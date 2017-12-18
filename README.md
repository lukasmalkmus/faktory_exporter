# lukasmalkmus/faktory_exporter

> A Faktory Exporter for Prometheus. - by **[Lukas Malkmus](https://github.com/lukasmalkmus)**

[![Travis Status][travis_badge]][travis]
[![CircleCI Status][circleci_badge]][circleci]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![GoDoc][docs_badge]][docs]
[![Docker Repository on Quay][quay_badge]][quay]
[![Docker Pulls][hub_badge]][hub]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

### Introduction

The *faktory_exporter* is a simple server that scrapes a configured Faktory
instance for stats by issuing the "INFO" command and exports them via HTTP for
Prometheus consumption.

### Usage

#### Installation

The easiest way to run the *faktory_exporter* is by grabbing the latest
binary from the [release page][release].

##### Building from source

This project uses [dep](https://github.com/golang/dep) for vendoring.

```bash
git clone https://github.com/lukasmalkmus/faktory_exporter
cd faktory_exporter
make
```

#### Using the exporter

```bash
./faktory_exporter [flags]
```

Help on flags:

```bash
./faktory_exporter --help
```

#### Using docker

```bash
docker run -d --rm -p9386:9386 quay.io/lukasmalkmus/faktory-exporter:latest
```

### Contributing

Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

### License

© Lukas Malkmus, 2017

Distributed under Apache License (`Apache License, Version 2.0`).

See [LICENSE](LICENSE) for more information.

[travis]: https://travis-ci.org/lukasmalkmus/faktory_exporter
[travis_badge]: https://travis-ci.org/lukasmalkmus/faktory_exporter.svg
[circleci]: https://circleci.com/gh/lukasmalkmus/faktory_exporter
[circleci_badge]: (https://circleci.com/gh/lukasmalkmus/faktory_exporter/tree/master.svg?style=shield)
[coverage]: https://coveralls.io/github/lukasmalkmus/faktory_exporter?branch=master
[coverage_badge]: https://coveralls.io/repos/github/lukasmalkmus/faktory_exporter/badge.svg?branch=master
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/faktory_exporter
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/faktory_exporter
[docs]: https://godoc.org/github.com/lukasmalkmus/faktory_exporter
[docs_badge]: https://godoc.org/github.com/lukasmalkmus/faktory_exporter?status.svg
[quay]: https://quay.io/repository/lukasmalkmus/faktory_exporter
[quay_badge]: https://quay.io/repository/lukasmalkmus/faktory_exporter/status
[hub]: https://hub.docker.com/r/lukasmalkmus/faktory-exporter
[hub_badge]: https://img.shields.io/docker/pulls/lukasmalkmus/faktory-exporter.svg
[release]: https://github.com/lukasmalkmus/faktory_exporter/releases
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/faktory_exporter.svg
[license]: https://opensource.org/licenses/Apache-2.0
[license_badge]: https://img.shields.io/badge/license-Apache-blue.svg