# Backend + System Design + AI Engineer Roadmap

Merged from the three roadmap.sh PDFs: **Backend**, **System Design**, **AI Engineer**.
Each item: `- [ ]` checkbox, **term**, a one-line plain-English meaning, and a resource.
Change `[ ]` to `[x]` when you finish an item.

Language for all exercises: **Go** (you can do everything below in Go).

**Suggested order:** Part 1 (Backend) -> Part 2 (System Design) -> Part 3 (AI Engineer).
Part 2 and Part 3 can start in parallel once Phase 1.4 (databases) is done.

Main references (keep open the whole time):
- roadmap.sh: https://roadmap.sh/backend | https://roadmap.sh/system-design | https://roadmap.sh/ai-engineer
- Book: *Designing Data-Intensive Applications* (DDIA): https://dataintensive.net
- System Design Primer: https://github.com/donnemartin/system-design-primer
- Go: https://go.dev/tour | https://gobyexample.com | https://go.dev/doc/effective_go

---

# PART 1: BACKEND

## Phase 1.1: Internet and web basics
- [ ] **How the internet works**: packets, IP addresses, routers. https://howdns.works and https://www.cloudflare.com/learning/
- [ ] **HTTP**: the request/response protocol of the web (methods, status codes, headers). https://developer.mozilla.org/en-US/docs/Web/HTTP
- [ ] **Domain name and DNS**: DNS turns `example.com` into an IP address. https://howdns.works
- [ ] **Hosting**: where your server program actually runs (VPS, cloud, shared). https://roadmap.sh/backend
- [ ] **Browsers and how they work**: what a browser does after you type a URL. https://developer.mozilla.org/en-US/docs/Web/Performance/How_browsers_work
- [ ] **TCP/UDP basics**: TCP = reliable ordered delivery, UDP = fast, no guarantees. https://beej.us/guide/bgnet/

## Phase 1.2: Language and tooling
- [ ] **Pick a language**: Go (already your choice). https://go.dev/tour
- [ ] **Go standard library `net/http`**: build a server without frameworks. https://pkg.go.dev/net/http
- [ ] **Git**: version control, track changes in code. https://git-scm.com/book
- [ ] **Repo hosting (GitHub / GitLab / Bitbucket)**: remote home for your Git repos, plus PRs and issues. https://docs.github.com/en/get-started
- [ ] **Go modules and project layout**: how Go manages dependencies. https://go.dev/doc/modules/managing-dependencies

## Phase 1.3: APIs
- [ ] **REST**: design APIs around resources (`GET /users/1`). https://restfulapi.net
- [ ] **JSON APIs**: send/receive data as JSON. https://pkg.go.dev/encoding/json
- [ ] **HATEOAS**: responses include links to what you can do next. https://restfulapi.net/hateoas/
- [ ] **OpenAPI specs**: a machine-readable description of your API. https://swagger.io/specification/
- [ ] **gRPC**: fast binary RPC (call a remote function) using Protocol Buffers. https://grpc.io/docs/languages/go/quickstart/
- [ ] **GraphQL**: client asks for exactly the fields it wants. https://graphql.org/learn/
- [ ] **SOAP**: older XML-based API style, just know it exists. https://roadmap.sh/backend
- [ ] **Real-time: Long polling, short polling, WebSockets, SSE**: ways to push updates to clients (SSE = server-sent events, one-way server to browser). https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API and https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events

## Phase 1.4: Databases
- [ ] **Relational DB (PostgreSQL first)**: tables, rows, SQL. https://www.postgresql.org/docs/current/ and https://www.pgexercises.com
- [ ] **MySQL / MariaDB / MS SQL / Oracle**: other relational DBs, skim the differences. https://roadmap.sh/backend
- [ ] **Normalization**: organize tables so data is not duplicated. https://www.guru99.com/database-normalization.html
- [ ] **ACID**: Atomic, Consistent, Isolated, Durable, the safety guarantees of transactions. https://www.postgresql.org/docs/current/tutorial-transactions.html
- [ ] **Transactions**: group steps so they all succeed or all fail. https://www.postgresql.org/docs/current/tutorial-transactions.html
- [ ] **Failure modes**: what happens to data on crash / network failure. DDIA chapters 7-8.
- [ ] **ORMs and query builders in Go**: map DB rows to code. https://sqlc.dev and https://gorm.io
- [ ] **Migrations**: versioned scripts that change the DB schema. https://github.com/golang-migrate/migrate
- [ ] **Database indexes**: lookup structure that makes queries fast. https://use-the-index-luke.com
- [ ] **N+1 problem**: 1 query for a list + N queries for each item, very slow. https://use-the-index-luke.com
- [ ] **Profiling performance (EXPLAIN)**: see how the DB runs a query. https://www.postgresql.org/docs/current/using-explain.html
- [ ] **Data replication**: copies of the DB on other machines. DDIA chapter 5.
- [ ] **Sharding strategies**: split data across machines. DDIA chapter 6.
- [ ] **CAP theorem**: during a network split you pick consistency or availability. https://jepsen.io/consistency
- [ ] **Scaling databases**: read replicas, partitioning, pooling. https://www.postgresql.org/docs/current/high-availability.html
- [ ] **NoSQL overview**: non-table databases, each fits a use case. https://roadmap.sh/backend
  - [ ] Document DB: MongoDB / CouchDB. https://www.mongodb.com/docs/manual/
  - [ ] Key-Value: Redis / DynamoDB. https://redis.io/docs/
  - [ ] Column DB: Cassandra. https://cassandra.apache.org/doc/latest/
  - [ ] Graph DB: Neo4j. https://neo4j.com/docs/getting-started/
  - [ ] Time-series DB: InfluxDB / TimescaleDB. https://docs.timescale.com
  - [ ] Realtime DB: Firebase / RethinkDB. https://firebase.google.com/docs/database
  - [ ] SQLite: tiny file-based SQL DB. https://www.sqlite.org/docs.html
- [ ] **Search engines: Elasticsearch / Solr**: full-text search. https://www.elastic.co/guide/index.html

## Phase 1.5: Caching
- [ ] **Client-side caching**: browser keeps copies. https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
- [ ] **CDN**: servers around the world that serve static files closer to users. https://www.cloudflare.com/learning/cdn/what-is-a-cdn/
- [ ] **Server-side caching (Redis, Memcached)**: keep hot data in memory. https://redis.io/docs/latest/develop/ and https://memcached.org

## Phase 1.6: Authentication and security
- [ ] **Authentication types**: Basic, Token, Cookie/session, JWT (a signed token carrying user info), OAuth (log in with another service), OpenID, SAML. https://jwt.io/introduction and https://www.oauth.com
- [ ] **Hashing algorithms**: MD5/SHA (fast, not for passwords); bcrypt/scrypt (slow on purpose, for passwords). https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
- [ ] **HTTPS and SSL/TLS**: encrypted HTTP. https://developer.mozilla.org/en-US/docs/Web/Security/Transport_Layer_Security
- [ ] **CORS**: browser rule on which sites may call your API. https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
- [ ] **CSP (Content Security Policy)**: header limiting what scripts a page can run. https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
- [ ] **OWASP Top 10**: most common web vulnerabilities. https://owasp.org/www-project-top-ten/
- [ ] **API security best practices**. https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html
- [ ] **Server security**. https://roadmap.sh/backend

## Phase 1.7: Testing
- [ ] **Unit testing**: test one function. https://go.dev/doc/tutorial/add-a-test
- [ ] **Integration testing**: test parts together (e.g. API + real DB). https://testcontainers.com/guides/getting-started-with-testcontainers-for-go/
- [ ] **Functional testing**: test the feature from the user's view. https://roadmap.sh/backend
- [ ] **Table-driven tests and benchmarks in Go**. https://go.dev/wiki/TableDrivenTests

## Phase 1.8: CI/CD and DevOps basics
- [ ] **CI/CD**: automatically test (CI) and deploy (CD) on every push. https://docs.github.com/en/actions
- [ ] **Containerization vs virtualization**: containers share the OS kernel; VMs run a full OS. https://docs.docker.com/get-started/docker-overview/
- [ ] **Docker**: package an app + dependencies into a container. https://docs.docker.com/get-started/
- [ ] **LXC**: lower-level Linux containers. https://linuxcontainers.org/lxc/introduction/
- [ ] **Kubernetes**: runs and manages many containers. https://kubernetes.io/docs/tutorials/kubernetes-basics/
- [ ] **Web servers: Nginx / Caddy / Apache / IIS**: front door that serves files and forwards requests. https://nginx.org/en/docs/beginners_guide.html and https://caddyserver.com/docs/
- [ ] **Basic Linux/infrastructure knowledge**. https://roadmap.sh/devops

## Phase 1.9: Architecture and design
- [ ] **GoF design patterns**: 23 classic reusable object-oriented solutions. https://refactoring.guru/design-patterns
- [ ] **Twelve-Factor App**: rules for cloud-friendly apps. https://12factor.net
- [ ] **Monolith vs microservices vs SOA vs serverless**: one big app vs many small services vs functions run on demand. https://martinfowler.com/articles/microservices.html
- [ ] **Service mesh**: a layer that manages service-to-service traffic. https://istio.io/latest/about/service-mesh/
- [ ] **Domain-Driven Design (DDD)**: structure code around the business domain. Book: *Domain-Driven Design Distilled* (Vaughn Vernon).
- [ ] **CQRS**: separate models for reading and writing. https://martinfowler.com/bliki/CQRS.html
- [ ] **Event sourcing**: store every change as an event, not just the final state. https://martinfowler.com/eaaDev/EventSourcing.html
- [ ] **Test-Driven Development (TDD)**: write the test first. https://martinfowler.com/bliki/TestDrivenDevelopment.html

## Phase 1.10: Message brokers
- [ ] **Message broker basics**: middleman that holds messages between services. https://www.rabbitmq.com/tutorials
- [ ] **RabbitMQ**: queue-based broker. https://www.rabbitmq.com/tutorials/tutorial-one-go
- [ ] **Kafka**: distributed log for high-volume streams. https://developer.confluent.io/courses/apache-kafka/events/

## Phase 1.11: Scale, resilience, observability
- [ ] **Types of scaling**: vertical (bigger machine) vs horizontal (more machines). https://github.com/donnemartin/system-design-primer
- [ ] **Mitigation strategies**: graceful degradation (partly work instead of fail), throttling (limit rate), backpressure (tell senders to slow down), load shifting, circuit breaker (stop calling a failing service). https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- [ ] **Migration strategies**: moving systems/data without downtime. https://martinfowler.com/bliki/StranglerFigApplication.html
- [ ] **Instrumentation / Monitoring / Telemetry**: add code to emit data / watch it / the data itself. https://opentelemetry.io/docs/concepts/
- [ ] **Metrics and logs**: Prometheus + Grafana. https://prometheus.io/docs/introduction/overview/
- [ ] **Go concurrency**: goroutines, channels, context, sync. https://go.dev/blog/pipelines and https://go.dev/blog/context

**Backend capstone (build it yourself):** a Go API with auth, Postgres, Redis cache, a RabbitMQ worker, Docker, CI, tests, and basic metrics.
- [ ] Capstone done

---

# PART 2: SYSTEM DESIGN

## Phase 2.1: Fundamentals
- [ ] **What is system design / how to approach it**: clarify requirements, estimate, draw components, discuss trade-offs. https://github.com/donnemartin/system-design-primer
- [ ] **Performance vs scalability**: slow for one user vs slow under many users. System Design Primer.
- [ ] **Latency vs throughput**: time per request vs requests per second. System Design Primer.
- [ ] **Availability vs consistency**: always answer vs always give the latest data. DDIA ch. 9.
- [ ] **CAP theorem (AP vs CP)**: AP = stay available, CP = stay consistent during a partition. https://jepsen.io/consistency

## Phase 2.2: Consistency and availability
- [ ] **Consistency patterns**: weak, eventual (all copies agree after a while), strong (always latest). https://jepsen.io/consistency
- [ ] **Availability patterns**: fail-over, active-active, active-passive. System Design Primer.
- [ ] **Replication**: master-slave, master-master. DDIA ch. 5.
- [ ] **Availability in numbers**: 99.9% (three 9s), 99.99% (four 9s), in parallel vs in sequence. https://sre.google/sre-book/table-of-contents/

## Phase 2.3: Network and traffic
- [ ] **DNS**. https://howdns.works
- [ ] **CDN: push vs pull**: you upload content vs CDN fetches on first request. https://www.cloudflare.com/learning/cdn/what-is-a-cdn/
- [ ] **Load balancers**: spread requests across servers. https://www.nginx.com/resources/glossary/load-balancing/
- [ ] **LB vs reverse proxy**: a reverse proxy sits in front of servers; LB is a reverse proxy focused on distribution. System Design Primer.
- [ ] **Load balancing algorithms**: round robin, least connections, hashing. System Design Primer.
- [ ] **Layer 4 vs Layer 7 LB**: L4 routes by IP/port, L7 routes by HTTP content. System Design Primer.
- [ ] **Horizontal scaling**. System Design Primer.
- [ ] **HTTP, TCP, UDP**. https://beej.us/guide/bgnet/
- [ ] **Communication styles: RPC, REST, gRPC, GraphQL**. https://grpc.io/docs/what-is-grpc/introduction/
- [ ] **Idempotent operations**: doing it twice has the same effect as once. https://developer.mozilla.org/en-US/docs/Glossary/Idempotent

## Phase 2.4: Application layer
- [ ] **Microservices**. https://microservices.io/patterns/
- [ ] **Service discovery**: how services find each other. https://microservices.io/patterns/server-side-discovery.html
- [ ] **Background jobs**: event-driven or schedule-driven, and how results are returned. https://roadmap.sh/system-design
- [ ] **Asynchronism**: don't make users wait. Task queues, message queues, back pressure. System Design Primer.

## Phase 2.5: Data layer
- [ ] **SQL vs NoSQL**. https://roadmap.sh/system-design
- [ ] **Replication, sharding, federation (split DB by function), denormalization (duplicate data for speed), SQL tuning**. System Design Primer + DDIA ch. 3, 5, 6.
- [ ] **NoSQL types**: key-value, document, wide-column, graph. DDIA ch. 2.
- [ ] **Caching strategies**: cache-aside, write-through, write-behind, refresh-ahead. System Design Primer.
- [ ] **Types of caching**: client, CDN, web server, database, application. System Design Primer.

## Phase 2.6: Performance antipatterns
Know each as "what mistake is this?":
- [ ] Busy Database / Busy Frontend / Chatty I/O / Extraneous Fetching
- [ ] Improper Instantiation / Monolithic Persistence / No Caching
- [ ] Noisy Neighbor / Synchronous I/O / Retry Storm
- Resource: https://learn.microsoft.com/en-us/azure/architecture/antipatterns/

## Phase 2.7: Monitoring
- [ ] **Health, availability, performance, security, usage monitoring**. https://sre.google/sre-book/monitoring-distributed-systems/
- [ ] **Instrumentation, visualization and alerts**. https://prometheus.io/docs/introduction/overview/

## Phase 2.8: Cloud design patterns (overview only, per roadmap)
Resource for all: https://learn.microsoft.com/en-us/azure/architecture/patterns/
- [ ] **Messaging**: Sequential Convoy, Scheduling Agent Supervisor, Queue-based Load Leveling, Publisher/Subscriber, Priority Queue, Pipes and Filters, Competing Consumers, Choreography, Claim Check, Async Request-Reply
- [ ] **Data management**: Valet Key, Static Content Hosting, Sharding, Materialized View, Index Table, Event Sourcing, CQRS, Cache-Aside
- [ ] **Design and implementation**: Strangler Fig, Sidecar, Leader Election, Gateway Routing / Offloading / Aggregation, External Config Store, Compute Resource Consolidation, Backends for Frontend, Anti-Corruption Layer, Ambassador
- [ ] **Reliability**: Deployment Stamps, Geodes, Health Endpoint Monitoring, Throttling, Bulkhead, Circuit Breaker, Compensating Transaction, Retry
- [ ] **Security**: Federated Identity, Gatekeeper, Valet Key

## Phase 2.9: Practice
- [ ] Read 5 case studies: https://aws.amazon.com/builders-library/
- [ ] Design on paper: URL shortener, rate limiter, chat app, news feed, file storage, video pipeline
- [ ] Do 10 mock designs, say trade-offs out loud
- [ ] Read DDIA fully (can span all phases)

**System design capstone:** pick one of your own projects and write a 2-page design doc (requirements, diagram, bottlenecks, failure cases).
- [ ] Capstone done

---

# PART 3: AI ENGINEER

## Phase 3.1: Introduction
- [ ] **What is an AI Engineer / AI vs ML Engineer**: AI engineers build products on top of existing models; ML engineers train models. https://roadmap.sh/ai-engineer
- [ ] **Core vocabulary**: LLM (large language model), inference (running a model), training, embeddings (text as numbers), vector DB, RAG (give the model documents to look at), prompt engineering, AI agent, AGI. https://huggingface.co/learn/llm-course
- [ ] **Intro video**: Karpathy, intro to LLMs. https://karpathy.ai/zero-to-hero.html
- [ ] **Using pre-trained models**: benefits and limits. https://huggingface.co/learn/llm-course

## Phase 3.2: Models and APIs
- [ ] **Popular models**: OpenAI, Anthropic Claude, Google Gemini, Mistral, Cohere. https://docs.anthropic.com | https://platform.openai.com/docs | https://ai.google.dev | https://docs.mistral.ai
- [ ] **Cloud platforms**: Azure AI, AWS SageMaker/Bedrock. https://aws.amazon.com/bedrock/
- [ ] **Capabilities and context length** (how much text it can read at once) and **knowledge cut-off dates**. Model docs above.
- [ ] **Chat completions API**: send messages, get a reply. https://platform.openai.com/docs/guides/text
- [ ] **Anthropic Messages API (Go SDK)**. https://docs.anthropic.com/en/api/overview
- [ ] **Tokens**: pieces of text the model reads; maximum tokens, token counting, pricing. https://platform.openai.com/tokenizer
- [ ] **Fine-tuning**: further training a model on your data. https://platform.openai.com/docs/guides/fine-tuning
- [ ] **Playground**: test prompts in the browser. https://platform.openai.com/playground

## Phase 3.3: Prompt engineering
- [ ] **Writing good prompts**: clear task, examples, output format. https://www.promptingguide.ai
- [ ] **Anthropic prompt guide**. https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/overview
- [ ] **Structured outputs / constraining outputs**. https://platform.openai.com/docs/guides/structured-outputs

## Phase 3.4: AI safety and ethics
- [ ] **Prompt injection**: hidden instructions in input that hijack the model. https://owasp.org/www-project-top-10-for-large-language-model-applications/
- [ ] **Bias and fairness; security and privacy concerns**. OWASP LLM Top 10 above.
- [ ] **Adversarial testing**: try to break your own app. OWASP LLM Top 10.
- [ ] **Moderation API; end-user IDs in prompts**. https://platform.openai.com/docs/guides/moderation
- [ ] **Robust prompts; know your use case; constrain inputs and outputs**. https://docs.anthropic.com/en/docs/test-and-evaluate/strengthen-guardrails/reduce-prompt-leak

## Phase 3.5: Open-source AI
- [ ] **Open vs closed models**. https://huggingface.co/learn/llm-course
- [ ] **Hugging Face Hub and tasks**: find and use models. https://huggingface.co/docs/hub/index
- [ ] **Inference SDK**: https://huggingface.co/docs/huggingface_hub/guides/inference
- [ ] **Transformers.js**: run models in the browser. https://huggingface.co/docs/transformers.js
- [ ] **Ollama**: run models on your own machine; models + SDK. https://github.com/ollama/ollama
- [ ] **Call Ollama from Go**. https://pkg.go.dev/github.com/ollama/ollama/api

## Phase 3.6: Embeddings and vector databases
- [ ] **What are embeddings**: a list of numbers that captures meaning, so similar texts are close. https://platform.openai.com/docs/guides/embeddings
- [ ] **Use cases**: semantic search, recommendations, anomaly detection, classification. Same link.
- [ ] **OpenAI embeddings API / models / pricing**. Same link.
- [ ] **Open-source embeddings**: Sentence Transformers, Hugging Face models. https://www.sbert.net
- [ ] **Vector databases: purpose**: store embeddings and find nearest ones fast. https://www.pinecone.io/learn/vector-database/
- [ ] **Pick one**: pgvector (fits your Postgres skills), Chroma, Qdrant, Weaviate, FAISS, LanceDB, Pinecone, Supabase, MongoDB Atlas. https://github.com/pgvector/pgvector | https://qdrant.tech/documentation/
- [ ] **Indexing embeddings and similarity search**. https://www.pinecone.io/learn/series/faiss/

## Phase 3.7: RAG
- [ ] **RAG use cases; RAG vs fine-tuning**. https://www.pinecone.io/learn/retrieval-augmented-generation/
- [ ] **Chunking**: split documents into pieces. https://www.pinecone.io/learn/chunking-strategies/
- [ ] **Embed -> store -> retrieve -> generate pipeline**. Same resources.
- [ ] **Implement with SDKs directly first, frameworks later** (LangChain, LlamaIndex, OpenAI Assistants, Replicate). https://python.langchain.com/docs/tutorials/rag/ | https://docs.llamaindex.ai
- [ ] **Build RAG in Go** with pgvector + an LLM API.

## Phase 3.8: AI agents
- [ ] **What agents are**: an LLM that decides which tools to call in a loop. https://www.anthropic.com/engineering/building-effective-agents
- [ ] **ReAct prompting**: model alternates reasoning and tool use. https://www.promptingguide.ai/techniques/react
- [ ] **Function / tool calling**. https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
- [ ] **Manual agent implementation**: write the loop yourself. Same Anthropic article.
- [ ] **OpenAI Assistants API**. https://platform.openai.com/docs/assistants/overview
- [ ] **MCP (Model Context Protocol)**: standard way to give models tools. https://modelcontextprotocol.io

## Phase 3.9: Multimodal AI
- [ ] **Use cases**: image/video understanding, image generation, audio, text-to-speech, speech-to-text. https://roadmap.sh/ai-engineer
- [ ] **Vision API**. https://platform.openai.com/docs/guides/images-vision
- [ ] **DALL-E / image generation API**. https://platform.openai.com/docs/guides/image-generation
- [ ] **Whisper (speech-to-text)**. https://platform.openai.com/docs/guides/speech-to-text
- [ ] **Hugging Face multimodal models**. https://huggingface.co/models
- [ ] **LangChain / LlamaIndex for multimodal apps**. Docs above.

## Phase 3.10: Tools
- [ ] **AI code editors and completion tools**. https://roadmap.sh/ai-engineer
- [ ] **Evaluation**: measure whether your AI feature works (test sets, scoring). https://docs.anthropic.com/en/docs/test-and-evaluate/develop-tests

**AI capstone:** a RAG app over your own notes with an agent that can call 2-3 tools, plus a small test set to measure answer quality.
- [ ] Capstone done

---

# Progress tracker
- [ ] Part 1 complete
- [ ] Part 2 complete
- [ ] Part 3 complete