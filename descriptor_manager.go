package vulkan

import (
	"sync"
)

// DescriptorPoolManager defines the DescriptorPoolManager type
// DescriptorPoolManager is a high-level utility that dynamically manages
// a growing collection of Vulkan descriptor pools. It eliminates the need
// to manually recreate pools when they run out of memory or get fragmented.
type DescriptorPoolManager struct {
	device         Device
	poolSizes      []DescriptorPoolSize
	maxSetsPerPool uint32
	flags          DescriptorPoolCreateFlags

	currentPool DescriptorPool
	usedPools   []DescriptorPool
	freePools   []DescriptorPool

	mu sync.Mutex
}

// NewDescriptorPoolManager creates a new DescriptorPoolManager
func NewDescriptorPoolManager(device Device, maxSetsPerPool uint32, poolSizes []DescriptorPoolSize, flags DescriptorPoolCreateFlags) (*DescriptorPoolManager, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if maxSetsPerPool == 0 {
		return nil, NewValidationError("maxSetsPerPool", "must be greater than 0")
	}
	if len(poolSizes) == 0 {
		return nil, NewValidationError("poolSizes", "cannot be empty")
	}

	return &DescriptorPoolManager{
		device:         device,
		poolSizes:      poolSizes,
		maxSetsPerPool: maxSetsPerPool,
		flags:          flags,
		usedPools:      make([]DescriptorPool, 0),
		freePools:      make([]DescriptorPool, 0),
	}, nil
}

// grabPool gets a free pool or creates a new one
func (m *DescriptorPoolManager) grabPool() (DescriptorPool, error) {
	if len(m.freePools) > 0 {
		pool := m.freePools[len(m.freePools)-1]
		m.freePools = m.freePools[:len(m.freePools)-1]
		return pool, nil
	}

	// Create a new pool
	pool, err := CreateDescriptorPool(m.device, &DescriptorPoolCreateInfo{
		Flags:     m.flags,
		MaxSets:   m.maxSetsPerPool,
		PoolSizes: m.poolSizes,
	})
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// AllocateDescriptorSets allocates one or more descriptor sets, creating new pools if necessary
func (m *DescriptorPoolManager) AllocateDescriptorSets(layouts []DescriptorSetLayout) ([]DescriptorSet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(layouts) == 0 {
		return nil, nil
	}

	if m.currentPool == nil {
		pool, err := m.grabPool()
		if err != nil {
			return nil, err
		}
		m.currentPool = pool
	}

	allocInfo := &DescriptorSetAllocateInfo{
		DescriptorPool: m.currentPool,
		SetLayouts:     layouts,
	}

	sets, err := AllocateDescriptorSets(m.device, allocInfo)
	// If we run out of memory or the pool gets fragmented, get a new pool and try again
	if err != nil {
		var isPoolError bool
		if vulkanErr, ok := err.(*VulkanError); ok {
			if vulkanErr.Result == ErrorOutOfPoolMemory || vulkanErr.Result == ErrorFragmentedPool {
				isPoolError = true
			}
		}

		if isPoolError {
			m.usedPools = append(m.usedPools, m.currentPool)
			// Clear currentPool immediately so a grabPool failure cannot leave
			// the same pool referenced both here and in usedPools (which would
			// lead to a double vkDestroyDescriptorPool in Destroy).
			m.currentPool = nil

			pool, grabErr := m.grabPool()
			if grabErr != nil {
				return nil, grabErr
			}
			m.currentPool = pool

			allocInfo.DescriptorPool = m.currentPool
			sets, err = AllocateDescriptorSets(m.device, allocInfo)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return sets, nil
}

// Reset resets all used pools and makes them available for reallocation
func (m *DescriptorPoolManager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentPool != nil {
		m.usedPools = append(m.usedPools, m.currentPool)
		m.currentPool = nil
	}

	var lastErr error
	newUsed := make([]DescriptorPool, 0, len(m.usedPools))

	for _, pool := range m.usedPools {
		err := ResetDescriptorPool(m.device, pool)
		if err != nil {
			lastErr = err
			newUsed = append(newUsed, pool)
		} else {
			m.freePools = append(m.freePools, pool)
		}
	}

	// Update used pools to only those that failed to reset
	m.usedPools = newUsed
	return lastErr
}

// Destroy destroys all Vulkan descriptor pools managed by this manager
func (m *DescriptorPoolManager) Destroy() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentPool != nil {
		DestroyDescriptorPool(m.device, m.currentPool)
		m.currentPool = nil
	}

	for _, pool := range m.usedPools {
		DestroyDescriptorPool(m.device, pool)
	}
	m.usedPools = nil

	for _, pool := range m.freePools {
		DestroyDescriptorPool(m.device, pool)
	}
	m.freePools = nil
}
