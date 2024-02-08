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
				GraceDuration:    0,
			},
		},
		{
			name: "change duration",
			want: Throttler{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 7,
				GraceDuration:    5,
			},
		},
		{
			name: "change both",
			want: Throttler{
				CacheOpt:         "new_cache_dir",
				ThrottleDuration: 8,
				GraceDuration:    0,
			},
		},
	}
	for _, tt := range tests {
		if tt.name == "change duration" {
			os.Setenv("THROTTLE_DURATION", "7")
			os.Setenv("THROTTLE_GRACE_SECONDS", "5")
		} else if tt.name == "change both" {
			os.Setenv("THROTTLE_DURATION", "8")
			os.Setenv("THROTTLE_GRACE_SECONDS", "0")
			os.Setenv("THROTTLE_DISKCACHE_DIR", "new_cache_dir")
		} else if tt.name == "default" {
			os.Setenv("THROTTLE_GRACE_SECONDS", "")
		}
		t.Run(tt.name, func(t *testing.T) {
			got := NewThrottler()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThrottler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThrottler_IsThrottledOrGraced(t *testing.T) {
	type fields struct {
		CacheOpt         string
		ThrottleDuration int
		GraceDuration    int
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
				GraceDuration:    0,
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
				GraceDuration:    0,
			},
			args: args{
				ocError: errors.New("test_throttling"),
			},
			want: true,
		},
		{
			name: "graced_true",
			fields: fields{
				CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
				ThrottleDuration: 5,
				GraceDuration:    25,
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
				GraceDuration:    tt.fields.GraceDuration,
			}
			if tt.name == "throttled_true" {
				if err := th.ThrottleError(tt.args.ocError); err != nil {
					t.Errorf("testing failed : %+v", err)
				}
			}
			if got := th.IsThrottledOrGraced(tt.args.ocError); got != tt.want {
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
		GraceDuration    int
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
				GraceDuration:    0,
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
				GraceDuration:    0,
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
			if tt.name == "default" && !th.IsThrottledOrGraced(tt.args.errObj) {
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
		graceDuration    int
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
				graceDuration:    0,
			},
			want: true,
		},
		{
			name: "Test_isOverThrottleDuration_false",
			args: args{
				cachedTime:       time.Now().Add(1 * time.Minute).Format(time.RFC3339), // 1 minute ahead of current < throtte duration 2
				throttleDuration: 2,
				graceDuration:    0,
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

func Test_isOverGraceDuration(t *testing.T) {
	type args struct {
		cachedTime       string
		throttleDuration int
		graceDuration    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_isOverGraceDuration_true",
			args: args{
				cachedTime:       time.Now().Add(-5 * time.Second).Format(time.RFC3339), // 2 sec after grace duration is over
				throttleDuration: 0,
				graceDuration:    3,
			},
			want: true,
		},
		{
			name: "Test_isOverGraceDuration_false",
			args: args{
				cachedTime:       time.Now().Add(2 * time.Second).Format(time.RFC3339), // still 8 sec left for grace duration
				throttleDuration: 0,
				graceDuration:    10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOverGraceDuration(tt.args.cachedTime, tt.args.graceDuration); got != tt.want {
				t.Errorf("isOverGraceDuration() = %v, want %v", got, tt.want)
			}
		})
	}

}
