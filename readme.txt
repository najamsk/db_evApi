1. build binay for linux from git bash 
	GOOS=linux GOARCH=amd64 go build -o main .

2. run ./deploy_production.bat or ./deploy_staging.bat from command line

3. upload eventapi_production_dockerimg.zip or eventapi_staging_dockerimg.zip on server where you want to deploy api

4. SSH server where image zip file uploaded

5. check api containers by listing using "docker ps" command

6. stop api containere usign "docker stop eventapi_production" 

7. Remove all unused containers, networks, images usinf "docker system prune -a" command . 

8. Load docker container from image using "docker load -i eventapi_production_dockerimg.zip"

9. Set volumes for documents
	docker run --rm -it --name eventapi_production -v /home/aghakhan/dockervolumes/eventvisor/api_live/uploads:/opt/uploads -d -p 4001:4000 eventapi_production