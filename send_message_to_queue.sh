echo "Sending Model2.json to Kafka"
cat model2.json | docker exec -i l0_kafka \
                          kafka-console-producer \
                          --topic orders \
                          --bootstrap-server localhost:9092
echo "Sent succesfully"
