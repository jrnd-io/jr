## JR and Locust distributed in Kubernetes

JR can run with  [Locust K6](https://locust.io)

### Prerequisites

- a Kubernetes cluster
- the `jr` image accessible either as built locally or from a registry
- [`helm`](https://helm.sh) client


###  Build the image

Build the customised `jr` and `locust` image with:

```
docker build -t locust-jr .
```

### Create the configmap for the jr locust library

Create a configmap  with:

```
kubectl create configmap locust-jr-lib --from-file jr.py
```


### Example with jr and the MongoDB Producer

Create a configmap with the MongoDB configuration and the locust python file:


```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mongodb-jr
data:
  mongoconfig.json: |
    {
      "mongo_uri": "mongodb+srv://...","
      "database": "[mongodb database]",
      "collection": "[mongodb collection]"
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: locust-jr
data:
  main.py: |
    from locust import task
    from lib.jr import JRUser

    class MyUser(JRUser):
        @task
        def jr_task(self):
            self.run_jr(["run", 
                        "net_device",
                        "--output","mongo",
                        "--mongoConfig", "/jrconfig/mongoconfig.json"])
```

Add the `deliveryhero` repo: 

```
https://github.com/deliveryhero/helm-charts/tree/master/stable/locust
```

Install the [Locust Helm chart](https://github.com/deliveryhero/helm-charts/tree/master/stable/locust) with:

```bash
helm install locust deliveryhero/locust \
   --set loadtest.name=jr-test-01 \
   --set loadtest.locust_locustfile_configmap=locust-jr \
   --set loadtest.locust_lib_configmap=locust-jr-lib \
   --set image.pullPolicy=Always \
   --set worker.image=registry.localhost:5000/locust-jr \
   --set worker.logLevel=DEBUG \
   --set worker.replicas=4  \
   --set extraConfigMaps.mongodb-jr="/jrconfig"
```

Expose the port `8089` of the locust master pod via port forward, connect with the browser to http://localhost:8089 and launch the test (ignore the Host field in the form)

