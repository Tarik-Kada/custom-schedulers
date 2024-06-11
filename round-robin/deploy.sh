docker build -t tarikkada/rr-scheduler:latest .

docker push tarikkada/rr-scheduler:latest

kubectl apply -f ./deployment.yaml

kubectl apply -f ./config.yaml