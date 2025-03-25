package apiserver

import (
	"database/sql"
	"net/http"

	"obscura.app/backend/pkg/logger"
	// Add your store package here
)

// Start initializes and starts the API server
func Start(cfg *Config) error {
	log, err := logger.NewLogger("cmd/app.log")

	if err != nil {
		return err
	}
	defer log.Close()

	db, err := newDB(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		return err
	}
	defer db.Close()

	// Initialize your store here
	// store := sqlstore.New(db)

	// Initialize your session store here
	// sessionStore := sessions.NewCookieStore([]byte(cfg.SessionKey))

	// Initialize your server here
	// srv := newServer(store, sessionStore)

	log.Info("Starting server on port: " + cfg.ServerConfig.Port)
	return http.ListenAndServe(cfg.Server.Port, nil) // Replace nil with your server
}

func newDB(cfg DatabaseConfig) (*sql.DB, error) {
	dsn := "" // Construct your DSN here using cfg
	db, err := sql.Open("postgre", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}