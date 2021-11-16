FROM python:3.9.1-slim
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN pip install --no-cache-dir  django gunicorn requests
EXPOSE 8000
RUN mkdir /pet-store
WORKDIR /pet-store
COPY . /pet-store


