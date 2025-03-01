# micro-ddns v0.0.3 Release Notes

v0.0.3 updated all dependencies to latest version.

## Changelog

- Updated dependencies ([#65](https://github.com/masteryyh/micro-ddns/pull/65))

## Supported DNS Providers

- [x] [Cloudflare](https://www.cloudflare.com/)
- [x] [AliCloud](https://www.aliyun.com/)
- [x] [DNSPod](https://www.dnspod.cn/)
- [x] [HuaweiCloud](https://www.huaweicloud.com/)
- [x] [JDCloud](https://www.jdcloud.com/)
- [x] [RFC2136 (experimental)](https://datatracker.ietf.org/doc/html/rfc2136)

## Supported Address Detection Methods

- [x] Interface
- [x] 3rd Party API
- [x] SSH

## Known Issues

- When using AliCloud DNS provider, you may encounter the following error: "You can’t finish this operation because the last operation has not been finished". You can try to wait for next DNS update and its not always reproducible.
