CREATE TABLE IF NOT EXISTS visits (
    id serial PRIMARY KEY,
    link_id integer NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address INET NOT NULL,
    referer TEXT DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_visits_link FOREIGN KEY (link_id) REFERENCES "links" (id)
);