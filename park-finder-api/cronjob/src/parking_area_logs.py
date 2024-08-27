from pymongo import UpdateOne
from bson import ObjectId

from . import db,scheduler_logs
reserve_col = db["parking_area"]



def update_open_area_status(parking_id):
    filter = {"_id": ObjectId(parking_id)}
    update = {"$set": {"open_status": True,"time_stamp_close":None},}

    print("open_status of :" +parking_id+ " change to True")
    if reserve_col.find_one(filter) is not None:
        reserve_col.update_one(filter,update)
    else:
        print("Not Found parking_area")
    scheduler_logs.remove_job("OPA"+parking_id)
