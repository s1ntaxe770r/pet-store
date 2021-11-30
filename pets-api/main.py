import os
from flask import Flask,request
from flask_sqlalchemy import SQLAlchemy
from flask.json import jsonify
import json
import threading
from dataclasses import dataclass
import pika
from pika import channel
from pika import connection



app =  Flask(__name__)
db = SQLAlchemy(app)
@dataclass  
class Pet(db.Model):
    # create pet id, name, category column
    pet_id: int
    name: str
    notes: str
    category: str
    pet_id = db.Column(db.Integer, primary_key = True)
    name = db.Column(db.String(100), nullable = False)
    notes = db.Column(db.String(50), nullable = False )
    category = db.Column(db.String(50), nullable = False )


RMQ_HOST = os.getenv('RMQ_HOST')
app.config["SQLALCHEMY_DATABASE_URI"] = 'sqlite:///pets.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
parameters = pika.URLParameters(RMQ_HOST)
db.create_all()
db.session.commit()
    


def serialize(message:Pet):
    return json.dumps(
        {
            'id': message.pet_id,
            'name': message.name,
            'notes': message.notes,
            'category': message.category,
        })
        

@app.post("/api/pets")
def create():
     data = request.json
     name =  data["name"]
     notes = data["notes"]
     category = data["category"]
     NewPet = Pet(name=name,notes=notes,category=category)
     db.session.add(NewPet)
     db.session.commit()
     conn = pika.BlockingConnection(parameters)
     channel  = conn.channel()
     channel.queue_declare("pets",durable=False)
     channel.basic_publish(exchange='', routing_key='pets', body=serialize(NewPet))
     print(" [x] Sent 'pet to RabbitMQ!'")
     return jsonify(
        {
            'id': NewPet.pet_id,
            'name': NewPet.name,
            'notes': NewPet.notes,
            'category': NewPet.category,
        }   
    )


@app.get("/api/pets")
def get_pets():
    r = db.session.query(Pet).all()
    print(r)
    return jsonify(r)


if __name__ == "__main__":
    app.run(host="0.0.0.0",port=6000)
