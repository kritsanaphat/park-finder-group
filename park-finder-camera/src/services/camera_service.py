from flask import jsonify
import requests
import os
from dotenv import load_dotenv
from PIL import Image
import cv2
import numpy as np

load_dotenv()

def capture_picture(Token, deviceSerial, channelNo):

    headers = {'Content-type': 'application/json', 'Token': Token}
    payload = dict(deviceSerial=deviceSerial, channelNo=channelNo)

    r = requests.post(f"{os.getenv('HOST_HikCentral')}/api/hccgw/resource/v1/device/capturePic",json=payload,headers=headers)
    
    if r.status_code == 200:
        return r.json()

def detection_car(pictureUrl):
    image = Image.open(requests.get(pictureUrl, stream=True).raw)
    # image = cv2.imread('testImg.png')
    image = image.resize((450,250))
    image_arr = np.array(image)

    grey = cv2.cvtColor(image_arr,cv2.COLOR_BGR2GRAY)
    blur = cv2.GaussianBlur(grey,(5,5),0)
    dilated = cv2.dilate(blur,np.ones((3,3)))
    kernel = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (2, 2))
    closing = cv2.morphologyEx(dilated, cv2.MORPH_CLOSE, kernel) 

    car_cascade_src = '../config/cars.xml'
    car_cascade = cv2.CascadeClassifier(car_cascade_src)

    # cars = car_cascade.detectMultiScale(image_arr, 1.1, 1)
    # cars_grey = car_cascade.detectMultiScale(grey, 1.1, 1)
    # cars_blur = car_cascade.detectMultiScale(blur, 1.1, 1)
    # cars_dilated = car_cascade.detectMultiScale(dilated, 1.1, 1)
    cars_closing = car_cascade.detectMultiScale(closing, 1.1, 1)

    return len(cars_closing)

