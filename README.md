# mibparser

This project is used to convert mib file to json format. <br/>
The created json format is in tree format. <br/>
This code takes the nodes id and the nodes parent id and finally creates an OID with this information <br/>
Informs the user if another mib file is needed for the created json format <br/>

### Requirements
There must be required mib files in the specified path <br/>

### Installation

```
go get github.com/limanmys/mib-parser-go
```

### Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/limanmys/mib-parser-go"
)

func main() {
    mibparserObject, err := mibparser.Load(mibparser.NewPath("./mib-files"))
    if err != nil {
        log.Fatalf("error when loading mibparser path")
    }
    mibTree, err := mibparserObject.GetJSONTree()
    if err != nil {
        log.Fatalf("error when parsing mib files, err " + err.Error())
    }
    fmt.Println("mibTree ", mibTree)

    mibObjects, err := mibparserObject.GetObjects()
    if err != nil {
        log.Fatalf("error when parsing mib files, err " + err.Error())
    }
    fmt.Println("mibObjects", mibObjects)
}


```
# Created by

<img src="https://avatars.githubusercontent.com/u/63673212?s=280&v=4" alt="Logo" width="150" height="150">

