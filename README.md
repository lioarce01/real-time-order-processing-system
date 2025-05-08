# Real-Time Order Processing System

This is a real-time order processing system designed to manage order creation and subsequent processing using modern technologies such as Kafka, PostgreSQL, and microservices. The system receives orders, stores them in the database, publishes events to a messaging system (Kafka), and then Kafka consumers process these events.

## Features

- **Order Creation**: The system allows you to create orders through a REST API.
- **Data Persistence**: Orders are stored in a PostgreSQL database.
- **Asynchronous Messaging**: Kafka is used to publish order creation events and allow other services to process them efficiently.
- **Scalability**: The architecture is designed to be scalable, handling a high volume of orders and events.

## Technologies Used

- **Golang**: Primary language used for backend and service implementation.
- **Kafka**: Messaging system used to publish and consume events.
- **PostgreSQL**: Relational database used to store orders and their details.
- **GORM**: Golang ORM used to interact with the PostgreSQL database.
- **Docker**: Tool for containerizing services and easing deployment.
- **Zap**: Structured and efficient logging library.
- **Sarama**: Kafka client for integration with the messaging system.
- **Clean Architecture**: Code structure based on Clean Architecture principles to maintain a clean, modular, and scalable system.

## Architecture

This project is structured following **Clean Architecture** principles, allowing the separation of different system layers, making the code more maintainable and testable. The architecture is divided into the following layers:

- **Domain**: Contains business entities and domain rules.
- **Application**: Contains the business logic of use cases, such as order creation.
- **Infrastructure**: Provides adapters to interact with external technologies like the database (PostgreSQL) and Kafka.
- **Interface**: Defines HTTP controllers and other mechanisms for interacting with the system.

## Workflow

1. **Order Creation**: Users can create an order via the REST API at the `/orders` endpoint. The request is processed, and the order is stored in the PostgreSQL database.
2. **Event Publishing**: After the order is saved, an event is published to Kafka to notify other services about the order creation.
3. **Event Consumption**: A Kafka consumer listens for the event and starts processing the order, which may include additional actions such as updating inventory or generating invoices.
