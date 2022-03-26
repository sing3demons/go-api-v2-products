ing:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/aws/deploy.yaml

build:
	docker build -t sing3demons/go-api-v2-product:0.0.4 .
push:
	docker push sing3demons/go-api-v2-product:0.0.4

go run :
	kubectl apply -f k8s/go-product

GET:
	curl http://kubernetes.docker.internal/api/v1/products