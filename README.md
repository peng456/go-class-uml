[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/peng456/goclassuml/parser) [![Go Report Card](https://goreportcard.com/badge/github.com/peng456/goclassuml)](https://goreportcard.com/report/github.com/peng456/goclassuml) [![codecov](https://codecov.io/gh/peng456/goclassuml/branch/master/graph/badge.svg)](https://codecov.io/gh/peng456/goclassuml) [![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/peng456/goclassuml.svg)](https://github.com/peng456/goclassuml/releases/)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 
[![DUMELS Diagram](https://www.dumels.com/api/v1/badge/23ff0222-e93b-4e9f-a4ef-4d5d9b7a5c7d)](https://www.dumels.com/diagram/23ff0222-e93b-4e9f-a4ef-4d5d9b7a5c7d) 
# goclassuml

goclassuml is an open-source tool developed to streamline the process of generating go class diagrams from Go source code. With goclassuml, developers can effortlessly visualize the structure and relationships within their Go projects, aiding in code comprehension and documentation. By parsing Go source code and producing diagrams svg .eg, goclassuml empowers developers to create clear and concise visual representations of their codebase architecture, package dependencies, and function interactions. This tool simplifies the documentation process and enhances collaboration among team members by providing a visual overview of complex Go projects. goclassuml is actively maintained and welcomes contributions from the Go community.


## Code of Conduct
Please, review the code of conduct [here](https://github.com/peng456/goclassuml/blob/master/CODE_OF_CONDUCT.md "here").

### Prerequisites
golang 1.17 or above

### Installing

```
go get github.com/peng456/goclassuml/parser
go install github.com/peng456/goclassuml/cmd/goclassuml@latest
```

This will install the command goclassuml in your GOPATH bin folder.

### Usage

```
goclassuml [-recursive] path/to/gofiles /tmp/gofiles.svg
```


Usage of goclassuml:
  -aggregate-private-members
        Show aggregations for private members. Ignored if -show-aggregations is not used.
  -hide-connections
        hides all connections in the diagram
  -hide-fields
        hides fields
  -hide-methods
        hides methods
  -ignore string
        comma separated list of folders to ignore
  -notes string
        Comma separated list of notes to be added to the diagram
  -output string
        output file path. If omitted, then this will default to standard output
  -recursive
        walk all directories recursively
  -show-aggregations
        renders public aggregations even when -hide-connections is used (do not render by default)
  -show-aliases
        Shows aliases even when -hide-connections is used
  -show-compositions
        Shows compositions even when -hide-connections is used
  -show-connection-labels
        Shows labels in the connections to identify the connections types (e.g. extends, implements, aggregates, alias of
  -show-implementations
        Shows implementations even when -hide-connections is used
  -show-options-as-note
        Show a note in the diagram with the none evident options ran with this CLI
  -title string
        Title of the generated diagram
  -hide-private-members
        Hides all private members (fields and methods)
```

#### Example
```
goclassuml $GOPATH/src/github.com/peng456/goclassuml/parser
```
```
// echoes

@mermaid
classDiagram

    class apiSubsetResizeNotifyController {
        <<struct>>
        - ctx *gin.Context

        + Req *pb_gen.ApiSubsetResizeNotifyReq
        + Resp interface

    }
    

    Handler --|> apiSubsetResizeNotifyController : implements



goclassuml $GOPATH/src/github.com/peng456/goclassuml/parser > ClassDiagram.mmd
// Generates a file ClassDiagram.mmd