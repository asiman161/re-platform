use actix_web::{Error, error, get, http::header::ContentType, HttpRequest, HttpResponse, Responder, web};
use actix_web_actors::ws;

use crate::chat;
use crate::core::storage;

#[get("/ping")]
pub async fn ping() -> impl Responder {
    HttpResponse::Ok().content_type(ContentType::plaintext()).body("pong")
}

#[get("/users")]
pub async fn users(pool: web::Data<storage::DbPgPool>) -> actix_web::Result<impl Responder> {
    let users = web::block(move || {
        let mut con = pool.get()?;

        storage::get_users(&mut con)
    })
        .await?
        .map_err(error::ErrorInternalServerError)?;

    Ok(HttpResponse::Ok().json(users))
}

#[get("/users/{user_id}")]
pub async fn get_user(
    pool: web::Data<storage::DbPgPool>,
    user_id: web::Path<i32>,
) -> actix_web::Result<impl Responder> {
    let user_uid = user_id.into_inner();

    // use web::block to offload blocking Diesel queries without blocking server thread
    let user = web::block(move || {
        // note that obtaining a connection from the pool is also potentially blocking
        let mut conn = pool.get()?;

        storage::find_user_by_id(&mut conn, user_uid)
    })
        .await?
        // map diesel query errors to a 500 error response
        .map_err(error::ErrorInternalServerError)?;

    Ok(match user {
        // user was found; return 200 response with JSON formatted user object
        Some(user) => HttpResponse::Ok().json(user),

        // user was not found; return 404 response with error message
        None => HttpResponse::NotFound().body(format!("No user found with UID: {user_uid}")),
    })
}

#[get("/room/{id}")]
pub async fn room(req: HttpRequest, stream: web::Payload, path: web::Path<u32>) -> Result<HttpResponse, Error> {
    let room_id = path.into_inner();
    ws::start(chat::session::new(room_id), &req, stream)
}
