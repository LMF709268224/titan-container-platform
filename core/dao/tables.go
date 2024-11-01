package dao

// Table creation SQL statements.
var cUserTable = `
    CREATE TABLE if not exists %s (
		account             VARCHAR(128)  NOT NULL UNIQUE,
		avatar              VARCHAR(128)  DEFAULT '',
		user_name           VARCHAR(128)  DEFAULT '',
		user_email          VARCHAR(128)  DEFAULT '',
		created_at          DATETIME      DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (account)
	) ENGINE=InnoDB COMMENT='user info';`

var cOrderTable = `
    CREATE TABLE if not exists %s (
		id           VARCHAR(128)  NOT NULL UNIQUE,
		account      VARCHAR(255)  NOT NULL,		
		cpu          INT           DEFAULT 0,
		ram       INT           DEFAULT 0,
		storage      INT           DEFAULT 0,
		duration       INT           DEFAULT 0,
		status       INT           DEFAULT 0,
		created_at   DATETIME      DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		KEY idx_account (account),
		KEY idx_status (status)
	) ENGINE=InnoDB COMMENT='order info';`
