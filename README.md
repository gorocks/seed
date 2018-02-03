seed
===

seed is a command-line tool to quick start  application.

## Requirements

- Go version >= 1.3.

## Installation

To install `seed` use the `go get` command:

```bash
go get github.com/Guazi-inc/seed
```

Then you can add `seed` binary to PATH environment variable in your `~/.bashrc` or `~/.bash_profile` file:

```bash
export PATH=$PATH:<your_main_gopath>/bin
```

> If you already have `seed` installed, updating `seed` is simple:

```bash
go get -u github.com/Guazi-inc/seed
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
seedVersion v0.0.1

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
seedVersion v0.0.1
seedVersion:0.0.1
2018/02/03 14:32:52 INFO     ▶ 0001 Creating application...
2018/02/03 14:32:52 SUCCESS  ▶ 0002 create dir:/go/src/github.com/Guazi-inc/seed/explame/
2018/02/03 14:32:52 SUCCESS  ▶ 0004 create file:/go/src/github.com/Guazi-inc/seed/explame/.gitgnore
2018/02/03 14:32:52 SUCCESS  ▶ 0008 create file:/go/src/github.com/Guazi-inc/seed/explame/README.md
2018/02/03 14:32:52 SUCCESS  ▶ 0009 create dir:/go/src/github.com/Guazi-inc/seed/explame/cmd/consumer/
2018/02/03 14:32:52 SUCCESS  ▶ 0011 create dir:/go/src/github.com/Guazi-inc/seed/explame/cmd/grpcserver/
2018/02/03 14:32:52 SUCCESS  ▶ 0013 create dir:/go/src/github.com/Guazi-inc/seed/explame/cmd/grpcweb/
2018/02/03 14:32:52 SUCCESS  ▶ 0015 create dir:/go/src/github.com/Guazi-inc/seed/explame/databases/
2018/02/03 14:32:52 SUCCESS  ▶ 0024 create file:/go/src/github.com/Guazi-inc/seed/explame/gip.yml
2018/02/03 14:32:52 SUCCESS  ▶ 0026 create file:/go/src/github.com/Guazi-inc/seed/explame/gometalinter.json
2018/02/03 14:32:52 SUCCESS  ▶ 0032 create dir:/go/src/github.com/Guazi-inc/seed/explame/model/
2018/02/03 14:32:52 SUCCESS  ▶ 0036 create file:/go/src/github.com/Guazi-inc/seed/explame/requirements.txt
2018/02/03 14:32:52 SUCCESS  ▶ 0037 create dir:/go/src/github.com/Guazi-inc/seed/explame/service/
2018/02/03 14:32:52 SUCCESS  ▶ 0039 New application successfully created!

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