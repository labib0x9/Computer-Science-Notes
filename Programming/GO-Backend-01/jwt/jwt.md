# JWT From First Principles

## 1. The problem: HTTP is stateless

Every HTTP request stands alone. The server has no memory of the previous request unless you build memory in. So after a user logs in with a username/password, how does the server know that the _next_ request — "GET /orders" — comes from that same authenticated user?

You need some token the client sends on every request that proves "I am the user who logged in." Two competing approaches exist to build that token: **sessions** and **JWT**. JWT didn't replace sessions because sessions are "wrong" — it exists because sessions have a specific weakness that matters at scale.

---

## 2. Session-based auth (the baseline)

```
1. POST /login {user, pass}
2. Server verifies credentials
3. Server creates a session record:
      session_id: "abc123" -> { user_id: 42, expires: ... }
   stored in Redis/DB/memory
4. Server responds:
      Set-Cookie: session_id=abc123; HttpOnly; Secure
5. Browser stores cookie, sends it automatically on every request
6. Server on each request:
      session = db.lookup("abc123")
      if session valid -> user_id = 42
```

The cookie itself is just a random opaque string. It carries **no information** — it's a pointer. The server must look it up in shared storage on every single request.

This is called **stateful authentication**: the source of truth for "who is this and are they still logged in" lives on the server.

**Consequence:** any server instance handling a request needs access to that shared session store. Fine for a monolith with one Redis. Gets awkward when you have many services (auth service, orders service, payments service) that each need to know "is this user logged in and who are they" — now they all need to hit the same session store, or you build a whole internal "verify this session" RPC.

---

## 3. Enter JWT: stateless authentication

The idea: instead of a random pointer to server-side data, hand the client a token that **contains the data itself**, cryptographically signed so it can't be tampered with.

```
1. POST /login {user, pass}
2. Server verifies credentials
3. Server creates a JWT:
      { user_id: 42, role: "admin", exp: 1720540800 }
   and SIGNS it with a secret key
4. Server responds:
      Authorization: Bearer eyJhbGciOi...
5. Client stores token, sends it in Authorization header
6. Server on each request:
      verify signature on the token (no DB lookup)
      if valid -> trust the claims inside it directly
```

This is **stateless authentication**: any server that knows the signing key (or public key) can verify the token and read the claims, with **zero shared storage, zero lookup**. This is the entire reason JWT exists — it lets you verify identity across many independent services without a shared session store.

That's the trade you're making:

- Session: small opaque token, server does a lookup, easy to revoke instantly, server owns the truth.
- JWT: self-contained token, server does math (crypto verify) instead of a lookup, hard to revoke early, the token itself is the truth until it expires.

Keep that tension in your head — nearly every JWT "gotcha" later (revocation, refresh tokens, short expiry) is the ecosystem re-inventing partial statefulness to compensate for what JWT gave up.

---

## 4. Anatomy of a JWT

A JWT is three Base64URL-encoded segments joined by dots:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MiIsImV4cCI6MTcyMDU0MDgwMH0.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ

┌──────── header ────────┐┌──────────── payload ────────────┐┌──── signature ────┐
eyJhbGciOiJIUzI1NiIs...    .eyJzdWIiOiI0MiIsImV4cCI6...       .TJVA95OrM7E2cB...
```

`header.payload.signature` — nothing more. Each of the first two segments is just **JSON, Base64URL-encoded** (not encrypted, not hashed — just encoded, like base64 of a string). Anyone can decode and read them. That's a critical point people miss: **a JWT is not secret**. It's tamper-evident, not confidential.

### 4.1 Header

```json
{ "alg": "HS256", "typ": "JWT" }
```

Says which algorithm was used to sign it, so the verifier knows how to check the signature.

### 4.2 Payload (claims)

```json
{
  "sub": "42",
  "role": "admin",
  "iat": 1720537200,
  "exp": 1720540800
}
```

This is the actual data — the "claims" about the user. Standard claim names in section 9.

### 4.3 Signature

```
signature = HMAC-SHA256(
    base64url(header) + "." + base64url(payload),
    secret_key
)
```

The signature is computed **over the exact bytes** of the first two encoded segments. This is the whole security mechanism: if anyone changes a single character of the header or payload, the signature no longer matches, and verification fails.

---

## 5. How it's created — step by step

Say the server wants to issue a token for user 42.

```
1. header  = {"alg":"HS256","typ":"JWT"}
2. payload = {"sub":"42","role":"admin","iat":1720537200,"exp":1720540800}

3. h = base64url(json(header))     -> "eyJhbGciOiJIUzI1NiIs..."
4. p = base64url(json(payload))    -> "eyJzdWIiOiI0MiIsImV4cCI6..."

5. signingInput = h + "." + p

6. sig = HMAC_SHA256(signingInput, secretKey)
7. s = base64url(sig)

8. token = h + "." + p + "." + s
```

That's it. No magic. Anyone with the secret key can run this same algorithm.

---

## 6. How it's verified — step by step

```
Given: incoming token string T, and the same secret key

1. split T into h, p, s  (by the dots)
2. recompute: expectedSig = base64url(HMAC_SHA256(h + "." + p, secretKey))
3. compare expectedSig to s using a CONSTANT-TIME comparison
      (never use ==, use crypto/subtle.ConstantTimeCompare in Go —
       naive == leaks timing info about how many bytes matched,
       which is a real (if slow) attack vector)
4. if mismatch -> reject, 401
5. if match -> decode payload p, check registered claims:
      - exp: is now < exp?  if not -> expired, reject
      - nbf: is now >= nbf? if not -> not yet valid, reject
      - iss/aud: do they match what THIS service expects?
6. if all checks pass -> trust payload, proceed as that user
```

Go-flavored pseudocode:

```go
func VerifyToken(token string, secret []byte) (Claims, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return Claims{}, ErrMalformed
    }
    header, payload, sig := parts[0], parts[1], parts[2]

    expected := hmacSHA256(header+"."+payload, secret)
    given := base64URLDecode(sig)

    if !hmac.Equal(expected, given) { // constant time
        return Claims{}, ErrBadSignature
    }

    var claims Claims
    json.Unmarshal(base64URLDecode(payload), &claims)

    if time.Now().Unix() > claims.Exp {
        return Claims{}, ErrExpired
    }
    if claims.Nbf != 0 && time.Now().Unix() < claims.Nbf {
        return Claims{}, ErrNotYetValid
    }
    if claims.Iss != expectedIssuer {
        return Claims{}, ErrBadIssuer
    }

    return claims, nil
}
```

Notice: **no database call anywhere in this function.** That's the entire value proposition of JWT in one code block.

---

## 7. Signing vs encryption — a distinction people blur constantly

- **Signing (JWS — JSON Web Signature)** — what 99% of "JWT" actually means. Payload is **readable by anyone**, but **tamper-proof**. Goal: integrity + authenticity, not secrecy.
- **Encryption (JWE — JSON Web Encryption)** — a different, less common JWT variant where the payload itself is encrypted, so only holders of the decryption key can read it. Structurally has 5 segments instead of 3.

**Rule of thumb:** never put secrets (passwords, SSNs, raw card numbers) in a signed-only JWT payload — it's base64, not encryption, anyone can read it with `atob()`. If you truly need confidentiality in the token, use JWE, or better — don't put sensitive data in the token at all, just put an ID and look the sensitive data up server-side when needed.

---

## 8. HS256 vs RS256 vs ES256

This is really a question of **symmetric vs asymmetric** signing.

### HS256 (HMAC + SHA-256) — symmetric

- **Same secret key** signs and verifies.
- Fast, simple.
- **Problem:** every service that needs to _verify_ tokens also has the power to _create_ them (it holds the same secret). Fine if you have one backend. Dangerous if you hand that key to multiple services or third parties — any one of them can now mint tokens impersonating any user.

```
sign:   HMAC(payload, secretKey)      -- needs secretKey
verify: HMAC(payload, secretKey)      -- needs SAME secretKey
```

### RS256 (RSA + SHA-256) — asymmetric

- **Private key signs, public key verifies.**
- The auth service holds the private key (keeps it secret). Every other microservice gets only the **public key** and can verify tokens but can never forge one.
- This is the natural fit for a distributed-systems setup like yours: auth-service signs, orders-service/payments-service/etc. each just embed the public key (or fetch it from a JWKS endpoint) and verify locally — no network call to auth-service needed per request.
- Slower than HMAC (RSA math is heavier), tokens are a bit bigger.

```
sign:   RSA_private_key.sign(payload)
verify: RSA_public_key.verify(payload, signature)   -- can't forge with only this
```

### ES256 (ECDSA + SHA-256, P-256 curve) — asymmetric, elliptic curve

- Same private/public split as RS256, but ECDSA keys are much smaller than RSA for equivalent security (256-bit EC key ≈ 3072-bit RSA key), so tokens are smaller and signing/verifying is faster.
- Preferred over RS256 in most new systems for this reason; RS256 mainly survives because of legacy/compatibility.

**Practical guidance:**
| | HS256 | RS256 | ES256 |
|---|---|---|---|
| Key type | 1 shared secret | keypair | keypair |
| Who can verify | anyone with the secret (= anyone who can also forge) | anyone with public key, cannot forge | anyone with public key, cannot forge |
| Good for | single backend, no external verifiers | multi-service / multi-party verification | same as RS256, smaller/faster, modern default |
| Watch out for | leaking the secret anywhere = game over | key management overhead | slightly less library support in old ecosystems |

Given you're building things like a toy Raft node / multi-service systems, **RS256 or ES256 is the right mental model** the moment more than one service needs to verify tokens independently — that's the whole reason asymmetric signing exists.

---

## 9. Standard (registered) claims

These are optional but standardized field names — using them lets every JWT library understand your token the same way.

| Claim | Meaning                                                  | Example            |
| ----- | -------------------------------------------------------- | ------------------ |
| `iss` | Issuer — who created/signed this token                   | `"auth.myapp.com"` |
| `sub` | Subject — who the token is about (usually user ID)       | `"42"`             |
| `aud` | Audience — who the token is intended for                 | `"orders-api"`     |
| `exp` | Expiration — Unix timestamp after which token is invalid | `1720540800`       |
| `iat` | Issued At — when the token was created                   | `1720537200`       |
| `nbf` | Not Before — token invalid until this time               | `1720537200`       |
| `jti` | JWT ID — unique ID for this specific token               | `"a1b2c3-..."`     |

`aud` and `iss` matter more than people initially think in a microservices setup: if `orders-service` only accepts tokens where `aud == "orders-service"`, a token minted for `chat-service` can't be replayed against `orders-service` even though it's cryptographically valid — this stops **token confusion** across services that all trust the same signing key.

`jti` is what makes revocation and refresh-token tracking possible (section 12) — it gives you something unique to blacklist or store, without needing to store the whole token.

---

## 10. Complete login/authentication flow

```
┌────────┐                                   ┌────────┐
│ Client │                                   │ Server │
└───┬────┘                                   └───┬────┘
    │  POST /login {email, password}             │
    │────────────────────────────────────────────>│
    │                                              │  verify credentials
    │                                              │  vs DB (bcrypt compare)
    │                                              │
    │                                              │  build claims:
    │                                              │  {sub, iss, aud, iat, exp}
    │                                              │  sign -> access_token
    │                                              │  sign -> refresh_token
    │                                              │  (see section 11)
    │  200 OK                                      │
    │  { access_token, refresh_token }             │
    │<─────────────────────────────────────────────│
    │                                              │
    │  GET /orders                                 │
    │  Authorization: Bearer <access_token>        │
    │────────────────────────────────────────────>│
    │                                              │  verify signature
    │                                              │  check exp/nbf/aud/iss
    │                                              │  extract sub=42 -> query
    │  200 OK {orders...}                          │
    │<─────────────────────────────────────────────│
```

Concrete HTTP:

```http
POST /login HTTP/1.1
Content-Type: application/json

{"email": "labib@example.com", "password": "hunter2"}
```

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MiIsImlzcyI6ImF1dGgubXlhcHAuY29tIiwiYXVkIjoib3JkZXJzLWFwaSIsImlhdCI6MTcyMDUzNzIwMCwiZXhwIjoxNzIwNTM4MTAwfQ.f9x...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MiIsImp0aSI6ImExYjJjMyIsImV4cCI6MTcyMTE0MjAwMH0.k2p..."
}
```

```http
GET /orders HTTP/1.1
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

If the token is bad/expired:

```http
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Bearer error="invalid_token"

{"error": "token expired"}
```

---

## 11. Access tokens vs refresh tokens

Short-lived tokens are safer (smaller window if stolen), but forcing a full re-login every 15 minutes is unusable. The fix: two tokens with different jobs.

|           | Access token                        | Refresh token                                                 |
| --------- | ----------------------------------- | ------------------------------------------------------------- |
| Purpose   | sent on every API request           | used only to get a new access token                           |
| Lifetime  | short (5–15 min typical)            | long (days–weeks)                                             |
| Sent to   | every resource server               | only the auth server's `/refresh` endpoint                    |
| Storage   | memory / short-lived                | HttpOnly cookie or secure storage, never localStorage ideally |
| If leaked | limited blast radius (expires soon) | serious — attacker can mint new access tokens repeatedly      |

Flow when access token expires:

```
┌────────┐                                    ┌────────┐
│ Client │                                    │ Server │
└───┬────┘                                    └───┬────┘
    │ GET /orders  (access_token expired)          │
    │──────────────────────────────────────────────>│
    │ 401 Unauthorized                              │
    │<────────────────────────────────────────────── │
    │                                                │
    │ POST /refresh  {refresh_token}                │
    │───────────────────────────────────────────────>│
    │                                    verify refresh_token
    │                                    check it's not revoked (jti lookup)
    │                                    issue NEW access_token
    │                                    (often: NEW refresh_token too — rotation)
    │ 200 {access_token, refresh_token}              │
    │<────────────────────────────────────────────── │
    │ retry GET /orders with new access_token        │
    │───────────────────────────────────────────────>│
```

**Refresh token rotation:** every time a refresh token is used, issue a brand new one and invalidate the old one (store its `jti` as "used"). If an old, already-rotated refresh token is ever presented again, that's a signal it was stolen and replayed — revoke the entire token family immediately. This is the standard mitigation against long-lived refresh token theft.

---

## 12. Expiration, rotation, revocation, logout — the part that fights JWT's core design

This is where the "stateless" promise gets complicated in practice, and it's worth being explicit about why.

**The core problem:** a JWT is valid until `exp`, full stop. The server that issued it has no way to reach out and invalidate it early, because nobody's _checking in_ with the server — that was the whole point. So "logout" and "revoke this specific token because the account was compromised" don't have a free, stateless answer.

Mitigations, roughly cheapest-and-weakest to most-robust-and-most-stateful:

1. **Short expiry on access tokens (5–15 min).** Cheapest mitigation — worst case exposure window is small. Doesn't solve immediate revocation, just bounds the damage.

2. **Client-side logout = delete the token.** Works for "this browser tab logs out," does nothing if the token was copied/stolen — it's still valid on the server until `exp`.

3. **Refresh token blacklist/allowlist (server-side, but only for refresh tokens, not access tokens).** Store `jti` of issued/revoked refresh tokens in Redis. On `/refresh`, check the blacklist. This re-introduces a small amount of statefulness, but only on the _low-frequency_ refresh path, not on every single API call — so you keep most of JWT's scaling benefit while regaining revocation for the token that matters most (since the refresh token can mint new access tokens indefinitely).

4. **Access token blacklist (full statefulness, defeats most of the point).** If you truly need instant revocation of an access token mid-flight (e.g., "kill this session right now, admin banned this user"), you're back to a lookup on every request — you've reintroduced sessions, just keyed by `jti` instead of a random session id. This is a legitimate design when instant revocation is a hard requirement (e.g., security-critical apps), but be honest that at that point you've given up the main advantage of JWT for that check.

5. **Short-lived access token + versioned/rotated signing key** is another lever: bump a `key_version` or a per-user `token_version` claim; on password change or forced logout-everywhere, increment `token_version` in the DB and reject tokens whose claim doesn't match current DB value. Still a DB check, but a very cheap one (single field), and it gives you "invalidate all sessions for this user" cheaply.

**The honest takeaway:** pure stateless JWT gives you _no_ revocation. Every real system re-adds a small amount of state at exactly the point where revocation matters (refresh tokens, or a version counter) while keeping the high-frequency access-token check itself lookup-free. Know that you are choosing a point on this stateless↔stateful spectrum, not "using JWT vs not."

---

## 13. Security best practices

- **Always verify the signature** before trusting anything in the payload — this sounds obvious but is the single most common real-world JWT bug (see `alg: none` below).
- **Pin the expected algorithm** on the verifying side. Don't let the token's `alg` header dictate which verification code path runs — an attacker who controls the header can otherwise trick a naive verifier.
- **Use RS256/ES256, not HS256, whenever more than one service needs to verify tokens.** HS256 means every verifier can also forge.
- **Short expiry on access tokens.** Long expiry on refresh tokens, but store/rotate them.
- **Store tokens carefully on the client:**
  - Access token: memory (JS variable) is safest for SPAs; avoid `localStorage` (readable by any injected script → XSS steals it).
  - Refresh token: `HttpOnly; Secure; SameSite=Strict` cookie so JS can't read it at all, mitigating XSS theft (CSRF becomes the concern instead — mitigate with `SameSite` + CSRF tokens for state-changing requests).
- **Always use HTTPS.** A JWT sent over plain HTTP is just handing your session to anyone on the network.
- **Validate `aud`/`iss` on every verify**, not just the signature — prevents token-confusion across services.
- **Keep the payload small and non-sensitive.** It's base64, not encrypted, and it's sent on every request — bloated payloads cost bandwidth on literally every call.
- **Constant-time signature comparison** (shown in section 6) to avoid timing side-channels.
- **Rotate signing keys periodically**, and support key IDs (`kid` header claim) so verifiers know which public key to use during rotation — lets you rotate without breaking in-flight tokens.

---

## 14. Common attacks and mitigations

**1. `alg: none` attack**
Some early JWT libraries let the token specify `"alg": "none"` and would skip signature verification entirely if you didn't explicitly pin the algorithm. An attacker crafts `{"alg":"none"}`, sets whatever payload they want, and the "signature" segment is just empty.
_Mitigation:_ explicitly require and check the algorithm on the verify side, never trust the header's declared `alg` to select behavior. Modern libraries fixed this by default, but always confirm.

**2. Algorithm confusion (RS256 → HS256 downgrade)**
If a server is configured to accept both RS256 and HS256, an attacker can take the server's known-public RSA **public** key, and use it as the HMAC **secret** to sign a forged HS256 token. Since the server naively treats "verify with whatever key I have" the same way regardless of algorithm, and HMAC verification with that "secret" succeeds, the forged token passes.
_Mitigation:_ never accept multiple algorithm families for the same key; hard-pin exactly one algorithm per verifier.

**3. Token theft via XSS**
If the access token sits in `localStorage`, any injected script can read `localStorage.getItem('token')` and exfiltrate it.
_Mitigation:_ keep tokens out of `localStorage`; use `HttpOnly` cookies for refresh tokens and short-lived in-memory tokens for access; sanitize inputs to prevent XSS in the first place.

**4. Token replay**
A stolen-but-valid token can be reused by the attacker until `exp`.
_Mitigation:_ short expiry, HTTPS everywhere to prevent interception, `jti` + refresh-rotation to detect reuse of an already-used refresh token (section 11).

**5. CSRF (if using cookies to carry the token)**
If a JWT rides in a plain cookie, browsers attach cookies automatically to cross-site requests, letting a malicious site trigger authenticated requests.
_Mitigation:_ `SameSite=Strict` or `Lax` cookies, plus a CSRF token for state-changing requests, or avoid cookies for the access token entirely (send it explicitly in the `Authorization` header instead, which cross-site forms can't do).

**6. Weak/guessable HMAC secret**
Some deployments use a short or default secret for HS256. Brute-forceable offline once an attacker has one valid token + signature to test guesses against.
_Mitigation:_ use a long, random, high-entropy secret (32+ bytes), or just use RS256/ES256 to sidestep the whole shared-secret problem.

**7. `kid` header injection**
The `kid` (key ID) header tells the verifier which key to use — if unsanitized, an attacker can craft a `kid` value that causes the server to look up a file path or run a malformed query, effectively controlling which "key" gets used to verify (sometimes tricking it into using attacker-supplied data as the key).
_Mitigation:_ treat `kid` as untrusted input; validate against a strict allowlist of known key IDs, never use it to build file paths or queries directly.

---

## 15. When to use JWT, and when not to

**Good fit:**

- Multiple independent services need to verify identity without a shared session store (microservices, exactly your kind of setup).
- Third-party API access delegation (OAuth2 access tokens are commonly JWTs).
- Mobile/native clients where cookie semantics are awkward.
- You can tolerate a short window where a revoked token is still technically valid (bounded by short expiry).

**Poor fit / reconsider:**

- Single monolith, one database, no cross-service verification needed — a session with a Redis lookup is simpler, gives instant revocation for free, and isn't meaningfully slower at reasonable scale. Don't reach for JWT by default just because it's trendy.
- You need **instant, hard revocation** as a core requirement (e.g., banking session kill-switch) — you'll end up rebuilding a server-side lookup anyway, so a session might just be more honest architecture.
- You're tempted to stuff large or sensitive data into the payload — that's a mismatch (base64 isn't encryption, and payload size hits every request).

---

## 16. Summary — the mental model to keep

- A session is a **reference** to server-side truth. A JWT **is** the truth, signed so it can't be forged.
- That single design choice — self-contained vs pointer — is the source of every property discussed above: no-lookup verification (the benefit), and painful revocation (the cost).
- Everything you'll see layered on top in real systems — refresh tokens, rotation, `jti` blacklists, `token_version` claims — exists to claw back a controlled, minimal amount of statefulness at the one point where you actually need it, without giving up the stateless win everywhere else.

For your Go services specifically: given you're already doing multi-service DDD work, RS256 (auth-service holds the private key, other services hold only the public key/JWKS) is the natural production choice over HS256 the moment a second service needs to verify tokens independently.
