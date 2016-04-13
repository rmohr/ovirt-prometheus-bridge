# ovirt-prometheus-bridge

`ovirt-prometheus-bridge` is a host autodiscovery service for Prometheus
targets. It can be used to query oVirt Engine for hosts. The result is stored
in a json file which Prometheus can use to find scrape targets. Prometheus will
then start collecting metrics from all hosts in your virtual datacenter where
[vdsm-prometheus](https://github.com/rmohr/vdsm-prometheus) is installed.

[![Build Status](https://travis-ci.org/rmohr/ovirt-prometheus-bridge.svg?branch=master)](https://travis-ci.org/rmohr/ovirt-prometheus-bridge)

In this example
```bash
docker run -e ENGINE_PASSWORD=engine -v $PWD:/targets rmohr/ovirt-prometheus-bridge -update-interval 60 -no-verify -engine-url=https://localhost:8443 -output /targets/targets.json
```
the service is querying the oVirt Engine API every 60 seconds and writes the
found hosts into the file `targets.json`.  The created file `targets.json`
looks like this:

```json
[
  {
    "targets": [
      "192.168.122.190:8181",
      "192.168.122.41:8181"
    ],
    "labels": {
      "cluster": "0294d770-70c0-4b99-a527-f8a4ff4de436"
    }
  }
]
```

Prometheus can monitor files like this and update its configuration whenever
this file changes.  Here is a sample configuration for prometheus:

```yaml
global:
  scrape_interval: 15s

  external_labels:
    monitor: 'ovirt'

scrape_configs:
  - job_name: 'prometheus'

    scrape_interval: 5s

    target_groups:
      - targets: ['localhost:9090']
        labels:
          group: 'prometheus'

  - job_name: 'vdsm'

    scrape_interval: 5s
    file_sd_configs:
      - names : ['/targets/*.json']
```

Here is another example Prometheus config which uses TLS to communicate with
[vdsm-prometheus](https://github.com/rmohr/vdsm-prometheus):

```yaml
global:
  scrape_interval: 15s

  external_labels:
    monitor: 'ovirt'

scrape_configs:
  - job_name: 'prometheus'

    scrape_interval: 5s

    target_groups:
      - targets: ['localhost:9090']
        labels:
          group: 'prometheus'

  - job_name: 'vdsm'

    scheme: 'https'
    tls_config:
      ca_file: '/etc/pki/prom/certs/cacert.pem'
      cert_file: '/etc/pki/prom/certs/promcert.pem'
      key_file: '/etc/pki/prom/certs/promkey.pem'
      insecure_skip_verify: false
    scrape_interval: 5s
    file_sd_configs:
      - names : ['/targets/*.json']
```

# Quick start

To quickly spawn ovirt-prometheus-bridge, Prometheus and Grafana you can use
the Docker compose file in this repository:

```bash
export HOSTIP=$OVIRT_ENGINE_IP
export ENGINE_PASSWORD=$OVIRT_ENGINE_PASSWORD
docker-compose up
```

Then add the Prometheus datasource to Grafana:

```bash
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" --data '{ "name":"oVirt", "type":"prometheus", "url":"http://prometheus:9090", "access":"proxy", "basicAuth":false }' http://admin:admin@localhost:3000/api/datasources
```

Prometheus will then listen on [localhost:9090](http://localhost:9090), Grafana
on [localhost:3000](http://localhost:3000) and ovirt-prometheus-bridge will
provide the scrape targets. The default credentials for Grafana are
`admin:admin`.
