# airport-subscription-exporter

Export usage info from airport subscription URL.

## Install

Download from [releases](https://github.com/fffonion/airport-subscription-exporter/releases) or run from docker

```
docker run -d -p 9991:9991 fffonion/airport-subscription-exporter
```

### Usage
Use the -h flag to see full usage:

```
$ airport-subscription-exporter -h
  -metrics.listen-addr string
        listen address for airport subscription exporter (default ":9991")
  -sub.update-interval string
        how long should exporter actually update subscription info (default "1h")
```

```
# HELP airport_last_refreshed Airport subscription info last refreshed time
# TYPE airport_last_refreshed gauge
airport_last_refreshed 1.726910091e+09
# HELP airport_userinfo_download Airport download traffic
# TYPE airport_userinfo_download counter
airport_userinfo_download 2.073741824e+07
# HELP airport_userinfo_expire Airport subscription expire time
# TYPE airport_userinfo_expire gauge
airport_userinfo_expire 1.742537346e+09
# HELP airport_userinfo_total Airport total subscription traffic
# TYPE airport_userinfo_total counter
airport_userinfo_total 1.073741824e+11
# HELP airport_userinfo_upload Airport upload traffic
# TYPE airport_userinfo_upload counter
airport_userinfo_upload 3.742537346e+04
```

## Sample prometheus config

```yaml
# scrape airport-subscription 
scrape_configs:
  - job_name: 'airport-subscription'
    static_configs:
    - targets: ["http://some-airport/api/sub/v1?token=1"]
      labels:
        instance: some-airport
    - targets: ["http://another-airport/api/sub/v1?token=1"]
      labels:
        instance: another-airport
    metrics_path: /scrape
    relabel_configs:
      - source_labels : [__address__]
        target_label: __param_target
      - target_label: __address__
        # IP of the exporter
        replacement: localhost:9991

# scrape airport-subscription-exporter itself
  - job_name: 'airport-subscription-exporter'
    static_configs:
      - targets:
        # IP of the exporter
        - localhost:9991
```

## Docker Build Instructions

Build for both `arm64` and `amd64`:
```
docker build -t <image-name>:latest-arm64 --platform linux/arm64 --build-arg GOARCH=arm64 .
docker build -t <image-name>:latest-amd64 --platform linux/amd64 --build-arg GOARCH=amd64 .
```

Merge them in one manifest:
```
docker manifest create <image-name>:latest --amend <image-name>:latest-arm64 --amend <image-name>:latest-amd64
docker manifest push <image-name>:latest
```
