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
      duration: '1m',
    };
    export default ()=> {
       console.log(exec.command("/usr/bin/jr", ["run",
                                                "net_device"
                                                ]))

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





