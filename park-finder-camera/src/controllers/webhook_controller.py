from flask import request, Response, json, Blueprint,jsonify
from dotenv import load_dotenv
from PIL import Image
from io import BytesIO

import requests
import os
import requests

from ..services.token_service import get_access_token
from ..services.webhook_service import check_exists_licence_plate
from ..services.camera_service import capture_picture


webhook = Blueprint("webhook", __name__)
load_dotenv()



@webhook.route('/check_licence_plate', methods=["POST"])
def check_licence_plate():
    try:
        data = request.get_json()
        if not data:
            return jsonify({'error': 'No JSON data received'}), 400
        
        module_code = data.get("module_code")
        customer_car_licence_plate = data.get("customer_license_plate")
        if not module_code or not customer_car_licence_plate:
            return jsonify({'error': 'Missing required data: module_code or customer_license_plate'}), 400

        exists = False
        token = get_access_token()
        access_token = token['data']['accessToken']

        pictureUrl = capture_picture(access_token, module_code, 1)
        print(pictureUrl)
        capture_url = pictureUrl['data']['captureUrl']

        response = requests.get(capture_url)
        response.raise_for_status()

        image_data = BytesIO(response.content)
        files = {'image': ('filename.jpg', image_data, 'image/jpeg')}
        headers = {'Apikey': os.getenv('AI_THAI_API_KEY')}
        
        response = requests.post(os.getenv('AI_THAI_HOST'), files=files, headers=headers)
        if response.text.strip() != "":
            print("Response is", response.json())
            exists = check_exists_licence_plate(customer_car_licence_plate, response.json())
        else:
            response_message = "Can't detect"
            return jsonify({'response': response_message}), 200

    except Exception as e:
        error_message = f"An error occurred: {e}"
        return jsonify({'error': error_message}), 500

    response_message = str(exists)
    return jsonify({'response': response_message}), 200


@webhook.route('/capture_camera', methods = ["POST"])
def capture_camera():
    module_code = request.args.get("module_code")
    try:
        token = get_access_token()
        access_token = token['data']['accessToken']
        pictureUrl = capture_picture(access_token, module_code, 1)
        capture_url = pictureUrl['data']['captureUrl']

    except Exception as e:
        error_message = f"An error occurred: {e}"
        response = jsonify({'error': error_message})
        response.status_code = 500 
        return response
        
    response = jsonify({'response': capture_url})
    response.status_code = 200  
    return response 

