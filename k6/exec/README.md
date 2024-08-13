## JR and K6 distributed in Kubernetes


JR can run with [xk6-exec](https://github.com/grafana/xk6-exec) with the [distributed K6](https://grafana.com/docs/k6/latest/set-up/set-up-distributed-k6/) in Kubernetes.

### Prerequisites

- a Kubernetes cluster
- the [K6 operator](https://grafana.com/docs/k6/latest/set-up/set-up-distributed-k6/) installed 
- the `jr` image accessible either as built locally or from a registry


### Build the image

From the `/k6/exec` folder launch:

```bash
docker build -t k6-jr .
```

The `k6-jr` image can launch a script that executes `jr` in the runner.
e.g.:

```javascript
import exec from 'k6/x/exec';

export default ()=> {
       console.log(exec.command("/usr/bin/jr", ["run",
                                                "net_device",
                                                "-n", "2",
                                                "-f", "100ms",
                                                "-d", "1m"
                                                ]))

}
```

### Create the CR

You can now create the CR from the k6 CRD to run the test.
The following example runs `jr` with a parallelism of `4` and `200` virtual users (e.g. 50 VUs each pod)


```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: jr-test
data:
  jr.js: |
    import exec from 'k6/x/exec';

    export const options = {
      vus: 20,
      duration: '5m',
    };
    export default ()=> {
      try{
          var output = exec.command("/usr/bin/jr", ["run",
                                                    "net_device",
                                                    ],{"continue_on_error": true});
      }catch (e) {
        console.log("ERROR: " + e);
        if (e.value && e.value.stderr) {
                console.log("STDERR: " + String.fromCharCode.apply(null, e.value.stderr))
        }
      }


    }    
---
apiVersion: k6.io/v1alpha1
kind: TestRun
metadata:
  name: k6-jr-example
spec:
  parallelism: 4
  script:
    configMap:
      name: jr-test
      file: jr.js
  runner:
    image: registry.localhost:5000/k6-jr:latest
```    


> Note: `k6` is responsible for running `jr` with the proper number pf virtual users, duration  and parallelism so `jr` in the script should be run emitting just one sample

To run the example:

```bash
kubectl apply -f test-jr-run.yaml
```

The k6 operator should launch 4 pods with `jr` runnning for the amount of time declared in the `options` variable of the `javascript`  test script.


### A more complex example with MongoDB producer  and Prometheus remote RW

In this example we will run a distributed test with `jr` using as producer `MongoDB` and writing the `k6` metrics to `Prometheus`. 

The prerequisites are:

- a reachable `MongoDB` cluster instance 
- a reachable `Prometheus` instance with the remote write receiver enabled 

The `yaml` (self explanatory) of the TestRun CR is the following:


```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: jr-configmap
data:
  mongoconfig.json: |
    {
      "mongo_uri": "mongodb+srv://<user>:<password>@...",
      "database": "<mongo database>",
      "collection": "<mongo collection>"
    }
  jr.js: |
    import exec from 'k6/x/exec';

    export const options = {
      vus: 20,
      duration: '5m',
    };
    export default ()=> {
      try{
          var output = exec.command("/usr/bin/jr", ["run",
                                                    "net_device",
                                                    "--output","mongo",
                                                    "--mongoConfig", "/jrconfig/mongoconfig.json"
                                                    ],{"continue_on_error": true});
      }catch (e) {
        console.log("ERROR: " + e);
        if (e.value && e.value.stderr) {
                console.log("STDERR: " + String.fromCharCode.apply(null, e.value.stderr))
        }
      }


    }
---
apiVersion: k6.io/v1alpha1
kind: TestRun
metadata:
  name: k6-jr-example
spec:
  parallelism: 4
  script:
    configMap:
      name: jr-configmap
      file: jr.js
  arguments: -o experimental-prometheus-rw
  runner:
    image: registry.localhost:5000/k6-jr:latest
    env: 
      - name: K6_PROMETHEUS_RW_SERVER_URL
        value: http://<prometheus endpoint>/api/v1/write
      - name: K6_PROMETHEUS_RW_TREND_STATS
        value: "p(95),p(99),min,max"
    volumeMounts:
    - name: config-volume
      mountPath: /jrconfig
    volumes:
    - name: config-volume
      configMap:
        name: jr-configmap
```





