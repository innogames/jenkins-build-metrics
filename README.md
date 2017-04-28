Jenkins Build Metrics
=====================

This simple go script will iterate over all jenkins jobs and will ask for the latest build. If
the build failed it will write a 1 to graphite with the predefined metric plus job name.

You can use this to have a build failed metric inside your grafana boards.

The script assume that jenkins serves his API over https and the graphite gateway listens on tcp. 

Example
-------

```
jenkins-build-metrcs -u foo -t bar -s jenkins.example.com
```

Documentation
-------------

* -u a user configured in jenkins
* -t the users token, this is not the password
* -s the host name of the jenkins server
* -p the prefix for the graphite metric: default is `backend.marketing.jenkins`
* -gp the port for your graphite metrics gateway: default `3002`
* -gh the host name of your graphite metrics gateway: default `127.0.0.1`
