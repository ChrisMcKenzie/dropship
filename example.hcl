# vim: set ft=hcl :
service "my-service" {
  sequentialUpdates = true
  checkInterval = "1s"
  
  preCommand = "echo hello world"
  postCommand = "echo hello world"

  artifact "rackspace" {
    bucket = "my-service"
    path = "final/blue/my-service.tar.gz"
    destination = "./usr/bin"
  }
}
