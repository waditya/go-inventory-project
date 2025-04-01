package main

import (
	"database/sql"
	"fmt"
)

// This file consits of database related methods

// We need Structure to process data returned as rows

// JSON tags are needed to help during encoding into JSON format while sending response

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {

	rows, err := db.Query("SELECT id, name, quantity, price from products")
	checkError(err)

	products := []product{}

	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)

		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil

}

// getProduct function will need product information hence we create it as a method of product struct

func (p *product) getProduct(db *sql.DB) error {

	query := fmt.Sprintf("SELECT name, quantity, price from products WHERE id=%v", p.ID)
	row := db.QueryRow(query) // Used for select query with single row output
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)

	if err != nil {
		return err
	}

	return nil

}
