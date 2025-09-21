prepare:
    go mod download
    go run github.com/hairyhenderson/gomplate/v4/cmd/gomplate -f btrbk.conf.tmpl -o btrbk.conf

mount:
    mountpoint -q ./mnt || sudo mount LABEL=BTRBK-manage-part ./mnt
