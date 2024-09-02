## Config File

```yaml
# This is an example config.yaml
# Using Cloudflare as DNS provider and a 3rd party API
# to detect public IP address of this network
ddns:
  # Use a unique name for your DDNS instance
  # You can add multiple instances here
  # For example one for v4 and one for v6
  - name: homelab
    domain: yourdomain.com
    # Use "@" in subdomain for zone apex
    # e.g. use DDNS for yourdomain.com itself
    subdomain: www
    # Or IPv6 for AAAA record
    stack: IPv4
    # Crontab expression
    # You can attach a timezone definition
    # by prepending TZ=<Your/Time_Zone>
    cron: "*/30 * * * *"
    # Here defines how should it detect
    # IP address of your network/device
    detection:
      # Or Interface to read IP address from
      # interfaces on your device
      type: ThirdParty
      api:
        url: https://api.ipify.org
        # You can specify a Json Path to extract
        # data from JSON response from API
      # jsonpath: ".ip"
    # type: Interface
    # interface:
    # Specify interface name here
    #   name: eth0
    provider:
      # Currently Cloudflare and AliCloud are supported
      name: Cloudflare
      cloudflare:
        apiToken: "<your-api-token>"
```

## Parameters

### DDNS related

| Name                               | Type   | Description                                                                                                                              |
|------------------------------------|--------|------------------------------------------------------------------------------------------------------------------------------------------|
| `ddns`                             | array  | Top level element for holding DDNS instances.                                                                                            |
| `ddns.name`                        | string | Name of the DDNS instance, cannot be same.                                                                                               |
| `ddns.domain`                      | string | Your domain name, without any subdomain.                                                                                                 |
| `ddns.subdomain`                   | string | Subdomain for this instance in punycode, use "@" for zone apex.                                                                          |
| `ddns.stack`                       | string | Use IPv4 or IPv6 address.                                                                                                                |
| `ddns.cron`                        | string | Crontab expression for how should the program arrange update operation. You can prepend `TZ=<Your/Time_Zone>` to specify your time zone. |
| `ddns.detection`                   | object | IP address detection method. Currently it can fetch address from 3rd-party HTTP API or read from network interface.                      |
| `ddns.detection.type`              | string | Detection method. Currently ThirdParty and Interface are supported.                                                                      |
| `ddns.detection.interface`         | object | Interface address detection specifications.                                                                                              |
| `ddns.detection.interface.name`    | string | Interface to read address from.                                                                                                          |
| `ddns.detection.api`               | object | Third-party API detection specification.                                                                                                 |
| `ddns.detection.api.url`           | string | 3rd-party API URL.                                                                                                                       |
| `ddns.detection.api.customHeaders` | object | (Optional) Custom headers that adds into requests to 3rd-party API.                                                                      |
| `ddns.detection.api.params`        | object | (Optional) Custom params that appends to API URL.                                                                                        |
| `ddns.detection.api.username`      | string | (Optional) API authentication username.                                                                                                  |
| `ddns.detection.api.password`      | string | (Optional) API authentication password.                                                                                                  |

### DNS provider related

| Name                                     | Type   | Description                                                                                                                |
|------------------------------------------|--------|----------------------------------------------------------------------------------------------------------------------------|
| `ddns.provider`                          | object | DNS provider specification.                                                                                                |
| `ddns.provider.name`                     | string | DNS provider name. Currently Cloudflare, AliCloud and DNSPod are supported.                                                |
| `ddns.provider.cloudflare`               | object | Credentials and settings for Cloudflare DNS provider.                                                                      |
| `ddns.provider.cloudflare.apiToken`      | string | Fine-grained API token for Cloudflare, recommended as this can limit permissions for a specific token.                     |
| `ddns.provider.cloudflare.globalApiKey`  | string | Global API key for Cloudflare, not recommended as this key has full power to access your Cloudflare account and resources. |
| `ddns.provider.cloudflare.email`         | string | Email of your Cloudflare account, must use with global API key.                                                            |
| `ddns.provider.alicloud`                 | object | Credentials and settings for AliCloud DNS provider.                                                                        |
| `ddns.provider.alicloud.accessKeyId`     | string | AccessKeyId of your AliCloud account. You can create a RAM sub user to limit permission of this access key.                |
| `ddns.provider.alicloud.accessKeySecret` | string | AccessKeySecret of your AliCloud account.                                                                                  |
| `ddns.provider.alicloud.regionId`        | string | Region ID of your resources.                                                                                               |
| `ddns.provider.alicloud.line`            | string | (Optional) Line of the DNS record.                                                                                         |
| `ddns.provider.dnspod`                   | object | Credentials and settings for DNSPod DNS provider.                                                                          |
| `ddns.provider.dnspod.secretId`          | string | SecretID of your Tencent Cloud account.                                                                                    |
| `ddns.provider.dnspod.secretKey`         | string | SecretKey of your Tencent Cloud account                                                                                    |
| `ddns.provider.dnspod.region`            | string | Region of your resources.                                                                                                  |
| `ddns.provider.dnspod.lineId`            | string | (Optional) ID of the line of your DNS record.                                                                              |
| `ddns.provider.huawei`                   | object | Credentials and settings for Huawei Cloud DNS provider.                                                                    |
| `ddns.provider.huawei.accessKey`         | string | Access key (AK) of the account.                                                                                            |
| `ddns.provider.huawei.secretAccessKey`   | string | Secret key (SK) of the account.                                                                                            |
| `ddns.provider.huawei.region`            | string | Region of resources and APIs.                                                                                              |
| `ddns.provider.jd`                       | object | Credentials and settings for JD Cloud DNS provider.                                                                        |
| `ddns.provider.jd.accessKey`             | string | Access key (AK) of the account.                                                                                            |
| `ddns.provider.jd.secretKey`             | string | Secret Key (SK) of the account.                                                                                            |
| `ddns.provider.jd.regionId`              | string | (Optional) Region ID of resources, leave empty for default value (cn-north-1).                                             |
| `ddns.provider.jd.viewId`                | number | (Optional) View ID of the record, leave empty for default value (-1).                                                      |