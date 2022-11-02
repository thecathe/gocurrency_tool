package scopemanager

// "github.com/thecathe/gocurrency_tool/analyser/log"

type StackOfIDs []ID

// creates a new empty ScopeIDStack and returns its pointer
func NewStackOfIDs() *StackOfIDs {
	return &StackOfIDs{}
}

// Receiver *ScopeIDStack
// Returns Scope ID from the top of the Stack, -1 if error.
func (stack *StackOfIDs) Peek() (ID, bool) {
	if size := (*stack).Size(); size > 0 {
		return (*stack)[size-1], true
	}

	return ID(""), false
}

// Returns x amount of IDs starting from the top of the stack.
// If x is above the size of the Stack, returns the whole stack.
func (stack *StackOfIDs) PeekX(x int) (*StackOfIDs, bool) {
	// log.DebugLog("ZZ Entering peekX: %d > %d\n", (*stack).Size(), x)
	if (*stack).Size() >= x {
		var _stack StackOfIDs = (*stack)[(*stack).Size()-x:]
		return &_stack, true
	}

	return stack, false
}

// Returns the Scope ID at the given Index, from 0.
func (stack *StackOfIDs) Get(index int) (ID, bool) {
	if (*stack).Size() > index {
		return (*stack)[index], true
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
		(*stack) = (*stack)[:(*stack).Size()-1]
	} else {
		// nothing to remove from stack
		return stack, false
	}
	return stack, true
}

// Receiver *ScopeIDStack
// Returns len() of ScopeIDStack
func (stack *StackOfIDs) Size() int {
	return len((*stack))
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
