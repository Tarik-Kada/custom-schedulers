FROM python:3.9-slim

WORKDIR /app
COPY scheduler.py /app

RUN pip install flask
RUN pip install waitress

CMD ["python", "scheduler.py"]
