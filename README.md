# k2hash_go

### Overview

**k2hash_go** implements a [k2hash](https://k2hash.antpick.ax/) client in golang.

### Install

Firstly you must install the [k2hash](https://k2hash.antpick.ax/) shared library.
```
$ curl -o- https://raw.github.com/yahoojapan/k2hash_go/master/utils/libk2hash.sh | bash
```
You can install **k2hash** library step by step from [source code](https://github.com/yahoojapan/k2hash). See [Build](https://k2hash.antpick.ax/build.html) for details.

After you make sure you set the [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) environment, download the **k2hash_go** package.
```
$ go get -u github.com/yahoojapan/k2hash_go
```

### Usage

Here is a simple example of **k2hash_go** which save a key and get it.

```golang
package main

import (
        "fmt"
        "os"

        "github.com/yahoojapan/k2hash_go/k2hash"
)

func SetAndGet() {
        // 1. Instantiate K2hash class
        file, _ := k2hash.NewK2hash("/tmp/test.k2h")
        defer file.Close()
        ok, err := file.Set("hello", "world")
        if ok != true {
                fmt.Fprintf(os.Stderr, "file.Set(hello, world) returned false, err %v\n", err)
        }
        // 2. Get
        // 2.1. no args
        val, err := file.Get("hello")
        if val == nil || err != nil {
                fmt.Fprintf(os.Stderr, "file.Get(hello) returned val %v err %v\n", val, err)
                return
        }
        fmt.Printf("val = %v, err = %v\n", val.String(), err)
}

func main() {
        SetAndGet()
}
```

### Development

Here is the step to start developing **k2hash_go**.

- Debian / Ubuntu

```bash
#!/bin/sh

sudo apt-get update -y && sudo apt-get install curl git -y && curl -s https://packagecloud.io/install/repositories/antpickax/stable/script.deb.sh | sudo bash
sudo apt-get install libfullock-dev k2hash-dev -y
go get github.com/yahoojapan/k2hash_go/k2hash

exit 0
```

- CentOS / Fedora

```bash
#!/bin/sh

sudo dnf makecache && sudo yum install curl git -y && curl -s https://packagecloud.io/install/repositories/antpickax/stable/script.rpm.sh | sudo bash
sudo dnf install libfullock-devel k2hash-devel -y
go get github.com/yahoojapan/k2hash_go/k2hash

exit 0
```

### Documents
  - [About K2HASH](https://k2hash.antpick.ax/)
  - [About AntPickax](https://antpick.ax/)

### License

MIT License. See the LICENSE file.

## AntPickax

[AntPickax](https://antpick.ax/) is 
  - an open source team in [Yahoo Japan Corporation](https://about.yahoo.co.jp/info/en/company/). 
  - a product family of open source software developed by [AntPickax](https://antpick.ax/).

