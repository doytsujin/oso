---
title: Quickstart
description: |
    Ready to get started? See oso in action, and walk through our quick tutorial
    for adding authorization to a simple web server.
weight: 1
---

## oso in 5 minutes

oso helps developers build authorization into their applications. If you’ve
never used oso before and want to see it in action, this guide is for you.
We’re going to walk through how to use oso to add authorization to a simple web
server.

{{< callout "Try it!" "green" >}}
To follow along, clone the {{% exampleGet "githubApp" %}}:

```bash
git clone {{% exampleGet "githubURL" %}}
```
{{< /callout >}}

## Run the server

Our sample application serves data about expenses submitted by users. The
sample application has three important files.

One file defines a simple `Expense` class and some sample data stored in a
map.

A second file has our HTTP server code, where we have defined a route handler
for `GET` requests to the path `/expenses/:id`. We’ve already added an
authorization check using the oso library to
control access to expense resources. You can learn more about how to add oso to
your application here.

The third file is the oso policy file, `expenses.polar`, and is currently
empty.

{{< callout "Try it!" "green" >}}

{{% exampleGet "installation" %}}

With the server running, open a
second terminal and make a request using cURL:

```bash
$ curl localhost:5050/expenses/1
Not Authorized!
```

You’ll get a “Not Authorized!” response because we haven’t added any rules to
our oso policy (in `expenses.polar`), and oso is deny-by-default.

{{< /callout >}}

Let’s start implementing our access control scheme by adding some rules to the
oso policy.

## Adding our first rule

oso rules are written in a declarative policy language called Polar. You can
include any kind of rule in a policy, but the oso library is designed to
evaluate allow rules, which specify the conditions that
allow an **actor** to perform an **action** on a **resource**.

{{< callout "Edit it!" "blue" >}}

In our policy file (`expenses.polar`), let's add a rule that allows
anyone with an email ending in `"@example.com"` to view all expenses:

```prolog
allow(actor: String, "GET", _expense: Expense) if
    actor.{{< exampleGet "endswith" >}}("@example.com");
```

Note that the call to {{< exampleGet "endswith" >}} is actually calling out to {{< exampleGet "endswithURL" >}} The actor value passed to oso is a string, and oso allows us to call methods on it.

{{< /callout >}}

The `Expense` and `String` terms following the colons in the head of the
rule are specializers, patterns that control rule
execution based on whether they match the supplied argument. This syntax
ensures that the rule will only be evaluated when the actor is a string and
the resource is an instance of the `Expense` class.

{{< callout "Try it!" "green" >}}
Once we've added our new rule and restarted the web server, every user with
an ``@example.com`` email should be allowed to view any expense:

```bash
$ curl -H "user: alice@example.com" localhost:5050/expenses/1
Expense(...)
```
{{< /callout >}}


Okay, so what just happened?

When we ask oso for a policy decision via `Oso.is_allowed()`, the oso engine
searches through its knowledge base to determine whether the provided
**actor**, **action**, and **resource** satisfy any **allow** rules. In the
above case, we passed in `"alice@example.com"` as the **actor**, `"GET"` as
the **action**, and the `Expense` object with `id=1` as the **resource**.
Since `"alice@example.com"` ends with `@example.com`, our rule is
satisfied, and Alice is allowed to view the requested expense.

{{< callout "Try it!" "green" >}}

If a user's email doesn't end in `"@example.com"`, the rule fails, and
they are denied access:

```bash
$ curl -H "user: alice@foo.com" localhost:5050/expenses/1
Not Authorized!
```
If you aren’t seeing the same thing, make sure you created your policy
correctly in `expenses.polar`.
{{< /callout >}}


## Using application data

We now have some basic access control in place, but we can do better.
Currently, anyone with an email ending in `@example.com` can see all expenses
– including expenses submitted by others.

{{< callout "Edit it!" "blue" >}}
Let's modify our existing rule such that users can only see their own
expenses:

```prolog
allow(actor: String, "GET", expense: Expense) if
    expense.{{< exampleGet "submitted_by" "submitted_by" >}} = actor;
```
{{< /callout >}}


Behind the scenes, oso looks up the `submitted_by` field on the provided
`Expense` instance and compares that value against the provided **actor**.
And just like that, an actor can only see an expense if they submitted it!

{{< callout "Try it!" "green" >}}

Alice can see her own expenses but not Bhavik's:

```bash
$ curl -H "user: alice@example.com" localhost:5050/expenses/1
Expense(...)
```

```bash
$ curl -H "user: alice@example.com" localhost:5050/expenses/3
Not Authorized!
```
{{< /callout >}}


Feel free to play around with the current policy and experiment with adding
your own rules!

For example, if you have `Expense` and `User` classes defined in your
application, you could write a policy rule in oso that says a `User` may
`"approve"` an `Expense` if they manage the `User` who submitted the
expense and the expense’s amount is less than $100.00:

```prolog
allow(approver: User, "approve", expense: Expense) if
    approver = expense.{{% exampleGet "submitted_by" "submitted_by" %}}.manager
    and expense.amount < 10000;
```

In the process of evaluating that rule, the oso engine would call back into the
application in order to make determinations that rely on application data, such
as:


* Which user submitted the expense in question?
* Who is their manager?
* Is their manager the user who’s attempting to approve the expense?
* Does the expense’s `amount` field contain a value less than $100.00?

For more on leveraging application data in an oso policy, check out
[Application Types](application-types).

