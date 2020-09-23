if [ -d deploy ]; then
	rm -rf deploy
fi
mkdir deploy

cp -r redis deploy/redis
cp -r nginx deploy/nginx
cp -r grafana deploy/grafana
cp -r go deploy/go
cp -r go deploy/mysql

