Quero que atue como um engenheiro senior de infraestrutura e toda interação seja em português do Brasil.
Estamos em uma instancia EC2 da AWS rodando Debian 13 e temos um NGINX instalado e configurado como proxy reverso para as seguintes aplicações listadas abaixo. o NGINX responde ao dominio www.gca.dev.br usando um certificado SSL da letsencrypt.
- toq_server rodando na porta 8000. uma servidor rest-api desenvolvido em go. é o proejto atual do vscode. o caminho gca.dev.br/app direciona para a api.
- grafana rodando em docker, vide docker compose, e o caminho grafana.gca.dev.br direciona para o grafana.
- swagger rodando em docker, vide docker compose, e o caminho swagger.gca.dev.br direciona para o swagger.
- jaeger rodando em docker, vide docker compose, e o caminho jaeger.gca.dev.br direciona para o jaeger.

estávamos enfrentando problemas de CORS e CSP, e para resolver isso foram acrescentados parametros no NGINX, entretanto, após esta configuração o swagger parou de funcionar.

A console do browser apresenta o erro listado abaixo. Ao carregar o swagger.gca.dev.br, a tela carrega uma barra de navegação, com o caminho https://api.gca.dev.br/docs/swagger.json e a mensagem Fetch error
Failed to fetch https://api.gca.dev.br/docs/swagger.json
Fetch error
Possible cross-origin (CORS) issue? The URL origin (https://api.gca.dev.br) does not match the page (https://swagger.gca.dev.br). Check the server returns the correct 'Access-Control-Allow-*' headers.

# analise as configurações do NGINX e sugira as alterações necessárias para resolver o problema do swagger, mantendo as outras aplicações funcionando corretamente. Não implemente nada sem minha autorização. Estamos na fase de diagnóstico.

Console do browser:
swagger-ui-bundle.js:2 [Report Only] Refused to apply inline style because it violates the following Content Security Policy directive: "style-src self unsafe-inline https:". Either the 'unsafe-inline' keyword, a hash ('sha256-RL3ie0nH+Lzz2YNqQN83mnU0J1ot4QL7b99vMdIX99w='), or a nonce ('nonce-...') is required to enable inline execution.

(index):1  Refused to load the image 'https://validator.swagger.io/validator?url=https%3A%2F%2Fapi.gca.dev.br%2Fdocs%2Fswagger.json' because it violates the following Content Security Policy directive: "img-src 'self' data: blob: https://*.amazonaws.com".

(index):1  Access to fetch at 'https://api.gca.dev.br/docs/swagger.json' from origin 'https://swagger.gca.dev.br' has been blocked by CORS policy: The 'Access-Control-Allow-Origin' header contains multiple values 'https://swagger.gca.dev.br, https://swagger.gca.dev.br', but only one is allowed. Have the server send the header with a valid value.
api.gca.dev.br/docs/swagger.json:1   Failed to load resource: net::ERR_FAILED
