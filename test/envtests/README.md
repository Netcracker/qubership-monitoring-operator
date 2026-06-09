# Envtests

## How to run it locally

### Run on the workstation

To configure your local environment you can use the `envtest`:

[https://book.kubebuilder.io/reference/envtest.html](https://book.kubebuilder.io/reference/envtest.html)

### Run locally inside the container

Another way to run integration tests locally use a container and run tests inside it if you don't want to
install all dependencies on your system.

Requirements to use this way:

* Clone repository with `monitoring-operator`
* Docker CE or other container runtime (podman, rancher-desktop, etc)
* The container with `kubebuilder`, `kubelet`, `kube-apiserver` and `etcd`

To run integration tests inside the container you need:

1. Navigate to `monitoring-operator` directory
2. Run the container with mounts:

    ```bash
    docker run -ti --rm -v ./:/operator <image>
    ```

3. After container will start, inside the container navigate to `/operator` dir:

    ```bash
    cd /operator
    ```

4. Run tests using the command:

    ```bash
    $GOPATH/bin/ginkgo --fail-on-pending --trace --race --compilers=2 --grace-period=10m -v ./bdd-tests/...
    ```
