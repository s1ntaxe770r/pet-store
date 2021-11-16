from django.shortcuts import render
import os
import requests

# Create your views here.


PETS_SERVICE = os.getenv("PETS_SERVICE_URL")


def index(request):
    data = requests.get(PETS_SERVICE)
    jsn_payload= data.json()
    if data.status_code != 200:
        return "Sorry we were unable to load this page right now"
    return render(request,"index.html",{'pets':jsn_payload})