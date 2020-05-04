job "putio-download" {
  datacenters = [ "dc1" ]
  type = "batch"
  periodic {
    cron = "@hourly"
    prohibit_overlap = true
  }
  group "putio" {
    restart {
      attempts = 0
    }
    task "download" {
      driver = "docker"
      config {
        image = "azakazkaran/putio-go-aria2"
        volumes = [
          "/home/azak/apps/nomad_0.11.1_linux_amd64/putio-download/download.yml:/app/config.yml"
        ]
      }
    }
    task "organize"{
      driver = "docker"
      config {
        image = "azakazkaran/putio-go-aria2"
        volumes = [
          "/home/azak/apps/nomad_0.11.1_linux_amd64/putio-download/organize.yml:/app/config.yml",
          "/mnt/data/unsorted:/mnt/data/unsorted"
        ]
      }
    }
  }
}
