from django.shortcuts import render
import os
import requests

# Create your views here.


INVENTORY_SERVICE = os.getenv("INVENTORY_SERVICE_URL")


def index(request):
    data = requests.get(INVENTORY_SERVICE)
    jsn_payload= data.json()
    pet_count = jsn_payload["total_pets"]
    print(jsn_payload)
    if data.status_code != 200:
        return "Sorry we were unable to load this page right now"
    return render(request,"index.html",{'inventory':pet_count})