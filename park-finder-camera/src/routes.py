from flask import Blueprint
from src.controllers.camera_controller import cameras
from src.controllers.webhook_controller import webhook
# main blueprint to be registered with application
api = Blueprint('api', __name__)

# register cameras with api blueprint
api.register_blueprint(cameras, url_prefix="/cameras")
api.register_blueprint(webhook, url_prefix="/webhook")