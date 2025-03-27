package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// SecBoolValue 避免nil
func SecBoolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolValue(false)
	}
	return types.BoolValue(*b)
}
