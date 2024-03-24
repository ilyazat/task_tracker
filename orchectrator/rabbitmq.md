### Message acknowledgment
RabbitMQ doesn't allow you to redefine an existing queue with different parameters and will return an error to any program that tries to do that.


### Durability


In this tutorial we will use manual message acknowledgements by passing a false for the "auto-ack" argument and then send a proper acknowledgment from the worker with d.Ack(false) (this acknowledges a single delivery), once we're done with a task.
```
err = ch.PublishWithContext(ctx,
  "",           // exchange
  q.Name,       // routing key
  false,        // mandatory
  false,
  <b>amqp.Publishing { </b>
    DeliveryMode: amqp.Persistent,
    ContentType:  "text/plain",
    Body:         []byte(body),
})
```
- Marking messages as persistent doesn't fully guarantee that a message won't be lost. Although it tells RabbitMQ to save the message to disk, there is still a short time window when RabbitMQ has accepted a message and hasn't saved it yet.
- If you need a stronger guarantee then you can use publisher confirms.
- durable stands for file saving data

### Fair Dispatch

Проблема: В ситуации, когда два четный воркер получает легкую работу, а нечетный -- тяжелую, то кто-то будет постоянно работать
и не будет свободен. Кролик просто распределяет сообщения, когда сообщение приходит в очередь, и он не смотрит на количество unacks

Решение: поставить prefetch count to 1 -- суть: "Не давай больше 1 сообщения единовременно воркеру, пока он не отдал аск, а отдай другому"

```
err = ch.Qos(
  1,     // prefetch count
  0,     // prefetch size
  false, // global
)
failOnError(err, "Failed to set QoS")
```

### Publish/Subscribe

У очереди концепция: одна задача -- один воркер. Pub/Sub -- это про большее количество потребителей.

The core idea in the messaging model in RabbitMQ is that the producer never sends any messages directly to a queue. Actually, quite often the producer doesn't even know if a message will be delivered to any queue at all.

Instead, the producer can only send messages to an exchange. An exchange is a very simple thing. On one side it receives messages from producers and the other side it pushes them to queues. The exchange must know exactly what to do with a message it receives. Should it be appended to a particular queue? Should it be appended to many queues? Or should it get discarded. The rules for that are defined by the exchange type.

```
err = ch.ExchangeDeclare(
  "logs",   // name
  "fanout", // type
  true,     // durable
  false,    // auto-deleted
  false,    // internal
  false,    // no-wait
  nil,      // arguments
)
```
We used `""` exchange above.

The fanout exchange is very simple. As you can probably guess from the name, it just broadcasts all the messages it receives to all the queues it knows. And that's exactly what we need for our logger.

### Temporary Queues
Being able to name a queue was crucial for us -- we needed to point the workers to the same queue. Giving a queue a name is important when you want to share the queue between producers and consumers.

In the amqp client, when we supply queue name as an empty string, we create a non-durable queue with a generated name:

```
q, err := ch.QueueDeclare(
  "",    // name
  false, // durable
  false, // delete when unused
  true,  // exclusive
  false, // no-wait
  nil,   // arguments
)
```

When the connection that declared it closes, the queue will be deleted because it is declared as exclusive.

### Bindings

Now we need to tell the exchange to send messages to our queue. **That relationship between exchange and a queue is called a binding.**

```
err = ch.QueueBind(
  q.Name, // queue name
  "",     // routing key
  "logs", // exchange
  false,
  nil,
)
```

Bindings can take an extra routing_key parameter. To avoid the confusion with a Channel.Publish parameter we're going to call it a binding key. This is how we could create a binding with a key:
```
err = ch.QueueBind(
        q.Name,    // queue name
        "black",   // routing key
        "logs",    // exchange
        false,
        nil
)
```
The meaning of a binding key depends on the exchange type.

### Direct Exchange

Our logging system from the previous tutorial broadcasts all messages to all consumers. We want to extend that to allow filtering messages based on their severity. 

We were using a fanout exchange, which doesn't give us much flexibility - it's only capable of mindless broadcasting.

**The routing algorithm behind a direct exchange is simple - a message goes to the queues whose binding key exactly matches the routing key of the message.**


####  Multiple Bindings

It is perfectly legal to bind multiple queues with the same binding key. In our example we could add a binding between X and Q1 with binding key black. In that case, the direct exchange will behave like fanout and will broadcast the message to all the matching queues. A message with routing key black will be delivered to both Q1 and Q2.


# Topics

Although using the direct exchange improved our system, it still has limitations - it can't do routing based on multiple criteria.

In our logging system we might want to subscribe to not only logs based on severity, but also based on the source which emitted the log. You might know this concept from the syslog unix tool, which routes logs based on both severity (info/warn/crit...) and facility (auth/cron/kern...).

### Topic Exchange

Messages sent to a topic exchange can't have an arbitrary routing_key - it must be a list of words, delimited by dots. 
The words can be anything, but usually they specify some features connected to the message. 
A few valid routing key examples: "stock.usd.nyse", "nyse.vmw", "quick.orange.rabbit". 
There can be as many words in the routing key as you like, up to the limit of 255 bytes.


a message sent with a particular routing key will be delivered to all the queues that are bound with a matching binding key. However there are two important special cases for binding keys:

    * (star) can substitute for exactly one word.
    # (hash) can substitute for zero or more words.


# RPC

## Callback Queue

A client sends a request message and a server replies with a response message. In order to receive a response we need to send a 'callback' queue address with the request. We can use the default queue. Let's try it:
```
q, err := ch.QueueDeclare(
  "",    // name
  false, // durable
  false, // delete when unused
  true,  // exclusive
  false, // noWait
  nil,   // arguments
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err = ch.PublishWithContext(ctx,
  "",          // exchange
  "rpc_queue", // routing key
  false,       // mandatory
  false,       // immediate
  amqp.Publishing{
    ContentType:   "text/plain",
    CorrelationId: corrId,
    ReplyTo:       q.Name,
    Body:          []byte(strconv.Itoa(n)),
})
```

- `persistent`: Marks a message as persistent (with a value of true) or transient (false). You may remember this property from the second tutorial.
- `content_type`: Used to describe the mime-type of the encoding. For example for the often used JSON encoding it is a good practice to set this property to: application/json.
- `reply_to`: Commonly used to name a callback queue.
- `correlation_id`: Useful to correlate RPC responses with requests.

здесь мы создаем отдельную колбек очередь на каждый рпс-запрос. Можно сделать эффективнее -- создать один коллбек на очередь
Создается проблема соотношения запрос и ответа на запрос.That's when the correlation_id property is used. 

We're going to set it to a unique value for every request. 

Later, when we receive a message in the callback queue we'll look at this property, and based on that we'll be able to match a response with a request. 

If we see an unknown correlation_id value, we may safely discard the message - it doesn't belong to our requests.



Our RPC will work like this:

- When the Client starts up, it creates an anonymous exclusive callback queue.
- For an RPC request, the Client sends a message with two properties: reply_to, which is set to the callback queue and correlation_id, which is set to a unique value for every request.
- The request is sent to an rpc_queue queue.
- The RPC worker (aka: server) is waiting for requests on that queue. When a request appears, it does the job and sends a message with the result back to the Client, using the queue from the reply_to field.
- The client waits for data on the callback queue. When a message appears, it checks the correlation_id property. If it matches the value from the request it returns the response to the application.

