Ao criar um listing do tipo casa em contrução não haverá sessão de fotos como nos outros property types. Entretanto o owner, durante o cadastro, poderá fazer upload da planta da casa em construção, um pdf e fotos de renderizaç~eos de como ficará, em jpeg.

A ideia original é que o owner pudesse usar o mesmo endpoint de upload de fotos, com variantes planta e render, adicionais aos foto horizontal, foto vertical e videos.

Entretanto isso pode complicar o processo visto que os endpoints atuais de upload de fotos são endpoints utilizados pelo Photographer, que faz upload de fotos já associadas a um listing publicado, e não durante o cadastro do listing.

Assim, verifique se os endpints `/listings/media/*` podem ser usados para esse propósito, ou se é necessário criar endpoints específicos para upload de planta e renders durante o cadastro do listing.

Caso seja necessário, ou seja recomendado, criar estes endpoints, defina a estrutura dos endpoints, os parâmetros necessários, e como o fluxo de upload funcionará para o owner durante o cadastro do listing de casa em construção.

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e crie um plano esta refatoração, com evidencia se sem suposições.