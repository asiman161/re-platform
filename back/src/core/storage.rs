use std::env;
use diesel::{Connection, PgConnection, QueryResult};
use crate::core::models::User;


pub struct Storage {
    dsn: String
}

pub fn new() -> Storage {
    Storage {
        dsn: env::var("DATABASE_URL").expect("DATABASE_URL must be set")
    }
}

impl Storage {
    pub fn get_users(self) -> QueryResult<Vec<User>> {
        use crate::schema::users::dsl::users;
        use diesel::prelude::*;

        users.load::<User>(&mut self.establish_connection())
    }

    fn establish_connection(self) -> PgConnection {
        PgConnection::establish(&self.dsn)
            .expect(&format!("Error connecting to {}", self.dsn))
    }
}
