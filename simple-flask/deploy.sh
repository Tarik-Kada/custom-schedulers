docker build -t tarikkada/scheduler-simple-flask:latest .

docker push tarikkada/scheduler-simple-flask:latest

kubectl apply -f ./deployment.yaml

kubectl apply -f ./config.yaml