package db

import (
    "database/sql"
    "fmt"
    "log"

    "vn.ghrm/internal/config"
    _ "github.com/lib/pq"
)

func SetupDatabase(cfg *config.Config) error {
    // Connect to PostgreSQL as the admin user
    adminDSN := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
        cfg.DbHost, cfg.DbPort, cfg.DbAdminUser, cfg.DbAdminDB)
    log.Printf("Admin DSN: %s", adminDSN)

    db, err := sql.Open("postgres", adminDSN)
    if err != nil {
        return fmt.Errorf("failed to connect to admin database: %v", err)
    }
    defer db.Close()

    // Test the connection
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping admin database: %v", err)
    }

    // Check if the hrm database already exists
    var exists bool
    query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
    err = db.QueryRow(query, cfg.DbName).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check if database exists: %v", err)
    }

    if exists {
        log.Printf("Database %s already exists, skipping creation", cfg.DbName)
    } else {
        // Create the hrm database if it doesn't exist
        _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DbName))
        if err != nil {
            return fmt.Errorf("failed to create database: %v", err)
        }
        log.Printf("Created database %s", cfg.DbName)
    }

    // Create the hrm user if it doesn't exist
    _, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", cfg.DbUser, cfg.DbPassword))
    if err != nil {
        if err.Error() != fmt.Sprintf(`pq: role "%s" already exists`, cfg.DbUser) {
            return fmt.Errorf("failed to create user: %v", err)
        }
        log.Printf("User %s already exists", cfg.DbUser)
    } else {
        log.Printf("Created user %s", cfg.DbUser)
    }

    // Grant privileges on the database to the hrm user
    _, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", cfg.DbName, cfg.DbUser))
    if err != nil {
        return fmt.Errorf("failed to grant database privileges: %v", err)
    }
    log.Printf("Granted database privileges to user %s on database %s", cfg.DbUser, cfg.DbName)

    // Connect to the hrm database as the admin user to grant schema privileges
    hrmDSN := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
        cfg.DbHost, cfg.DbPort, cfg.DbAdminUser, cfg.DbName)
    hrmDB, err := sql.Open("postgres", hrmDSN)
    if err != nil {
        return fmt.Errorf("failed to connect to hrm database: %v", err)
    }
    defer hrmDB.Close()

    // Grant USAGE and CREATE privileges on the public schema to the hrm user
    _, err = hrmDB.Exec(fmt.Sprintf("GRANT USAGE, CREATE ON SCHEMA public TO %s", cfg.DbUser))
    if err != nil {
        return fmt.Errorf("failed to grant schema privileges: %v", err)
    }
    log.Printf("Granted USAGE and CREATE privileges on public schema to user %s", cfg.DbUser)

    return nil
}