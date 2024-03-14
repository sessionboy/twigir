
# 项目部署
rsync -cavzP ./docker-compose.yml root@47.99.243.195:/home/dgraph/
# ssh root@47.99.243.195 "\
# cd /home/projects/tantu/web; \
# yarn install ; \
# yarn build ; \
# sh deploy.sh; \
# "
