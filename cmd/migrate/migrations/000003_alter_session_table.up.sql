ALTER TABLE sessions 
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT unique_user_session UNIQUE (user_id);