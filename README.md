# Grafana Dashboard Generator


## Building

```shell
go mod tidy
go build .
```

## Running

```shell
./gdgenerator -config ./examples/generic.yaml --manifests-directory ./dashboard-resources --manifests
grafanactl resources push --path ./dashboard-resources/generic-dashboard.json
```

## More Information

* https://github.com/grafana/grafana-foundation-sdk
* https://github.com/grafana/grafanactl
* https://github.com/grafana/dashboards-as-code-workshop/tree/main/part-one-golang-starter
