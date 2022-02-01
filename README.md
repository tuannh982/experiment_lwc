Experiment LWC
===

Just learning how container works

## Quickstarts

Quickstarts with alpine as base image
```shell
rm -rf rootfs/
mkdir -p rootfs
tar -xf alpine-minirootfs-3.15.0-x86_64.tar.gz -C ./rootfs
```

then,

```shell
sudo ./experiment_lwc
```

## Build

```shell
go build -trimpath
```

## Usage

- Run `./experiment_lwc -h` for command usage
- Run `sudo ./experiment_lwc` to run
  - Note: this program must be executed as root (needed for cgroups)