from pymongo import UpdateOne
import requests
import os
from . import db,scheduler_logs
reserve_col = db["log_reserve"]

def CR_cancel_reserve(order_id):
    filter = {"order_id": order_id}
    update = {"$set": {"status": "Cancel",},}

    print("Oder Id:"+order_id+"has benn cancel")
    if reserve_col.find_one(filter) is not None:
        reserve_col.update_one(filter,update)
    else:
        print("Not Found order_id")
    scheduler_logs.remove_job("CR_"+order_id)

def CER_cancel_extend_reserve(order_id):
    filter = {"order_id": order_id}
    update = {
        "$set": {"status": "Parking", "is_extend": False},
        "$inc": {"hour_end": -1}
    }
    print("Oder Id:"+order_id+"extendd has benn cancel")
    if reserve_col.find_one(filter) is not None:
        reserve_col.update_one(filter,update)
    else:
        print("Not Found order_id")
    scheduler_logs.remove_job("CER_"+order_id)

def CRIA_cancel_reserve_in_advance(order_id):
    url =  'http://'+os.getenv("HOST")+'/webhook/internal/notification/confirm_reserve_in_advance_notification/cancel'
    myobj = {'order_id': order_id[0],'email':order_id[1]}
    print("-------------------------\n")
    print("Request to",url)
    x = requests.get(url, params = myobj)
    print("Response is",x.text)
    print("-------------------------\n")
    scheduler_logs.remove_job("CRIA_"+order_id[0]+","+order_id[1])

def BTOR_timeout_reserve_job(order_id):
    url =  'http://'+os.getenv("HOST")+'/webhook/internal/cronjob/before_timeout_reserve'
    myobj = {'order_id': order_id}
    print("-------------------------\n")
    print("Request to",url)
    x = requests.post(url, params = myobj)
    print("Response is",x.text)
    print("-------------------------\n")
    scheduler_logs.remove_job("BTOR_"+order_id)

def TOR_timeout_reserve_job(order_id):
    url =  'http://'+os.getenv("HOST")+'/webhook/internal/cronjob/timeout_reserve'
    myobj = {'order_id': order_id}
    print("-------------------------\n")
    print("Request to",url)
    x = requests.post(url, params = myobj)
    print("Response is",x.text)
    print("-------------------------\n")
    scheduler_logs.remove_job("TOR_"+order_id)

def ATOR_timeout_reserve_job(order_id):
    url =  'http://'+os.getenv("HOST")+'/webhook/internal/cronjob/after_timeout_reserve'
    myobj = {'order_id': order_id}
    print("-------------------------\n")
    print("Request to",url)
    x = requests.post(url, params = myobj)
    print("Response is",x.text)
    print("-------------------------\n")
    scheduler_logs.remove_job("ATOR_"+order_id)



