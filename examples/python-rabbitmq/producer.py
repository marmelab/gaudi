#!/usr/bin/env python
import pika
import os

import logging

logging.basicConfig(format='%(levelname)s:%(message)s', level=logging.CRITICAL)
logging.getLogger('pika').setLevel(logging.CRITICAL)

connection = pika.BlockingConnection(pika.ConnectionParameters(os.environ.get('RABBITMQ_PORT_5672_TCP_ADDR'), 5672, '/'))
channel = connection.channel()

channel.queue_declare(queue='hello')

channel.basic_publish(exchange='',
                      routing_key='hello',
                      body='Hello World!')

print " [x] Sent 'Hello World!'"

connection.close()
