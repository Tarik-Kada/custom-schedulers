docker build -t tarikkada/static-scheduler:latest .

docker push tarikkada/static-scheduler:latest

kubectl apply -f ./deployment.yaml

kubectl apply -f ./config.yaml