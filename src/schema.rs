// @generated automatically by Diesel CLI.

diesel::table! {
    users (id) {
        id -> Int4,
        username -> Text,
        pass_hash -> Text,
        salt -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}
