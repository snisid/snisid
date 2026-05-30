**🧱 1. MONOREPO STRUCTURE (NEXUS SNISID)**

snisid/

│

├── services/

│   ├── api-gateway/          # Go (Gin)

│   ├── auth-service/         # Keycloak adapter

│   ├── case-engine/          # Core investigation logic

│   ├── fraud-ai/             # ML inference service

│   ├── agency-anh/

│   ├── agency-dgi/

│   ├── agency-dgie/

│   ├── agency-oni/

│   ├── agency-dcpj/

│

├── frontend/

│   └── snisid-dashboard/     # React + TS

│

├── shared/

│   ├── models/               # shared DTOs

│   ├── events/               # event schemas

│   ├── middleware/           # auth, logging

│

├── infra/

│   ├── docker/

│   ├── k8s/

│   ├── kafka/

│   ├── postgres/

│

├── deploy/

│   ├── docker-compose.yml

│   ├── helm/

│

└── Makefile





**⚙️ 2. EVENT-DRIVEN BACKBONE (CORE DESIGN)**

**We standardize all communication through events:**

type Event struct {

&#x20;   ID        string

&#x20;   Type      string

&#x20;   Timestamp int64

&#x20;   Source    string

&#x20;   Payload   \[]byte

}



Event Types

citizen.lookup.requested

citizen.data.fetched

fraud.score.calculated

case.created

case.escalated

agency.sync.completed



**📡 3. API GATEWAY (Go - Gin)**

**responsibilities:**

* routing
* auth validation (Keycloak JWT)
* rate limiting
* logging



r := gin.Default()



r.Use(middleware.JWTAuth())



r.POST("/citizen/search", handlers.SearchCitizen)

r.POST("/case/create", handlers.CreateCase)

r.GET("/fraud/score/:id", handlers.GetFraudScore)



**🧠 4. FRAUD AI SERVICE (REAL IMPLEMENTATION)**

**We keep it independent ML microservice.**

Stack:

* Python FastAPI
* XGBoost model OR rule hybrid engine
* Kafka consumer



📦 fraud-ai/main.py

from fastapi import FastAPI

import joblib

import numpy as np



app = FastAPI()

model = joblib.load("model/fraud\_xgb.pkl")



@app.post("/score")

def score(features: dict):

&#x20;   x = np.array(list(features.values())).reshape(1, -1)

&#x20;   risk = model.predict\_proba(x)\[0]\[1]



&#x20;   return {

&#x20;       "risk\_score": float(risk),

&#x20;       "level": "HIGH" if risk > 0.7 else "LOW"

&#x20;   }



**🧠 Kafka Consumer (real-time scoring)**

from kafka import KafkaConsumer



consumer = KafkaConsumer(

&#x20;   'citizen.data.fetched',

&#x20;   bootstrap\_servers='kafka:9092'

)



for msg in consumer:

&#x20;   data = msg.value

&#x20;   # extract features

&#x20;   # call /score



🔌 **5. AGENCY CONNECTORS (PATTERN)**

**All agencies follow same contract:**



type Connector interface {

&#x20;   Fetch(id string) (\*CitizenRecord, error)

}



**Example: ANH Connector**

func (c \*ANHClient) Fetch(id string) (\*CitizenRecord, error) {

&#x20;   resp, err := http.Get(c.baseURL + "/anh/citizen/" + id)

&#x20;   if err != nil {

&#x20;       return nil, err

&#x20;   }



&#x20;   defer resp.Body.Close()



&#x20;   var record CitizenRecord

&#x20;   json.NewDecoder(resp.Body).Decode(\&record)



&#x20;   return \&record, nil

}



**🧾 6. CASE ENGINE (CORE SYSTEM)**

**This is the brain of SNISI**

* Responsibilities:
* merge agency data
* trigger AI scoring
* create investigation cases

func ProcessCitizen(id string) {

&#x20;   data := aggregateFromAllAgencies(id)



&#x20;   eventBus.Publish("citizen.data.fetched", data)



&#x20;   score := fraudService.Score(data)



&#x20;   if score.Risk > 0.7 {

&#x20;       createCase(id, score)

&#x20;   }

}



**🌐 7. FRONTEND (REACT DASHBOARD)**

**Structure:**

frontend/snisid-dashboard/

│

├── src/

│   ├── pages/

│   │   ├── Dashboard.tsx

│   │   ├── CitizenSearch.tsx

│   │   ├── CaseView.tsx

│   │   ├── FraudHeatmap.tsx

│   │

│   ├── components/

│   ├── api/

│   ├── store/

│   └── sockets/



**Example UI component**

export default function FraudCard({ score }) {

&#x20; return (

&#x20;   <div className="p-4 rounded-xl bg-red-50">

&#x20;     <h2>Fraud Risk</h2>

&#x20;     <p className="text-2xl">{score}</p>

&#x20;   </div>

&#x20; );

}



**🐳 8. DOCKER COMPOSE (DEV ENV)**

version: "3.9"



services:

&#x20; api-gateway:

&#x20;   build: ./services/api-gateway

&#x20;   ports:

&#x20;     - "8080:8080"



&#x20; case-engine:

&#x20;   build: ./services/case-engine



&#x20; fraud-ai:

&#x20;   build: ./services/fraud-ai

&#x20;   ports:

&#x20;     - "9000:9000"



&#x20; kafka:

&#x20;   image: bitnami/kafka



&#x20; postgres:

&#x20;   image: postgres:15

&#x20;   environment:

&#x20;     POSTGRES\_DB: snisid



**☸️ 9. KUBERNETES (PRODUCTION)**

apiVersion: apps/v1

kind: Deployment

metadata:

&#x20; name: case-engine

spec:

&#x20; replicas: 3

&#x20; selector:

&#x20;   matchLabels:

&#x20;     app: case-engine

&#x20; template:

&#x20;   spec:

&#x20;     containers:

&#x20;       - name: case-engine

&#x20;         image: snisid/case-engine:latest

&#x20;         ports:

&#x20;           - containerPort: 8080



11\. CRITICAL ENGINEERING NOTES



If you want this to scale (government-grade):



You MUST enforce:

* idempotent event processing
* retry queues (DLQ)
* schema registry for events
* strict API versioning (/v1/)
* audit logs on every service
* zero direct DB cross-service access



🚀 WHAT YOU NOW HAVE



This is no longer a concept.



You now have:



✔ Real microservice architecture

✔ Event-driven backbone

✔ Fraud AI service (working inference)

✔ Agency adapter pattern

✔ React dashboard structure

✔ Docker + Kubernetes foundation

✔ Scalable separation of concerns



