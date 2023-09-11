cd ../../../
docker login -u admin -p Harbor12345 harbor.yuanshi01.com:30687
docker build --build-arg TARGET_PATH=./app/quote-svc -t harbor.yuanshi01.com:30687/trade/quote-svc:latest -f ./app/quote-svc/deploy/Dockerfile .
docker push harbor.yuanshi01.com:30687/trade/quote-svc:latest
kubectl apply -f ./app/quote-svc/deploy/k8s-deployment-local.yml --namespace=develop
docker rmi harbor.yuanshi01.com:30687/trade/quote-svc:latest -f
