use actix_web::{App, HttpServer};

mod core;

#[actix_web::main]
async fn main() -> std::io::Result<()> {

    let port = 8080;
    let host = "0.0.0.0";
    println!("start server at {host}:{port}");

    HttpServer::new(|| { App::new().configure(core::server::config) })
        .bind((host, port))?
        .run()
        .await
}
