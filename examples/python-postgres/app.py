import psycopg2
import os

try:
    conn = psycopg2.connect("dbname='project' user='docker' host='" + os.environ.get('DB_PORT_5432_TCP_ADDR') + "' password='docker'")
except:
    print ("Unable to connect to the database")


cur = conn.cursor()

# Create table
cur.execute("""CREATE TABLE IF NOT EXISTS book (title varchar(100) NOT NULL);""")

# Add an user
values = {"title":"Harry Potter and the Half-Blood Prince"}
cur.execute("INSERT INTO book (title) VALUES (%(title)s);", values)

# Fetch the user
cur.execute("SELECT * from book")
rows = cur.fetchall()
for row in rows:
    print ("- ", row[0])
