# start minikube
minikube start --vm-driver=virtualbox

# get information about the cluster
kubectl cluster-info
kubectl get cs # component-info

kubectl create -f file.yml

kubectl delete pod name
kubectl delete svc|service name

# list pods
kubectl get pods
# list services
kubectl get svc|services

kube describe svc name

# scaling
kubectl scale rc web --replicas=10