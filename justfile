default:
    just pix
pix:
    mkdir -p bin
    go build -o bin/pix ./cmd/pix
wallpaper-finder:
    mkdir -p bin
    go build -o bin/wallpaper-finder ./cmd/wallpaper-finder
