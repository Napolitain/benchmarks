use axum::{routing::get, Json, Router};
use serde::Serialize;

#[derive(Serialize)]
struct Message {
    message: &'static str,
}

async fn hello() -> Json<Message> {
    Json(Message {
        message: "Hello, World!",
    })
}

#[tokio::main(flavor = "current_thread")]
async fn main() {
    let app = Router::new().route("/", get(hello));

    println!("Rust Axum server listening on :8080 (single-threaded)");
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
