kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/aws/deploy.yaml
kubectl apply -f k8s/postgres/
kubectl apply -f k8s/go-product/