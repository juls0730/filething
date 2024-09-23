# Filething

**dependencies**:

- bun
- go

## Getting started

To run filething, run

```BASH
go generate
go build -tags netgo -ldflags=-s
DB_HOST=localhost:5432 DB_NAME=filething DB_USER=postgres STORAGE_PATH=data ./filething
```

### Contributing

To run filething in dev mode with a hot reloading Ui server and auto rebuilding backend server, run

```BASH
DB_HOST=localhost:5432 DB_NAME=filething DB_USER=postgres STORAGE_PATH=data CompileDaemon --build="go build -tags netgo,dev -ldflags=-s" --command=./filething --exclude-dir=data/ --exclude-dir=ui/ --graceful-kill
```
