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
  after "graphite-event" {
    host = "http://my-graphite-server"
    tags = "my-service deployment"
    what = "deployed to {{.Hostname}}"
    data = "{{.Hash}}"
  }

  # Run this command after the update finishes
  after "script" {
    command = "initctl my-service start"
  }
}

