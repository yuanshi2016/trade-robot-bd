cd ../../../
docker login -u admin -p Harbor12345 harbor.yuanshi01.com:30687
docker build --build-arg TARGET_PATH=./app/usercenter-svc -t harbor.yuanshi01.com:30687/trade/usercenter-svc:latest -f ./app/usercenter-svc/deploy/Dockerfile .
docker push harbor.yuanshi01.com:30687/trade/usercenter-svc:latest
kubectl apply -f ./app/usercenter-svc/deploy/k8s-deployment-local.yml --namespace=develop
docker rmi harbor.yuanshi01.com:30687/trade/usercenter-svc:latest -f
