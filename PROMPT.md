o arquivo de rotas em `/codigos/go_code/toq_server/internal/adapter/left/http/routes/routes.go` está poluido com diversas rotas criadas e não implementadas, devolvendo `func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }`.

É necessário excluir estas rotas não implementadas para melhorar a clareza do código e evitar confusões futuras.

Esta remoção deve incluir esqueletos nos handlers e quaisquer referências associadas a essas rotas.

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e crie um plano para resolver a causa raiz identificada, com evidencia se sem suposições.