docker build -t tarikkada/ext-kube-sched:latest .
docker push tarikkada/ext-kube-sched:latest

kubectl apply -f deployment.yaml
kubectl apply -f config.yaml
