package database

type DBConnection struct {
}

func Connect() *DBConnection {
	// cfg := config.Config()
	// conn, err := pgx.NewConnPool(cfg)
	// if err != nil {
	// 	logging.Fatal(fmt.Sprintf("Unable to connect to database: %v\n", err))
	// }
	// return conn

	return nil
}

func (dbConn *DBConnection) Close() {

}
