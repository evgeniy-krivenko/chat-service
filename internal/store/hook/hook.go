// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
)

// The ChatFunc type is an adapter to allow the use of ordinary
// function as Chat mutator.
type ChatFunc func(context.Context, *store.ChatMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f ChatFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.ChatMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.ChatMutation", m)
}

// The FailedJobFunc type is an adapter to allow the use of ordinary
// function as FailedJob mutator.
type FailedJobFunc func(context.Context, *store.FailedJobMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f FailedJobFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.FailedJobMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.FailedJobMutation", m)
}

// The JobFunc type is an adapter to allow the use of ordinary
// function as Job mutator.
type JobFunc func(context.Context, *store.JobMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f JobFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.JobMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.JobMutation", m)
}

// The MessageFunc type is an adapter to allow the use of ordinary
// function as Message mutator.
type MessageFunc func(context.Context, *store.MessageMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f MessageFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.MessageMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.MessageMutation", m)
}

// The ProblemFunc type is an adapter to allow the use of ordinary
// function as Problem mutator.
type ProblemFunc func(context.Context, *store.ProblemMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f ProblemFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.ProblemMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.ProblemMutation", m)
}

// The ProfileFunc type is an adapter to allow the use of ordinary
// function as Profile mutator.
type ProfileFunc func(context.Context, *store.ProfileMutation) (store.Value, error)

// Mutate calls f(ctx, m).
func (f ProfileFunc) Mutate(ctx context.Context, m store.Mutation) (store.Value, error) {
	if mv, ok := m.(*store.ProfileMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *store.ProfileMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, store.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m store.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m store.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m store.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op store.Op) Condition {
	return func(_ context.Context, m store.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m store.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m store.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m store.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk store.Hook, cond Condition) store.Hook {
	return func(next store.Mutator) store.Mutator {
		return store.MutateFunc(func(ctx context.Context, m store.Mutation) (store.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, store.Delete|store.Create)
func On(hk store.Hook, op store.Op) store.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, store.Update|store.UpdateOne)
func Unless(hk store.Hook, op store.Op) store.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) store.Hook {
	return func(store.Mutator) store.Mutator {
		return store.MutateFunc(func(context.Context, store.Mutation) (store.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []store.Hook {
//		return []store.Hook{
//			Reject(store.Delete|store.Update),
//		}
//	}
func Reject(op store.Op) store.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []store.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...store.Hook) Chain {
	return Chain{append([]store.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() store.Hook {
	return func(mutator store.Mutator) store.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...store.Hook) Chain {
	newHooks := make([]store.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
