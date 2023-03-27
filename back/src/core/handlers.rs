use actix_web::{Error, get, http::header::ContentType, HttpRequest, HttpResponse, Responder, web};
use actix_web_actors::ws;

use crate::core::storage;
use crate::chat;

#[get("/ping")]
pub async fn ping() -> impl Responder {
    HttpResponse::Ok().content_type(ContentType::plaintext()).body("pong")
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

#[get("/room/{id}")]
pub async fn room(req: HttpRequest, stream: web::Payload, path: web::Path<u32>) -> Result<HttpResponse, Error> {
    let room_id = path.into_inner();
    ws::start(chat::session::new(room_id), &req, stream)
}
