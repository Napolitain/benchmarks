from fastapi import FastAPI
from fastapi.responses import JSONResponse
import uvicorn

app = FastAPI(default_response_class=JSONResponse)

@app.get("/")
async def hello():
    return {"message": "Hello, World!"}

if __name__ == "__main__":
    print("FastAPI server listening on :8080")
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8080,
        log_level="critical",
        access_log=False,
    )
