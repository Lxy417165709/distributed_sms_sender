
docker rm zookeeper -f
docker rm kafka -f

docker pull zookeeper
docker pull wurstmeister/kafka

docker run --name zookeeper -p 2181:2181 -d zookeeper
docker run --name kafka -p 9092:9092  -e KAFKA_ZOOKEEPER_CONNECT=120.26.162.39:2181 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://120.26.162.39:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -d wurstmeister/kafka
