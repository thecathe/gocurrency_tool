package scopemanager

type StackOfIDs []ID

// creates a new empty ScopeIDStack and returns its pointer
func NewStackOfIDs() *StackOfIDs {
	return &StackOfIDs{}
}

// Receiver *ScopeIDStack
// Returns Scope ID from the top of the Stack, -1 if error.
func (stack *StackOfIDs) Peek() (ID, bool) {
	if (*stack).Size() >= 0 {
		return (*stack)[(*stack).Size()], true
	}

	return ID(""), false
}

// Receiver *ScopeIDStack
func (stack *StackOfIDs) Push(scope_id ID) *StackOfIDs {
	(*stack) = append((*stack), scope_id)
	return stack
}

// Receiver *ScopeIDStack
func (stack *StackOfIDs) Pop() (*StackOfIDs, bool) {
	if (*stack).Size() >= 0 {
		// remove last id
		(*stack) = (*stack)[:(*stack).Size()]
	} else {
		// nothing to remove from stack
		return stack, false
	}
	return stack, true
}

// Receiver *ScopeIDStack
// Returns len() of ScopeIDStack -1
func (stack *StackOfIDs) Size() int {
	return len((*stack)) - 1
}

// Returns a pointer for a new Stack with the IDs in reverse order.
func (stack *StackOfIDs) Reverse() *StackOfIDs {
	var reversed StackOfIDs = *NewStackOfIDs()

	// Add each from the end of the stack
	for i := 0; i < (*stack).Size(); i++ {
		reversed = append(reversed, (*stack)[(*stack).Size()-i])
	}

	return &reversed
}
