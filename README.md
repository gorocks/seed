seed
===

seed is a command-line tool to quick start  application.

## Installation

To install `seed` use the `go get` command:

```bash
go get github.com/Guazi-inc/seed
```

## Basic commands

seed provides a variety of commands which can be helpful at various stages of development. The top level commands include:

```
    version     Prints the current seed version
    new         Creates a  app for template
    httptest    set up a http server for test

```

### seed version

To display the current version of `seed`, `seedgo` and `go` installed on your machine:

```bash
$ seed version
seedVersion:0.0.1
├── GoVersion : go1.9
├── GOOS      : linux
├── GOARCH    : amd64
├── NumCPU    : 2
├── GOPATH    : /Users/user/.go
├── GOROOT    : /usr/local/go
├── Compiler  : gc
└── Date      : Saturday, 3 Feb 2018
```
For more information on the usage, run `seed help version`.

### seed new

To create a new seedgo web application:

```bash
$ seed new my-web-app -tp="template/path"
seedVersion:0.0.1
2018/02/04 09:07:49 [INFO]    : Creating application...
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/.gitgnore
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/.gitlab-ci.yml
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/README.md
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/consumer/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/consumer/main.go
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/grpcserver/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/grpcserver/main.go
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/grpcweb/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/cmd/grpcweb/main.go
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/databases/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/databases/init-tables.py
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/databases/init.sql
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/fixtures/apply/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/fixtures/apply/user.yml
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/gometalinter.json
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/model/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/model/user.go
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/model/user_test.go
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/requirements.txt
2018/02/04 09:07:49 [SUCCESS] : create dir:/$GOPATH/src/github.com/Guazi-inc/seed/explame/service/
2018/02/04 09:07:49 [SUCCESS] : create file:/$GOPATH/src/github.com/Guazi-inc/seed/explame/service/preaudit-service.go
2018/02/04 09:07:49 [SUCCESS] : New application successfully created!


```

For more information on the usage, run `seed help new`.

## Help

To print more information on the usage of a particular command, use `seed help <command>`.

For instance, to get more information about the `run` command:

```bash
$ seed help new
USAGE
  Seed new -n=[appname] -tp=[template path]

OPTIONS
  -g=finance
      this application belong which group
  
  -n
      set a name for application
  
  -tn=eipis-apply
      template name,use which template
  
  -tp
      template path
  
DESCRIPTION
  Creates a  application for the given app name and template in the current directory.
```