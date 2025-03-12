package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SupplyChain struct {
	contractapi.Contract
}

type Product struct {
	ProductID      string            `json:"ProductID"`
	CurrentStatus  string            `json:"CurrentStatus"`
	StatusHistory  map[string]string `json:"StatusHistory"`
}

func (s *SupplyChain) RegisterProduct(ctx contractapi.TransactionContextInterface, productID string, status string) error {
	exists, err := s.ProductExists(ctx, productID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("product %s already exists", productID)
	}
	timestamp := time.Now().Format(time.RFC3339)
	product := Product{
		ProductID:     productID,
		CurrentStatus: status,
		StatusHistory: map[string]string{timestamp: status},
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(productID, productJSON)
}

func (s *SupplyChain) UpdateStatus(ctx contractapi.TransactionContextInterface, productID string, newStatus string) error {
	product, err := s.QueryProduct(ctx, productID)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	product.CurrentStatus = newStatus
	product.StatusHistory[timestamp] = newStatus

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(productID, productJSON)
}

func (s *SupplyChain) QueryProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from ledger: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("product %s does not exist", productID)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *SupplyChain) ListAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (s *SupplyChain) ProductExists(ctx contractapi.TransactionContextInterface, productID string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return false, fmt.Errorf("failed to read from ledger: %v", err)
	}

	return productJSON != nil, nil
}
