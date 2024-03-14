
# 生产环境发布
go build
rsync -cavzP ./Dockerfile root@10.10.110.14:/home/twigir/server
rsync -cavzP ./docker-compose.yaml root@10.10.110.14:/home/twigir/server
rsync -cavzP ./server root@10.10.110.14:/home/twigir/server
# rsync -cavzP --delete-after ./ --exclude-from='.rsync-exclude' root@10.10.110.14:/home/twigir/server
ssh root@10.10.110.14 "\
cd /home/twigir/server; \
docker-compose up -d; \
docker images | grep none | awk '{print $3}' | xargs docker rmi \
"
