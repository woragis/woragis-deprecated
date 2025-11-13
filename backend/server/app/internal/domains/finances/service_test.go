package finances

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNormalizeAmountRequiresExchangeRate(t *testing.T) {
	svc := &service{}
	_, _, err := svc.normalizeAmount(100, "EUR", "USD", 0)
	if err == nil {
		t.Fatalf("expected error when exchange rate missing")
	}
}

func TestCashflowProjectionIncludesTemplates(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	userID := uuid.New()

	// Seed two months of actual data in normalized USD.
	month1 := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	month2 := time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)

	repo.transactions = append(repo.transactions,
		&Transaction{ID: uuid.New(), UserID: userID, Type: TransactionTypeIncome, NormalizedAmount: 2000, Amount: 2000, Currency: "USD", BaseCurrency: "USD", OccurredAt: month1},
		&Transaction{ID: uuid.New(), UserID: userID, Type: TransactionTypeExpense, NormalizedAmount: 800, Amount: 800, Currency: "USD", BaseCurrency: "USD", OccurredAt: month1},
		&Transaction{ID: uuid.New(), UserID: userID, Type: TransactionTypeIncome, NormalizedAmount: 2100, Amount: 2100, Currency: "USD", BaseCurrency: "USD", OccurredAt: month2},
		&Transaction{ID: uuid.New(), UserID: userID, Type: TransactionTypeExpense, NormalizedAmount: 900, Amount: 900, Currency: "USD", BaseCurrency: "USD", OccurredAt: month2},
	)

	tpl, err := NewRecurringTemplate(userID, "Salary", TransactionTypeIncome, "compensation", "", 500, "USD", "USD", FrequencyMonthly, 1, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating template: %v", err)
	}
	tpl.NormalizedAmount = 500
	repo.templates = append(repo.templates, *tpl)

	projection, err := svc.CashflowProjection(context.Background(), CashflowProjectionRequest{UserID: userID, PastMonths: 2, FutureMonths: 2})
	if err != nil {
		t.Fatalf("unexpected error building projection: %v", err)
	}

	if len(projection.Months) != 4 {
		t.Fatalf("expected 4 months in projection, got %d", len(projection.Months))
	}

	future := projection.Months[len(projection.Months)-1]
	if future.ProjectedIncome <= 0 {
		t.Fatalf("expected projected income to include template contribution")
	}
}

func TestImportTransactionsCSV(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, nil)

	userID := uuid.New()

	csv := "type,category,description,amount,currency,base_currency,exchange_rate,occurred_at,tags\n" +
		"income,sales,January retainer,1000,USD,USD,,2025-01-05T00:00:00Z,client-a;Retainer\n" +
		"expense,software,SAAS license,200,EUR,USD,1.1,2025-01-06T00:00:00Z,operations"

	payload := ImportTransactionsRequest{
		UserID:          userID,
		Format:          ImportFormatCSV,
		BaseCurrency:    "USD",
		DefaultCurrency: "USD",
		Contents:        []byte(csv),
	}

	txs, err := svc.ImportTransactions(context.Background(), payload)
	if err != nil {
		t.Fatalf("unexpected error importing csv: %v", err)
	}

	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}

	if len(repo.transactions) != 2 {
		t.Fatalf("repository should contain imported transactions")
	}

	if repo.transactions[1].BaseCurrency != "USD" {
		t.Fatalf("expected base currency normalization to USD")
	}
}

type mockRepository struct {
	transactions []*Transaction
	templates    []RecurringTemplate
}

func newMockRepository() *mockRepository {
	return &mockRepository{transactions: []*Transaction{}, templates: []RecurringTemplate{}}
}

func (m *mockRepository) CreateTransaction(_ context.Context, tx *Transaction) error {
	m.transactions = append(m.transactions, cloneTransaction(tx))
	return nil
}

func (m *mockRepository) UpdateTransaction(_ context.Context, tx *Transaction) error {
	for idx, existing := range m.transactions {
		if existing.ID == tx.ID {
			m.transactions[idx] = cloneTransaction(tx)
			return nil
		}
	}
	return nil
}

func (m *mockRepository) GetTransaction(_ context.Context, userID, id uuid.UUID) (*Transaction, error) {
	for _, tx := range m.transactions {
		if tx.ID == id && tx.UserID == userID {
			return cloneTransaction(tx), nil
		}
	}
	return nil, NewDomainError(ErrCodeNotFound, ErrTransactionNotFound)
}

func (m *mockRepository) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	return m.QueryTransactions(ctx, TransactionFilter{UserID: userID, From: from, To: to})
}

func (m *mockRepository) QueryTransactions(_ context.Context, filter TransactionFilter) ([]Transaction, error) {
	result := []Transaction{}

	for _, tx := range m.transactions {
		if tx.UserID != filter.UserID {
			continue
		}

		if !filter.From.IsZero() && tx.OccurredAt.Before(filter.From) {
			continue
		}
		if !filter.To.IsZero() && tx.OccurredAt.After(filter.To) {
			continue
		}

		result = append(result, *cloneTransaction(tx))
	}

	return result, nil
}

func (m *mockRepository) BulkCreateTransactions(_ context.Context, txs []*Transaction) error {
	for _, tx := range txs {
		m.transactions = append(m.transactions, cloneTransaction(tx))
	}
	return nil
}

func (m *mockRepository) BulkUpdateCategory(_ context.Context, _ uuid.UUID, _ []uuid.UUID, _ string) error {
	return nil
}

func (m *mockRepository) BulkUpdateType(_ context.Context, _ uuid.UUID, _ []uuid.UUID, _ TransactionType) error {
	return nil
}

func (m *mockRepository) BulkDeleteTransactions(_ context.Context, _ uuid.UUID, _ []uuid.UUID) error {
	return nil
}

func (m *mockRepository) SetArchived(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return nil
}

func (m *mockRepository) SetRecurring(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return nil
}

func (m *mockRepository) SetEssential(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return nil
}

func (m *mockRepository) AggregateSummary(_ context.Context, userID uuid.UUID, _, _ time.Time) (Summary, error) {
	summary := Summary{}

	for _, tx := range m.transactions {
		if tx.UserID != userID {
			continue
		}
		if tx.Type == TransactionTypeIncome {
			summary.IncomeTotal += tx.NormalizedAmount
			summary.BaseCurrency = tx.BaseCurrency
		} else if tx.Type == TransactionTypeExpense {
			summary.ExpenseTotal += tx.NormalizedAmount
			summary.BaseCurrency = tx.BaseCurrency
		}
	}

	summary.SavingsAllocation = summary.IncomeTotal * 0.5
	return summary, nil
}

func (m *mockRepository) CreateTemplate(_ context.Context, template *RecurringTemplate) error {
	m.templates = append(m.templates, *template)
	return nil
}

func (m *mockRepository) UpdateTemplate(_ context.Context, template *RecurringTemplate) error {
	for idx, existing := range m.templates {
		if existing.ID == template.ID {
			m.templates[idx] = *template
			return nil
		}
	}
	return nil
}

func (m *mockRepository) DeleteTemplate(_ context.Context, _ uuid.UUID, templateID uuid.UUID) error {
	filtered := make([]RecurringTemplate, 0, len(m.templates))
	for _, tpl := range m.templates {
		if tpl.ID != templateID {
			filtered = append(filtered, tpl)
		}
	}
	m.templates = filtered
	return nil
}

func (m *mockRepository) GetTemplate(_ context.Context, userID, templateID uuid.UUID) (*RecurringTemplate, error) {
	for _, tpl := range m.templates {
		if tpl.ID == templateID && tpl.UserID == userID {
			copy := tpl
			return &copy, nil
		}
	}
	return nil, NewDomainError(ErrCodeNotFound, ErrTemplateNotFound)
}

func (m *mockRepository) ListTemplates(_ context.Context, userID uuid.UUID) ([]RecurringTemplate, error) {
	result := []RecurringTemplate{}
	for _, tpl := range m.templates {
		if tpl.UserID == userID {
			result = append(result, tpl)
		}
	}
	return result, nil
}

func (m *mockRepository) MonthlyTotals(_ context.Context, userID uuid.UUID, from, to time.Time) ([]MonthlyTotal, error) {
	monthly := map[string]MonthlyTotal{}

	for _, tx := range m.transactions {
		if tx.UserID != userID {
			continue
		}
		if !from.IsZero() && tx.OccurredAt.Before(from) {
			continue
		}
		if !to.IsZero() && tx.OccurredAt.After(to) {
			continue
		}

		year, month, _ := tx.OccurredAt.Date()
		key := monthlyKey(year, int(month))
		en := monthly[key]
		en.Year = year
		en.Month = int(month)
		if tx.Type == TransactionTypeIncome {
			en.IncomeTotal += tx.NormalizedAmount
		} else {
			en.ExpenseTotal += tx.NormalizedAmount
		}
		monthly[key] = en
	}

	result := make([]MonthlyTotal, 0, len(monthly))
	for _, value := range monthly {
		result = append(result, value)
	}

	return result, nil
}

func cloneTransaction(tx *Transaction) *Transaction {
	copy := *tx
	return &copy
}

func TestParseCSVTransactionsTrimsTags(t *testing.T) {
	csv := "type,category,description,amount,currency,base_currency,exchange_rate,occurred_at,tags\n" +
		"income,consulting,,100,USD,USD,,2025-01-02T00:00:00Z, Tag-One ; TAG-two "
	data := base64.StdEncoding.EncodeToString([]byte(csv))

	repo := newMockRepository()
	svc := NewService(repo, nil)

	userID := uuid.New()
	payload := ImportTransactionsRequest{
		UserID:          userID,
		Format:          ImportFormatCSV,
		BaseCurrency:    "USD",
		DefaultCurrency: "USD",
		Contents:        decodeBase64(t, data),
	}

	txs, err := svc.ImportTransactions(context.Background(), payload)
	if err != nil {
		t.Fatalf("unexpected error importing csv: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected single transaction")
	}

	if len(txs[0].Tags) != 2 {
		t.Fatalf("expected tags to be normalized")
	}
}

func decodeBase64(t *testing.T, encoded string) []byte {
	t.Helper()
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}
	return decoded
}
