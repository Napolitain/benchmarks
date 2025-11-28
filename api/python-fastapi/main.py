from fastapi import FastAPI
from fastapi.responses import ORJSONResponse
import uvicorn

app = FastAPI(default_response_class=ORJSONResponse)

@app.get("/")
async def hello():
    return {"message": "Hello, World!"}

if __name__ == "__main__":
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8080,
        log_level="critical",
        access_log=False,
        loop="uvloop",
    )
    print("FastAPI server listening on :8080 (single-threaded)")
