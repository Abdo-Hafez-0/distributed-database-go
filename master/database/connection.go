package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql" 
    "github.com/joho/godotenv"
)

// DB holds the global database connection instance.
var DB *sql.DB

// Config holds database connection parameters loaded from environment.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LoadConfig reads database configuration from environment variables.
// It attempts to load a .env file first; missing variables are fatal.
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[database] No .env file found — using system environment variables")
	}

	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Name == "" {
		log.Fatal("[database] Missing required environment variables: DB_HOST, DB_PORT, DB_USER, DB_NAME")
	}

	return cfg
}

// DSN builds the MySQL Data Source Name string from Config.
func (c Config) DSN() string {
	// Format: user:password@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

// Connect initialises the global DB connection pool using environment config.
// Call this once at application startup.
func Connect() (*sql.DB, error) {
	cfg := LoadConfig()
	dsn := cfg.DSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("[database] sql.Open failed: %w", err)
	}

	// Verify the connection is actually reachable.
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("[database] ping failed (check credentials / host): %w", err)
	}

	// Sensible connection-pool defaults — tune per workload.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	DB = db
	log.Printf("[database] Connected to MySQL at %s:%s/%s", cfg.Host, cfg.Port, cfg.Name)
	return db, nil
}

// MustConnect is like Connect but calls log.Fatal on error.
// Convenient for main() where a DB connection is non-negotiable.
func MustConnect() *sql.DB {
	db, err := Connect()
	if err != nil {
		log.Fatalf("[database] Could not connect: %v", err)
	}
	return db
}

// Close gracefully closes the global DB connection pool.
func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("[database] Error closing connection: %v", err)
		}
	}
}