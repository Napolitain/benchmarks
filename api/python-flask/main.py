from flask import Flask
import orjson

app = Flask(__name__)

@app.route("/")
def hello():
    response = app.response_class(
        response=orjson.dumps({"message": "Hello, World!"}),
        status=200,
        mimetype='application/json'
    )
    return response

if __name__ == "__main__":
    from gunicorn.app.base import BaseApplication

    class StandaloneApplication(BaseApplication):
        def __init__(self, app, options=None):
            self.options = options or {}
            self.application = app
            super().__init__()

        def load_config(self):
            for key, value in self.options.items():
                if key in self.cfg.settings and value is not None:
                    self.cfg.set(key.lower(), value)

        def load(self):
            return self.application

    options = {
        "bind": "0.0.0.0:8080",
        "workers": 1,
    }
    print("Flask server listening on :8080 (gunicorn, 1 worker)")
    StandaloneApplication(app, options).run()
