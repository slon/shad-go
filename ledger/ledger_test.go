package ledger_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/ledger"
	"gitlab.com/slon/shad-go/pgfixture"
)

func TestLedger(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	dsn := pgfixture.Start(t)

	ctx := context.Background()

	l0, err := ledger.New(ctx, dsn)
	require.NoError(t, err)
	defer func() { _ = l0.Close() }()

	t.Run("SimpleCommands", func(t *testing.T) {
		checkBalance := func(account ledger.ID, amount ledger.Money) {
			b, err := l0.GetBalance(ctx, account)
			require.NoError(t, err)
			require.Equal(t, amount, b)
		}

		require.NoError(t, l0.CreateAccount(ctx, "a0"))
		checkBalance("a0", 0)

		require.Error(t, l0.CreateAccount(ctx, "a0"))

		require.NoError(t, l0.Deposit(ctx, "a0", ledger.Money(100)))
		checkBalance("a0", 100)

		require.NoError(t, l0.Withdraw(ctx, "a0", ledger.Money(50)))
		checkBalance("a0", 50)

		require.ErrorIs(t, l0.Withdraw(ctx, "a0", ledger.Money(100)), ledger.ErrNoMoney)

		require.NoError(t, l0.CreateAccount(ctx, "a1"))

		require.NoError(t, l0.Transfer(ctx, "a0", "a1", ledger.Money(40)))
		checkBalance("a0", 10)
		checkBalance("a1", 40)

		require.ErrorIs(t, l0.Transfer(ctx, "a0", "a1", ledger.Money(50)), ledger.ErrNoMoney)
	})

	t.Run("ErroneousCases", func(t *testing.T) {
		checkBalance := func(account ledger.ID, amount ledger.Money) {
			b, err := l0.GetBalance(ctx, account)
			require.NoError(t, err)
			require.Equal(t, amount, b)
		}
		require.NoError(t, l0.CreateAccount(ctx, "b0"))
		checkBalance("b0", 0)

		require.NoError(t, l0.Deposit(ctx, "b0", ledger.Money(100)))
		checkBalance("b0", 100)

		require.Error(t, l0.Deposit(ctx, "b0", ledger.Money(-100)))
		checkBalance("b0", 100)

		require.Error(t, l0.Withdraw(ctx, "b0", ledger.Money(-50)))
		checkBalance("b0", 100)

		require.Error(t, l0.Transfer(ctx, "b0", "b999", ledger.Money(50)))
		checkBalance("b0", 100)

		require.Error(t, l0.Transfer(ctx, "b999", "b0", ledger.Money(50)))
		checkBalance("b0", 100)

		require.NoError(t, l0.CreateAccount(ctx, "b999"))
		require.NoError(t, l0.Deposit(ctx, "b999", ledger.Money(200)))
		checkBalance("b999", 200)

		require.Error(t, l0.Transfer(ctx, "b0", "b999", ledger.Money(-50)))
		checkBalance("b0", 100)
		checkBalance("b999", 200)

		require.NoError(t, l0.Transfer(ctx, "b0", "b999", ledger.Money(50)))
		checkBalance("b0", 50)
		checkBalance("b999", 250)

		require.Error(t, l0.Deposit(ctx, "c0", ledger.Money(100)))
		require.Error(t, l0.Withdraw(ctx, "c0", ledger.Money(100)))

		_, err := l0.GetBalance(ctx, "c0")
		require.Error(t, err)

	})

	t.Run("Transactions", func(t *testing.T) {
		const nAccounts = 10
		const initialBalance = 5

		var accounts []ledger.ID
		for i := 0; i < nAccounts; i++ {
			id := ledger.ID(fmt.Sprint(i))
			accounts = append(accounts, id)

			require.NoError(t, l0.CreateAccount(ctx, id))
			require.NoError(t, l0.Deposit(ctx, id, initialBalance))
		}

		var wg sync.WaitGroup
		done := make(chan struct{})

		spawn := func(action func() error) {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-done:
						return

					default:
						if err := action(); err != nil {
							if !errors.Is(err, ledger.ErrNoMoney) {
								t.Errorf("operation failed: %v", err)
								return
							}
						}
					}
				}
			}()
		}

		for i := 0; i < nAccounts; i++ {
			i := i

			account := accounts[i]
			next := accounts[(i+1)%len(accounts)]
			prev := accounts[(i+len(accounts)-1)%len(accounts)]

			spawn(func() error {
				balance, err := l0.GetBalance(ctx, account)
				if err != nil {
					return err
				}

				if balance < 0 {
					return fmt.Errorf("%q balance is negative", account)
				}

				return nil
			})

			spawn(func() error {
				return l0.Transfer(ctx, account, next, 1)
			})

			spawn(func() error {
				return l0.Transfer(ctx, account, prev, 1)
			})
		}

		time.Sleep(time.Second * 10)
		close(done)
		wg.Wait()

		var total ledger.Money
		for i := 0; i < nAccounts; i++ {
			amount, err := l0.GetBalance(ctx, accounts[i])
			require.NoError(t, err)

			total += amount
		}

		require.Equal(t, ledger.Money(initialBalance*nAccounts), total)
	})
}
