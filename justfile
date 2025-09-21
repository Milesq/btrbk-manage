mount:
    mountpoint -q ./mnt || sudo mount LABEL=BTRBK-manage-part ./mnt

test:
    go fmt ./...
    go mod tidy
    go test -cover -timeout=1s -race ./...

prepare:
    go mod download
    go run github.com/hairyhenderson/gomplate/v4/cmd/gomplate -f btrbk.conf.tmpl -o btrbk.conf
