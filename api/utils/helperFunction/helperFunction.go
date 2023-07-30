package helperfunction

// Helper function to check if a uint32 value exists in a slice.
func ContainsUint32(slice []uint32, val uint32) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
