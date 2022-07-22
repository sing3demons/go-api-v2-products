kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.3.0/deploy/static/provider/aws/deploy.yaml
kubectl apply -f k8s/postgres/ -f k8s/go-product/