package concurrent_aggregator

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/order"
	"github.com/medunes/go-kata/01-context-cancellation-concurrency/01-concurrent-aggregator/profile"

	"golang.org/x/sync/errgroup"
)

type Option func(ua *UserAggregator)

func WithLogger(l *slog.Logger) Option {
	return func(ua *UserAggregator) {
		ua.logger = l
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(ua *UserAggregator) {
		ua.timeout = timeout
	}
}

type UserAggregator struct {
	orderService   order.Service
	profileService profile.Service
	timeout        time.Duration
	logger         *slog.Logger
}

func NewUserAggregator(os order.Service, ps profile.Service, opts ...Option) *UserAggregator {
	ua := &UserAggregator{orderService: os, profileService: ps}
	for _, opt := range opts {
		opt(ua)
	}
	return ua
}

type AggregatedProfile struct {
	Name string
	Cost float64
}

func (ua *UserAggregator) Aggregate(ctx context.Context, id int) ([]*AggregatedProfile, error) {
	var (
		localCtx context.Context
		cancel   context.CancelFunc
		pr       *profile.Profile
		or       []*order.Order
		au       []*AggregatedProfile
		g        *errgroup.Group
	)
	if ua.timeout > 0 {
		localCtx, cancel = context.WithTimeout(ctx, ua.timeout)
	} else {
		localCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	g, localCtx = errgroup.WithContext(localCtx)
	g.Go(func() error {
		var err error
		pr, err = ua.profileService.Get(localCtx, id)
		if err != nil {
			return fmt.Errorf("profile not found")
		}
		return nil
	})
	g.Go(func() error {
		var err error
		or, err = ua.orderService.GetAll(localCtx, id)
		if err != nil {
			return fmt.Errorf("profile not found")
		}
		return nil

	})

	err := g.Wait()
	if err != nil {
		return nil, err
	}
	for _, o := range or {
		if o.UserId == id {
			au = append(au, &AggregatedProfile{pr.Name, o.Cost})
		}
	}
	return au, nil
}
