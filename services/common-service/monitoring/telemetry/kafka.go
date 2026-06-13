package telemetry

// KafkaCarrier implements propagation.TextMapCarrier for Kafka message headers.
// This allows trace context to be propagated through Kafka messages,
// so a consumer can continue the same trace as the producer.
//
// Producer usage (inject trace context into Kafka headers):
//
//	headers := telemetry.KafkaCarrier{}
//	telemetry.GetPropagator().Inject(ctx, headers)
//	msg := &kafka.Message{
//	    Headers: headers.ToKafkaHeaders(),
//	    Value:   payload,
//	}
//	producer.Produce(msg, nil)
//
// Consumer usage (extract trace context from Kafka headers):
//
//	carrier := telemetry.KafkaCarrierFromHeaders(msg.Headers)
//	ctx = telemetry.GetPropagator().Extract(ctx, carrier)
//	// now ctx contains the trace from the producer

type KafkaCarrier map[string]string

// Get returns the value for a given key.
func (c KafkaCarrier) Get(key string) string {
	return c[key]
}

// Set sets the value for a given key.
func (c KafkaCarrier) Set(key string, value string) {
	c[key] = value
}

// Keys returns all keys in the carrier.
func (c KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// KafkaHeader represents a single Kafka message header.
// This is framework-agnostic — works with confluent-kafka-go, sarama, segmentio/kafka-go.
type KafkaHeader struct {
	Key   string
	Value []byte
}

// ToKafkaHeaders converts the carrier to a slice of KafkaHeader.
// Map to your Kafka library's header type accordingly.
//
// confluent-kafka-go:
//
//	for _, h := range carrier.ToKafkaHeaders() {
//	    msg.Headers = append(msg.Headers, kafka.Header{Key: h.Key, Value: h.Value})
//	}
//
// segmentio/kafka-go:
//
//	for _, h := range carrier.ToKafkaHeaders() {
//	    msg.Headers = append(msg.Headers, kafkago.Header{Key: h.Key, Value: h.Value})
//	}
func (c KafkaCarrier) ToKafkaHeaders() []KafkaHeader {
	headers := make([]KafkaHeader, 0, len(c))
	for k, v := range c {
		headers = append(headers, KafkaHeader{Key: k, Value: []byte(v)})
	}
	return headers
}

// KafkaCarrierFromHeaders creates a KafkaCarrier from a slice of KafkaHeader.
// Use this on the consumer side to extract trace context.
func KafkaCarrierFromHeaders(headers []KafkaHeader) KafkaCarrier {
	carrier := KafkaCarrier{}
	for _, h := range headers {
		carrier[h.Key] = string(h.Value)
	}
	return carrier
}
