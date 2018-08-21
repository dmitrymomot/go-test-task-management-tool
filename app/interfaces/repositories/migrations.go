package repositories

// Migrate runs database migrations
func Migrate(db DbHandler) error {
	q := `CREATE TABLE dbname.tasks(
	    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	    title VARCHAR(250) NULL DEFAULT NULL,
	    description TEXT(1000) DEFAULT NULL,
	    status VARCHAR(20) NOT NULL,
		created_at DATETIME NOT NULL,
		completed_at DATETIME DEFAULT NULL,
	    PRIMARY KEY(id),
	    INDEX(status)
	) ENGINE = InnoDB CHARSET = utf8 COLLATE utf8_general_ci;`

	_, err := db.Execute(q)
	return err
}
