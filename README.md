# Dropship

Dropship is a simple tool for installing and updating artifacts from CDN.

## Features

- Automatically performs md5sum checks of artifact that is on server and remote
and will download automatically
- Distributed sequential updates
- Multiple Artifact Repository Support

## Installation

To install on ubuntu do the following:

```
echo "deb http://dl.bintray.com/chrismckenzie/deb trusty main" >> /etc/apt/sources.list
sudo apt-get update
sudo apt-get install dropship
```

## Configuration

To setup dropship you will need to add/update the following files.

First you will need to tell dropship how to connect to your artifact repository
so you will need to uncomment out the desired repo and fill in its options.

_/etc/dropship.d/dropship.hcl_
```hcl
# vim: set ft=hcl :
# Location that service config will be read from
service_path = "/etc/dropship.d/services"

# Rackspace Repo Config
# =====================
rackspace {
  user = "<your-rackspace-user>"
  key = "<your-rackspace-key>"
  region = "<rackspace-region>"
}
```

You will then have to create a file in the services directory of dropship. this 
will tell dropship how to check and install you artifact. You can have multiple
`service` definitions in one file or multiple files.

_/etc/dropship.d/services/my-service.hcl_
```hcl
# vim: set ft=hcl :
service "my-service" {
  # Use a semaphore to update one machine at a time
  sequentialUpdates = true

  # Check for updates every 10s
  checkInterval = "10s"

  # Run this command before update starts
  before "script" {
    command = "initctl my-service stop"
  }

  # Artifact defines what repository to use (rackspace) and where 
  # your artifact live on that repository
  artifact "rackspace" {
    bucket = "my-container"
    path = "my-service.tar.gz"
    destination = "./test/dest"
  }

  # After successful update send an event to graphite
  # this allows you to show deploy annotations in tools like grafana
  # 
  # The graphite hook will automatically add this services name into the 
  # graphite tags. You also have access to all of the services meta data
  # like Name, "current hash", hostname.
  after "graphite-event" {
    host = "http://<my-graphite-server>"
    tags = "deployment"
    what = "deployed to {{.Name}} on {{.Hostname}}"
    data = "{{.Hash}}"
  }

  # Run this command after the update finishes
  after "script" {
    command = "initctl my-service start"
  }
}
```

## Roadmap

- [X] Hooks
- [ ] Support for Amazon S3, and FTP
- [ ] Support for different file types deb, rpm, file _(currently only tar.gz)_
- [ ] Reporting system
- [ ] Redis, etcd for semaphore
