use std::env;

use diesel::{Connection, PgConnection, QueryResult, r2d2};

use crate::core::models::User;
use diesel::prelude::*;

use r2d2_redis::{r2d2 as redis_r2d2, RedisConnectionManager};

pub type DbPgPool = r2d2::Pool<r2d2::ConnectionManager<PgConnection>>;

pub struct Storage {
    database_dsn: String,
}

pub fn new() -> Storage {
    Storage {
        database_dsn: env::var("DATABASE_DSN").expect("DATABASE_DSN must be set"),
    }
}

type DbError = Box<dyn std::error::Error + Send + Sync>;

// pub fn get_users2(con: &mut PgConnection) -> QueryResult<Vec<User>> {
pub fn get_users(con: &mut PgConnection) -> Result<Vec<User>, DbError> {
    use crate::schema::users::dsl::users;

    Ok(users.load::<User>(con)?)
}

/// Run query using Diesel to find user by uid and return it.
pub fn find_user_by_id(
    conn: &mut PgConnection,
    user_id: i32,
) -> Result<Option<User>, DbError> {
    use crate::schema::users::dsl::*;

    let user = users
        .filter(id.eq(user_id))
        .first::<User>(conn)
        .optional()?;

    Ok(user)
}

impl Storage {
    pub fn get_users(self) -> QueryResult<Vec<User>> {
        use crate::schema::users::dsl::users;
        use diesel::prelude::*;

        users.load::<User>(&mut self.establish_connection())
    }

    fn establish_connection(self) -> PgConnection {
        PgConnection::establish(&self.database_dsn)
            .expect(&format!("Error connecting to {}", self.database_dsn))
    }
}

pub fn init_pg_pool() -> DbPgPool {
    let dsn = env::var("DATABASE_DSN").expect("DATABASE_DSN must be set");
    let manager = r2d2::ConnectionManager::<PgConnection>::new(dsn);
    r2d2::Pool::builder().build(manager).expect("database URL should be valid dsn to db")
}

pub fn init_rd_pool(dsn: &str) -> redis_r2d2::Pool<RedisConnectionManager> {
    let manager = RedisConnectionManager::new(dsn).expect("Invalid connection URL");
    redis_r2d2::Pool::builder()
        .build(manager)
        .expect("failed to connect to Redis")
}
