## micro-ddns
micro-ddns is a tool that can update your DNS record
dynamically based on current IP address, support multiple
IP address detection methods and DNS providers.

## Quick Start
```yaml
# /etc/micro-ddns/config.yaml
ddns:
  - name: homelab
    domain: yourdomain.com
    subdomain: test
    stack: IPv4
    cron: "*/30 * * * *"
    detection:
      type: ThirdParty
      api:
        url: https://api.ipify.org
    provider:
      name: Cloudflare
      cloudflare:
        apiToken: "<redacted>"
```

```
$ micro-ddns run --config /etc/micro-ddns/config.yaml
```

OR you can run as a container:

```
$ docker run --name ddns -d -v /path/to/config.yaml:/etc/micro-ddns/config.yaml masteryyh/micro-ddns:alpine
```

## License
This project is licensed under the Apache License 2.0. For more details, see the LICENSE file in the repository.

Copyright (c) 2024 masteryyh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
