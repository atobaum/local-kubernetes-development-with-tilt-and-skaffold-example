load('ext://restart_process', 'docker_build_with_restart')

k8s_context("k3d-local-dev")

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main ./'

local_resource(
    'deploy',
    'python3 record-start-time.py',
)

local_resource(
  "build",
  compile_cmd,
  deps=["./main.go"],
  resource_deps=["deploy"]
)

docker_build_with_restart(
  "myapp",
  ".",
  entrypoint=["/app/main"],
  dockerfile="./tilt/Dockerfile-live-update",
  only = [
    "./build"
  ],
  live_update=[
    sync("./build/main", "/app/main"),
  ]
)

# docker_build(
#   "myapp",
#   ".",
# )

k8s_yaml("k8s/deployment.yaml")
k8s_resource("myapp", port_forwards=8080)