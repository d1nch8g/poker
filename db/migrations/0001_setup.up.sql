CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  verified BOOLEAN NOT NULL,
  is_admin BOOLEAN NOT NULL,
  encrypted_wallet TEXT UNIQUE NOT NULL,
  passwhash TEXT UNIQUE NOT NULL
);