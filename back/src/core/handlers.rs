use actix_web::{Responder, get};

#[get("/ping")]
pub async fn ping() -> impl Responder {
    format!("pong")
}
