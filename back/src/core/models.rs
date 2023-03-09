use chrono::NaiveDateTime;
use diesel::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Debug, Queryable, Serialize, Deserialize)]
pub struct User {
    pub id: i32,
    pub username: String,
    pub password_hash: String,
    pub password_salt: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}
