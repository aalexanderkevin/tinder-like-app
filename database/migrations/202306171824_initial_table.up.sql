CREATE TABLE users (
	id VARCHAR (255) PRIMARY KEY,
	name VARCHAR (200) NOT NULL,
	email VARCHAR (200) NOT NULL,
  	password_salt varchar(200) NOT NULL,
  	password varchar(200) NOT NULL,
  	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
