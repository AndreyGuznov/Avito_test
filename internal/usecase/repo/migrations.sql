CREATE TABLE users (
	id int not NULL,
	balance float,
	primary key (id)
);

CREATE TABLE transanctions (
    user_id INT not NULL,
	amount float NOT null,
	operation text not null,
    date TIMESTAMP (3) not null,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);