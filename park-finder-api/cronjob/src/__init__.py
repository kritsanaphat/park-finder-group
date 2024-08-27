import os
from http import HTTPStatus

import pytz
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.date import DateTrigger
from dotenv import load_dotenv
from flask import Flask, jsonify, request
from flask_cors import CORS
from flask_cors import CORS
from pymongo import MongoClient
from . import job

status = load_dotenv()
if not status:
    print("Cannot load environment file.")


# Load env
mongodb_uri = os.getenv("MONGODB_URI")

# Setup scheduler
tz = pytz.timezone("Asia/Bangkok")
scheduler = BackgroundScheduler(timezone=tz)

# Setup flask
app = Flask(__name__)
cors = CORS()

# Setup db
mongodb_client = MongoClient(mongodb_uri)
db = mongodb_client.get_database(os.getenv("DATABASE_NAME"))
scheduler_col = db["log_scheduler"]

def fetch_scheduler():
    result = list(scheduler_col.find())
    job_list = []
    for i in result:
        job_list.append({i["job_id"]:i["trigger"]})
        scheduler_col.delete_one({"job_id": i["job_id"]})


    job.fetch_scheduler_logs(job_list)
    print("Fetch All Scheduler in Mongodb",job_list)
    result = list(scheduler_col.find())




def create_app() -> Flask:
    # Config cors
    cors.init_app(app)

    # Register blueprint
    from .job import job_api
    app.register_blueprint(job_api)

    # Show all route
    print("All routes")
    for route in app.url_map.iter_rules():
        print(f"{route.methods} {route}")

    #fetch all scheduler logs in db
    fetch_scheduler()


    # Start schedule
    scheduler.start()

    return app

__all__ = [
    "create_app",
]

