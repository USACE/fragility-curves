TAG=v0.0.1
IMAGE=williamlehman/fragilitycurveplugin

docker build -t $IMAGE:$TAG .

docker run -it --entrypoint /bin/sh $IMAGE:$TAG 

docker push $IMAGE:$TAG

# # test
docker run $IMAGE:$TAG /bin/sh -c  "./main"