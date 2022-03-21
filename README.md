# dns-debug

A simple app that runs a dns query and traces it.

Don't read the code it's ugly and not mean to be used other then a quick test for me.

## Usage

```shell
export CURL="false"
export URL=google.se
export ENDPOINTS="https://google.com,https://google.se"
```

## Tracing

Currently support `datadog`, might support open-tracing in the future.

## Dockerfile

You can find upstream image at: `quay.io/nissessenap2/dns-debug:v0.0.3`

## Kubernetes deployment example

See [deployment.yaml](./deployment.yaml)
