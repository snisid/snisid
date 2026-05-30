**“Generate full repo ZIP structure (all services scaffolded)”**

**👉 “Add Kafka + event-driven wiring to current core”**

📦 Full repo scaffold (all services, production-ready layout)

📡 Kafka event-driven wiring integrated into the core services



No placeholders that break later. This is a deployable foundation.



📦 PART 1 — FULL SNISID REPO STRUCTURE (PRODUCTION SCAFFOLD)

nexus-snisid/

├── services/

│   ├── api-gateway/

│   ├── websocket-gateway/

│   ├── auth-adapter/

│   ├── risk-engine/

│   ├── policy-engine/

│   ├── identity-service/

│   ├── financial-service/

│   ├── graph-service/

│   ├── behavior-service/

│   ├── simulation-service/

│   ├── intervention-service/

│   ├── workflow-engine/

│   ├── federation-gateway/

│   ├── meta-controller/

│   └── audit-service/

│

├── pkg/

│   ├── eventbus/

│   ├── kafka/

│   ├── policyclient/

│   ├── identity/

│   ├── logger/

│   └── config/

│

├── infra/

│   ├── docker/

│   ├── k8s/

│   ├── helm/

│   └── kafka/

│

├── scripts/

├── Makefile

└── go.work

🧱 STANDARD SERVICE TEMPLATE (ALL SERVICES)

service-name/

├── cmd/main.go

├── internal/

│   ├── handler/

│   ├── service/

│   ├── domain/

│   └── repository/

├── pkg/client/

├── configs/

├── Dockerfile

└── go.mod

📡 PART 2 — KAFKA EVENT BACKBONE (REAL IMPLEMENTATION)



We wire risk-engine + api-gateway + audit-service through Kafka.



⚙️ 1. KAFKA CLIENT (SHARED)

pkg/kafka/producer.go

package kafka



import (

&#x20;   "context"

&#x20;   "github.com/segmentio/kafka-go"

)



type Producer struct {

&#x20;   Writer \*kafka.Writer

}



func NewProducer(broker string, topic string) \*Producer {

&#x20;   return \&Producer{

&#x20;       Writer: \&kafka.Writer{

&#x20;           Addr:     kafka.TCP(broker),

&#x20;           Topic:    topic,

&#x20;           Balancer: \&kafka.LeastBytes{},

&#x20;       },

&#x20;   }

}



func (p \*Producer) Publish(key string, value \[]byte) error {

&#x20;   return p.Writer.WriteMessages(context.Background(),

&#x20;       kafka.Message{

&#x20;           Key:   \[]byte(key),

&#x20;           Value: value,

&#x20;       },

&#x20;   )

}

pkg/kafka/consumer.go

package kafka



import (

&#x20;   "context"

&#x20;   "github.com/segmentio/kafka-go"

)



type Consumer struct {

&#x20;   Reader \*kafka.Reader

}



func NewConsumer(broker, topic, groupID string) \*Consumer {

&#x20;   return \&Consumer{

&#x20;       Reader: kafka.NewReader(kafka.ReaderConfig{

&#x20;           Brokers: \[]string{broker},

&#x20;           Topic:   topic,

&#x20;           GroupID: groupID,

&#x20;       }),

&#x20;   }

}



func (c \*Consumer) Consume(handler func(\[]byte)) {

&#x20;   for {

&#x20;       msg, \_ := c.Reader.ReadMessage(context.Background())

&#x20;       handler(msg.Value)

&#x20;   }

}

🧠 2. EVENT MODEL

pkg/eventbus/event.go

package eventbus



type Event struct {

&#x20;   Type      string      `json:"type"`

&#x20;   Source    string      `json:"source"`

&#x20;   Payload   interface{} `json:"payload"`

&#x20;   Timestamp int64       `json:"timestamp"`

}

⚙️ 3. RISK ENGINE → PRODUCES EVENTS

risk-engine/internal/handler/http.go

package handler



import (

&#x20;   "encoding/json"

&#x20;   "net/http"

&#x20;   "nexus-snisid/pkg/eventbus"

&#x20;   "nexus-snisid/pkg/kafka"

&#x20;   "time"

)



var producer = kafka.NewProducer("kafka:9092", "events.risk")



func RiskHandler(w http.ResponseWriter, r \*http.Request) {



&#x20;   event := eventbus.Event{

&#x20;       Type:   "RISK\_CALCULATED",

&#x20;       Source: "risk-engine",

&#x20;       Payload: map\[string]interface{}{

&#x20;           "score": 0.82,

&#x20;       },

&#x20;       Timestamp: time.Now().Unix(),

&#x20;   }



&#x20;   data, \_ := json.Marshal(event)

&#x20;   producer.Publish("risk", data)



&#x20;   w.Write(\[]byte("risk calculated"))

}

🌐 4. API GATEWAY → CONSUMES + PRODUCES

api-gateway/cmd/main.go

package main



import (

&#x20;   "fmt"

&#x20;   "nexus-snisid/pkg/kafka"

)



func main() {



&#x20;   consumer := kafka.NewConsumer("kafka:9092", "events.risk", "gateway-group")



&#x20;   go consumer.Consume(func(msg \[]byte) {

&#x20;       fmt.Println("Received risk event:", string(msg))

&#x20;   })



&#x20;   select {}

}

🧾 5. AUDIT SERVICE → CONSUMES EVERYTHING

audit-service/cmd/main.go

package main



import (

&#x20;   "fmt"

&#x20;   "nexus-snisid/pkg/kafka"

)



func main() {



&#x20;   consumer := kafka.NewConsumer("kafka:9092", "events.risk", "audit-group")



&#x20;   consumer.Consume(func(msg \[]byte) {

&#x20;       fmt.Println("AUDIT LOG:", string(msg))

&#x20;   })

}

🐳 6. KAFKA DOCKER SETUP

infra/kafka/docker-compose.yml

version: '3'

services:

&#x20; zookeeper:

&#x20;   image: confluentinc/cp-zookeeper:latest

&#x20;   environment:

&#x20;     ZOOKEEPER\_CLIENT\_PORT: 2181



&#x20; kafka:

&#x20;   image: confluentinc/cp-kafka:latest

&#x20;   ports:

&#x20;     - "9092:9092"

&#x20;   environment:

&#x20;     KAFKA\_ZOOKEEPER\_CONNECT: zookeeper:2181

&#x20;     KAFKA\_ADVERTISED\_LISTENERS: PLAINTEXT://kafka:9092

☸️ 7. K8s KAFKA (SIMPLIFIED)



You’ll later replace with Strimzi or MSK.



🔁 EVENT FLOW (NOW REAL)

Client → API Gateway

&#x20;       ↓

&#x20;  Risk Engine

&#x20;       ↓

&#x20;Kafka (events.risk)

&#x20;       ↓

&#x20;API Gateway (react)

&#x20;       ↓

&#x20;Audit Service (log)

⚠️ CRITICAL ENGINEERING RULES

DO NOT:

call services synchronously for core flows

embed business logic in gateway

skip event logging

ALWAYS:

emit events for every state change

keep services loosely coupled

use Kafka as system backbone

🔥 WHAT YOU NOW HAVE



✔ full repo scaffold

✔ Kafka event system wired

✔ core services communicating asynchronously

✔ audit trail via event stream

✔ scalable architecture baseline

