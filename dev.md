docker --debug buildx build -t sucicada/su-home-okaeri:latest .
DOCKER_HOST=tcp://mint.sucicada.me docker images | grep su-home-okaeri


GOTRACEBACK=all DOCKER_HOST=mint.sucicada.me docker --debug  buildx build -t sucicada/su-home-okaeri:latest  .


GOTRACEBACK=all DOCKER_HOST=mint.sucicada.me docker  --debug build -t sucicada/su-home-okaeri:latest --progress=plain .