package database

import (
	"backend/internal/models"
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
