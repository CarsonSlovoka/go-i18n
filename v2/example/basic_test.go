package example

import "testing"

func TestDemoBasic(t *testing.T) {
	// DemoBasic() // This example is not suitable for the test because it involves interaction with the user. If you want to test, then you can call this function on your project and run it.
}


func TestDemoUseFunc2Render(t *testing.T) {
	demoUseFunc2Render()
}

func TestDemoZeroOneTwoFewManyOther(t *testing.T) {
	demoZeroOneTwoFewManyOther()
}

func TestDemoJSON(t *testing.T) {
	demoJSON()
}
