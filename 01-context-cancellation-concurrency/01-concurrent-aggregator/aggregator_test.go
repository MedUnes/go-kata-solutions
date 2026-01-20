package concurrent_aggregator

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/order"
	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	type Input struct {
		profileService    ProfileServiceMock
		orderService      OrderServiceMock
		searchedProfileId int
		timeout           time.Duration
	}
	type Expected struct {
		aggregatedProfiles []*AggregatedProfile
		err                error
	}
	type TestCase struct {
		name     string
		input    Input
		expected Expected
	}
	basicProfiles := []*profile.Profile{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
		{Id: 3, Name: "Charlie"},
		{Id: 4, Name: "Dave"},
		{Id: 5, Name: "Eva"},
	}
	basicOrders := []*order.Order{
		{Id: 1, UserId: 1, Cost: 100.0},
		{Id: 2, UserId: 1, Cost: 20.6},
		{Id: 3, UserId: 3, Cost: 30.79},
	}

	testCases := []TestCase{
		{
			"profiles fetched in one second, orders fetched in two seconds, timeout is ",
			Input{
				ProfileServiceMock{
					1 * time.Second,
					nil,
					basicProfiles,
				},
				OrderServiceMock{
					2 * time.Second,
					nil,
					basicOrders,
				},
				1,
				5 * time.Second,
			},
			Expected{
				[]*AggregatedProfile{{"Alice", 100.0}},
				nil,
			},
		},
	}
	for _, tc := range testCases {

		t.Run("Test aggregator function", func(t *testing.T) {
			ctx := context.Background()
			u := NewUserAggregator(
				tc.input.orderService,
				tc.input.profileService,
				func(ua *UserAggregator) {
					ua.timeout = tc.input.timeout
				},
				func(ua *UserAggregator) {
					ua.logger = slog.Default()
				},
			)
			aggregatedProfiles, err := u.Aggregate(ctx, 1)

			require.Equal(t, tc.expected.err, err)
			require.GreaterOrEqual(t, len(aggregatedProfiles), len(tc.expected.aggregatedProfiles))
			if len(aggregatedProfiles) > 0 {
				ap := aggregatedProfiles[0]
				require.IsType(t, &AggregatedProfile{}, ap)
				assert.Equal(t, tc.expected.aggregatedProfiles[0].Name, ap.Name)
				assert.Equal(t, tc.expected.aggregatedProfiles[0].Cost, ap.Cost)
			}

		})
	}
}
