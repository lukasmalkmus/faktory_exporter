# lukasmalkmus/faktory_exporter

> A Faktory Exporter for Prometheus. - by **[Lukas Malkmus](https://github.com/lukasmalkmus)**

[![Travis Status][travis_badge]][travis]
[![Go Report][report_badge]][report]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

> [!IMPORTANT]  
> This project is not actively maintained. As explained in [this issue](https://github.com/lukasmalkmus/faktory_exporter/issues/2) any contribution will be reviewed and potentially accepted, however it is not regularly tested. A very similar [Python Faktory exporter](https://github.com/prospect-server/faktory-prometheus-exporter) has been developed and is, as of 2024, actively used.

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
make docker
docker run -d --rm -p9386:9386 faktory-exporter:master
```

### Contributing

Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

### License

© Lukas Malkmus, 2017

Distributed under Apache License (`Apache License, Version 2.0`).

See [LICENSE](LICENSE) for more information.

[travis]: https://travis-ci.org/lukasmalkmus/faktory_exporter
[travis_badge]: https://travis-ci.org/lukasmalkmus/faktory_exporter.svg
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/faktory_exporter
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/faktory_exporter
[release]: https://github.com/lukasmalkmus/faktory_exporter/releases
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/faktory_exporter.svg
[license]: https://opensource.org/licenses/Apache-2.0
[license_badge]: https://img.shields.io/badge/license-Apache-blue.svg
