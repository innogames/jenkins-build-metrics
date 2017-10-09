Jenkins Build Metrics
=====================

This simple go program will iterate over all Jenkins jobs and will ask for the latest build. If
the build failed, it will write a `1` to graphite with the predefined metric plus job name.

You can use this to have a build failed metric inside your grafana boards.

The script assumes that jenkins serves his API over https and the graphite gateway listens on tcp.


Compiling The Binary
--------------------

Native compilation (when compile platform equals target platform) can be executed by just running `make`

For cross-platform compilation, please override `GOOS` and `GOARCH` environment variables, e.g. `GOOS=linux GOARCH=amd64 make`


Running The Tool
----------------

```
jenkins-build-metrics -u user -t token -s jenkins.example.com
```

Arguments
---------

* `-u` a user configured in jenkins
* `-t` the users token, this is not the password
* `-s` the host name of the jenkins server
* `-p` the prefix for the graphite metric: default is `backend.marketing.jenkins`
* `-gp` the port for your graphite metrics gateway: default `3002`
* `-gh` the host name of your graphite metrics gateway: default `127.0.0.1`

