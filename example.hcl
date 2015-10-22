# vim: set ft=hcl :
service "data-service-api" {
  name = "data-service-api-blue"
  sequentialUpdate = true
  restartService = true
  checkInterval = "1s"

  artifact "rackspace" {
    bucket = "data-service"
    path = "final/blue/data-service.tar.gz"
    type = "application/x-gzip"
    destination = "./sites/"
  }
}
