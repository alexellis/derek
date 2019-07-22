export TAG="latest-dev"
if [ $1 ] ; then
  TAG=$1
fi

echo "Building: $DOCKER_NS/derek:$TAG"

docker build -t $DOCKER_NS/derek:$TAG . -f Dockerfile
