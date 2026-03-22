-- users
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY UNIQUE,
    nickname TEXT NOT NULL UNIQUE,
    birthday TEXT NOT NULL,
    gender TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    profile_image TEXT DEFAULT NULL,
    session_id TEXT UNIQUE,
    session_created_at TEXT DEFAULT NULL,
    session_expired_at TEXT DEFAULT NULL
);