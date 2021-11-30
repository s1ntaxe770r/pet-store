import json
import re
from django.http.response import HttpResponse
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

def create_pet(request):
    print(request.method)
    print(request.POST.get('name'))
    print(request.POST.get('category'))
    if request.POST:
        req = request.POST.dict()
        name = req.get("name") 
        notes  =  req.get("notes")
        category = req.get("category")
        headers = {'Content-type': 'application/json'}
        # post message to pet service
        resp = requests.post(PETS_SERVICE,headers=headers,json={"name":name,"notes":notes,"category":category})
        resp_data = resp.json()
        if resp.status_code != 200:
            return HttpResponse("unable to create pet")
        return HttpResponse(f'<h1>created pet {resp_data["id"]}</h1>')
    return render(request,"create.html")

    