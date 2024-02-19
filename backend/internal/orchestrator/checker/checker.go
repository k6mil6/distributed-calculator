package checker

import (
	"context"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"log/slog"
	"time"
)

type SubexpressionGetter interface {
	TakenAt(context context.Context) ([]model.Subexpression, error)
	MakeNonTaken(context context.Context, id int) error
	MakeBeingChecked(context context.Context, id int) error
	GetById(context context.Context, id int) (model.Subexpression, error)
}

type Checker struct {
	subexpressionGetter SubexpressionGetter

	checkInterval time.Duration
	logger        *slog.Logger
}

func New(subexpressionGetter SubexpressionGetter, checkInterval time.Duration, logger *slog.Logger) *Checker {
	return &Checker{
		subexpressionGetter: subexpressionGetter,
		checkInterval:       checkInterval,
		logger:              logger,
	}
}

func (c *Checker) Start(context context.Context) {
	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()

	c.logger.Info("checker started")

	for {
		select {
		case <-ticker.C:
			c.Check(context)
		case <-context.Done():
			return
		}
	}
}

func (c *Checker) Check(ctx context.Context) {
	subexpressions, err := c.subexpressionGetter.TakenAt(ctx)
	if err != nil {
		c.logger.Error("failed to get subexpressions", err)
		return
	}

	if len(subexpressions) == 0 {
		return
	}

	for _, subexpression := range subexpressions {
		if subexpression.IsDone {
			continue
		}

		go func(subexpression model.Subexpression) {
			if err := c.subexpressionGetter.MakeBeingChecked(ctx, subexpression.ID); err != nil {
				c.logger.Error("failed to make subexpression being checked", err)
				return
			}

			checkTime := subexpression.TakenAt.Add(time.Duration(subexpression.Timeout) * time.Second).Add(2 * time.Minute)

			c.logger.Info("checking subexpression", "subexpression_id", subexpression.ID, "check_time", checkTime)
			timer := time.NewTimer(time.Until(checkTime))
			defer timer.Stop()
			<-timer.C

			subexp, err := c.subexpressionGetter.GetById(ctx, subexpression.ID)
			if err != nil {
				c.logger.Error("failed to get subexpression", err)
				return
			}

			if !subexp.IsDone {
				opCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				if err := c.subexpressionGetter.MakeNonTaken(opCtx, subexp.ID); err != nil {
					c.logger.Error("failed to make subexpression non-taken", err)
				} else {
					c.logger.Info("made subexpression non-taken", "subexpression_id", subexp.ID)
				}
			}

			c.logger.Info("checked subexpression", "subexpression_id", subexp.ID)
		}(subexpression)
	}
}
