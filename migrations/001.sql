CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS content;
CREATE SCHEMA IF NOT EXISTS guides;
CREATE SCHEMA IF NOT EXISTS chat;

SET search_path TO auth, content, guides, chat, public;

CREATE TYPE auth.role AS ENUM ('admin', 'worker');

CREATE TABLE IF NOT EXISTS auth.client (
    id UUID PRIMARY KEY,
    trade_union_id text,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth."user" (
    id UUID PRIMARY KEY,
    role auth.role NOT NULL DEFAULT 'worker',
    login TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth.user_info (
    user_id UUID PRIMARY KEY,
    first_name TEXT,
    second_name TEXT,
    patronymic TEXT DEFAULT NULL,
    image_url TEXT DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth."user"(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS auth.invite_token (
    id SERIAL PRIMARY KEY,
    content UUID NOT NULL,
    used BOOLEAN DEFAULT false,
    role auth.role NOT NULL
);

CREATE TABLE IF NOT EXISTS content.news (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS content.projects (
    id UUID PRIMARY KEY,
    intro TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS content.documents (
    id UUID PRIMARY KEY,
    url TEXT,
    title TEXT,
    type TEXT
);


CREATE TYPE guides.guide_type AS ENUM ('worker', 'student');

CREATE TABLE IF NOT EXISTS guides.guides (
    id SERIAL PRIMARY KEY,
    type guides.guide_type NOT NULL,
    title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS guides.themes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    guide_id INTEGER REFERENCES guides.guides(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat.chat (
    id UUID PRIMARY KEY,
    title TEXT
);

CREATE TABLE IF NOT EXISTS chat.chat_users (
    chat_id UUID,
    user_id UUID,
    UNIQUE(chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES chat.chat(id)
);

CREATE TABLE IF NOT EXISTS chat.messages (
    id UUID PRIMARY KEY,
    content TEXT,
    user_id UUID,
    chat_id UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (chat_id) REFERENCES chat.chat(id)
);
