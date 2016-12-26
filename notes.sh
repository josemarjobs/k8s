# gox \
#   -os="linux" \
#   -arch="amd64" \
#   -output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" app

# start minikube
minikube start --vm-driver=virtualbox

# get information about the cluster
kubectl cluster-info
kubectl get cs # component-info

# create components
kubectl create -f file.yml
# delete them
kubectl delete pod name
kubectl delete svc|service name
# list pods
kubectl get pods
# list services
kubectl get svc|services
# describe
kubectl describe svc name
kubectl describe pod name
# scaling
kubectl scale rc web --replicas=10

# expose
kubectl expose pod mypod --type=NodePort # to the outside world
kubectl expose pod mypod --type=ClusterIP # to the cluster