---
date: '2021-01-07T02:46:33.217Z'
docname: using/examples/inheritance
images: {}
path: /using-examples-inheritance
title: Resources with Inheritance
---

# Resources with Inheritance

Some applications have many resources that should have similar authorization
rules applied to them.  This is a common scenario in workflow driven
applications that have different user types and a large number of resources.

A common set of rules will apply to many resources, with some exceptions.  In
this guide, we will cover various ways for modeling this scenario with oso.

## Setup

In this example, we will consider a hypothetical EMR (electronic medical record)
application.  We’ll discuss a few resources:


* **Order**: A record of an action that medical staff will perform on a patient


* **Test**: A diagnostic test that will be performed on a patient.


* **Lab**: A lab test that will be performed on a patient.

These resources are all examples of different types of patient data.

## Basic Policy

Let’s start by considering a basic policy controlling access to these three
resources:

```
allow(actor: Actor, "read", resource: Order) if
    actor.role = "medical_staff" and
    actor.treated(resource.patient);

allow(actor: Actor, "read", resource: Test) if
    actor.role = "medical_staff" and
    actor.treated(resource.patient);

allow(actor: Actor, "read", resource: Lab) if
    actor.role = "medical_staff" and
    actor.treated(resource.patient);
```

Let’s take a look at the first rule in the policy. This `allow` rule permits
an actor to perform the `"read"` action on an `Order` if:


1. The actor’s `role` property is equal to `"medical_staff"`.


2. The application method `treated(patient)`

    on the actor returns true for the patient of the `resource`.

Note the head of the rule.  Each argument uses a type specializer to
ensure this rule only applies to certain types of resources and actor.  This
rule indicates that the `actor` parameter must be an instance of the `Actor`
class and the `resource` parameter must be an instance of the `Order` class.

This policy meets our goal above. We have expressed the same rule for the three
types of patient data, but it is a bit repetitive.  Let’s try to improve it.

## Using a Rule to Express Common Behavior

Our policy doesn’t just need to contain `allow` rules.  We can write any
rules we’d like, and compose them as needed to express our policy!

```
can_read_patient_data(actor, "read", resource) if
    actor.role = "medical_staff" and
    actor.treated(resource.patient);

allow(actor: Actor, "read", resource) if
    can_read_patient_data(actor, "read", resource);
```

Now, we’ve taken the repeated logic and expressed it as the
`can_read_patient_data` rule.  When the `allow` rule is evaluated,
it will check if `can_read_patient_data` is true. The policy is much shorter!

Unfortunately, we’ve lost one property of our last policy: the specializers.
This rule would be evaluated for any type of resource, not just our three
examples of patient data above. That’s not what we want.

## Bringing Back Specializers

We can combine this idea with our first policy to make sure only our three
patient data resources use the `can_read_patient_data` rule.

```
allow(actor: Actor, "read", resource: Order) if
    can_read_patient_data(actor, "read", resource);

allow(actor: Actor, "read", resource: Test) if
    can_read_patient_data(actor, "read", resource);

allow(actor: Actor, "read", resource: Lab) if
    can_read_patient_data(actor, "read", resource);
```

Now, we still have three rules, but the body isn’t repeated anymore.

## One Rule to Rule Them All

We haven’t talked about the application side of this yet.  So far, we’ve assumed
`Order`, `Lab`, and `Test` are application classes.

Here’s how our application classes might be implemented:

Python

```
@oso.polar_class
class PatientData:
    def __init__(self, patient):
        self.patient = patient


@oso.polar_class
class Lab(PatientData):
    pass


@oso.polar_class
class Order(PatientData):
    pass


@oso.polar_class
class Test(PatientData):
    pass
```

Ruby

```
class PatientData
  attr_accessor :patient

  def initialize(patient:)
    @patient = patient
  end
end

OSO.register_class(PatientData)

class Lab < PatientData
end

OSO.register_class(Lab)

class Order < PatientData
end

OSO.register_class(Order)

class Test < PatientData
end

OSO.register_class(Test)
```

Java

Java example coming soon.

Node.js

```
class PatientData {
  constructor(patient) {
    this.patient = patient;
  }
}

oso.registerClass(PatientData);

class Lab extends PatientData {}

oso.registerClass(Lab);

class Order extends PatientData {}

oso.registerClass(Order);

class Test extends PatientData {}

oso.registerClass(Test);
```

We used inheritance to capture some of the common
functionality needed (storing the patient).  In a real application these
would probably be ORM models.

We can use the same idea to shorten our policy even further!

```
allow(actor: Actor, "read", resource: PatientData) if
    actor.role = "medical_staff" and
    actor.treated(resource.patient);
```

Now, this `allow` rule will be evaluated for any instance that is a subclass
of `PatientData`.  Polar understands the class inheritance structure when
selecting rules to evaluate!

<!-- TODO: include when groups are back
Working with groups
-------------------

This worked well for us, but remember this is just an example.  Not all
applications may have encoded relationships this way.  Maybe when we wrote our
code we didn't create a ``PatientData`` class, and just implemented ``Lab``,
``Order`` and ``Test`` separately.  We still want to treat them as one concept
in our policy, but don't want to change our application.

Polar includes a ``group`` construct that can be used for exactly this purpose:

.. literalinclude:: /examples/inheritance/05-group.polar
   :language: polar
   :emphasize-lines: 1

The :ref:`group` declaration creates a new type in our Polar file called
``PatientData``.  This time, we tell Polar that ``Lab``, ``Order`` and ``Test``
are part of this group.  We can write our rule in the same way as before. -->
## Summary

In this guide, we saw an example of an application policy that could result in
significant repetition.   We tried out a few strategies for representing common
policy, but using it across many resource types.  First, we wrote a custom rule
that moved duplicated logic into one place.  Then we used specializers and
application types to condense our policy even further.

<!-- TODO
Finally, we saw how groups & inheritance can both be exploited to
write flexible policies that accurately model our application's domain model. -->
