# DockPush

Push Docker images into your Diego-powered Cloud Foundry.

![DockPush](https://i0.wp.com/deafwhale.com/pics/beached11.jpg)

## Installation

```
$ go get github.com/cloudcredo/dockpush
$ cf install-plugin $GOPATH/bin/dockpush
```

Please note - DockPush requires Cloud Foundry CLI 6.8.0 or later.

## Usage

```
cf dockpush docker-image run-command app-name

cf dockpush cloudfoundry/inigodockertest:latest /dockerapp docker-test

cf dp -m=512 -i=2 -d=1200 cloudfoundry/inigodockertest:latest /dockerapp docker-test
-m: memory limit (in MB) for container, default 1024
-i: number of instances, default 1
-d: disk space limit (in MB) for container, default 1024
```

## Where can I get a Diego-powered Cloud Foundry?

 https://github.com/cloudfoundry-incubator/diego-release#deploying-diego-to-a-local-bosh-lite-instance

 If you'd like a production-scale Diego-powered Cloud Foundry please do [get in contact.](http://www.cloudcredo.com/contact-us/)

## A message from our sponsors

[CloudCredo](http://www.cloudcredo.com) would like to hear from you.
