# rke2diff

`rke2diff` - Diff Rancher RKE2 releases! 🚀

**Notes:**
* This tool uses the GitHub API to fetch releases. The API is rate-limited to 60 requests for unauthenticated requests. More info about rate limiting can be found [here](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28).

## Install

### From source

1. `git clone https://github.com/mikejoh/rke2diff.git`
2. `cd rke2diff`
3. `make build`
4. `make install` (assumes that `~/.local/bin` is used)

### Download and run

1. Download (using `v0.1.3` as an example):
```
curl -LO https://github.com/mikejoh/rke2diff/releases/download/0.1.3/rke2diff_0.1.3_linux_amd64.tar.gz
```
2. Unpack:
```
tar xzvf rke2diff_0.1.3_linux_amd64.tar.gz
```
3. Run:
```
./rke2diff -version
```

## Usage

Example:
```
rke2diff -rke2 v1.28.9+rke2r1 -rke2 v1.29.5+rke2r1

┌──────────────────────────────────┬─────────────────────────────┬─────────────────────────────┐
│ NAME                             │ V1.28.9+RKE2R1 (2024-04-29) │ V1.29.5+RKE2R1 (2024-05-22) │
├──────────────────────────────────┼─────────────────────────────┼─────────────────────────────┤
│ Containerd                       │ v1.7.11-k3s2                │ v1.7.11-k3s2                │
│ CoreDNS                          │ v1.11.1                     │ v1.11.1                     │
│ Etcd                             │ v3.5.9-k3s1                 │ v3.5.9-k3s1                 │
│ Helm-controller                  │ v0.15.9                     │ v0.15.9                     │
│ Ingress-Nginx                    │ nginx-1.9.6-hardened1 ❌    │ nginx-1.9.6-hardened1 ❌    │
│ Kubernetes                       │ v1.28.9                     │ v1.29.5 🚀                  │
│ Metrics-server                   │ v0.7.1                      │ v0.7.1                      │
│ Runc                             │ v1.1.12                     │ v1.1. ❌                    │
│ harvester-cloud-provider         │ 0.2.300                     │ 0.2.300                     │
│ harvester-csi-driver             │ 0.1.1700                    │ 0.1.1700                    │
│ rancher-vsphere-cpi              │ 1.7.001                     │ 1.7.001                     │
│ rancher-vsphere-csi              │ 3.1.2-rancher400            │ 3.1.2-rancher400            │
│ rke2-calico                      │ v3.27.300                   │ v3.27.300                   │
│ rke2-calico-crd                  │ v3.27.002                   │ v3.27.002                   │
│ rke2-canal                       │ v3.27.3-build2024042301     │ v3.27.3-build2024042301     │
│ rke2-cilium                      │ 1.15.400                    │ 1.15.500 🚀                 │
│ rke2-coredns                     │ 1.29.002                    │ 1.29.002                    │
│ rke2-ingress-nginx               │ 4.9.100                     │ 4.9.100                     │
│ rke2-metrics-server              │ 3.12.002                    │ 3.12.002                    │
│ rke2-snapshot-controller         │ 1.7.202                     │ 1.7.202                     │
│ rke2-snapshot-controller-crd     │ 1.7.202                     │ 1.7.202                     │
│ rke2-snapshot-validation-webhook │ 1.7.302                     │ 1.7.302                     │
├──────────────────────────────────┼─────────────────────────────┼─────────────────────────────┤
│ ❌ = VERSION STRING INVALID      │                             │                             │
│ 🚀 = BUMPED                      │                             │                             │
└──────────────────────────────────┴─────────────────────────────┴─────────────────────────────┘
```
