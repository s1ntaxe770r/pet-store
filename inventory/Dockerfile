FROM python:3.9.1-slim 

WORKDIR /app 

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1


RUN pip install --upgrade pip
RUN pip install --no-cache-dir flask redis pika gunicorn

COPY . .

# start server with 4 workers 

CMD ["python","/app/inventory/main.py"]    