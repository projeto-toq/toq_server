É necessário que existam rotinas de limpeza para
    `toq_db.device_tokens` sugira uma política de retenção, configurável em env.yaml
    `toq_db.sessions` sugeira uma política de retenção, configurável em env.yaml
    `toq_db.holiday_calendar_dates` anteriores a 1 anos, configurável em env.yaml
    `toq_db.media_processing_jobs`sugira uma política de retenção, configurável em env.yaml
    `toq_db.photographer_agenda_entries`anteriores a 1 anos, configurável em env.yaml
    `toq_db.photographer_photo_session_bookings` anteriores a 1 anos, configurável em env.yaml

Talvez já existam go routines de limpeza em `/codigos/go_code/toq_server/internal/core/go_routines` para alguns desses itens, ou talvez existam partes do código que já façam isso, mas é necessário revisar todo o código e as configurações atuais para garantir que todas essas tabelas tenham rotinas de limpeza adequadas, com políticas de retenção configuráveis via env.yaml.

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e só então proponha o plano de correção.