# Paseto

> PASETO (Platform-Agnostic SEcurity TOkens) is a specification and reference implementation for secure stateless
> tokens.

## Key Differences between PASETO and JWT

Unlike JSON Web Tokens (JWT), which gives developers more than enough rope with which to hang themselves, PASETO only
allows secure operations. JWT gives you "algorithm agility", PASETO gives you "versioned protocols". It's incredibly
unlikely that you'll be able to use PASETO in an insecure way.

> Caution: Neither JWT nor PASETO were designed for stateless session management. PASETO is suitable for tamper-proof
> cookies, but cannot prevent replay attacks by itself.