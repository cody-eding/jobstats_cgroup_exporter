# cgroup Prometheus exporter
[![GitHub release](https://img.shields.io/github/v/release/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter?include_prereleases&sort=semver)](https://github.com/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter/releases/latest)
![GitHub All Releases](https://img.shields.io/github/downloads/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter/total)

# Jobstats-Compatible cgroup Prometheus Exporter

The `jobstats_cgroup_exporter` produces metrics from v2 cgroups specifically for use with [Jobstats](https://princetonuniversity.github.io/jobstats/), a tool used for monitoring resource utilization on Slurm HPC clusters.

This exporter is a modified adaptation of treydock’s excellent [cgroup_exporter](https://github.com/treydock/cgroup_exporter), on which the [Princeton University fork](https://github.com/plazonic/cgroup_exporter) is also based. Unlike the Princeton fork, `jobstats_cgroup_exporter` is based on a later version of treydock’s exporter and fully supports v2 cgroups. Unlike treydock’s exporter, this adaptation **only** supports v2 cgroups. It also removes the ability to track other cgroups or processes outside of Slurm. This reduces the amount of code required and focuses this build solely on providing the metrics needed for Jobstats.

By default, this exporter listens on port `9306`, and all metrics are exposed via the `/metrics` endpoint.

# Usage

The exporter scans `/system.slice/slurmstepd.scope` by default. This can be overridden by providing a comma-separated list of paths via the `--config.paths` flag.

For example, if Slurm is compiled to support multiple `slurmd` instances and your cgroup paths are:

`/sys/fs/cgroup/system.slice/<nodename>_slurmstepd.scope`

you must pass:

`--config.paths=/system.slice/<nodename>_slurmstepd.scope`

replacing `<nodename>` with the host’s `slurmd` `NodeName`.

## Install

Download the [latest release](https://github.com/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter/releases)

## Build from source

To produce the `cgroup_exporter` binaries:

```
make build
```

Or

```
go get github.com/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter
```

## Process metrics

```
setcap cap_sys_ptrace=eip /usr/bin/cgroup_exporter
```

## Metrics

Example of metrics exposed by this exporter with default settings:

```
# HELP cgroup_cpu_info Information about the cgroup CPUs
# TYPE cgroup_cpu_info gauge
cgroup_cpu_info{cgroup="/system.slice/slurmstepd.scope/job_223239",cpus="0,2,4,6,8,10,12,14,16,18,20,22,24,26,28",jobid="223239"} 1
cgroup_cpu_info{cgroup="/system.slice/slurmstepd.scope/job_223344",cpus="30,32,34,36,38,40,42,44",jobid="223344"} 1
# HELP cgroup_cpu_system_seconds Cumalitive CPU system seconds for cgroup
# TYPE cgroup_cpu_system_seconds gauge
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 854.241433
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 853.312255
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 47.884299
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 45.165054
# HELP cgroup_cpu_total_seconds Cumalitive CPU total seconds for cgroup
# TYPE cgroup_cpu_total_seconds gauge
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 178389.268434
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 178382.956806
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 140.809977
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 137.528931
# HELP cgroup_cpu_user_seconds Cumalitive CPU user seconds for cgroup
# TYPE cgroup_cpu_user_seconds gauge
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 177535.027
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 177529.644551
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 92.925678
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 92.363877
# HELP cgroup_cpus Number of CPUs in the cgroup
# TYPE cgroup_cpus gauge
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 15
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 0
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 8
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 0
# HELP cgroup_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, goversion from which cgroup_exporter was built, and the goos and goarch for the build.
# TYPE cgroup_exporter_build_info gauge
cgroup_exporter_build_info{branch="",goarch="amd64",goos="linux",goversion="go1.25.7",revision="unknown",tags="unknown",version=""} 1
# HELP cgroup_info User slice information
# TYPE cgroup_info gauge
cgroup_info{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",uid="209096162",username="hpcuser1"} 1
cgroup_info{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",uid="750909302",username="hpcuser2"} 1
# HELP cgroup_memory_cache_bytes Memory cache used in bytes
# TYPE cgroup_memory_cache_bytes gauge
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 6.0682293248e+10
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 6.0682092544e+10
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 9.428992e+07
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 9.4035968e+07
# HELP cgroup_memory_fail_count Memory fail count
# TYPE cgroup_memory_fail_count gauge
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 0
# HELP cgroup_memory_rss_bytes Memory RSS used in bytes
# TYPE cgroup_memory_rss_bytes gauge
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 3.247726592e+09
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 3.233476608e+09
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 1.95446784e+09
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 1.948413952e+09
# HELP cgroup_memory_total_bytes Memory total given to cgroup in bytes
# TYPE cgroup_memory_total_bytes gauge
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 6.442450944e+10
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 1.8446744073709552e+19
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 3.4359738368e+10
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 1.8446744073709552e+19
# HELP cgroup_memory_used_bytes Memory used in bytes
# TYPE cgroup_memory_used_bytes gauge
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 6.393001984e+10
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 6.3915569152e+10
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 2.04875776e+09
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 2.04244992e+09
# HELP cgroup_memsw_total_bytes Swap total given to cgroup in bytes
# TYPE cgroup_memsw_total_bytes gauge
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 0
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 1.8446744073709552e+19
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 0
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 1.8446744073709552e+19
# HELP cgroup_memsw_used_bytes Swap used in bytes
# TYPE cgroup_memsw_used_bytes gauge
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239",jobid="223239",step="",task=""} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223239/step_batch/user/task_0",jobid="223239",step="batch",task="0"} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344",jobid="223344",step="",task=""} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223344/step_batch/user/task_0",jobid="223344",step="batch",task="0"} 0
# HELP cgroup_uid Uid number of user running this job
# TYPE cgroup_uid gauge
cgroup_uid{jobid="223239",username="hpcuser1"} 2.09096162e+08
cgroup_uid{jobid="223344",username="hpcuser2"} 7.50909302e+08
```