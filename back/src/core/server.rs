use actix_web::{web};
use crate::core::handlers;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api")
            .service(handlers::ping)
            .service(handlers::users)
    );
}
