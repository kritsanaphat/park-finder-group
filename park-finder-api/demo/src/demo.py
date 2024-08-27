import gspread
from oauth2client.service_account import ServiceAccountCredentials
from datetime import datetime,timedelta
import pytz
import time

from pymongo.mongo_client import MongoClient
from flask import Flask, jsonify,Blueprint

app = Flask(__name__)

demo_api = Blueprint("demo_api", __name__, url_prefix="/demo_api")

@demo_api.route('/redeem_from_park_finder')
def redeem_from_park_finder():
    response = {'barcode_url':'https://parkingadmindata.s3.ap-southeast-1.amazonaws.com/ppt/360_F_255979498_vewTRAL5en9T0VBNQlaDBoXHlCvJzpDl.jpg?response-content-disposition=inline&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEEAaDmFwLXNvdXRoZWFzdC0xIkgwRgIhAJ5coXVM8W6crkPno7FpUBHRBBaB8l3U3FLYUuxyuEc8AiEAyFoVEAb8c85bTtR65A5FZZ0Xb2DWxmtj6UejWr6SeTsq5AIIGRAAGgwyMTExMjU3MzI0OTYiDE8fQRlR3aSWEhz%2BwyrBAp27Sfk6N%2BgptEs1TtxAxj2q0XhN6XGxxXiUb4AYBfDHYbmBKTXDH7aIqduAbAExr%2FxHdidFQp%2B9%2FjxfXfYaWGxeGVhDXDPtetTdQ%2FM8gei63lQiuSBZNHqfcUHQkwnew0ZLYGLLtYDTDP6ZJQ5b0GeRmN6g2i39GbcAxhJD5bclPFxq9X6M2L3W0WaGWdQWg3qeXB%2B7M0La3ZaShhS1TP7lVbjmXik9EAAnv8VKOUyx8MEX%2Fig3ekvdfCkxd1PhBz5956ycMrC2yQ2dl8OnKr2ZvjAmIGVZVELiMyLbllluIKUK%2BrQkjbyBqt%2BhZUupIrUCT%2FRKMFGOW%2FVQ9fPiCEykjKKcFiVs3%2Ft6cEwszmVm0kypIi3hhncipYacNx5E3oV5L42MkAxpt3nqnqjTJ9HyKGlihfSOkyW2Qd9K6d25DDDroZmuBjqyAlP5XquiDZplMgWW2dbRqDbS3VorCBrRz1lGbaGgwutMjCuW%2FeukIcq8K3WPjQgODL032bRazVao9hsqxHuRE5Purk8RWDKAw0E46niQofZ%2Fvvvym01u5Lg%2FaMkk9ZnuYr%2F8kmXPTg89axIgk1k7yPLFXFhICHYTv0wgeJ17M1DAmbSBY3bXIlSWYifj2SKwFOkvNEonB5qCk2QxaQ1WY9w6Vk2ZtmVWpWlxe8aPQXquBUd44N1cUs8AcIZVB3I3vTZNFbtjjPnkyStum9KOK8uBC%2F8i1XiY5RDFTEu%2F%2BCSLHgrff0dnQ0HRkiNxn4bSrZ6nfMwC%2F03h6leSUEP7CYiWkBOwN9w%2BBeUI7y5XIkbNstEt5CHpveWxo%2BbyXQUxpGcTVnsnSZxil1kVVvuISQd7gw%3D%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20240209T162224Z&X-Amz-SignedHeaders=host&X-Amz-Expires=300&X-Amz-Credential=ASIATCKATASIFVVGEI36%2F20240209%2Fap-southeast-1%2Fs3%2Faws4_request&X-Amz-Signature=ba53f6e84550329f1ebe8479914185c279a55a23a9f29dd1f1d7c75ec814093d'}
    return jsonify(response)  

@demo_api.route('/')
def hello():
    response = {'message': 'Hello, World! This is your Flask API.'}
    return jsonify(response)