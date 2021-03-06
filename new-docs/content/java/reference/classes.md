---
title: Application Data
---

## Working with Java Types

oso’s Java authorization library lets you write policy rules over Java objects directly.
This document explains how different types of Java objects can be used in oso policies.

<!-- TODO: make below note reference correct doc -->
**NOTE**: More detailed examples of working with application classes can be found in Policy Examples.

### Class Instances

You may pass an instance of any Java class into oso and access its methods
and fields from your policy; see Application Types.

Java instances can be constructed from within an oso policy using the
New operator:

```
new User("alice@example.com")
```

To construct instances of a Java class, the class must be **registered**
using the `registerClass()` method:

```
registerClass(User.class)
```

If you want to refer to the class using another name from within a policy,
you may supply an alias:

```
registerClass(Person.class, "User")
```

At instantiation time, oso will search the list returned by
[Class.getConstructors()](https://docs.oracle.com/javase/10/docs/api/java/lang/Class.html#getConstructors())
for a constructor that is applicable to the supplied positional constructor
arguments. For example, given the Polar expression `new
User("alice@example.com")`, oso will search for a `Constructor` with one
parameter compatible with `String.class`, e.g.:

```
public User(String username) { ... }
```

Applicability is determined using [Class.isAssignableFrom(Class<?> cls)](https://docs.oracle.com/javase/10/docs/api/java/lang/Class.html#isAssignableFrom(java.lang.Class)),
which allows arguments that are instances of subclasses or implementations
of interfaces to properly match the constructor’s parameter types.

<!-- TODO: link to javadoc above -->
### Numbers and Booleans

Polar supports integer and floating point real numbers, as well as booleans (see Primitive Types).

**NOTE**: Java primitives may be passed into oso, but numbers and booleans created in an oso policy will be
converted to [autoboxed](https://docs.oracle.com/javase/tutorial/java/data/autoboxing.html) Integer, Float, and Boolean types respectively.

This means that methods called from oso must have autoboxed argument types. E.g.,

```
class Foo {
   public static unboxed(int a, int b) {
      // ...
   }
   public static boxed(Integer a, Integer b) {
      // ...
   }
}
```

The `boxed()` method may be called from a policy,
but attempting to call `unboxed()` will fail.

### Strings

Java Strings are mapped to Polar Strings. Java’s String methods may be accessed from policies:

```
allow(actor, action, resource) if actor.username.endsWith("example.com");
```

```
public class User {
   public String username;

   public User(String username) {
      this.username = username;
   }

   public static void main(String[] args) {
      User user = new User("alice@example.com");
      assert oso.isAllowed(user, "foo", "bar");
   }
}
```

### Lists and Arrays

Java [Arrays](https://docs.oracle.com/javase/tutorial/java/nutsandbolts/arrays.html) *and* objects that implement the [List](https://docs.oracle.com/javase/10/docs/api/java/util/List.html) interface are
mapped to Polar Lists. Java’s `List` methods may be accessed from policies:

```
allow(actor, action, resource) if actor.groups.contains("HR");
```

```
public class User {
   public List<String> groups;

   public User(List<String> groups) {
      this.groups = groups;
   }

   public static void main(String[] args) {
      User user = new User(List.of("HR", "payroll"));
      assert oso.isAllowed(user, "foo", "bar");
   }
}
```

Note that the `isAllowed()` call would also succeed if `groups` were an array.

**WARNING**: Polar does not support methods that mutate lists in place. E.g. `add()` will have no effect on
a list in Polar.

Likewise, lists constructed in Polar may be passed into Java methods:

```
allow(actor, action, resource) if actor.has_groups(["HR", "payroll"]);
```

```
   public boolean hasGroups(List<String> groups) {
      for(String g : groups) {
         if (!this.groups.contains(g))
            return false;
      }
      return true;
   }

   public static void main(String[] args) {
      User user = new User(List.of("HR", "payroll"));
      assert oso.isAllowed(user, "foo", "bar");
   }
```

Java methods like [`List.get`](https://docs.oracle.com/javase/10/docs/api/java/util/List.html#get(int)) may be used for random access to
list elements, but there is currently no Polar syntax for that is
equivalent to the Java expression `user.groups[1]`. To access
the elements of a list without using a method, you may iterate
over it with In (List Membership) or destructure it with
pattern matching.

### Maps

Java objects that implement the [Map](https://docs.oracle.com/javase/10/docs/api/java/util/Map.html) interface
are mapped to Polar Dictionaries:

```
allow(actor, action, resource) if actor.roles.project1 = "admin";
```

```
public class User {
   public Map<String, String> roles;

   public User(Map<String, String> roles) {
      this.roles = roles;
   }

   public static void main(String[] args) {
      User user = new User(Map.of("project1", "admin"));
      assert oso.isAllowed(user, "foo", "bar");
   }
}
```

Likewise, dictionaries constructed in Polar may be passed into Java methods.

### Enumerations

You may iterate over a Java [Enumeration](https://docs.oracle.com/javase/10/docs/api/java/util/Enumeration.html)
(or anything that can be converted to one, such as a `Collection` or `Iterable`)
using the Polar In (List Membership) operator:

```
allow(actor, action, resource) if "payroll" in actor.getGroups();
```

```
   public List<String> getGroups() {
      return List.of("HR", "payroll");
   }

   public static void main(String[] args) {
      User user = new User(Map.of("project1", "admin"));
      assert oso.isAllowed(user, "foo", "bar");
   }
```

### `null`

The Java `null` reference is registered as the Polar constant nil.
If a Java method can return `null`, you may want to compare the result
to `nil`:

```
allow(actor, action, resource) if actor.getOptional() != nil;
```

```
public class User {
    ...

    public Thing getOptional() {
        if someCondition() {
            return new Thing();
        } else {
            return null;
        }
    }
 }
```

<!-- ### Summary

### Java → Polar Types Summary

| Java type

 | Polar type

 |     |
 | --- ||  |  |  |  |  ||  |  |  |
| int/Integer

                                                    | Integer

                                                                                 |
| float/Float

                                                    | Float

                                                                                   |
| double/Double

                                                  | Float

                                                                                   |
| boolean/Boolean

                                                | Boolean

                                                                                 |
| List

                                                           | List

                                                                                    |
| Array

                                                          | List

                                                                                    |
| Map

                                                            | Dictionary

                                                                              |
| String

                                                         | String

                                                                                  | -->
