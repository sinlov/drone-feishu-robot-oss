[![ci](https://github.com/sinlov/drone-feishu-robot-oss/workflows/ci/badge.svg?branch=main)](https://github.com/sinlov/drone-feishu-robot-oss/actions/workflows/ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/sinlov/drone-feishu-robot-oss?label=go.mod)](https://github.com/sinlov/drone-feishu-robot-oss)
[![GoDoc](https://godoc.org/github.com/sinlov/drone-feishu-robot-oss?status.png)](https://godoc.org/github.com/sinlov/drone-feishu-robot-oss)
[![goreportcard](https://goreportcard.com/badge/github.com/sinlov/drone-feishu-robot-oss)](https://goreportcard.com/report/github.com/sinlov/drone-feishu-robot-oss)

[![docker hub version semver](https://img.shields.io/docker/v/sinlov/drone-feishu-robot-oss?sort=semver)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss/tags?page=1&ordering=last_updated)
[![docker hub image size](https://img.shields.io/docker/image-size/sinlov/drone-feishu-robot-oss)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss)
[![docker hub image pulls](https://img.shields.io/docker/pulls/sinlov/drone-feishu-robot-oss)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss/tags?page=1&ordering=last_updated)

[![GitHub license](https://img.shields.io/github/license/sinlov/drone-feishu-robot-oss)](https://github.com/sinlov/drone-feishu-robot-oss)
[![codecov](https://codecov.io/gh/sinlov/drone-feishu-robot-oss/branch/FE-new-build-workflow/graph/badge.svg)](https://codecov.io/gh/sinlov/drone-feishu-robot-oss)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/sinlov/drone-feishu-robot-oss)](https://github.com/sinlov/drone-feishu-robot-oss/tags)
[![github release](https://img.shields.io/github/v/release/sinlov/drone-feishu-robot-oss?style=social)](https://github.com/sinlov/drone-feishu-robot-oss/releases)

## for what

- this project used to drone CI

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/sinlov/drone-feishu-robot-oss)](https://github.com/sinlov/drone-feishu-robot-oss/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

## Pipeline Settings (.drone.yml)

`1.x`

- notify-build-failure

```yaml
steps:
  - name: notify-failure-feishu-robot
    # image: sinlov/drone-feishu-robot-oss:latest
    # pull: if-not-exists
    image: sinlov/drone-feishu-robot-oss:latest
    settings:
      # debug: true # plugin debug switch
      # ntp_target: "pool.ntp.org" # if not set will not sync
      # timeout_second: 10 # default 10
      feishu_webhook:
        # https://docs.drone.io/pipeline/environment/syntax/#from-secrets
        from_secret: feishu_group_bot_token
      feishu_secret:
        from_secret: feishu_group_secret_bot
      feishu_msg_title: "Drone CI Notification" # default [Drone CI Notification]
      # let notification card change more info see https://open.feishu.cn/document/ukTMukTMukTM/uAjNwUjLwYDM14CM2ATN
      feishu_enable_forward: true
      drone_system_admin_token: # non-essential parameter 1.5.0+
        from_secret: drone_system_admin_token
      # ignore last success by distance
      feishu_ignore_last_success_by_admin_token_distance: 1 # if distance is 0 will not ignore, use 1 will let notify build change to success
    when:
      event: # https://docs.drone.io/pipeline/exec/syntax/conditions/#by-event
        - promote
        - rollback
        - push
        - pull_request
        - tag
      status: # only support failure/success, both open will send anything
        - failure
        # - success
```

# Features

- more see [features/README.md](features/README.md)

## env

- minimum go version: go 1.18
- change `go 1.18`, `^1.18`, `1.18.10` to new go version

### libs

| lib                                        | version |
|:-------------------------------------------|:--------|
| https://github.com/stretchr/testify        | v1.8.4  |
| https://github.com/sebdah/goldie           | v2.5.3  |

- more see [go.mod](go.mod)

# dev

## depends

in go mod project

```bash
# warning use private git host must set
# global set for once
# add private git host like github.com to evn GOPRIVATE
$ go env -w GOPRIVATE='github.com'
# use ssh proxy
# set ssh-key to use ssh as http
$ git config --global url."git@github.com:".insteadOf "https://github.com/"
# or use PRIVATE-TOKEN
# set PRIVATE-TOKEN as gitlab or gitea
$ git config --global http.extraheader "PRIVATE-TOKEN: {PRIVATE-TOKEN}"
# set this rep to download ssh as https use PRIVATE-TOKEN
$ git config --global url."ssh://github.com/".insteadOf "https://github.com/"

# before above global settings
# test version info
$ git ls-remote -q https://github.com/sinlov/drone-feishu-robot-oss.git

- test code

add env then test

```bash
export PLUGIN_MSG_TYPE=post \
  export PLUGIN_WEBHOOK=7138d7b3-abc
```

```bash
make test
```

- see help

```bash
make dev
```

update main.go file set env then and run

```bash
export PLUGIN_MSG_TYPE= \
  export PLUGIN_WEBHOOK= \
  export DRONE_REPO=sinlov/drone-feishu-robot-oss \
  export DRONE_REPO_NAME=drone-feishu-robot-oss \
  export DRONE_REPO_NAMESPACE=sinlov \
  export DRONE_REMOTE_URL=https://github.com/sinlov/drone-feishu-robot-oss \
  export DRONE_REPO_OWNER=sinlov \
  export DRONE_COMMIT_AUTHOR=sinlov \
  export DRONE_COMMIT_AUTHOR_AVATAR=  \
  export DRONE_COMMIT_AUTHOR_EMAIL=sinlovgmppt@gmail.com \
  export DRONE_COMMIT_BRANCH=main \
  export DRONE_COMMIT_LINK=https://github.com/sinlov/drone-feishu-robot-oss/commit/68e3d62dd69f06077a243a1db1460109377add64 \
  export DRONE_COMMIT_SHA=68e3d62dd69f06077a243a1db1460109377add64 \
  export DRONE_COMMIT_REF=refs/heads/main \
  export DRONE_COMMIT_MESSAGE="mock message commit" \
  export DRONE_STAGE_STARTED=1674531206 \
  export DRONE_STAGE_FINISHED=1674532106 \
  export DRONE_BUILD_STATUS=success \
  export DRONE_BUILD_NUMBER=1 \
  export DRONE_BUILD_LINK=https://drone.xxx.com/sinlov/drone-feishu-robot-oss/1 \
  export DRONE_BUILD_EVENT=push \
  export DRONE_BUILD_STARTED=1674531206 \
  export DRONE_BUILD_FINISHED=1674532206
```

- then run

```bash
make run
```

## docker

```bash
# then test build as test/Dockerfile
$ make dockerTestRestartLatest
# if run error
# like this error
# err: missing webhook, please set webhook
#  fix env settings then test

# see run docker fast
$ make dockerTestRunLatest

# clean test build
$ make dockerTestPruneLatest

# see how to use
$ docker run --rm sinlov/drone-feishu-robot-oss:latest -h
```

## use

- use to replace
  `sinlov/drone-feishu-robot-oss` to you code
