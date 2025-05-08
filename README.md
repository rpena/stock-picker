# stock-picker
web service that looks up a fixed number of closing prices of a specific stock


## Steps
These steps are specific for minikube, but should be easy to setup in other kubernetes environments

- Open a terminal and if needed, install minikube for your specific environment
- Run `minikube start` - it'll take a little bit of time but should complete in a few min.
- Run `minikube addons enable ingress` - this will be needed to access the server from your local environment to your service running inside minikube
- Run `eval $(minikube docker-env)` - this will setup docker to point to your minikube environment so the docker image will be sent there
- Clone the github repo locally and in a terminal navigate to the root directory - .../stock-picker
- Log into docker and make sure you can run docker images to see the minikube docker images before continuing
- Run `make docker-build` - to create the docker image and provide it to minikube
- cd to deploy dir, and deploy the service
    ```
    cd deploy
    kubectl apply -f configmap.yaml
    kubectl apply -f secret.yaml
    kubectl apply -f service.yaml
    kubectl apply -f deployment.yaml
    kubectl apply -f ingress.yaml
    ```
- Check the pod is running - `kubectl get pods`
- Edit your /etc/hosts file to include the following
    127.0.0.1    stockpickerping.com
- In another terminal run `minikube tunnel` to create a tunnel from your localhost to ingress
- In a terminal run `curl http://stockpickerping.com/stock/` to get the output requested
- You can also check the pod logs to see additional logging - `kubectl logs stock-picker-deployment-XXX`
