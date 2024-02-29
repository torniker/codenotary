package accounting

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/k0kubun/pp/v3"

	"github.com/ddosify/go-faker/faker"
	"github.com/spf13/cobra"
	"github.com/torniker/codenotary/immudb"
	"github.com/torniker/codenotary/model"
)

const baseURL = "http://0.0.0.0:5656"

func Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create-accounting-collection",
		Short: "Create Collection",
		Run:   createAccountingCollection,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get-accounting-collection",
		Short: "Get Collection",
		Run:   getAccountingCollection,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "drop-accounting-collection",
		Short: "Drop Collection",
		Run:   dropAccountingCollection,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "add-accounting-document",
		Short: "Add Random Accounting Document",
		Run:   addAccountingDocument,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "replace-accounting-document",
		Short: "Replace Accounting Document by ID",
		Run:   replaceAccountingDocument,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "search-accounting-document",
		Short: "Search Accounting Documents",
		Run:   searchAccountingDocument,
	})
}

func searchAccountingDocument(cmd *cobra.Command, args []string) {
	resp, err := http.Get(fmt.Sprintf("%s/api/accounting", baseURL))
	if err != nil {
		log.Fatal(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(body))
}

func addAccountingDocument(cmd *cobra.Command, args []string) {
	faker := faker.NewFaker()
	index := faker.IntBetween(0, 1)
	types := []string{"sending", "receiving"}
	accounting := model.Accounting{
		Number:  faker.RandomBankAccount(),
		Name:    faker.RandomCompanyName(),
		IBAN:    faker.RandomBankAccountIban(),
		Address: faker.RandomAddressStreetAddress(),
		Amount:  faker.RandomInt(),
		Type:    types[index],
	}
	err := immudb.New().CreateDocument(context.Background(), "accounting", accounting)
	if err != nil {
		log.Fatal(err)
		return
	}
	pp.Print(accounting)
}

func replaceAccountingDocument(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("ID is required")
		return
	}
	faker := faker.NewFaker()
	index := faker.IntBetween(0, 1)
	types := []string{"sending", "receiving"}
	accounting := model.Accounting{
		Number:  faker.RandomBankAccount(),
		Name:    faker.RandomCompanyName(),
		IBAN:    faker.RandomBankAccountIban(),
		Address: faker.RandomAddressStreetAddress(),
		Amount:  faker.RandomInt(),
		Type:    types[index],
	}
	query := immudb.Query{
		Expressions: []immudb.Expression{
			{
				FieldComparisons: []immudb.FieldComparison{
					{Field: "id", Operator: immudb.OperatorEqual, Value: args[0]},
				},
				OrderBy: []immudb.OrderBy{},
			},
		},
	}
	err := immudb.New().ReplaceDocument(context.Background(), "accounting", accounting, query)
	if err != nil {
		log.Fatal(err)
		return
	}
	accounting.ID = args[0]
	pp.Print(accounting)
}

func createAccountingCollection(cmd *cobra.Command, args []string) {
	collection := immudb.Collection{
		Name:        "accounting",
		IDFieldName: "id",
		Fields: []immudb.Field{
			{Name: "account_number", Type: "STRING"},
			{Name: "account_name", Type: "STRING"},
			{Name: "iban", Type: "STRING"},
			{Name: "address", Type: "STRING"},
			{Name: "amount", Type: "INTEGER"},
			{Name: "type", Type: "STRING"},
		},
		Indexes: []immudb.Index{
			{
				Fields: []string{
					"account_number",
				},
				IsUnique: true,
			},
		},
	}
	err := immudb.New().CollectionCreate(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}
	pp.Print(collection)
}

func getAccountingCollection(cmd *cobra.Command, args []string) {
	collection, err := immudb.New().Collection(context.Background(), "accounting")
	if err != nil {
		log.Fatal(err)
	}
	pp.Print(collection)
}

func dropAccountingCollection(cmd *cobra.Command, args []string) {
	err := immudb.New().DropCollection(context.Background(), "accounting")
	if err != nil {
		log.Fatal(err)
		return
	}
	pp.Print("accounting deleted successfully")
}
