from flask import request, Response, json, Blueprint
import requests
import os
from dotenv import load_dotenv
from ..services.token_service import get_access_token
from ..services.camera_service import capture_picture, detection_car

cameras = Blueprint("cameras", __name__)

load_dotenv()

@cameras.route('/getpicture', methods = ["POST"])
def get_picture():

    token = get_access_token()
    access_token = token['data']['accessToken']

    pictureUrl = capture_picture(access_token , "L38082195", 1)
    
    return Response(
        response=pictureUrl['data']['captureUrl'],
        status=200,
        mimetype='application/json'
    )

@cameras.route('/checkcar', methods = ["POST"])
def check_car():

    token = get_access_token()
    access_token = token['data']['accessToken']

    pictureUrl = capture_picture(access_token , "L38082195", 1)

    car = detection_car(pictureUrl)
    
    return Response(
        response=car,
        status=200,
        mimetype='application/json'
    )
    