#!/usr/bin/env python
import pika
import os

connection = pika.BlockingConnection(pika.ConnectionParameters(os.environ.get('RABBITMQ_PORT_5672_TCP_ADDR'), 5672, '/'))
channel = connection.channel()

channel.queue_declare(queue='hello')

def callback(ch, method, properties, body):
    print " [x] Received %r" % (body,)


channel.basic_consume(callback,
                      queue='hello',
                      no_ack=True)

print " [*] Waiting for messages..."
channel.start_consuming()
