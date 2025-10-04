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
    test -f ./mnt/@/data || echo data > ./mnt/@/data

bck: mount
    sudo btrbk run -c ./btrbk.conf

clear target="meta":
    #!/usr/bin/env bash
    if [ "{{target}}" = "meta" ]; then
        sudo btrfs subvolume delete ./mnt/@snaps/.meta/*/** 2>/dev/null || true
        rm -rf ./mnt/@snaps/.meta/*
    elif [ "{{target}}" = "all" ]; then
        just clear
        sudo btrfs subvolume delete ./mnt/@snaps/*/** 2>/dev/null || true
        sudo btrfs subvolume delete ./mnt/@snaps/* 2>/dev/null || true
        rm -rf ./mnt/@snaps/***
    else
        echo "Invalid target: {{target}}. Use 'meta' or 'all'."
        exit 1
    fi
