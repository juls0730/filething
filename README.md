# Filething

**dependencies**:

- bun
- go

## Getting started

To run filething, run

```BASH
bun --cwd=./ui install
RENDERING_MODE=static bun --bun --cwd=./ui run generate
go build -tags netgo -ldflags=-s
DB_HOST=localhost:5432 DB_NAME=filething DB_USER=postgres STORAGE_PATH=data ./filething
```

Or if you want to run filething with SSR (you will need node on the target server), run

```BASH
bun --cwd=./ui install
bun --bun --cwd=./ui run build
go build -tags netgo,ssr -ldflags=-s
DB_HOST=localhost:5432 DB_NAME=filething DB_USER=postgres STORAGE_PATH=data ./filething
```

### Contributing

To run filething in dev mode with a hot reloading Ui server and auto rebuilding backend server, run

```BASH
DB_HOST=localhost:5432 DB_NAME=filething DB_USER=postgres STORAGE_PATH=data CompileDaemon --build="go build -tags netgo,dev -ldflags=-s" --command=./filething --exclude-dir=data/ --exclude-dir=ui/ --graceful-kill
```
