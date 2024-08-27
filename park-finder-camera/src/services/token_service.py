from flask import jsonify
import requests
import os
from dotenv import load_dotenv

load_dotenv()

def get_access_token():
    appKey = os.getenv("HikCentral_appKey")
    secretKey = os.getenv("HikCentral_secretKey")

    headers = {'Content-type': 'application/json'}
    payload = dict(appKey=appKey, secretKey=secretKey)

    r = requests.post(f"{os.getenv('HOST_HikCentral')}/api/hccgw/platform/v1/token/get",json=payload,headers=headers)
    
    if r.status_code == 200:
        return r.json()