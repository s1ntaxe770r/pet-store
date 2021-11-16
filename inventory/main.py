import datetime
from flask import Flask, jsonify
from flask.globals import request
import json
from pika import connection
import redis 
from os import getenv
from typing import List
import pika
import threading


redis_url:str = getenv("REDIS_URL")
rmq:str =  getenv("RMQ_HOST")

app = Flask(__name__)
client = redis.from_url(redis_url)


def main():
    parameters = pika.URLParameters(rmq)
    connection = pika.SelectConnection(parameters)
    connection.ioloop.start()
    channel = connection.channel()
    channel.queue_declare(queue='pets')
   
    def callback(ch, method, properties, body):
        print(" [x] Received %r" % body)
        data =  json.loads(body)
        print(data)
        pet_id = data["name"]
        print("[*] inserted pet id"+str(pet_id)+" into redis")
        pet_no = client.get("total_pets")
        if pet_no == None:
            client.set("total_pets",1)
        else:
            client.set("total_pets",int(pet_no)+1)
     
    channel.basic_consume(queue='pets', on_message_callback=callback, auto_ack=True)

    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()


@app.get('/inventory/stats')
def stats():
    total = client.get("total_pets")
    return jsonify(
        {"total_pets":int(total)})

consumer_thread = threading.Thread(target=main)
consumer_thread.start()

if __name__ == '__main__':
    app.run(host="0.0.0.0",port=7000)