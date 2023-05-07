use actix_cors::Cors;
use actix_web::{web};

use crate::core::handlers;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api")
            .service(handlers::ping)
            .service(handlers::users)
            .service(handlers::get_user)
            .service(handlers::room)
    );
}

pub fn cors() -> Cors {
    let cors = Cors::default()
        .allow_any_origin()
        .allow_any_method()
        .supports_credentials();

    return cors;
}
