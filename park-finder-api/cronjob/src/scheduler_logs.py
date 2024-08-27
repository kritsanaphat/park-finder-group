from pymongo import UpdateOne

from . import db
scheduler_col = db["log_scheduler"]


def save_job(job_id,trigger):
    data = {
        "job_id" : job_id,
        "trigger" : trigger,}
    result = scheduler_col.insert(data)
    print(result)

def remove_job(job_id):
    print(job_id)
    filter = {
        "job_id" : job_id,
        }
    result = scheduler_col.delete_one(filter)
    print(result)

