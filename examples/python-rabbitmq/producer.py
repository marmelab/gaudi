import pika, os, time, logging

# Connect to RabbitMQ
logging.basicConfig(format='%(LEVELNAME)s:%(message)s', level=logging.CRITICAL)
logging.getLogger('pika').setLevel(logging.CRITICAL)

connection = pika.BlockingConnection(pika.ConnectionParameters(os.environ.get('RABBITMQ_PORT_5672_TCP_ADDR'), 5672, '/', connection_attempts=100))
channel = connection.channel()

channel.queue_declare(queue='hello')

i = 0

while 1:

	body = 'Hello World #' + str(i)
	channel.basic_publish(exchange='',
                      routing_key='hello',
                      body=body)

	print " [x] Sent '" + body + "'"
	i += 1
	
	time.sleep(1)

connection.close()
