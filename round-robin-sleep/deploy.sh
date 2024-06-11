docker build -t tarikkada/rr-sleep-scheduler:latest .

docker push tarikkada/rr-sleep-scheduler:latest

kubectl apply -f ./deployment.yaml

kubectl apply -f ./config.yaml