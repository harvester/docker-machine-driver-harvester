docker-machine-driver-harvester
========
[![Build Status](https://drone-publish.rancher.io/api/badges/harvester/docker-machine-driver-harvester/status.svg)](https://drone-publish.rancher.io/harvester/docker-machine-driver-harvester)

The [Harvester](https://github.com/harvester/harvester) machine driver for Docker.


## Development

### Building
```bash
make
```

The binary is placed in the `bin` directory.

The compressed binary is placed in the `dist/artifacts` directory.


## Usage

Put the binary to your $PATH directory

```bash
docker-machine create --driver harvester
```
