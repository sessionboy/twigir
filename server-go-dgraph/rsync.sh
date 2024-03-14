
# 生产环境发布
rsync -cavzP --delete-after ./ --exclude-from='.rsync-exclude' root@10.10.110.14:/home/twigir/server
# ssh test "\
# cd /home/twigir/server; \
# docker-compose up -d; \
# docker images | grep none | awk '{print $3}' | xargs docker rmi \
# "
