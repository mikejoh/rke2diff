# rke2diff

`rke2diff` - Diff Rancher RKE2 releases! 🚀

**Notes:**

* This tool uses the GitHub API to fetch releases. The API is rate-limited to 60 requests for unauthenticated requests. More info about rate limiting can be found [here](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28).

## Install

### From source

1. `git clone https://github.com/mikejoh/rke2diff.git`
2. `cd rke2diff`
3. `make build` (places the compiled binary in `./build/`
4. `make install` (copies the compiled binary to `~/.local/bin`)

### Download and run

1. Download (using `v0.1.3` as an example):

```bash
curl -LO https://github.com/mikejoh/rke2diff/releases/download/0.1.3/rke2diff_0.1.3_linux_amd64.tar.gz
```

2. Unpack:

```bash
tar xzvf rke2diff_0.1.3_linux_amd64.tar.gz
```

3. Run:

```bash
./rke2diff -version
```

## Usage

```
rke2diff -h
Usage of ./build/rke2diff:
  -per-page int
     Skip release candidate releases. (default 100)
  -pick
     Interactive release picker.
  -releases
     Show all releases.
  -rke2 value
     RKE2 version to compare, can be set multiple times.
  -skip-rc
     Skip release candidate releases. (default true)
  -version
     Print the version number.
```

Compare to version of RKE2:

```bash
rke2diff -rke2 v1.28.11+rke2r1 -rke2 v1.31.0+rke2r1
┌──────────────────────────────────┬──────────────────────────────┬─────────────────────────────┐
│ NAME                             │ V1.28.11+RKE2R1 (2024-07-01) │ V1.31.0+RKE2R1 (2024-09-04) │
├──────────────────────────────────┼──────────────────────────────┼─────────────────────────────┤
│ containerd                       │ v1.7.17-k3s1                 │ v1.7.20-k3s1 🚀             │
│ coredns                          │ v1.11.1                      │ v1.11.1                     │
│ etcd                             │ v3.5.13-k3s1                 │ v3.5.13-k3s1                │
│ helm-controller                  │ v0.15.10                     │ v0.16.3 🚀                  │
│ ingress-nginx                    │ v1.10.1-hardened1            │ v1.10.4-hardened2 🚀        │
│ kubernetes                       │ v1.28.11                     │ v1.31.0 🚀                  │
│ metrics-server                   │ v0.7.1                       │ v0.7.1                      │
│ runc                             │ v1.1.12                      │ v1.1.13 🚀                  │
│ harvester-cloud-provider         │ 0.2.400                      │ 0.2.600 🚀                  │
│ harvester-csi-driver             │ 0.1.1700                     │ 0.1.1800 🚀                 │
│ rancher-vsphere-cpi              │ 1.7.001                      │ 1.9.000 🚀                  │
│ rancher-vsphere-csi              │ 3.1.2-rancher400             │ 3.3.1-rancher100 🚀         │
│ rke2-calico                      │ v3.27.300                    │ v3.28.100 🚀                │
│ rke2-calico-crd                  │ v3.27.002                    │ v3.28.100 🚀                │
│ rke2-canal                       │ v3.28.0-build2024062503      │ v3.28.1-build2024082701 🚀  │
│ rke2-cilium                      │ 1.15.500                     │ 1.16.101 🚀                 │
│ rke2-coredns                     │ 1.29.002                     │ 1.29.004 🚀                 │
│ rke2-ingress-nginx               │ 4.10.101                     │ 4.10.401 🚀                 │
│ rke2-metrics-server              │ 3.12.002                     │ 3.12.002                    │
│ rke2-snapshot-controller         │ 1.7.202                      │ 1.7.202                     │
│ rke2-snapshot-controller-crd     │ 1.7.202                      │ 1.7.202                     │
│ rke2-snapshot-validation-webhook │ 1.7.302                      │ 1.7.302                     │
├──────────────────────────────────┼──────────────────────────────┼─────────────────────────────┤
│ ❌ = VERSION STRING INVALID      │                              │                             │
│ 🚀 = BUMPED                      │                              │                             │
└──────────────────────────────────┴──────────────────────────────┴─────────────────────────────┘
```
