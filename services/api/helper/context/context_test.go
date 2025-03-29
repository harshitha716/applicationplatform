package apicontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextVariableManagement(t *testing.T) {
	// Start with an empty context
	ctx := context.Background()

	// Test adding a variable
	ctx = AddCtxVariableToCtx(ctx, "testKey", "testValue")
	value := getCtxVariableFromCtx(ctx, "testKey")
	assert.Equal(t, "testValue", value, "AddCtxVariableToCtx should add a variable to the context")

	// Test adding multiple variables
	ctx = AddCtxVariableToCtx(ctx, "numberKey", 42)
	ctx = AddCtxVariableToCtx(ctx, "boolKey", true)

	// Test getting all context variables
	allVars := GetAllCtxVariablesFromCtx(ctx)
	assert.Len(t, allVars, 3, "Should have 3 variables in context")
	assert.Equal(t, "testValue", allVars["testKey"], "First variable should be correct")
	assert.Equal(t, 42, allVars["numberKey"], "Second variable should be correct")
	assert.Equal(t, true, allVars["boolKey"], "Third variable should be correct")

	// Test removing a variable
	ctx = RemoveCtxVariableFromCtx(ctx, "numberKey")
	remainingVars := GetAllCtxVariablesFromCtx(ctx)
	assert.Len(t, remainingVars, 2, "Should have 2 variables after removal")
	assert.Nil(t, getCtxVariableFromCtx(ctx, "numberKey"), "Removed variable should return nil")

	// Test removing a non-existent variable
	ctx = RemoveCtxVariableFromCtx(ctx, "nonExistentKey")
	assert.Len(t, GetAllCtxVariablesFromCtx(ctx), 2, "Removing non-existent key should not change context")

	// Test getting a non-existent variable
	assert.Nil(t, getCtxVariableFromCtx(ctx, "nonExistentKey"), "Getting non-existent key should return nil")

	// Test context with no variables
	emptyCtx := context.Background()
	assert.Empty(t, GetAllCtxVariablesFromCtx(emptyCtx), "Empty context should return empty map")
	assert.Nil(t, getCtxVariableFromCtx(emptyCtx, "anyKey"), "Getting key from empty context should return nil")
}

func TestContextOverwrite(t *testing.T) {
	// Test overwriting an existing variable
	ctx := context.Background()
	ctx = AddCtxVariableToCtx(ctx, "key", "initialValue")
	ctx = AddCtxVariableToCtx(ctx, "key", "updatedValue")

	value := getCtxVariableFromCtx(ctx, "key")
	assert.Equal(t, "updatedValue", value, "Variable should be overwritten")
}

func TestContextTypes(t *testing.T) {
	// Test storing different types of values
	ctx := context.Background()

	// Primitive types
	ctx = AddCtxVariableToCtx(ctx, "string", "hello")
	ctx = AddCtxVariableToCtx(ctx, "int", 123)
	ctx = AddCtxVariableToCtx(ctx, "float", 3.14)
	ctx = AddCtxVariableToCtx(ctx, "bool", true)

	// Complex types
	type TestStruct struct {
		Name string
		Age  int
	}
	testStruct := TestStruct{Name: "John", Age: 30}
	ctx = AddCtxVariableToCtx(ctx, "struct", testStruct)

	// Slice and map
	testSlice := []string{"a", "b", "c"}
	testMap := map[string]int{"x": 1, "y": 2}
	ctx = AddCtxVariableToCtx(ctx, "slice", testSlice)
	ctx = AddCtxVariableToCtx(ctx, "map", testMap)

	// Verify all types
	assert.Equal(t, "hello", getCtxVariableFromCtx(ctx, "string"))
	assert.Equal(t, 123, getCtxVariableFromCtx(ctx, "int"))
	assert.Equal(t, 3.14, getCtxVariableFromCtx(ctx, "float"))
	assert.Equal(t, true, getCtxVariableFromCtx(ctx, "bool"))
	assert.Equal(t, testStruct, getCtxVariableFromCtx(ctx, "struct"))
	assert.Equal(t, testSlice, getCtxVariableFromCtx(ctx, "slice"))
	assert.Equal(t, testMap, getCtxVariableFromCtx(ctx, "map"))
}
