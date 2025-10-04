mount:
    mountpoint -q ./mnt || sudo mount LABEL=BTRBK-manage-part ./mnt

build:
    go build -o bin/ ./cmd/...

run cmd:
    go run ./cmd/btrbk-{{cmd}}/main.go

test:
    go fmt ./...
    go mod tidy
    go test -cover -timeout=1s -race ./...

prepare: mount
    go mod download
    go run github.com/hairyhenderson/gomplate/v4/cmd/gomplate -f btrbk.conf.tmpl -o btrbk.conf

    sudo btrfs subvolume create ./mnt/@ ./mnt/@home ./mnt/@snaps
    sudo chown -R $(id -u):$(id -g) ./mnt
    touch ./mnt/.gitkeep
    echo data > ./mnt/@/data

bck: mount
    sudo btrbk run -c ./btrbk.conf
