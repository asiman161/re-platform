use actix_web::{Responder, get, HttpResponse, http::header::ContentType};
use crate::core::storage;

#[get("/ping")]
pub async fn ping() -> impl Responder {
    format!("pong")
}

#[get("/users")]
pub async fn users() -> impl Responder {
    match storage::new().get_users() {
        Ok(users) => {
            let body = serde_json::to_string(&users).unwrap();
            HttpResponse::Ok().content_type(ContentType::json()).body(body)
        }
        Err(err) => HttpResponse::InternalServerError().body(err.to_string())
    }
}
