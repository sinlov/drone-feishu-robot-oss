[![go-ubuntu](https://github.com/sinlov/drone-feishu-robot-oss/workflows/go-ubuntu/badge.svg?branch=main)](https://github.com/sinlov/drone-feishu-robot-oss/actions)
[![GoDoc](https://godoc.org/github.com/sinlov/drone-feishu-robot-oss?status.png)](https://godoc.org/github.com/sinlov/drone-feishu-robot-oss/)
[![GoReportCard](https://goreportcard.com/badge/github.com/sinlov/drone-feishu-robot-oss)](https://goreportcard.com/report/github.com/sinlov/drone-feishu-robot-oss)
[![codecov](https://codecov.io/gh/sinlov/drone-feishu-robot-oss/branch/main/graph/badge.svg)](https://codecov.io/gh/sinlov/drone-feishu-robot-oss)
[![docker version semver](https://img.shields.io/docker/v/sinlov/drone-feishu-robot-oss?sort=semver)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss/tags?page=1&ordering=last_updated)
[![docker image size](https://img.shields.io/docker/image-size/sinlov/drone-feishu-robot-oss)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss)
[![docker pulls](https://img.shields.io/docker/pulls/sinlov/drone-feishu-robot-oss)](https://hub.docker.com/r/sinlov/drone-feishu-robot-oss/tags?page=1&ordering=last_updated)

## for what

- this project used to drone CI

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

### cli tools to init project fast

```
$ curl -L --fail https://raw.githubusercontent.com/sinlov/drone-feishu-robot-oss/main/drone-feishu-robot-oss
# let drone-feishu-robot-oss file folder under $PATH
$ chmod +x drone-feishu-robot-oss
# see how to use
$ drone-feishu-robot-oss -h
```