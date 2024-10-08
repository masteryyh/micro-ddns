## Config File

Below is an example configuration file written in YAML, but you can also write in JSON.

Use `--config` or `-c` option in `run` command to specify configuration file path.

```yaml
# This is an example config.yaml
# Using Cloudflare as DNS provider and ipify API
# to detect public IP address
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
    # Choose a detection method defined by yourself
    detectionRef: api
    # Choose a DNS provider
    providerRef: cloudflare

# Here defines how should it detect
# IP address of your network/device
detection:
  # Name must be unique
  - name: api
    api:
      url: https://api.ipify.org
      # You can specify a Json Path to extract
      # data from JSON response from API
      # jsonpath: ".ip"
    # interface:
    # Specify interface name here
    #   name: eth0

provider:
  # For supported DNS providers, check the detailed document below
  - name: cloudflare
    cloudflare:
      apiToken: "<your-api-token>"
```

## Parameters

### DDNS fields

| Name                               | Type   | Description                                                                                                                              |
|------------------------------------|--------|------------------------------------------------------------------------------------------------------------------------------------------|
| `ddns`                             | array  | Top level element for holding DDNS instances.                                                                                            |
| `ddns.name`                        | string | Name of the DDNS instance, cannot be same.                                                                                               |
| `ddns.domain`                      | string | Your domain name, without any subdomain.                                                                                                 |
| `ddns.subdomain`                   | string | Subdomain for this instance in punycode, use "@" for zone apex.                                                                          |
| `ddns.stack`                       | string | Use IPv4 or IPv6 address.                                                                                                                |
| `ddns.cron`                        | string | Crontab expression for how should the program arrange update operation. You can prepend `TZ=<Your/Time_Zone>` to specify your time zone. |

### Address detection fields

| Name                                     | Type    | Description                                                                                                                                                                 |
|------------------------------------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `detection.name`                     | string  | Address detection specification name, must be unique.                                                                 |
| `detection.interface`         | object | Interface address detection specifications.                                                                                              |
| `detection.interface.name`    | string | Interface to read address from.                                                                                                          |
| `detection.api`               | object | Third-party API detection specification.                                                                                                 |
| `detection.api.url`           | string | 3rd-party API URL.                                                                                                                       |
| `detection.api.customHeaders` | object | (Optional) Custom headers that adds into requests to 3rd-party API.                                                                      |
| `detection.api.params`        | object | (Optional) Custom params that appends to API URL.                                                                                        |
| `detection.api.username`      | string | (Optional) API authentication username.                                                                                                  |
| `detection.api.password`      | string | (Optional) API authentication password.                                                                                                  |

### DNS provider fields

| Name                                     | Type    | Description                                                                                                                                                                 |
|------------------------------------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `provider.name`                     | string  | DNS provider specification name, must be unique.                                                                 |
| `provider.cloudflare`               | object  | Credentials and settings for Cloudflare DNS provider.                                                                                                                       |
| `provider.cloudflare.apiToken`      | string  | Fine-grained API token for Cloudflare, recommended as this can limit permissions for a specific token. Conflict with `globalApiKey` and `email`.                            |
| `provider.cloudflare.globalApiKey`  | string  | Global API key for Cloudflare, not recommended as this key has full permission to access your Cloudflare account and resources. Use with `email`. Conflict with `apiToken`. |
| `provider.cloudflare.email`         | string  | Email of your Cloudflare account. Use with `globalApiKey`. Conflict with `apiToken`.                                                                                        |
| `provider.alicloud`                 | object  | Credentials and settings for AliCloud DNS provider.                                                                                                                         |
| `provider.alicloud.accessKeyId`     | string  | AccessKeyId of your AliCloud account. You can create a RAM sub user to limit permission of this access key.                                                                 |
| `provider.alicloud.accessKeySecret` | string  | AccessKeySecret of your AliCloud account.                                                                                                                                   |
| `provider.alicloud.line`            | string  | (Optional) Line of the DNS record.                                                                                                                                          |
| `provider.dnspod`                   | object  | Credentials and settings for DNSPod DNS provider.                                                                                                                           |
| `provider.dnspod.secretId`          | string  | SecretID of your Tencent Cloud account.                                                                                                                                     |
| `provider.dnspod.secretKey`         | string  | SecretKey of your Tencent Cloud account                                                                                                                                     |
| `provider.dnspod.lineId`            | string  | (Optional) ID of the line of your DNS record.                                                                                                                               |
| `provider.huawei`                   | object  | Credentials and settings for Huawei Cloud DNS provider.                                                                                                                     |
| `provider.huawei.accessKey`         | string  | Access key (AK) of the account.                                                                                                                                             |
| `provider.huawei.secretAccessKey`   | string  | Secret key (SK) of the account.                                                                                                                                             |
| `provider.huawei.region`            | string  | Region of resources and APIs.                                                                                                                                               |
| `provider.jd`                       | object  | Credentials and settings for JD Cloud DNS provider.                                                                                                                         |
| `provider.jd.accessKey`             | string  | Access key (AK) of the account.                                                                                                                                             |
| `provider.jd.secretKey`             | string  | Secret Key (SK) of the account.                                                                                                                                             |
| `provider.jd.viewId`                | number  | (Optional) View ID of the record, leave empty for default value (-1).                                                                                                       |
| `provider.rfc2136`                  | object  | Credentials and settings for [RFC 2136](https://www.ietf.org/rfc/rfc2136.txt) compatible DNS provider.                                                                      |
| `provider.rfc2136.address`          | string  | IP address or domain name of DNS server.                                                                                                                                    |
| `provider.rfc2136.port`             | number  | (Optional) Port of DNS server. Leave empty for default value (53).                                                                                                          |
| `provider.rfc2136.useTcp`           | boolean | (Optional) Specify if TCP should be used when communicating with DNS server. By default it uses UDP.                                                                        |
| `provider.rfc2136.tsig`             | object  | (Optional) Information about [RFC 2845](https://www.ietf.org/rfc/rfc2845.txt) TSIG authentication.                                                                          |
| `provider.rfc2136.tsig.keyName`     | string  | Name of TSIG key.                                                                                                                                                           |
| `provider.rfc2136.tsig.key`         | string  | TSIG key value. Should be a base64 encoded string.                                                                                                                          |
| `provider.rfc2136.gssTsig`          | object  | (Optional) Information about [RFC 3645](https://www.ietf.org/rfc/rfc3645.txt) GSS-TSIG authentication. Widely used in Windows Server DNS secured DNS update.                |
| `provider.rfc2136.gssTsig.domain`   | string  | Domain of directory service. Used in Kerberos authentication.                                                                                                               |
| `provider.rfc2136.gssTsig.username` | string  | Username used in Kerberos authentication.                                                                                                                                   |
| `provider.rfc2136.gssTsig.password` | string  | Password used in Kerberos authentication.                                                                                                                                   |
