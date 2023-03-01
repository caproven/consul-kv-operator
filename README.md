# consul-kv-operator

## Run Consul
Consul can be installed in k8s through a helm chart. Run the below command(s):

```bash
# Only needed first time for machine
helm repo add hashicorp https://helm.releases.hashicorp.com

# Install on k8s cluster
helm install consul hashicorp/consul --create-namespace --namespace consul --values consul-config.yaml
```

This installs Consul into its own namespace, `consul`.

## Access Consul

### Direct Access

You can access the Consul CLI by exec'ing into the Consul server pod.

```bash
kubectl exec --stdin --tty consul-consul-server-0 --namespace consul -- /bin/sh
```

Here, you have access to the CLI. Run `consul help` for more info.

### Remote Access

Consul can also be accessed through an [API](https://www.consul.io/api-docs) or a [CLI](https://www.consul.io/downloads).

Port-forward the Consul service to your local machine:

```bash
kubectl port-forward -n consul svc/consul-server 8500:8500 &
```

Example usage of API:

* Lookup a value

    ```bash
    curl -k http://localhost:8500/v1/kv/<key>
    ```

* Write a value

    ```bash
    curl \
    --request PUT \
    --data '<value>' \
    http://localhost:8500/v1/kv/<key>
    ```

* Authenticated

    ```bash
    curl \
    --header "X-Consul-Token: <token>"
    <rest of request>
    ```

Example usage of CLI:

* Lookup a value

    ```bash
    consul kv get <key>
    ```

* Write a value

    ```bash
    consul kv put <key> <value>
    ```

* Authenticated

    ```bash
    consul <command> -token='<token>'
    ```

## Example

Make sure Consul is deployed & port-forwarded. You need access to Consul before the below example may be followed.

1. Write some values to Consul

    ```bash
    consul kv put foo "foo_val"
    consul kv put bar "bar_val"
    ```

2. Run the operator

    ```bash
    make install
    make run
    ```

3. Apply the sample KVSecret custom resource

    ```bash
    kubectl apply -k config/samples/
    ```

4. Check that the secret `foobar-secret` contains the expected content from Consul

5. Update a Consul value

    ```bash
    consul kv put bar "new_bar_val"
    ```

6. Check that the secret `foobar-secret` was updated. May take a few seconds for the operator to refresh the secret
