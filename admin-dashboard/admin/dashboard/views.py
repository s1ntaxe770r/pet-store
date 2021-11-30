from django.shortcuts import render
import os
import requests

# Create your views here.


INVENTORY_SERVICE = os.getenv("INVENTORY_SERVICE_URL")
CATEGORY_SERVICE = os.getenv("CATEGORY_SERVICE_URL")


def index(request):
    data = requests.get(INVENTORY_SERVICE)
    category_data = requests.get(CATEGORY_SERVICE+"/pets/categories/reptile")
    jsn_payload= data.json()
    category_json = category_data.json()
    reptiles = category_json["data"]
    pet_count = jsn_payload["total_pets"]

    print(jsn_payload)
    if data.status_code != 200:
        return "Sorry we were unable to load this page right now"
    return render(request,"index.html",{'inventory':pet_count},{'reptiles':reptiles})