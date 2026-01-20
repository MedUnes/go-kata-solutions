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
	simulatedError    error
	simulatedProfiles []*profile.Profile
}
type OrderServiceMock struct {
	simulatedDuration time.Duration
	simulatedError    error
	simulatedOrders   []*order.Order
}

func (ps ProfileServiceMock) Get(ctx context.Context, id int) (*profile.Profile, error) {
	ticker := time.NewTicker(ps.simulatedDuration)
	fmt.Println("Simulating profile search..")
	select {
	case <-ticker.C:
		if ps.simulatedError != nil {
			return nil, fmt.Errorf("simulated profile search error")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if ps.simulatedError != nil {
		return nil, fmt.Errorf("simulated profile service error")
	}
	for _, p := range ps.simulatedProfiles {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, nil
}
func (os OrderServiceMock) GetAll(ctx context.Context, userId int) ([]*order.Order, error) {
	ticker := time.NewTicker(os.simulatedDuration)
	fmt.Println("Simulating orders search..")
	select {
	case <-ticker.C:
		if os.simulatedError != nil {
			return nil, fmt.Errorf("simulated orders search error")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var userOrders []*order.Order
	for _, o := range os.simulatedOrders {
		if o.UserId == userId {
			userOrders = append(userOrders, o)
		}
	}
	return userOrders, nil

}
func TestAggregate(t *testing.T) {
	t.Run("Test aggregator function", func(t *testing.T) {
		ctx := context.Background()
		orderService := &OrderServiceMock{
			simulatedDuration: 2 * time.Second,
			simulatedOrders: []*order.Order{
				{Id: 1, UserId: 1, Cost: 100.0},
				{Id: 2, UserId: 1, Cost: 20.6},
				{Id: 3, UserId: 3, Cost: 30.79},
			},
			simulatedError: nil,
		}
		profileService := &ProfileServiceMock{
			simulatedDuration: 1 * time.Second,
			simulatedProfiles: []*profile.Profile{
				{Id: 1, Name: "Alice"},
				{Id: 2, Name: "Bob"},
				{Id: 3, Name: "Charlie"},
				{Id: 4, Name: "Dave"},
				{Id: 5, Name: "Eva"},
			},
			simulatedError: nil,
		}
		u := NewUserAggregator(
			orderService,
			profileService,
			func(ua *UserAggregator) {
				ua.timeout = 3 * time.Second
			},
			func(ua *UserAggregator) {
				ua.timeout = 3 * time.Second
			},
		)
		aggregatedProfiles, err := u.Aggregate(ctx, 1)

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(aggregatedProfiles), 1)

		aggregatedProfile := aggregatedProfiles[0]
		assert.NotNil(t, aggregatedProfile)

		assert.NotNil(t, aggregatedProfile)
		assert.Equal(t, "Alice", aggregatedProfile.Name)
		assert.Equal(t, 100.0, aggregatedProfile.Cost)
	})
}
