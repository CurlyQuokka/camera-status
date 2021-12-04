# camera-status

This is simple Go program that helps me to monitor health of my CCTV camera. It simply checks the camera daemon status, if the recordings are up-to-date and if there is sufficient space available in the filesystem, then serves that info via HTTP server.

## build

```
go build -o camera-status cmd/camera-status/main.go
```

## usage

```
./camera-status <http_port>
```

