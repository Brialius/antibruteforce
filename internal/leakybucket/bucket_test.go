package leakybucket

import (
	"context"
	"testing"
	"time"
)

var ctx = context.Background()

func TestBucket(t *testing.T) {
	tests := []struct {
		name                string
		rateLimit           uint64
		duration            time.Duration
		inactiveCycles      uint64
		resetAt             int
		timeoutAfter        time.Duration
		requests            int
		timeBetweenRequests time.Duration
		wantTimeout         bool
		want                int
	}{
		{
			name:                "TestBucket10",
			rateLimit:           10,
			duration:            1 * time.Second,
			inactiveCycles:      2,
			resetAt:             0,
			timeoutAfter:        5 * time.Second,
			requests:            50,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                10,
		},
		{
			name:                "TestBucket100",
			rateLimit:           100,
			duration:            1 * time.Second,
			inactiveCycles:      2,
			resetAt:             0,
			timeoutAfter:        5 * time.Second,
			requests:            500,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                100,
		},
		{
			name:                "TestBucket1000",
			rateLimit:           1000,
			duration:            5 * time.Second,
			inactiveCycles:      2,
			resetAt:             0,
			timeoutAfter:        20 * time.Second,
			requests:            5000,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                1000,
		},
		{
			name:                "TestBucket10WithReset",
			rateLimit:           10,
			duration:            1 * time.Second,
			inactiveCycles:      2,
			resetAt:             6,
			timeoutAfter:        5 * time.Second,
			requests:            50,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                15,
		},
		{
			name:                "TestBucket100WithReset",
			rateLimit:           100,
			duration:            1 * time.Second,
			inactiveCycles:      2,
			resetAt:             51,
			timeoutAfter:        5 * time.Second,
			requests:            500,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                150,
		},
		{
			name:                "TestBucket1000WithReset",
			rateLimit:           1000,
			duration:            5 * time.Second,
			inactiveCycles:      2,
			resetAt:             501,
			timeoutAfter:        20 * time.Second,
			requests:            5000,
			timeBetweenRequests: 0,
			wantTimeout:         false,
			want:                1500,
		},
		{
			name:                "TestBucket10UnderRate",
			rateLimit:           10,
			duration:            1 * time.Second,
			inactiveCycles:      2,
			resetAt:             0,
			timeoutAfter:        20 * time.Second,
			requests:            50,
			timeBetweenRequests: 150 * time.Millisecond,
			wantTimeout:         false,
			want:                50,
		},
		{
			name:                "TestBucket10Timeout",
			rateLimit:           10,
			duration:            1 * time.Second,
			inactiveCycles:      10,
			resetAt:             0,
			timeoutAfter:        10 * time.Second,
			requests:            50,
			timeBetweenRequests: 0,
			wantTimeout:         true,
			want:                10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBucket(tt.name, tt.rateLimit, tt.duration, tt.inactiveCycles)
			var passed int
			for i := 1; i <= tt.requests; i++ {
				if i == tt.resetAt {
					b.ResetLimit(ctx)
				}
				if b.CheckLimit(ctx) {
					passed++
				}
				time.Sleep(tt.timeBetweenRequests)
			}

			timeout := time.After(tt.timeoutAfter)
			select {
			case <-timeout:
				if !tt.wantTimeout {
					t.Error("Bucket didn't finished in time")
				}
			case <-b.Inactive(ctx):
			}

			if passed != tt.want {
				t.Errorf("passed = %v, want %v", passed, tt.want)
			}
		})
	}
}
