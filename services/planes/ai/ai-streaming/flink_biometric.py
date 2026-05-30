from pyflink.common import WatermarkStrategy, Encoder
from pyflink.common.serialization import JsonRowDeserializationSchema, JsonRowSerializationSchema
from pyflink.common.typeinfo import Types
from pyflink.datastream import StreamExecutionEnvironment
from pyflink.datastream.connectors.kafka import FlinkKafkaConsumer, FlinkKafkaProducer

def biometric_pipeline():
    env = StreamExecutionEnvironment.get_execution_environment()
    env.set_parallelism(1)

    # Kafka Configuration
    kafka_props = {'bootstrap.servers': 'kafka:9092', 'group.id': 'flink-biometric-group'}
    
    # Deserialization Schemas
    photo_schema = Types.ROW_NAMED(['identityId', 'timestamp'], [Types.STRING(), Types.STRING()])
    deepfake_schema = Types.ROW_NAMED(['identityId', 'is_deepfake', 'confidence'], [Types.STRING(), Types.BOOLEAN(), Types.FLOAT()])

    # Consumers
    photo_consumer = FlinkKafkaConsumer("photo.ingested", JsonRowDeserializationSchema.builder().type_info(photo_schema).build(), kafka_props)
    deepfake_consumer = FlinkKafkaConsumer("deepfake.detected", JsonRowDeserializationSchema.builder().type_info(deepfake_schema).build(), kafka_props)

    photo_stream = env.add_source(photo_consumer)
    deepfake_stream = env.add_source(deepfake_consumer)

    # Simple Join logic (Internal Flink state handles the match)
    verified_stream = photo_stream.join(deepfake_stream) \
        .where(lambda x: x[0]) \
        .equal_to(lambda y: y[0]) \
        .window(TumblingEventTimeWindows.of(Time.minutes(1))) \
        .apply(lambda x, y: (x[0], not y[1], y[2])) # identityId, is_verified, confidence

    # Producer for verified events
    serialization_schema = JsonRowSerializationSchema.builder().with_type_info(
        Types.ROW_NAMED(['identityId', 'verified', 'confidence'], [Types.STRING(), Types.BOOLEAN(), Types.FLOAT()])).build()
    
    kafka_producer = FlinkKafkaProducer("face.verified", serialization_schema, kafka_props)
    verified_stream.add_sink(kafka_producer)

    env.execute("SNISID Biometric Verification Pipeline")

if __name__ == "__main__":
    biometric_pipeline()
