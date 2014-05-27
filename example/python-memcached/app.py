import memcache
import os

mc = memcache.Client([os.environ.get('MEMCACHED_PORT_11211_TCP_ADDR')], debug=0)

mc.set("fou", "barre")
value = mc.get("fou")

print("value is :", value)
