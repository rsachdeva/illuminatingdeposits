# Illuminating Deposits Public Facing REST API

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# Docker Compose Deployment
 
### To start all services:
#### docker-compose -f ./deploy/compose/docker-compose.api.yml up --build

The --build option is there for any code changes.

### Then Migrate and set up seed data:
#### export COMPOSE_IGNORE_ORPHANS=True
#### docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

##### To view logs of running services in a separate terminal:
###### docker-compose -f ./deploy/compose/docker-compose.api.yml logs -f --tail 1  

### Shutdown 

#### docker-compose -f ./deploy/compose/docker-compose.api.yml down
#### docker-compose -f ./deploy/compose/docker-compose.seed.yml down

#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

# Push Images to Docker Hub

docker build -t rsachdeva/illuminatingdeposits.api:v0.1 -f ./build/Dockerfile.api .  
docker push rsachdeva/illuminatingdeposits.api:v0.1 (as an example)  
docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed .  
docker push rsachdeva/illuminatingdeposits.seed:v0.1 (as an example)  

# Kubernetes Deployment - WIP

kubectl apply -f deploy/kubernetes/zipkin-deployment.yaml   
kubectl apply -f deploy/kubernetes/zipkin-service.yaml  

kubectl apply -f deploy/kubernetes/ic-traefik-lb.yaml  
kubectl apply -f deploy/kubernetes/ingress.yaml  

Access Traefik Dashboard at [http://localhost:3000/dashboard/#/](http://localhost:3000/dashboard/#/)   

Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io/zipkin/](http://zipkin.127.0.0.1.nip.io/zipkin)  

### Shutdown

kubectl delete -f deploy/kubernetes/.

# HTTP Client Requests:
See cmd/httpclient/editorsupport/HealthCRUD.http for examples.
Use dev env for localhost or change for prod if running web service at different IP address

(Development is WIP)