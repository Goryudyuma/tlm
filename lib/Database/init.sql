CREATE TABLE IF NOT EXISTS account (
	id INTEGER NOT NULL AUTO_INCREMENT,
	parent INTEGER NOT NULL,
	token VARCHAR(255),
	accesstoken VARCHAR(255),
	accesstokensecret VARCHAR(255),
	lastlogin DATETIME,
	createdAt DATETIME,
	PRIMARY KEY(id),
	INDEX(parent)
);

CREATE TABLE IF NOT EXISTS query (
	id INTEGER NOT NULL,
	accountid INTEGER NOT NULL,
	query LONGTEXT NOT NULL,
	failcount INTEGER NOT NULL DEFAULT 0,
	INDEX(failcount),
	FOREIGN KEY (accountid) REFERENCES account(id)
);