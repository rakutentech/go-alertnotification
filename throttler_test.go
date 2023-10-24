package alertnotification

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/GitbookIO/diskache"
)

func TestNewThrottler(t *testing.T) {
	tests := []struct {
		name string
		want Throttler
	}{
		{
			name: "default",
			want: Throttler{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
		},
		{
			name: "change duration",
			want: Throttler{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 7,
			},
		},
		{
			name: "change both",
			want: Throttler{
				CacheOpt:         "new_cache_dir",
				ThrottleDuration: 8,
			},
		},
	}
	for _, tt := range tests {
		if tt.name == "change duration" {
			os.Setenv("THROTTLE_DURATION", "7")
		} else if tt.name == "change both" {
			os.Setenv("THROTTLE_DURATION", "8")
			os.Setenv("THROTTLE_DISKCACHE_DIR", "new_cache_dir")
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := NewThrottler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThrottler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThrottler_IsThrottled(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
	}
	type args struct {
		ocError error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "default",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
			args: args{
				ocError: errors.New("test_throttling"),
			},
			want: false,
		},
		{
			name: "throttled_true",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
			args: args{
				ocError: errors.New("test_throttling"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Throttler{
				CacheOpt:         tt.fields.CacheOpt,
				ThrottleDuration: tt.fields.ThrottleDuration,
			}
			if tt.name == "throttled_true" {
				if err := th.ThrottleError(tt.args.ocError); err != nil {
					t.Errorf("testing failed : %+v", err)
				}

			}
			if got := th.IsThrottled(tt.args.ocError); got != tt.want {
				t.Errorf("Throttler.IsThrottled() = %v, want %v", got, tt.want)
			}
			err := th.CleanThrottlingCache()
			if err != nil {
				t.Errorf("Cannot clean after test. %+v", err)
			}

		})
	}
}
func TestThrottler_IsThrottledGraced(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
	}
	type args struct {
		ocError error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "default",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
			args: args{
				ocError: errors.New("test_throttling"),
			},
			want: true,
		},
	}

	os.Setenv("THROTTLE_GRACE_SECONDS", "10")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Throttler{
				CacheOpt:         tt.fields.CacheOpt,
				ThrottleDuration: tt.fields.ThrottleDuration,
			}
			if got := th.IsThrottled(tt.args.ocError); got != tt.want {
				t.Errorf("Throttler.IsThrottled() = %v, want %v", got, tt.want)
			}
			err := th.CleanThrottlingCache()
			if err != nil {
				t.Errorf("Cannot clean after test. %+v", err)
			}

		})
	}
}
func TestThrottler_IsThrottledOverGraced(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
	}
	type args struct {
		ocError error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "default",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
			args: args{
				ocError: errors.New("test_throttling"),
			},
			want: false,
		},
	}

	os.Setenv("THROTTLE_GRACE_SECONDS", "0")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Throttler{
				CacheOpt:         tt.fields.CacheOpt,
				ThrottleDuration: tt.fields.ThrottleDuration,
			}
			if got := th.IsThrottled(tt.args.ocError); got != tt.want {
				t.Errorf("Throttler.IsThrottled() = %v, want %v", got, tt.want)
			}
			err := th.CleanThrottlingCache()
			if err != nil {
				t.Errorf("Cannot clean after test. %+v", err)
			}

		})
	}
}

func TestThrottler_ThrottleError(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
	}
	type args struct {
		errObj error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
			},
			args: args{
				errObj: errors.New("test_throttling"),
			},
			wantErr: false,
		},
		{
			name: "test_error",
			fields: fields{
				CacheOpt:         "/no_permission_dir",
				ThrottleDuration: 5,
			},
			args: args{
				errObj: errors.New("test_throttling"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			th := &Throttler{
				CacheOpt:         tt.fields.CacheOpt,
				ThrottleDuration: tt.fields.ThrottleDuration,
			}

			if err := th.ThrottleError(tt.args.errObj); (err != nil) != tt.wantErr {
				t.Errorf("Throttler.ThrottleError() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "default" && !th.IsThrottled(tt.args.errObj) {
				t.Errorf("Throttler.ThrottleError() error = %v, wantErr %v", errors.New("throttling failed"), tt.wantErr)
			}
			if !tt.wantErr {
				err := th.CleanThrottlingCache()
				if err != nil {
					t.Errorf("Cannot clean after test. %+v", err)
				}
			}

		})
	}
}

func TestThrottler_getDiskCache(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
	}
	cachePart := fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME"))
	opts := diskache.Opts{
		Directory: cachePart,
	}
	dc, err := diskache.New(&opts)
	if err != nil {
		t.Errorf("Throttler.getDiskCache() error = %v", err)
		return
	}

	tests := []struct {
		name    string
		fields  fields
		want    *diskache.Diskache
		wantErr bool
	}{
		{
			name: "TestThrottler_getDiskCache_success",
			fields: fields{
				CacheOpt:         cachePart,
				ThrottleDuration: 5,
			},
			want:    dc,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Throttler{
				CacheOpt:         tt.fields.CacheOpt,
				ThrottleDuration: tt.fields.ThrottleDuration,
			}
			got, err := th.getDiskCache()
			if (err != nil) != tt.wantErr {
				t.Errorf("Throttler.getDiskCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Throttler.getDiskCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOverThrottleDuration(t *testing.T) {
	type args struct {
		cachedTime       string
		throttleDuration int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_isOverThrottleDuration_true",
			args: args{
				cachedTime:       time.Now().Add(-3 * time.Minute).Format(time.RFC3339), // -3 minutes => pass 2 minutes durations
				throttleDuration: 2,
			},
			want: true,
		},
		{
			name: "Test_isOverThrottleDuration_false",
			args: args{
				cachedTime:       time.Now().Add(1 * time.Minute).Format(time.RFC3339), // 1 minute ahead of current < throtte duration 2
				throttleDuration: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOverThrottleDuration(tt.args.cachedTime, tt.args.throttleDuration); got != tt.want {
				t.Errorf("isOverThrottleDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
