package ex00

import "testing"

func TestGetElement(t *testing.T) {
	var arr = []int{10, 20, 30, 40, 50}

	val, err := getElement(arr, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != 30 {
		t.Errorf("Expected element at index 2 to be 30, got %d", val)
	}

	_, err = getElement(arr, 5)
	expectedErr := "index out of bounds"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got %v", expectedErr, err)
	}

	_, err = getElement([]int{}, 0)
	expectedErr = "empty slice"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got %v", expectedErr, err)
	}

	_, err = getElement(arr, -1)
	expectedErr = "index out of bounds"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got %v", expectedErr, err)
	}
}
