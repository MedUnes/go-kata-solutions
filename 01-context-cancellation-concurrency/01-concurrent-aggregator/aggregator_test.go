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
	simulateError     bool
	simulatedProfiles []*profile.Profile
}
type OrderServiceMock struct {
	simulatedDuration time.Duration
	simulateError     bool
	simulatedOrders   []*order.Order
}

func (ps *ProfileServiceMock) Get(ctx context.Context, id int) (*profile.Profile, error) {

	fmt.Println("Simulating profile search..")
	select {
	case <-time.After(ps.simulatedDuration):
		if ps.simulateError {
			return nil, fmt.Errorf("simulated profile search error")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	for _, p := range ps.simulatedProfiles {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, nil
}
func (os *OrderServiceMock) GetAll(ctx context.Context, userId int) ([]*order.Order, error) {
	fmt.Println("Simulating orders search..")
	var emptyOrders []*order.Order
	select {
	case <-time.After(os.simulatedDuration):
		if os.simulateError {
			return emptyOrders, fmt.Errorf("simulated orders search error")
		}
	case <-ctx.Done():
		return emptyOrders, ctx.Err()
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
		profileService        *ProfileServiceMock
		orderService          *OrderServiceMock
		searchedProfileId     int
		timeout               time.Duration
		parentContextDuration time.Duration
	}
	type Expected struct {
		aggregatedProfiles []*AggregatedProfile
		err                bool
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
	//var emptyProfiles []*profile.Profile
	//var emptyOrders []*order.Order
	var emptyAggregateProfiles []*AggregatedProfile
	testCases := []TestCase{
		{
			"no errors, profile service in time, order service in time, more than one order",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					20 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{2, 1, 20.6},
						{3, 3, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				[]*AggregatedProfile{{"Alice", 100.0}, {"Alice", 20.6}},
				false,
			},
		},
		{
			"no errors, profile service in time, order service in time, one order",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					20 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 3, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				[]*AggregatedProfile{{"Alice", 100.0}},
				false,
			},
		},
		{
			"no errors, profile service in time, order service in time, no orders",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					20 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 2, 100.0},
						{3, 3, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				false,
			},
		},
		{
			"no errors, profile service in time, order service timeout, more than one order",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					120 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no errors, profile service in time, order service timeout, one order",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					120 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 2, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no errors, profile service in time, order service timeout, no orders",
			Input{
				&ProfileServiceMock{
					10 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					120 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 3, 100.0},
						{3, 2, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no errors, profile service timeout, order service in time, more than one order",
			Input{
				&ProfileServiceMock{
					120 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					20 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no service timeout, no errors, more than one order",
			Input{
				&ProfileServiceMock{
					120 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					150 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				0 * time.Millisecond,
				0,
			},
			Expected{
				[]*AggregatedProfile{{"Alice", 100.0}, {"Alice", 30.79}},
				false,
			},
		},
		{
			"no service timeout, profile error, more than one order",
			Input{
				&ProfileServiceMock{
					120 * time.Millisecond,
					true,
					basicProfiles,
				},
				&OrderServiceMock{
					150 * time.Millisecond,
					false,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				0 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no service timeout, order error, more than one order",
			Input{
				&ProfileServiceMock{
					70 * time.Millisecond,
					false,
					basicProfiles,
				},
				&OrderServiceMock{
					100 * time.Millisecond,
					true,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				0 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"no service timeout, profile error, order error, profile and order take same time",
			Input{
				&ProfileServiceMock{
					120 * time.Millisecond,
					true,
					basicProfiles,
				},
				&OrderServiceMock{
					120 * time.Millisecond,
					true,
					[]*order.Order{
						{1, 1, 100.0},
						{3, 1, 30.79},
					},
				},
				1,
				0 * time.Millisecond,
				0,
			},
			Expected{
				emptyAggregateProfiles,
				true,
			},
		},
		{
			"order service error propagates correctly",
			Input{
				&ProfileServiceMock{10 * time.Millisecond, false, basicProfiles},
				&OrderServiceMock{10 * time.Millisecond, true, nil}, // Error here
				1,
				100 * time.Millisecond,
				0,
			},
			Expected{
				nil,
				true,
			},
		},
		{
			"timeout error propagates correctly",
			Input{
				&ProfileServiceMock{200 * time.Millisecond, false, basicProfiles},
				&OrderServiceMock{10 * time.Millisecond, false, nil},
				1,
				50 * time.Millisecond,
				0,
			},
			Expected{
				nil,
				true,
			},
		},
		{
			"parent context cancellation propagates correctly",
			Input{
				&ProfileServiceMock{200 * time.Millisecond, false, basicProfiles},
				&OrderServiceMock{10 * time.Millisecond, false, nil},
				1,
				500 * time.Millisecond,
				100 * time.Millisecond,
			},
			Expected{
				nil,
				true,
			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.input.parentContextDuration > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.input.parentContextDuration)
				defer cancel()
			}
			u := NewUserAggregator(
				tc.input.orderService,
				tc.input.profileService,
				WithTimeout(tc.input.timeout),
				WithLogger(slog.Default()),
			)
			aggregatedProfiles, err := u.Aggregate(ctx, tc.input.searchedProfileId)

			if tc.expected.err {
				require.Error(t, err)

			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected.aggregatedProfiles, aggregatedProfiles)

		})
	}
}
