# Finite fields

## Definition
A finite fields is defined as a finite set of elements and two operations $+$ (addition) and &bull; (multiplication) that satisfy the following:
1. If $a$ and $b$ are in the set, $a + b$ and $a$ &bull; $b$ are in the set. (*Closed property*) 
2. $0$ exists and has the property $a + 0 = a$. (*Additive identity*)
3. $1$ exists and has the property $a \cdot 1 = a$. (*Multiplicative identity* )
4. If $a$ is in the set, $-a$ is in the set, which is defined as the value that makes $a + (-a) = 0$. (*Additive inverse*)
5. If $a$ is in the set and is not $0$, $a^{-1}$ is in the set, which is defined as the value that makes $a \cdot a^{-1} = 1$. (*Multiplicative inverse* )

For example, the set $\{-1, 0, 1\}$ is closed under normal multiplication. However, it is not closed under normal addition.

To make these sets closed, we can define addition and multiplication in a particular way to make these sets closed.

## Finite field sets
A finite field set is a set $$F_{p}=\{0, 1, 2, ..., p-1\}$$ where $p$ is the *order* or *size* of the finite field set. (Here, the numbers 0, 1, 2 etc stand for the elements in the field, but are not necessarily their natural numbers.) Some examples of finite fields are $F_{11}$, $F_{17}$ and $F_{983}$. Note, that $p$ is a prime number in each of these examples. Later, it will become clear why.

## Defining the arithmatic 
Addition and subtraction in finite field $F_{p}$ will be defined as $$a +_{f} b = (a + b) \mod p$$ and $$a -_{f} b = (a - b) \mod p$$ where $a, b \in F_{p}$. This makes addition and subtraction of finite field elements closed.

Multiplication in finite field $F_{p}$ will be defined as $$a \cdot_{f} b = (a \cdot b) \mod p$$. This makes multiplication of finite field elements closed.

### Division
Addition, subtraction, multiplication and exponentiation of field elements still feel quite natural because it is not much different from their operations with natural numbers. However, the intuition that we have from working with natural numbers doesn't help with understanding division of field elements. Fermat's Little Theorem is used to compute division.

## Fermat's Little Theorem
### Theorem
Fermat's Little Theorem is a fundamental result in number theory that states the following:

If $p$ is a prime number and $a$ is an integer not divisible by $p$, then $a^{(p-1)}$ is congruent to $1$ modulo $p$, which can be written as:

$$a^{(p-1)} ≡ 1 \mod p$$

### Proof

TODO

### Application in finite field division
Using Fermat's Little Theorem, we know $$b^{(p-1)}=1$$ because $p$ is prime. Thus: $$b^{-1}=b^{-1}\cdot _f 1=b^{-1}\cdot _f b^{(p-1)} = b^{(p-2)}$$ or $$b^{-1}=b^{(p-2)}$$

In $F_{19}$ this means that $b^{18}=1$, which means that $b^{-1}=b^{17}$ for all $b > 0$.

So in $F_{19}$:
$$2/7 = 2 \cdot 7^{19-2} = 2 \cdot 7^{17} = 3$$

## Extended Euclidean Algorithm
Another way to look at division is using the extended Euclidean algorithm, which is used to compute the greatest common divider of two integers $a$ and $b$. It also gives $x$ and $y$ such that $$a\cdot x + b\cdot y = \gcd(a,b)$$. 

If $n \in F _{p}$, then $\gcd(n, p) = 1$ because $p$ is prime. Plugging this into the extended Euclidean algorithm gives
$x\cdot n + y\cdot p = 1$. Since we are using modulo arithmetic here, we can write 

$$x\cdot n = 1\mod p$$

Going back to our original question, we wanted to know the multiplicative inverse of $n^{-1}\in F _{p}$ such that $n \cdot n^{-1} = 1$. As you can see, the extend Euclidean algorithm gives exactly that number, $n^{-1} = x \mod p$.

