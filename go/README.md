# go

This is the documentation for the GO implementation of the tool.


## Requirements

* Go: 1.18


## Build

Change into the project folder with a CLI and build the application via

```go build pomplanter.go pomreader.go```

or just

```go build .```



## Run

Change into the project folder with a CLI and start the application either by

```go run pomplanter.go pomreader.go {pomFileName}```

or 

```go run . {pomFileName}```

or 

```.\pomplanter.exe {pomFileName}```

in case the project is already build.

The PlantUML content will be send to the console output and could be copied to a PlantUML viewer easily.