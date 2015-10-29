# vim: set ft=hcl :
service "my-service" {
  name = "my-service-blue"
  sequentialUpdate = true
  checkInterval = "1s"
  
  postCommand = "echo hello world"

  artifact "rackspace" {
    bucket = "my-service"
    path = "final/blue/my-service.tar.gz"
    type = "application/x-gzip"
    destination = "./usr/bin"
  }
}

service "my-service" {
  name = "my-service-blue"
  sequentialUpdate = true
  checkInterval = "1s"
  
  postCommand = "echo hello world"

  artifact "rackspace" {
    bucket = "my-service"
    path = "final/green/my-service.tar.gz"
    type = "application/x-gzip"
    destination = "./usr/bin"
  }
}
