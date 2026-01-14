Atualmente o processo de processamento de medias, confrome descrito em :
- `/codigos/go_code/toq_server/docs/aws_media_processing_implementation_summary.md`
- `/codigos/go_code/toq_server/docs/aws_media_processing_useful_commands.md`
- `/codigos/go_code/toq_server/docs/media_processing_guide.md`

Deveria, caso o usuário apague uma media através do endpoint `DELETE /listings/media/delete`, apagar a media do do diretório `raw` e do diretório `processed` do bucket S3.

Creio que isso não está acontecendo hoje.

Busque todas as informações que precisa consultando as configurações reais da AWS e o código, não confiando na documentação, para ter certeza da causa raiz, e só então proponha o plano de correção.