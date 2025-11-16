package database

import (
	"backend/internal/models"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProductDetails operations

func (db *DB) InsertProduct(name string, ccal, fats, proteins, carbs int64) (*models.ProductDetails, error) {
	product := &models.ProductDetails{}
	
	query := `
		INSERT INTO product_details (name, ccal, fats, proteins, carbs)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uuid, name, ccal, fats, proteins, carbs
	`
	
	err := db.QueryRow(query, name, ccal, fats, proteins, carbs).Scan(
		&product.UUID,
		&product.Name,
		&product.Ccal,
		&product.Fats,
		&product.Proteins,
		&product.Carbs,
	)
	
	if err != nil {
		return nil, err
	}
	
	return product, nil
}

func (db *DB) GetProductByUUID(productUUID uuid.UUID) (*models.ProductDetails, error) {
	product := &models.ProductDetails{}
	
	query := `
		SELECT uuid, name, ccal, fats, proteins, carbs
		FROM product_details
		WHERE uuid = $1
	`
	
	err := db.QueryRow(query, productUUID).Scan(
		&product.UUID,
		&product.Name,
		&product.Ccal,
		&product.Fats,
		&product.Proteins,
		&product.Carbs,
	)
	
	if err != nil {
		return nil, err
	}
	
	return product, nil
}

// Record operations

func (db *DB) InsertRecord(productUUID uuid.UUID, amount int64, login string) (*models.Record, error) {
	record := &models.Record{}
	
	query := `
		INSERT INTO records (product_uuid, amount, login)
		VALUES ($1, $2, $3)
		RETURNING uuid, product_uuid, amount, login, created_at
	`
	
	err := db.QueryRow(query, productUUID, amount, login).Scan(
		&record.UUID,
		&record.ProductUUID,
		&record.Amount,
		&record.Login,
		&record.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

func (db *DB) GetRecordByUUID(recordUUID uuid.UUID) (*models.Record, error) {
	record := &models.Record{}
	
	query := `
		SELECT uuid, product_uuid, amount, login, created_at
		FROM records
		WHERE uuid = $1
	`
	
	err := db.QueryRow(query, recordUUID).Scan(
		&record.UUID,
		&record.ProductUUID,
		&record.Amount,
		&record.Login,
		&record.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

func (db *DB) GetRecordsByLoginAndTimeRange(login string, startTime, endTime time.Time) ([]*models.Record, error) {
	query := `
		SELECT uuid, product_uuid, amount, login, created_at
		FROM records
		WHERE login = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
	`
	
	rows, err := db.Query(query, login, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*models.Record
	for rows.Next() {
		record := &models.Record{}
		err := rows.Scan(
			&record.UUID,
			&record.ProductUUID,
			&record.Amount,
			&record.Login,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return records, nil
}

// UserPreferences operations

func (db *DB) GetUserPreferences(login string) (*models.UserPreferences, error) {
	prefs := &models.UserPreferences{}
	
	query := `
		SELECT login, noon, lang
		FROM user_preferences
		WHERE login = $1
	`
	
	err := db.QueryRow(query, login).Scan(
		&prefs.Login,
		&prefs.Noon,
		&prefs.Lang,
	)
	
	if err != nil {
		return nil, err
	}
	
	return prefs, nil
}

func (db *DB) UpsertUserNoon(login string, noon time.Time) error {
	query := `
		INSERT INTO user_preferences (login, noon)
		VALUES ($1, $2)
		ON CONFLICT (login)
		DO UPDATE SET noon = EXCLUDED.noon
	`
	
	_, err := db.Exec(query, login, noon)
	return err
}

func (db *DB) UpsertUserLang(login string, lang string) error {
	query := `
		INSERT INTO user_preferences (login, lang)
		VALUES ($1, $2)
		ON CONFLICT (login)
		DO UPDATE SET lang = EXCLUDED.lang
	`
	
	_, err := db.Exec(query, login, lang)
	return err
}

// UserCommonItem operations

func (db *DB) InsertUserCommonItem(login, path, name string, productUUID *uuid.UUID) (*models.UserCommonItem, error) {
	item := &models.UserCommonItem{}
	
	query := `
		INSERT INTO user_common_items (login, path, name, product_uuid)
		VALUES ($1, $2::ltree, $3, $4)
		RETURNING uuid, login, path, name, product_uuid, created_at
	`
	
	err := db.QueryRow(query, login, path, name, productUUID).Scan(
		&item.UUID,
		&item.Login,
		&item.Path,
		&item.Name,
		&item.ProductUUID,
		&item.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return item, nil
}

func (db *DB) GetUserCommonItemsByLogin(login string) ([]*models.UserCommonItem, error) {
	query := `
		SELECT uuid, login, path, name, product_uuid, created_at
		FROM user_common_items
		WHERE login = $1
		ORDER BY path
	`
	
	rows, err := db.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []*models.UserCommonItem
	for rows.Next() {
		item := &models.UserCommonItem{}
		err := rows.Scan(
			&item.UUID,
			&item.Login,
			&item.Path,
			&item.Name,
			&item.ProductUUID,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return items, nil
}

func (db *DB) GetUserCommonItemsByLoginAndPath(login, pathPattern string) ([]*models.UserCommonItem, error) {
	query := `
		SELECT uuid, login, path, name, product_uuid, created_at
		FROM user_common_items
		WHERE login = $1 AND path ~ $2
		ORDER BY path
	`
	
	rows, err := db.Query(query, login, pathPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []*models.UserCommonItem
	for rows.Next() {
		item := &models.UserCommonItem{}
		err := rows.Scan(
			&item.UUID,
			&item.Login,
			&item.Path,
			&item.Name,
			&item.ProductUUID,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return items, nil
}

func (db *DB) GetUserCommonItemsAtLevel(login, parentPath string) ([]*models.UserCommonItem, error) {
	var query string
	var rows *sql.Rows
	var err error
	
	if parentPath == "" {
		query = `
			SELECT uuid, login, path, name, product_uuid, created_at
			FROM user_common_items
			WHERE login = $1 AND nlevel(path) = 1
			ORDER BY path
		`
		rows, err = db.Query(query, login)
	} else {
		query = `
			SELECT uuid, login, path, name, product_uuid, created_at
			FROM user_common_items
			WHERE login = $1 AND path ~ $2 AND nlevel(path) = $3
			ORDER BY path
		`
		pattern := parentPath + ".*{1}"
		level := strings.Count(parentPath, ".") + 2
		rows, err = db.Query(query, login, pattern, level)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []*models.UserCommonItem
	for rows.Next() {
		item := &models.UserCommonItem{}
		err := rows.Scan(
			&item.UUID,
			&item.Login,
			&item.Path,
			&item.Name,
			&item.ProductUUID,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return items, nil
}
