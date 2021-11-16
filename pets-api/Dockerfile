FROM python:3.9.1-slim 

RUN mkdir /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1


RUN pip install --upgrade pip
RUN pip install --no-cache-dir flask flask-sqlalchemy pika gunicorn
COPY .  /app


 
EXPOSE 6000

CMD ["python","/app/pets-api/main.py"]    