job "putio-go-aria2" {
  datacenters = [ "dc1" ]
  group "download" {
    restart {
      attempts = 0
    }
    reschedule {
      interval       = "1h"
      delay          = "30s"
      delay_function = "exponential"
      max_delay      = "120s"
      unlimited      = false
    }
    task "download" {
      driver = "docker"
      config {
        image = "azakazkaran/putio-go-aria2"
        volumes = [
          "/home/azak/putio-go-aria2.yml:/app/config.yml"
        ]
      }
    }
  }
}
