# Named Flags

## Introduction

This simple package aims to solve a very specific problem: convert an integer to and from a struct with fields of type bool and vice-versa. Or, in other words, store a set of named flags, and translate them to and from an integer and a user-friendly format where you access boolean flag values through named attributes/fields.

You might think that "Named Flags" is a bad or stupid name for this concept, and you might be right. If you can think of a better name then you must be good at naming things, and I'm very happy for you.


## Background

Inspired by the talk [Making a Text Adventure Parser](https://www.youtube.com/watch?v=II3O1CJA-x8) by Evan Wright (see [his GitHub](https://github.com/evancwright/)) wherein he describes storing game object flags in an integer field in a table, where the bit at a given position references a property/capability/attribute of that game object. I had initially implemented something to achieve this kind of conversion in Python leveraging the hyper-flexibility that the language provides, but I chose to abandon doing the project in Python because I wanted to push myself to start learning Go.

If you are reading this at some point in the far-flung future and I have published a text adventure game written in Go, then just know that this repo was where it all began.


## Notes on structure and conversion

The conversion to int populates the bits from right to left, so the first (0th) field of your `struct` will correspond to the 0th power of 2 in the int value, the 1st being 2^1, etc.

E.g. if you have the following struct
```go
type GameObjectFlags struct {Room, Container, Actor, Openable, Open, Lockable, Locked, Visible bool}
```

and you convert an instance of it to an int using `ToInt`, then it will be represented with the following bits:

```
128     64     32       16   8        4     2         1
Visible Locked Lockable Open Openable Actor Container Room
```

## Basic Usage

Given a `struct` of your own making (let's call it `example`) with only `bool` fields, you can instantiate it by calling `namedflags.FromInt[example](7)`. We use the type parameter notation (`[example]`) to tell `FromInt` what type we want to instantiate. If we later want to dump our flags to an integer we can just call `namedflags.ToInt(myExample)`.

```go
import (
	"fmt"

	"github.com/matthew-hoad/namedflags"
)

type example struct (A, B, C bool)

func main() {
	var a example
	var err error
	var b uint

	a, err = namedflags.FromInt[example](7)
	fmt.Println(a)
	// {true true true}

	a.C = false

	// no need to pass the type parameter here
	// because it is inferred from `a`
	b, err = namedflags.ToInt(a)
	fmt.Println(b)
	// 3
}
```


## Example Use Case

Lets say we're tracking which members of an adventuring party are alive.

```go
package lotr

import (
	"github.com/matthew-hoad/namedflags"
)

type Fellowship struct {
	Frodo, Sam, Merry, Pippin, Gandalf,
	Aragorn, Legolas, Gimli, Boromir bool
}

func (f *Fellowship) Moria() {f.Gandalf = false}
func (f *Fellowship) AmonHen() {f.Boromir = false}
func (f *Fellowship) Mirkwood() {f.Gandalf = true}

func main() {
	// Start with a full party
	fellowship, _ := namedflags.FromInt[Fellowship](511)

	fellowship.Moria()
	// fellowship.Gandalf -> false
	
	fellowship.AmonHen()
	// fellowship.Boromir -> false
	
	fellowship.Mirkwood()
	// fellowship.Gandalf -> true
	
	intValue, _ := namedflags.ToInt(fellowship)
	// intValue -> 255
}
```

Next, let's say we're storing the game state in a database using the [GORM package](https://github.com/go-gorm/gorm) where `fellowship` might be just one field in the game state struct. We can implement a field serializer and deserializer for translating between the value in the database (`uint`) and in our code (`Fellowship`). The following is a naïve implementation based on the code from GORM's docs [here](https://gorm.io/docs/serializer.html).

```go
// method for setting the value from the database
// NOTE: here we use the pointer receiver notation because we want to set the value
func (f *Fellowship) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	*f, err = namedflags.FromInt(dbValue.(type))
	return err
}

// method for converting the value to the database
// NOTE: don't need pointer receiver notation here because we are not modifying the value
func (f Fellowship) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	intValue, err = namedflags.ToInt(f)
	return intValue, err
}
```
