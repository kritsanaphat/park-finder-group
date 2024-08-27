from flask import Flask
from flask_cors import CORS
from http import HTTPStatus

app = Flask(__name__)
cors = CORS()

def create_app() -> Flask:
    # Config cors
    cors.init_app(app)

    # Register blueprint
    from src.routes import api
    app.register_blueprint(api, url_prefix = "/api")

    # Show all route
    print("All routes")
    for route in app.url_map.iter_rules():
        print(f"{route.methods} {route}")

    #fetch all scheduler logs in db

    return app

__all__ = [
    "create_app",
]

