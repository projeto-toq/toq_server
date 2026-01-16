O sistema de agendamento de sessão de fotos do Listing, descrita em `/codigos/go_code/toq_server/docs/photo_session_flow.md`, necessita 2 refatorações:

1) A busca de slots disponíveis de fotógrafo, via endpoint `GET /listings/photo-session/slots` permite apenas escolher dia e período (manhã/tarde), sem definir um horário fixo, o que gera imprecisão na agenda.
    1.2) Cliente solicitou alterar o tempo padrão de sessão de fotos para 2 horas, que deve ser por uma variável de env.yaml, e pode ser em qualquer horário do dia.
    1.3) A busca por slots deve retornar todos os horários disponíveis em blocos de 2 horas, dentro da disponibilidade do fotografo.

2) Endpoint: `POST /api/v2/listings/photo-session/reserve` para Reserva de um slot pelo Owner, tem como resposta poucos dados de identificação.
    2.1) Necessário popular endpoint com dados/foto do fotógrafo.

3) Atualize o documento `docs/photo_session_flow.md` para refletir as mudanças feitas no código.


Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e só então proponha o plano de correção.