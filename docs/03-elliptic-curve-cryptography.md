# Elliptic Curve Cryptography
So far we have only shown elliptic curves and the point addition on $\R^{2}$. That worked, because $\R$ is a field (you can apply all the properties of a field to it).

The elliptic curve can also be made over a finite field. We can write the set of elements on the elliptic curve on in the real domain as 

$$\{(x,y)\in \R^{2} \mid y^{2} = x^{3} + ax + b, 4a^{3} + 27b^{2} \neq 0 \} \cup \{0\}$$

The elliptic curve restricted over $\mathbb{F}_{p}$ simple becomes

$$\{(x,y)\in \mathbb{F}_{p}^{2} \mid y^{2} = x^{3} + ax + b \mod p, 4a^{3} + 27b^{2} \neq 0 \mod p\} \cup \{0\}$$

Where $0$ is the point at infinity and $a$ and $b$ two integers in $\mathbb{F}_{p}$.

[What previously was a continuous curve is now a set of disjoint points in the $xy$-plane.][ref1]

[ref1]: https://andrea.corbellini.name/2015/05/23/elliptic-curve-cryptography-finite-fields-and-discrete-logarithms/

## Point addition
Point addition in over the finite field is a bit different from point additon over the real domain. Because, what does it mean for three points $P$, $Q$ and $R$ to be aligned in $\mathbb{F}_{p}$.

Three points are said to be aligned if there is a line in $\mathbb{F}_{p}$ that connects them. Since the line in $\mathbb{F}_{p}$ is *wrapped around* $\mathbb{F}_{p}$, 

![Alt text](images/point-addition-mod-p.png)

In point addition over finite fields, the inverse of $P$, or $-P$, is $(x_{P}, -y_{p} \mod p))$.

