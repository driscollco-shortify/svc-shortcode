repoName = "services"
serviceName = "shortcode"
project = "shortify-pro"
region = "europe-west1"

.PHONY: build
build:
	@env GOOS=linux CGO_ENABLED=0 go build -o service ../main.go

.PHONY: docker
docker:
	@docker build -t $(serviceName):latest .
	@docker tag $(serviceName):latest $(region)-docker.pkg.dev/$(project)/$(repoName)/$(serviceName):latest
	@docker push $(region)-docker.pkg.dev/$(project)/$(repoName)/$(serviceName):latest
	@rm -f service

.PHONY: deploy
deploy:
	@gcloud run deploy $(serviceName) --project $(project) --image $(region)-docker.pkg.dev/$(project)/$(repoName)/$(serviceName) \
	--platform managed --region $(region) --port 80 --cpu 2 --memory 2Gi --allow-unauthenticated --max-instances 30 --concurrency 50 \
	--set-env-vars "Cloudflare_Credentials_ApiKey=$(CLOUDFLARE_API_KEY),Grafana_Logs_UserId=$(GRAFANA_LOGS_USERID),\
	Grafana_Logs_AccessToken=$(GRAFANA_LOGS_ACCESS_TOKEN),Grafana_Metrics_UserId=$(GRAFANA_METRICS_USERID),\
	Grafana_Metrics_Key=$(GRAFANA_METRICS_KEY),Grafana_Traces_AuthKey=$(GRAFANA_TRACES_AUTHKEY),\
	Profiler_Host=$(GRAFANA_PROFILER_HOST),Profiler_Username=$(GRAFANA_PROFILER_USERNAME),Profiler_Password=$(GRAFANA_PROFILER_PASSWORD)\
	SafeSite_ApiKey=$(SAFESITE_APIKEY),Service_Handlers_Create_Jwt_PublicKey=$(JWT_KEY_PUBLIC),\
	Service_Handlers_Create_Jwt_PrivateKey=$(JWT_KEY_PRIVATE)"