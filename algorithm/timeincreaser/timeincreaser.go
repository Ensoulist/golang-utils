package timeincreaser

type IncCountType interface {
	~int64 | ~float64 | ~int | ~float32 | ~int32
}

// Automatically increases value by IncCount() at each IncInterval() until reaching Max() limit, or Min() limit if IncCount() is negative.
// GetData() and SetData() specify where the data is stored, typically in the data pointed to by the param parameter.
// The data to be set and retrieved includes: the timestamp of the last settlement and the count at the last settlement.
type ITimeIncreaser[T IncCountType, K IncCountType] interface {
	Max(param any) K
	Min(param any) K
	IncInterval(param any) int64
	IncCount(param any) T

	GetData(param any) (K, int64)
	SetData(count K, ts int64, param any)
}

// The IncreaserGet function calculates the current value of the counter, the timestamp,
// and the incremented count at the given time point 'now'.
// The 'increaser' parameter is an object implementing the ITimeIncreaser interface,
// used to get and set the counter data.
// The 'now' parameter is the current timestamp used for calculating the increment.
// The 'param' parameter is an arbitrary type used to pass additional information.
// The return values include: the current value of the counter, the updated timestamp,
// and the incremented count.
func IncreaserGet[T, K IncCountType](increaser ITimeIncreaser[T, K], now int64, param any) (K, int64, K) {
	baseCount, baseTs := increaser.GetData(param)
	incInterval := increaser.IncInterval(param)
	if incInterval <= 0 {
		return baseCount, baseTs, 0
	}
	incCount := increaser.IncCount(param)

	addRound := (now - baseTs) / incInterval
	if addRound <= 0 {
		return baseCount, baseTs, 0
	}
	addCount := K(T(addRound) * incCount)
	addTs := baseTs + addRound*incInterval
	afterCount := baseCount + addCount

	maxLimit := increaser.Max(param)
	if baseCount > maxLimit {
		maxLimit = baseCount
	}
	if afterCount > maxLimit {
		afterCount = maxLimit
	}
	minLimit := increaser.Min(param)
	if afterCount < minLimit {
		afterCount = minLimit
	}
	return afterCount, addTs, addCount
}

// The IncreaserAdd function is used to add a specified count to the counter.
// The 'increaser' parameter is an object implementing the ITimeIncreaser interface,
// used to get and set the counter data.
// The 'count' parameter is the value to be added to the counter.
// The 'now' parameter is the current timestamp used for calculating the increment.
// The 'dryrun' parameter is a boolean indicating whether to simulate the addition without actually updating the data.
// The 'param' parameter is an arbitrary type used to pass additional information.
// The return values include: the updated counter value, a boolean indicating whether the addition was successful, and the automatically added count.
func IncreaserAdd[T, K IncCountType](increaser ITimeIncreaser[T, K], count K, now int64, dryrun bool, param any) (K, bool, K) {
	minLimit := increaser.Min(param)
	nowCount, addTs, autoAddCount := IncreaserGet(increaser, now, param)
	afterCount := nowCount + count
	if afterCount < minLimit {
		return nowCount, false, 0
	}

	if dryrun {
		return afterCount, true, autoAddCount
	}

	maxLimit := increaser.Max(param)
	if afterCount >= maxLimit || nowCount >= maxLimit {
		addTs = now
	}

	increaser.SetData(afterCount, addTs, param)
	return afterCount, true, autoAddCount
}

// IncreaserSettle function finalizes the counter value at a given timestamp.
// The 'increaser' parameter is an object implementing the ITimeIncreaser interface,
// used to get and set the counter data.
// The 'now' parameter is the current timestamp used for calculating the increment.
// The 'param' parameter is an arbitrary type used to pass additional information.
// The return values include: the finalized counter value, the timestamp of the settlement,
// and the automatically added count.
func IncreaserSettle[T, K IncCountType](increaser ITimeIncreaser[T, K], now int64, param any) (K, int64, K) {
	_, _, autoAddCount := IncreaserAdd(increaser, 0, now, false, param)
	count, ts := increaser.GetData(param)
	return count, ts, autoAddCount
}
