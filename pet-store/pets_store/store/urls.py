from os import name
from django.urls import path
from . import views

urlpatterns = [
    path('', views.index),
    path('create/',views.create_pet,name="create")
]

 