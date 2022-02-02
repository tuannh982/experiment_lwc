Experiment LWC
===

Just learning how container works

## Prerequesites

install `cgdelete`

update iptables FORWARD policy

```shell
iptables --policy FORWARD ACCEPT
```

## Quickstarts

Quickstarts with alpine as base image
```shell
mkdir -p rootfs/alpine-minirootfs-3.15.0-x86_64
tar -xf ./rootfs_store/alpine-minirootfs-3.15.0-x86_64.tar.gz -C ./rootfs/alpine-minirootfs-3.15.0-x86_64
```

then, start 2 containers

[Container1]
```shell
sudo ./experiment_lwc --opt run --config samrunple_config.yaml
```

[Container2]
```shell
sudo ./experiment_lwc --opt run --config sample_config_2.yaml
```

those 2 containers can communicate with each others, like this

[Container2]
```shell
/ # ping 10.10.10.2
PING 10.10.10.2 (10.10.10.2): 56 data bytes
64 bytes from 10.10.10.2: seq=0 ttl=64 time=0.091 ms
64 bytes from 10.10.10.2: seq=1 ttl=64 time=0.043 ms
^C
--- 10.10.10.2 ping statistics ---
2 packets transmitted, 2 packets received, 0% packet loss
round-trip min/avg/max = 0.043/0.067/0.091 ms
```

## Build

```shell
go build -trimpath
```

## Usage

```shell
./experiment_lwc --help
./experiment_lwc --opt run --help
./experiment_lwc --opt inspect --help
./experiment_lwc --opt cleanup --help
./experiment_lwc --opt network --help
```

## Cleanup

Run `cleanup.sh` to clean the database and container directory

## TODOs

- Container can connect to internet
- Add seccomp filter to container process