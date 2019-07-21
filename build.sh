export TAG="latest-dev"
if [ $1 ] ; then
  TAG=$1
fi

echo "Building: alexellis/derek:$TAG"

docker build -t alexellis/derek:$TAG . -f Dockerfile

