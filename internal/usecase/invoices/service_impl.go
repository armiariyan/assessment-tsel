package invoices

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"

	"context"

	"github.com/armiariyan/bepkg/database/mysql"
)

type service struct {
	db                       *mysql.SQLDB
	customersRepository      repositories.CustomersRepository
	itemsRepository          repositories.ItemsRepository
	invoicesRepository       repositories.InvoicesRepository
	invoiceItemsRepository   repositories.InvoiceItemsRepository
	invoiceSummaryRepository repositories.InvoiceSummaryRepository
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *mysql.SQLDB) *service {
	s.db = db
	return s
}

func (s *service) SetCustomersRepository(repo repositories.CustomersRepository) *service {
	s.customersRepository = repo
	return s
}

func (s *service) SetItemsRepository(repo repositories.ItemsRepository) *service {
	s.itemsRepository = repo
	return s
}

func (s *service) SetInvoicesRepository(repo repositories.InvoicesRepository) *service {
	s.invoicesRepository = repo
	return s
}

func (s *service) SetInvoiceItemsRepository(repo repositories.InvoiceItemsRepository) *service {
	s.invoiceItemsRepository = repo
	return s
}

func (s *service) SetInvoiceSummaryRepository(repo repositories.InvoiceSummaryRepository) *service {
	s.invoiceSummaryRepository = repo
	return s
}

func (s *service) Validate() InvoicesService {
	if s.db == nil {
		panic("db is nil")
	}

	if s.customersRepository == nil {
		panic("customersRepository is nil")
	}

	if s.itemsRepository == nil {
		panic("itemsRepository is nil")
	}

	if s.invoicesRepository == nil {
		panic("invoicesRepository is nil")
	}

	if s.invoiceItemsRepository == nil {
		panic("invoiceItemsRepository is nil")
	}

	if s.invoiceSummaryRepository == nil {
		panic("invoiceSummaryRepository is nil")
	}

	return s
}

func (s *service) GetListInvoices(ctx context.Context, req GetListInvoicesRequest) (resp constants.DefaultResponse, err error) {
	queryParams := entities.InvoiceListParams{
		Limit:     req.Limit,
		Page:      req.Page,
		InvoiceID: req.InvoiceID,
		Subject:   req.Subject,
		Customer:  req.Customer,
		Status:    req.Status,
		IssueDate: req.IssueDate,
		DueDate:   req.DueDate,
	}

	invoices, count, err := s.invoicesRepository.FindAllAndCountWithCondition(ctx, queryParams)
	if err != nil {
		log.Error(ctx, "failed find all and count invoices with condition", queryParams, err)
		err = fmt.Errorf("something went wrong [0]")
		return
	}

	var results []GetListInvoicesResponse
	for _, v := range invoices {
		status := "unpaid"
		if v.IsPaid {
			status = "paid"
		}

		results = append(results, GetListInvoicesResponse{
			InvoiceID:    v.InvoiceID,
			Subject:      v.Subject,
			CustomerName: v.Name,
			TotalItems:   v.TotalItems,
			IssueDate:    v.IssueDate,
			DueDate:      v.DueDate,
			Status:       status,
		})

	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: results,
			PaginationData: constants.PaginationData{
				Page:        req.Page,
				Limit:       req.Limit,
				TotalPages:  uint(math.Ceil(float64(count) / float64(req.Limit))),
				TotalItems:  uint(count),
				HasNext:     req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
				HasPrevious: req.Page > 1,
			},
		},
		Errors: make([]string, 0),
	}

	return
}

func (s *service) GetDetailInvoice(ctx context.Context, uniqueInvoiceID string) (resp constants.DefaultResponse, err error) {
	// * INVOICE
	invoice, err := s.invoicesRepository.FindByUniqueInvoiceID(ctx, uniqueInvoiceID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find invoice by id %s", uniqueInvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data invoice not found")
			return
		}

		err = fmt.Errorf("something went wrong [0]")
		return
	}

	// * INVOICE SUMMARY
	invoiceSummary, err := s.invoiceSummaryRepository.FindByInvoiceID(ctx, invoice.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find invoice_summary by invoice id with unique invoice id %s", uniqueInvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data invoice summary not found")
			return
		}

		err = fmt.Errorf("something went wrong [1]")
		return
	}

	// * INVOICE ITEMS
	invoiceItems, err := s.invoiceItemsRepository.FindByInvoiceID(ctx, invoice.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find invoice_items by invoice id with unique invoice id %s", uniqueInvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data invoice items not found")
			return
		}

		err = fmt.Errorf("something went wrong [0]")
		return
	}
	var itemIDs []int
	for _, v := range invoiceItems {
		itemIDs = append(itemIDs, v.ItemID)
	}
	// * ITEMS
	items, err := s.itemsRepository.FindByIDs(ctx, itemIDs)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find items by invoice id with unique invoice id %s", uniqueInvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data items not found")
			return
		}

		err = fmt.Errorf("something went wrong [1]")
		return
	}

	mapItemsIdAndName := make(map[int]string)
	for _, item := range items {
		mapItemsIdAndName[item.ID] = item.Name
	}

	// * CUSTOMER
	customer, err := s.customersRepository.FindByID(ctx, invoice.CustomerID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find customer with customer id %d when get detail invoice", invoice.CustomerID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data customer not found")
			return
		}

		err = fmt.Errorf("something went wrong [1]")
		return
	}

	var responseItems []Item
	for _, v := range invoiceItems {
		responseItems = append(responseItems, Item{
			Name:      mapItemsIdAndName[v.ItemID],
			Quantity:  v.Quantity,
			UnitPrice: v.UnitPrice,
			Amount:    v.Amount,
		})
	}

	result := GetDetailInvoiceResponse{
		InvoiceID: invoice.InvoiceID,
		Subject:   invoice.Subject,
		IssueDate: invoice.IssueDate,
		DueDate:   invoice.DueDate,
		Customer: Customer{
			Name:     customer.Name,
			Address:  customer.Address,
			City:     customer.City,
			Country:  customer.Country,
			Postcode: customer.Postcode,
		},
		Items: responseItems,
		Summary: Summary{
			TotalItems: invoiceSummary.TotalItems,
			SubTotal:   invoiceSummary.Subtotal,
			Tax:        invoiceSummary.Tax,
			GrandTotal: invoiceSummary.GrandTotal,
			Status:     "unpaid",
		},
	}

	if invoiceSummary.IsPaid {
		result.Summary.Status = "paid"
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    result,
		Errors:  make([]string, 0),
	}

	return
}

func (s *service) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (resp constants.DefaultResponse, err error) {
	// * validate given customer ID exists
	customer, err := s.customersRepository.FindByID(ctx, req.CustomerID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find customers by id %d", req.CustomerID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data customers not found")
			return
		}

		err = fmt.Errorf("something went wrong [0]")
		return
	}

	// * validate given items id exists
	var itemIDs []int
	for _, v := range req.InvoiceItems {
		itemIDs = append(itemIDs, v.ItemID)
	}

	items, err := s.itemsRepository.FindByIDs(ctx, itemIDs)
	if err != nil {
		log.Error(ctx, "failed find by IDs from items", err)
		err = fmt.Errorf("something went wrong [1]")
		return
	}

	if len(items) != len(itemIDs) {
		err = fmt.Errorf("there is invalid items")
		return
	}

	// * get last id for invoice ID format
	lastID, err := s.invoicesRepository.FindLastDataID(ctx)
	if err != nil {
		log.Error(ctx, "failed find last data ID from invoices", err)
		err = fmt.Errorf("something went wrong [2]")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(ctx, "failed to begin tx", err)
		err = fmt.Errorf("something went wrong [3]")
		tx.Rollback()
		return
	}

	var (
		issueDate, _ = time.Parse("2006-01-02", req.IssueDate) // * parse string into time.Time
		dueDate, _   = time.Parse("2006-01-02", req.DueDate)   // * parse string into time.Time
		invoiceID    = fmt.Sprintf("%04d", (lastID + 1))
	)

	// * MANAGE INVOICE
	invoices := entities.Invoice{
		InvoiceID:  invoiceID,
		IssueDate:  issueDate,
		DueDate:    dueDate,
		Subject:    req.Subject,
		CustomerID: customer.ID,
	}
	res, err := s.invoicesRepository.InsertWithTx(ctx, invoices, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed insert with tx into invoices with ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [4]")
		tx.Rollback()
		return
	}

	countRows, err := res.RowsAffected()
	if countRows < 1 || err != nil {
		log.Error(ctx, "SQL error rows affected invoices insert with tx", err)
		err = fmt.Errorf("something went wrong [5]")
		tx.Rollback()
		return
	}

	// * MANAGE INVOICE SUMMARY
	latestInvoiceID, _ := res.LastInsertId()
	invoiceSummary := entities.InvoiceSummary{
		InvoiceID:  int(latestInvoiceID),
		TotalItems: len(req.InvoiceItems),
		Subtotal:   req.InvoiceSummary.SubTotal,
		Tax:        req.InvoiceSummary.Tax,
		GrandTotal: req.InvoiceSummary.GrandTotal,
	}

	res, err = s.invoiceSummaryRepository.InsertWithTx(ctx, invoiceSummary, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed insert with tx into invoice summary with invoice ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [6]")
		tx.Rollback()
		return
	}

	countRows, err = res.RowsAffected()
	if countRows < 1 || err != nil {
		log.Error(ctx, "SQL error rows affected invoice summary insert with tx", err)
		err = fmt.Errorf("something went wrong [7]")
		tx.Rollback()
		return
	}

	// * MANAGE INVOICE ITEMS
	var invoiceItems []entities.InvoiceItems
	for _, v := range req.InvoiceItems {
		invoiceItems = append(invoiceItems, entities.InvoiceItems{
			InvoiceID: int(latestInvoiceID),
			ItemID:    v.ItemID,
			Quantity:  v.Quantity,
			UnitPrice: v.UnitPrice,
			Amount:    v.Amount,
		})
	}

	res, err = s.invoiceItemsRepository.BulkInsertWithTx(ctx, invoiceItems, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed bulk insert with tx into invoice items with invoice ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [8]")
		tx.Rollback()
		return
	}

	tx.Commit()

	countRows, err = res.RowsAffected()
	if countRows < int64(len(invoiceItems)) || err != nil {
		log.Error(ctx, "SQL error rows affected invoice summary insert with tx", err)
		err = fmt.Errorf("something went wrong [9]")
		tx.Rollback()
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    struct{}{},
		Errors:  make([]string, 0),
	}

	return
}

func (s *service) EditInvoice(ctx context.Context, req EditInvoiceRequest) (resp constants.DefaultResponse, err error) {
	// * validate given invoice ID exists
	invoice, err := s.invoicesRepository.FindByUniqueInvoiceID(ctx, req.InvoiceID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find invoice by id %s", req.InvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("invoice not found")
			return
		}

		err = fmt.Errorf("something went wrong [0]")
		return
	}

	// * validate given customer ID exists
	customer, err := s.customersRepository.FindByID(ctx, req.CustomerID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find customers by id %d", req.CustomerID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("data customers not found")
			return
		}

		err = fmt.Errorf("something went wrong [1]")
		return
	}

	// * validate given items id exists
	var itemIDs []int
	for _, v := range req.InvoiceItems {
		itemIDs = append(itemIDs, v.ItemID)
	}

	items, err := s.itemsRepository.FindByIDs(ctx, itemIDs)
	if err != nil {
		log.Error(ctx, "failed find by IDs from items", err)
		err = fmt.Errorf("something went wrong [2]")
		return
	}

	if len(items) != len(itemIDs) {
		err = fmt.Errorf("there is invalid items")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(ctx, "failed to begin tx", err)
		err = fmt.Errorf("something went wrong [3]")
		return
	}

	var (
		issueDate, _ = time.Parse("2006-01-02", req.IssueDate) // * parse string into time.Time
		dueDate, _   = time.Parse("2006-01-02", req.DueDate)   // * parse string into time.Time
	)

	// * MANAGE INVOICE
	invoices := entities.Invoice{
		ID:         invoice.ID,
		InvoiceID:  req.InvoiceID,
		IssueDate:  issueDate,
		DueDate:    dueDate,
		Subject:    req.Subject,
		CustomerID: customer.ID,
	}
	_, err = s.invoicesRepository.UpdateWithTx(ctx, invoices, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed update with tx into invoices with ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [4]")
		tx.Rollback()
		return
	}

	// * MANAGE INVOICE SUMMARY
	invoiceSummary := entities.InvoiceSummary{
		InvoiceID:  invoice.ID,
		TotalItems: len(req.InvoiceItems),
		Subtotal:   req.InvoiceSummary.SubTotal,
		Tax:        req.InvoiceSummary.Tax,
		GrandTotal: req.InvoiceSummary.GrandTotal,
	}

	_, err = s.invoiceSummaryRepository.UpdateWithTx(ctx, invoiceSummary, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed update with tx into invoice summary with invoice ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [6]")
		tx.Rollback()
		return
	}

	// * MANAGE INVOICE ITEMS
	existingItems, err := s.invoiceItemsRepository.FindByInvoiceID(ctx, invoice.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find by invoice ID from invoice items with invoice ID %s", invoices.InvoiceID), err)
		err = fmt.Errorf("something went wrong [8]")
		tx.Rollback()
		return
	}

	existingItemsMap := make(map[int]entities.InvoiceItems)
	for _, item := range existingItems {
		existingItemsMap[item.ItemID] = item
	}

	for _, v := range req.InvoiceItems {
		if existingItem, ok := existingItemsMap[v.ItemID]; ok {
			existingItem.Quantity = v.Quantity
			existingItem.UnitPrice = v.UnitPrice
			existingItem.Amount = v.Amount

			_, errUpdate := s.invoiceItemsRepository.UpdateWithTx(ctx, existingItem, tx)
			if errUpdate != nil {
				log.Error(ctx, fmt.Sprintf("failed update with tx into invoice items with invoice ID %s", invoices.InvoiceID), errUpdate)
				err = fmt.Errorf("something went wrong [9]")
				tx.Rollback()
				return
			}
		} else {
			newItem := entities.InvoiceItems{
				InvoiceID: invoice.ID,
				ItemID:    v.ItemID,
				Quantity:  v.Quantity,
				UnitPrice: v.UnitPrice,
				Amount:    v.Amount,
			}

			_, errUpdate := s.invoiceItemsRepository.InsertWithTx(ctx, newItem, tx)
			if errUpdate != nil {
				log.Error(ctx, fmt.Sprintf("failed insert with tx into invoice items with invoice ID %s", invoices.InvoiceID), errUpdate)
				err = fmt.Errorf("something went wrong [10]")
				tx.Rollback()
				return
			}
		}
	}

	for _, existingItem := range existingItems {
		found := false
		for _, newItem := range req.InvoiceItems {
			if existingItem.ItemID == newItem.ItemID {
				found = true
				break
			}
		}
		if !found {
			_, err = s.invoiceItemsRepository.DeleteWithTx(ctx, existingItem.ID, tx)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("failed delete with tx into invoice items with invoice ID %s", invoices.InvoiceID), err)
				err = fmt.Errorf("something went wrong [11]")
				tx.Rollback()
				return
			}
		}
	}

	tx.Commit()

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    struct{}{},
		Errors:  make([]string, 0),
	}

	return
}

func (s *service) DeleteDetailInvoice(ctx context.Context, uniqueInvoiceID string) (resp constants.DefaultResponse, err error) {
	// * validate given invoice ID exists
	invoice, err := s.invoicesRepository.FindByUniqueInvoiceID(ctx, uniqueInvoiceID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find invoice by id %s", uniqueInvoiceID), err)

		if err == sql.ErrNoRows {
			err = fmt.Errorf("invoice not found")
			return
		}

		err = fmt.Errorf("something went wrong [0]")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(ctx, "failed to begin tx", err)
		err = fmt.Errorf("something went wrong [1]")
		tx.Rollback()
		return
	}

	res, err := s.invoicesRepository.DeleteByUniqueInvoiceIDWithTx(ctx, invoice.InvoiceID, tx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed delete with tx from invoices with unique invoice ID %s", invoice.InvoiceID), err)
		err = fmt.Errorf("something went wrong [2]")
		tx.Rollback()
		return
	}

	countRows, err := res.RowsAffected()
	if countRows < 1 || err != nil {
		log.Error(ctx, "SQL error rows affected invoices delete with tx", err)
		err = fmt.Errorf("something went wrong [3]")
		tx.Rollback()
		return
	}

	tx.Commit()

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    fmt.Sprintf("success delete invoice data with id %s", uniqueInvoiceID),
		Errors:  make([]string, 0),
	}

	return
}
