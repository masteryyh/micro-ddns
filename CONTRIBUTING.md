# How to contribute

First of all, thanks for reading this and trying to contribute to this repo.

## Testing

As we currently don't have a working unit testing or integrated testing guideline or instance, a full test is recommended before committing, and we still recommend committing with a working unit test.

## Committing

Please fork this repo first, create and commit to your own forked repo's branch, and open a Pull Request that we can review and merge your branch.

Please include a short and precise description about what you have done in the commit message.

We strongly recommend use GPG to sign your commit, and add a sign-off message after your commit message. We enabled vigilant mode so we can make sure that all commits are committed by yourself.

```
$ git commit -s -S -m "what-i-have-done"
```

For now we don't force contributors to sign their commits, but in the future we might refuse Pull Requests that commits are not signed.

## Coding conventions

1. We indent using Tab, you can choose any tab size you want when you browsing or writing code, but always choose Tab for indent;

2. All Go source files must add an Apache 2.0 license header comment, you can add your name and email under `Contributors` section like this:
```
/*
Copyright © 2024 masteryyh <yyh991013@163.com>

Contributors:
- your-self <your-email@email.com>
- another-contributor <blahblah@some-email.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
```

3. All Go source files about to commit must have `go vet` and `go fmt` processed.

## Important resources

### Committing with GPG signature and sign-off message

1. Check out the [Git's official documents](https://git-scm.com/docs/git-commit).

2. What is, and how to enable vigilant mode in your GitHub account: [Displaying verification statuses for all of your commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/displaying-verification-statuses-for-all-of-your-commits).

3. How to generate or use your existing GPG key:

    - [Generating a new GPG key](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key)
    
    - [Checking for existing GPG keys](https://docs.github.com/en/authentication/managing-commit-signature-verification/checking-for-existing-gpg-keys)

    - [Adding a GPG key to your GitHub account](https://docs.github.com/en/authentication/managing-commit-signature-verification/adding-a-gpg-key-to-your-github-account)

    - [Telling Git about your signing key](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key)

    - [Signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)

### DNS provider documentations

- AliCloud:

    - API: https://www.alibabacloud.com/help/en/dns/api-alidns-2015-01-09-overview
    - SDK: https://next.api.aliyun.com/api-tools/sdk/Alidns?version=2015-01-09&language=java-async-tea&tab=primer-doc

- Cloudflare:

    - API: https://developers.cloudflare.com/api/
    - SDK: https://developers.cloudflare.com/fundamentals/api/reference/sdks/

- DNSPod:

    - API: https://cloud.tencent.com/document/api/1427/56194
    - SDK: https://cloud.tencent.com/document/sdk/Go

- HuaweiCloud:

    - API: https://console.huaweicloud.com/apiexplorer/#/openapi/DNS/doc
    - SDK: https://console.huaweicloud.com/apiexplorer/#/sdkcenter/DNS?lang=Go

- JDCloud:

    - API: https://docs.jdcloud.com/cn/jd-cloud-dns/api/introduction
    - SDK: https://docs.jdcloud.com/cn/sdk/go

- RFC2136:

    As RFC2136 is actually a standard, not a DNS provider, we can only provide IETF documentation about this: https://datatracker.ietf.org/doc/html/rfc2136