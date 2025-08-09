
Sequência recomendada (JWT Access + Refresh) com boas práticas:

Login inicial

Client -> POST /auth/login (credenciais + device_id opcional).
Server valida credenciais, cria uma sessão (session_id) e persiste: user_id, device_id, refresh_jti, hash(refresh_token), exp, rotation_counter, ip, user_agent.
Server gera:
Access Token (TTL curto: 5–15 min) – contém jti, sub, aud, iat, exp, scopes mínimos.
Refresh Token (TTL longo: 7–30 dias) – inclui jti diferente, session_id, rotation_counter.
Entrega:
Web: Refresh em cookie HttpOnly + Secure + SameSite=Strict/Lax; Access no corpo ou header (client guarda só em memória).
Mobile: Refresh em secure storage (Keychain / Keystore); Access só em memória (não persistir).
Uso normal de APIs

Client envia apenas Access Token no header Authorization: Bearer <access>.
Não enviar Refresh Token nas chamadas normais.
Antes do vencimento / erro 401 por expiração

Se faltam <60s (ou ao receber 401 expirado) iniciar refresh.
Client -> POST /auth/refresh.
Envia somente o Refresh Token (cookie ou body/json).
Opcional: envia o jti do último Access (para reforçar detecção de reuse) + device_id.
Validação do refresh no servidor

Verifica assinatura, exp, issuer/audience.
Busca sessão pelo refresh_jti (hash) ou session_id.
Checa: não revogado, não usado (rotações anteriores), device/ip (se política exigir), rotation_counter coerente.
(Opcional) Reuse detection: se token já marcado usado => revoga toda a sessão e força re-login.
Rotação

Gera novo Access Token (novo jti).
Gera novo Refresh Token (novo jti e rotation_counter+1).
Marca antigo refresh como usado/revogado (store: previous_jti, used_at).
Atualiza metadados (last_refresh_at, ip, user_agent).
Resposta do refresh

Retorna novo Access + novo Refresh (substitui completamente os anteriores).
Web: set cookie HttpOnly novamente (overwrite).
Mobile: sobrescreve secure storage.
Concurrency / múltiplos refresh

Aceitar só o primeiro. Segundo simultâneo com o mesmo refresh antigo = reuse => invalidar sessão.
Implementar lock otimista (checar estado ‘active’ antes de rotacionar).
Logout

Client -> POST /auth/logout.
Server marca sessão e refresh atuais como revogados; opcional: adiciona access_jti à deny list até expirar.
Revogação / segurança extra

Maintain deny list de jti de Access (só até exp deles) se precisar revogar antes do exp natural.
Opcional: throttle refresh endpoint (ex: max 3/min por sessão) para mitigar brute force.
Armazenamento seguro

Nunca colocar Refresh em localStorage/sessionStorage.
Mobile: usar Keychain (iOS) / EncryptedSharedPreferences ou Keystore (Android).
Web: sempre cookie HttpOnly (evita XSS). Usar SameSite + CSRF token se SameSite=Lax/None.
Claims mínimas
Access: user_id (sub), roles/scopes, jti, exp curto.
Refresh: evitar dados sensíveis; usar jti + session_id; manter detalhes server-side.
Expirações
Access curto (5–15 min).
Refresh rotacionado pode ter sliding window (renova até limite absoluto, ex: 30 dias).
Expiração absoluta da sessão (hard cap) força novo login.
Detecção de comprometimento
Se reuse detectado: revogar sessão, invalidar todos tokens, log de segurança, exigir reautenticação forte.
Tratamento no cliente
Fila de requisições pendentes: se 401 expirado, pausar, fazer refresh uma vez, refazer pendentes.
Evitar múltiplos refresh concorrentes (usar mutex/singleflight client-side).
Fluxo resumido (happy path): Login -> (Access curto + Refresh longo) -> Chamadas com Access -> Expira -> POST /auth/refresh (envia Refresh) -> Server valida + rotaciona -> Novo Access + Refresh -> Continua.

Motivos para não tornar tudo síncrono com refresh silencioso a cada request:

Aumenta superfície de ataque (refresh token mais exposto).
Aumenta carga do endpoint de refresh.
Perde benefício de TTL curto.
Quando usar método síncrono especial:

Cenários onde a confirmação da entrega (ex: push crítico) depende de retorno imediato – use método sync opcional; caso contrário padrão assíncrono.
Principais erros a evitar:

Reusar refresh token sem rotação.
Guardar refresh em localStorage.
Access TTL muito longo.
Não registrar metadados (ip, device) por sessão.
Não detectar reuse.
Essencial: Rotação + Detecção de Reuso + Armazenamento seguro + TTL curto de Access.