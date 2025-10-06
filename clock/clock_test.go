package clock

import (
	"testing"
	"time"
)

func TestRealClocker_Now(t *testing.T) {
	t.Parallel()

	clocker := RealClocker{}
	
	// 現在時刻を取得
	now1 := clocker.Now()
	time.Sleep(1 * time.Millisecond) // 少し待機
	now2 := clocker.Now()
	
	// 2回の呼び出しで異なる時刻が返されることを確認
	if now1.Equal(now2) {
		t.Errorf("RealClocker.Now() returned same time twice: %v", now1)
	}
	
	// 現在時刻に近い時刻が返されることを確認（1秒以内）
	expected := time.Now()
	diff := now2.Sub(expected)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("RealClocker.Now() returned time too far from now: got %v, expected around %v", now2, expected)
	}
}

func TestFixedClocker_Now(t *testing.T) {
	t.Parallel()

	clocker := FixedClocker{}
	
	// 固定時刻を取得
	now1 := clocker.Now()
	time.Sleep(1 * time.Millisecond) // 少し待機
	now2 := clocker.Now()
	
	// 2回の呼び出しで同じ時刻が返されることを確認
	if !now1.Equal(now2) {
		t.Errorf("FixedClocker.Now() returned different times: %v, %v", now1, now2)
	}
	
	// 期待される固定時刻と一致することを確認
	expected := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	if !now1.Equal(expected) {
		t.Errorf("FixedClocker.Now() got %v, want %v", now1, expected)
	}
}

func TestFixedClocker_WithCustomTime(t *testing.T) {
	t.Parallel()

	// カスタム時刻を設定
	customTime := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	clocker := FixedClocker{Time: customTime}
	
	// カスタム時刻が返されることを確認
	now := clocker.Now()
	// FixedClockerは常に固定時刻を返すため、Timeフィールドは無視される
	expected := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	if !now.Equal(expected) {
		t.Errorf("FixedClocker.Now() with custom time got %v, want %v", now, expected)
	}
}

func TestClocker_Interface(t *testing.T) {
	t.Parallel()

	// RealClockerがClockerインターフェースを実装していることを確認
	var _ Clocker = RealClocker{}
	
	// FixedClockerがClockerインターフェースを実装していることを確認
	var _ Clocker = FixedClocker{}
}

func TestFixedClocker_Consistency(t *testing.T) {
	t.Parallel()

	clocker := FixedClocker{}
	
	// 複数回呼び出しても同じ時刻が返されることを確認
	times := make([]time.Time, 10)
	for i := 0; i < 10; i++ {
		times[i] = clocker.Now()
	}
	
	// すべての時刻が同じであることを確認
	for i := 1; i < len(times); i++ {
		if !times[0].Equal(times[i]) {
			t.Errorf("FixedClocker.Now() returned different times at call %d: %v vs %v", i, times[0], times[i])
		}
	}
}

func TestRealClocker_TimeProgression(t *testing.T) {
	t.Parallel()

	clocker := RealClocker{}
	
	// 複数回呼び出して時刻が進行することを確認
	times := make([]time.Time, 5)
	for i := 0; i < 5; i++ {
		times[i] = clocker.Now()
		time.Sleep(1 * time.Millisecond)
	}
	
	// 時刻が進行していることを確認
	for i := 1; i < len(times); i++ {
		if !times[i].After(times[i-1]) {
			t.Errorf("RealClocker.Now() time did not progress: %v should be after %v", times[i], times[i-1])
		}
	}
}
