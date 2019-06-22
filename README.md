# Static file filter for Golang httpd
Implementation of http.FileSystem interface to enable filtering of static files served with http.FileServer

## Installation

Install:

```shell
go get -u github.com/ahakkila/staticfs
```

Import:

```go
import "github.com/ahakkila/staticfs"
```

## Quickstart

```go

type SFSOption int

const (
	DirectoryFilter SFSOption 
	DotFileFilter
)

type StaticFileSystem struct {
	fs      http.FileSystem
	options SFSOption
}

func NewStaticFileSystem(path string, options SFSOption) StaticFileSystem 
// path - path to static resource root
// options - binary OR of desired Filter options

```

* currently supports two options
  * **DirectoryFilter** - do not display directory contents, directory URLs pointing to resources not containing index.html will result in Not Found
  * **DotFileFilter** - if any segment of the requested URL path is dot prefixed, the request will result in Not Found, e.g.
    * /path/to/*.htaccess* - Not Found
    * /path/to/*.git*/subdirectory/resource - Not Found
    * /*.secret* - Not Found
    * /request/to/*..*/relative/url - Not Found

```go
sfs := staticfs.NewStaticFileSystem("static", (staticfs.DirectoryFilter | staticfs.DotFileFilter))
http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(sfs)))
```

* 

