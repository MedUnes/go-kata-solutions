package concurrent_aggregator

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/order"
	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/profile"

	"github.com/stretchr/testify/assert"
)

type ProfileServiceMock struct {
	simulatedDuration time.Duration
	simulateError     bool
}
type OrderServiceMock struct {
	simulatedDuration time.Duration
	simulateError     bool
}

func (ps ProfileServiceMock) Get(ctx context.Context, id int) (*profile.Profile, error) {
	ticker := time.NewTicker(ps.simulatedDuration)
	fmt.Println("Simulating profile search..")
	select {
	case <-ticker.C:
		if ps.simulateError {
			return nil, fmt.Errorf("simulated profile search error")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if ps.simulateError {
		return nil, fmt.Errorf("simulated profile service error")
	}
	profiles := []profile.Profile{
		{Id: id, Name: "Alice"},
		{Id: id, Name: "Bob"},
		{Id: id, Name: "Charlie"},
		{Id: id, Name: "Dave"},
		{Id: id, Name: "Eva"},
	}
	for _, p := range profiles {
		if p.Id == id {
			return &p, nil
		}
	}
	return nil, nil
}
func (os OrderServiceMock) GetAll(ctx context.Context, userId int) ([]*order.Order, error) {
	ticker := time.NewTicker(os.simulatedDuration)
	fmt.Println("Simulating orders search..")
	select {
	case <-ticker.C:
		if os.simulateError {
			return nil, fmt.Errorf("simulated orders search error")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	orders := []*order.Order{
		{Id: 1, UserId: 1, Cost: 100.0},
		{Id: 2, UserId: 1, Cost: 20.6},
		{Id: 3, UserId: 3, Cost: 30.79},
	}
	var userOrders []*order.Order
	for _, o := range orders {
		if o.UserId == userId {
			userOrders = append(userOrders, o)
		}
	}
	return orders, nil

}
func TestAggregate(t *testing.T) {
	t.Run("Test aggregator function", func(t *testing.T) {
		ctx := context.Background()
		orderService := &OrderServiceMock{
			simulatedDuration: 10 * time.Millisecond,
			simulateError:     false,
		}
		profileService := &ProfileServiceMock{
			simulatedDuration: 10 * time.Millisecond,
			simulateError:     false,
		}
		u := NewUserAggregator(
			orderService,
			profileService,
			func(ua *UserAggregator) {
				ua.timeout = 3 * time.Millisecond
			},
			func(ua *UserAggregator) {
				ua.timeout = 3 * time.Millisecond
			},
		)
		aggregatedProfiles, err := u.Aggregate(ctx, 1)

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(aggregatedProfiles), 1)

		aggregatedProfile := aggregatedProfiles[0]
		assert.NotNil(t, aggregatedProfile)

		assert.NotNil(t, aggregatedProfile)
		assert.Equal(t, "Alice", aggregatedProfile.Name)
		assert.Equal(t, 20.6, aggregatedProfile.Cost)
	})
}
