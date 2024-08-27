import os
from http import HTTPStatus

import pytz
from dotenv import load_dotenv
from flask import Flask, jsonify, request
from flask_cors import CORS
from flask_cors import CORS



status = load_dotenv()
if not status:
    print("Cannot load environment file.")


# Setup flask
app = Flask(__name__)
cors = CORS()

def create_app() -> Flask:
    # Config cors
    cors.init_app(app)

    # Register blueprint
    from .demo import demo_api
    app.register_blueprint(demo_api)

    # Show all route
    print("All routes")
    for route in app.url_map.iter_rules():
        print(f"{route.methods} {route}")

    return app

__all__ = [
    "create_app",
]
