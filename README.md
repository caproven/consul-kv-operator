# consul-kv-operator
// TODO(user): Add simple overview of use/purpose

## Run Consul
Consul can be installed in k8s through a helm chart. Run the below command:

```bash
helm install consul hashicorp/consul  --create-namespace --namespace consul --version "0.43.0" --values consul-config.yaml
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

Consul can also be accessed through an [API](https://www.consul.io/api-docs).

Port-forward the Consul service to your local machine:

```bash
kubectl port-forward svc/consul-consul-server 8500:8500 &
```

Example usage:

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
    --header "X-Consul-Token: <consul token>"
    <rest of request>
    ```