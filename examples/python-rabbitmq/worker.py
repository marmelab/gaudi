import pika, os, time, logging

# Connect to RabbitMQ
logging.basicConfig(format='%(LEVELNAME)s:%(message)s', level=logging.CRITICAL)
logging.getLogger('pika').setLevel(logging.CRITICAL)

connection = pika.BlockingConnection(pika.ConnectionParameters(os.environ.get('RABBITMQ_PORT_5672_TCP_ADDR'), 5672, '/', connection_attempts=100))
channel = connection.channel()

channel.queue_declare(queue='hello')

print " [*] Waiting for messages..."

# Print received messages
def callback(ch, method, properties, body):
    print " [x] Received %r" % (body,)

channel.basic_consume(callback,
                      queue='hello',
                      no_ack=True)

channel.start_consuming()
