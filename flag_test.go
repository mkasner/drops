package drops

import (
	"fmt"
	"testing"
)

func TestFlags(t *testing.T) {
	const (
		Flag1 Flag = 1 << iota
		Flag2
		Flag3
	)

	var flagset Flag
	var flag Flag
	// setting flags
	fmt.Println("Setting flags...")
	flag = Flag1
	flagset = SetFlag(flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag1 %b %t\n", flagset, IsFlag(flagset, flag))
	if !IsFlag(flagset, flag) {
		t.Fatalf("Flag1 not set")
	}

	flag = Flag2
	flagset = SetFlag(flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag2 %b %t\n", flagset, IsFlag(flagset, flag))
	if !IsFlag(flagset, flag) {
		t.Fatalf("Flag2 not set")
	}

	flag = Flag3
	flagset = SetFlag(flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag3 %b %t\n", flagset, IsFlag(flagset, flag))
	if !IsFlag(flagset, flag) {
		t.Fatalf("Flag3 not set")
	}

	fmt.Println("Unsetting flags...")
	flag = Flag3
	flagset = SetFlag(flag)
	flagset = UnsetFlag(flagset, flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag3 %b %t\n", flagset, IsFlag(flagset, flag))
	if IsFlag(flagset, flag) {
		t.Fatalf("Flag3 should be not set")
	}

	flag = Flag1
	flagset = SetFlag(flag)
	flagset = UnsetFlag(flagset, flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag1 %b %t\n", flagset, IsFlag(flagset, flag))
	if IsFlag(flagset, flag) {
		t.Fatalf("Flag1 should be not set")
	}

	// setting multiple flags
	flag = Flag1
	flagset = 0
	flagset = AddFlag(flagset, Flag1, Flag2, Flag3)
	if !IsFlag(flagset, Flag2) {
		t.Fatalf("Flag2 should be set")
	}
	flagset = UnsetFlag(flagset, flag)
	fmt.Println(flagset, flag)
	fmt.Printf("flag1 %b %t\n", flagset, IsFlag(flagset, flag))
	if IsFlag(flagset, flag) {
		t.Fatalf("Flag1 should be not set")
	}

	// setting flag multiple times
	fmt.Println("Setting flags multiple times")
	flag = Flag1
	flagset = 0
	flagset = AddFlag(flagset, flag)
	flagset = AddFlag(flagset, flag)
	if !IsFlag(flagset, flag) {
		t.Fatalf("Flag1 should be set")
	}
	flagset = UnsetFlag(flagset, flag)
	flagset = UnsetFlag(flagset, flag)
	if IsFlag(flagset, flag) {
		t.Fatalf("Flag1 should not be set")
	}

}
