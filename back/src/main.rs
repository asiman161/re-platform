use std::env;
use actix_web::{App, HttpServer};
use dotenv::dotenv;
use back::core::server;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    let port = env::var("PORT")
        .expect("PORT must be set")
        .parse::<u16>().expect("PORT must be a positive integer");
    let host = "0.0.0.0";
    println!("start server at {host}:{port}");

    HttpServer::new(|| {
        App::new()
            .wrap(server::cors())
            .configure(server::config)
    })
        .bind((host, port))?
        .run()
        .await
}
