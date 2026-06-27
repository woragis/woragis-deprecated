# Queue Declarations: RabbitMQ Patterns

## Overview
How we declare queues in RabbitMQ with proper configuration (durability, dead letter queues, etc.).

## Key Points

### Queue Configuration
- Durable queues (survive broker restart)
- Dead letter exchange
- Message TTL (if needed)
- Priority (if needed)

### Exchange Configuration
- Topic exchange for routing
- Dead letter exchange for failures
- Durable exchanges

## Implementation Details

### Exchange Declaration
```go
err = ch.ExchangeDeclare(
    "woragis.tasks",
    "topic",
    true,  // durable
    false, // auto-deleted
    false, // internal
    false, // no-wait
    nil,
)
```

### Queue Declaration
```go
args := amqp.Table{
    "x-dead-letter-exchange":    "woragis.dlx",
    "x-dead-letter-routing-key": queueName + ".failed",
}
queue, err := ch.QueueDeclare(
    queueName,
    true,  // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    args,  // arguments
)
```

### Queue Binding
```go
err = ch.QueueBind(
    queueName,
    routingKey,
    "woragis.tasks",
    false,
    nil,
)
```

## Benefits
- Durable queues
- Dead letter queues
- Proper routing
- Fault tolerance

## Challenges
- Configuration complexity
- Need to understand RabbitMQ
- Queue management
- Monitoring

## Lessons Learned
- Durable queues essential
- Dead letter queues crucial
- Proper routing important
- Monitoring helps

## Future Improvements
- Queue metrics
- Queue monitoring dashboard
- Automatic queue creation
- Queue optimization
