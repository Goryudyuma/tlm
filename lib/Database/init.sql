CREATE TABLE IF NOT EXISTS account (
	id INTEGER NOT NULL,
	parent INTEGER NOT NULL,
	token VARCHAR(255),
	lastlogin DATETIME,
	createdAt DATETIME,
	PRIMARY KEY(id),
	INDEX(parent)
);

CREATE TABLE IF NOT EXISTS query (
	id INTEGER NOT NULL,
	accountid INTEGER NOT NULL,
	query LONGTEXT,
	failcount INTEGER,
	INDEX(failcount),
	FOREIGN KEY (accountid) REFERENCES account(id)
);