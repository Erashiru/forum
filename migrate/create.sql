CREATE TABLE IF NOT EXISTS posts (
    post_id  INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    user_id INTEGER
);

CREATE TABLE IF NOT EXISTS users (
    user_id  INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE CONSTRAINT users_uc_username,
    email TEXT NOT NULL UNIQUE CONSTRAINT users_uc_email,
    hash_password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		comment TEXT NOT NULL,
		created_at DATE NOT NULL
	);

CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER,
		session_token TEXT,
		expires_at TIME
	);	

CREATE TABLE IF NOT EXISTS reactions (
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		reaction_status INTEGER DEFAULT 0 NOT NULL
	);

CREATE TABLE IF NOT EXISTS comment_reactions (
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		reaction_status INTEGER DEFAULT 0 NOT NULL
	);

CREATE TABLE IF NOT EXISTS categories (
		post_id INTEGER NOT NULL,
		name TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS tags (
		tags_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS requests (
		request_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL UNIQUE,
		created_at DATE NOT NULL
	);

CREATE TABLE IF NOT EXISTS reports (
		report_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		report_text TEXT NOT NULL,
		post_id INTEGER NOT NULL,
		created_at DATE NOT NULL
	);

CREATE TABLE IF NOT EXISTS roles (
		role_id INTEGER PRIMARY KEY AUTOINCREMENT,
		role TEXT NOT NULL,
		user_id INTEGER NOT NULL
	);