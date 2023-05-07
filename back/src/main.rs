use std::{env};
use actix_web::{App, HttpServer, web};
use dotenv::dotenv;
use back::core::{server, storage};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    let port = env::var("PORT")
        .expect("PORT must be set")
        .parse::<u16>().expect("PORT must be a positive integer");
    let host = "0.0.0.0";

    fetch_redis();




    let dsn = env::var("REDIS_DSN").expect("REDIS_DSN must be set");

    println!("start server at {host}:{port}");


    let pg_pool = storage::init_pg_pool();

    use r2d2_redis::{r2d2 as redis_r2d2, RedisConnectionManager};
    let manager = RedisConnectionManager::new(dsn.as_str()).unwrap();

    let rd_pool = redis_r2d2::Pool::builder()
        .build(manager)
        .unwrap();

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pg_pool.clone()))
            .app_data(web::Data::new(rd_pool.clone()))
            .wrap(server::cors())
            .configure(server::config)
    })
        .bind((host, port))?
        .run()
        .await
}

fn connect_redis(dsn: &str) -> redis::Connection {
    let con = redis::Client::open(dsn).expect("Invalid connection URL")
        .get_connection()
        .expect("failed to connect to Redis");

    return con
}

fn fetch_redis()  {
    // connect to redis
    let dsn = env::var("REDIS_DSN").expect("REDIS_DSN must be set");
    let mut con = connect_redis(dsn.as_str());

    let bar: String = redis::cmd("SET")
        .arg("foo")
        .arg("bar")
        .query(&mut con)
        .expect("failed to execute GET for 'foo'");
    println!("value for 'foo' = {}", bar);
}
