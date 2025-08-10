O rpc PushOptIn (PushOptInRequest) returns (PushOptInResponse) foi alterado e não mas envia device token com parametro, já que ele é sempre enviado no login.
Assim, mantendo a arquitetura hexagonal já implementada, as boas praticas de progamação GO, criando a documetnação de função e a documentação interna para explicar a logica e não criando nanhum tipo de mock:
altere o handler, serviços e outras funcções necessárias para que ao chamar o rpc seja ajustado o campo opt_status da tabela users para 1
