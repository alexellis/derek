export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

docker build -t alexellis/derek:$eTAG . -f Dockerfile && \
docker create --name derek alexellis/derek:$eTAG