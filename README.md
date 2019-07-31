# sspanel_checkin

## Requirements
- [Golang](https://golang.org/doc/install) with `$GOPATH` set to your preferred directory

## Build
1. Get source code

```
$ git clone https://github.com/iczc/ss-panel-checkin.git $GOPATH/src/github.com/iczc/ss-panel-checkin
```
or
```
$ go get github.com/iczc/ss-panel-checkin
```

2. Compile the project

```
cd $GOPATH/src/github.com/iczc/ss-panel-checkin
go build
```

## Usage

```
./ss-panel-checkin -host <host> -email <email> -passwd <password>
```