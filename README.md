# Jobstats-Compatible cgroup Prometheus Exporter

[![GitHub release](https://img.shields.io/github/v/release/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter?include_prereleases&sort=semver)](https://github.com/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter/releases/latest)
![GitHub All Releases](https://img.shields.io/github/downloads/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter/total)

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

To produce the `jobstats_cgroup_exporter` binaries:

```
make build
```

Or

```
go get github.com/GVSU-Advanced-Research-Computing/jobstats_cgroup_exporter
```

## Process metrics

The exporter must be able to read system process information:

```
setcap cap_sys_ptrace=eip /<path>/<to>/jobstats_cgroup_exporter
```

If running as a systemd service:

```
AmbientCapabilities=CAP_SYS_PTRACE
```

## Metrics

Example of metrics exposed by this exporter with default settings:

```
# HELP cgroup_cpu_info Information about the cgroup CPUs
# TYPE cgroup_cpu_info gauge
cgroup_cpu_info{cgroup="/system.slice/slurmstepd.scope/job_223478",cpus="0,4,8,12,16,20,24,28",jobid="223478"} 1
cgroup_cpu_info{cgroup="/system.slice/slurmstepd.scope/job_223482",cpus="1,2,5,6,9,10,13,14,18,22,26,30,32,34,36",jobid="223482"} 1
# HELP cgroup_cpu_system_seconds Cumalitive CPU system seconds for cgroup
# TYPE cgroup_cpu_system_seconds gauge
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 30836.499999
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 30824.641046
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 804.392829
cgroup_cpu_system_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 798.801398
# HELP cgroup_cpu_total_seconds Cumalitive CPU total seconds for cgroup
# TYPE cgroup_cpu_total_seconds gauge
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 671488.627949
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 671471.974181
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 165449.390505
cgroup_cpu_total_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 165441.533265
# HELP cgroup_cpu_user_seconds Cumalitive CPU user seconds for cgroup
# TYPE cgroup_cpu_user_seconds gauge
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 640652.127949
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 640647.333135
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 164644.997675
cgroup_cpu_user_seconds{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 164642.731867
# HELP cgroup_cpus Number of CPUs in the cgroup
# TYPE cgroup_cpus gauge
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 8
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 0
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 15
cgroup_cpus{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 0
# HELP cgroup_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, goversion from which cgroup_exporter was built, and the goos and goarch for the build.
# TYPE cgroup_exporter_build_info gauge
cgroup_exporter_build_info{branch="HEAD",goarch="amd64",goos="linux",goversion="go1.24.0",revision="565e497926b1dc69936b5b98bc3f16d5a9265f6c",tags="unknown",version="0.0.3-alpha"} 1
# HELP cgroup_info User slice information
# TYPE cgroup_info gauge
cgroup_info{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",uid="750920098",username="hpcuser1"} 1
cgroup_info{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",uid="209096162",username="hpcuser2"} 1
# HELP cgroup_memory_cache_bytes Memory cache used in bytes
# TYPE cgroup_memory_cache_bytes gauge
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 1.3164015616e+10
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 1.3163266048e+10
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 6.1180706816e+10
cgroup_memory_cache_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 6.1180125184e+10
# HELP cgroup_memory_fail_count Memory fail count
# TYPE cgroup_memory_fail_count gauge
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 0
cgroup_memory_fail_count{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 0
# HELP cgroup_memory_rss_bytes Memory RSS used in bytes
# TYPE cgroup_memory_rss_bytes gauge
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 1.1304026112e+10
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 1.1299147776e+10
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 2.997366784e+09
cgroup_memory_rss_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 2.988306432e+09
# HELP cgroup_memory_total_bytes Memory total given to cgroup in bytes
# TYPE cgroup_memory_total_bytes gauge
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 3.4359738368e+10
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 1.8446744073709552e+19
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 6.442450944e+10
cgroup_memory_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 1.8446744073709552e+19
# HELP cgroup_memory_used_bytes Memory used in bytes
# TYPE cgroup_memory_used_bytes gauge
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 2.4468041728e+10
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 2.4462413824e+10
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 6.41780736e+10
cgroup_memory_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 6.4168431616e+10
# HELP cgroup_memsw_total_bytes Swap total given to cgroup in bytes
# TYPE cgroup_memsw_total_bytes gauge
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 0
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 1.8446744073709552e+19
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 0
cgroup_memsw_total_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 1.8446744073709552e+19
# HELP cgroup_memsw_used_bytes Swap used in bytes
# TYPE cgroup_memsw_used_bytes gauge
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478",jobid="223478",step="",task=""} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223478/step_batch/user/task_0",jobid="223478",step="batch",task="0"} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482",jobid="223482",step="",task=""} 0
cgroup_memsw_used_bytes{cgroup="/system.slice/slurmstepd.scope/job_223482/step_batch/user/task_0",jobid="223482",step="batch",task="0"} 0
# HELP cgroup_uid Uid number of user running this job
# TYPE cgroup_uid gauge
cgroup_uid{jobid="223478",username="hpcuser1"} 7.50920098e+08
cgroup_uid{jobid="223482",username="hpcuser2"} 2.09096162e+08
```