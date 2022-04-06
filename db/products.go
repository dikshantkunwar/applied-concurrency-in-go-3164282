package db

import (
	"fmt"
	"sort"
	"sync"

	"github.com/applied-concurrency-in-go/models"
	"github.com/applied-concurrency-in-go/utils"
)

type ProductDB struct {
	products sync.Map
}

// NewProducts creates a new empty product DB
func NewProducts() (*ProductDB, error) {
	p := &ProductDB{}
	// load start position
	if err := utils.ImportProducts(&p.products); err != nil {
		return nil, err
	}
	return p, nil
}

// Exists checks whether a product with a give id exists
func (p *ProductDB) Exists(id string) error {
	_, ok := p.products.Load(id)
	if !ok {
		return fmt.Errorf("no product found for %s product id", id)
	}
	return nil
}

// Find returns a given product if one exists
func (p *ProductDB) Find(id string) (models.Product, error) {
	prod, ok := p.products.Load(id)
	if !ok {
		return models.Product{}, fmt.Errorf("no product found for id %s", id)
	}

	return toProduct(prod), nil
}

// Upsert creates or updates a product in the orders DB
func (p *ProductDB) Upsert(prod models.Product) {
	p.products.Store(prod.ID, prod)
}

// FindAll returns all products in the system
func (p *ProductDB) FindAll() []models.Product {
	var allProducts []models.Product

	p.products.Range(func(_, value interface{}) bool {
		allProducts = append(allProducts, toProduct(value))
		return true
	})

	sort.Slice(allProducts, func(i, j int) bool {
		return allProducts[i].ID < allProducts[j].ID
	})
	return allProducts
}

func toProduct(prod interface{}) models.Product {
	product, ok := prod.(models.Product)
	if !ok {
		panic(fmt.Errorf("error casting %v to order", prod))
	}
	return product
}
