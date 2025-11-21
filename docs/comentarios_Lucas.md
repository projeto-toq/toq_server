## Na pergunta “De quem é o imóvel?”, não existe uma opção específica para casos em que o corretor/imobiliária tem o imóvel em gestão ou exclusividade, impossibilitando diferenciar esses cadastros especiais.
Estes são os tipos atuais:
1;property_owner;1;myself;Myself;;1
2;property_owner;2;spouse;Spouse;;1
3;property_owner;3;parents;Parents;;1
4;property_owner;4;grandparents;GrandParents;;1
5;property_owner;5;children;Children;;1
6;property_owner;6;uncles;Uncles;;1
7;property_owner;7;siblings;Siblings;;1

quais quer incluir e com qual descrição: exemplo
8;property_owner;7;agency;"Agency Exclusive";;1
7;property_owner;7;realtor;"Realtor Exclusive";;1


## Na etapa de dimensões/comodidades do apartamento, há um campo de “Área edificável” além da metragem já definida, o que é redundante e confunde o usuário.

Talvez seja `Área não edificável`  deveria er em casa/tereno residencial apenas. Na base temos:

landSize                float64
nonBuildable            float64
buildable               float64

## Na parte em que o usuário informa onde aceita permuta, é perguntado bairro, cidade e UF, sendo que o nível de detalhe de bairro é desnecessário nesse momento.
Certeza que quer remover? EU aceito uma permuta da cidade de São Paulo nos jardins ou Aclimação, mas nÃo aceito no Jardim Miriam. Só São Paulo, deixa muita margem para discussão.

## No agendamento de fotógrafo, o sistema permite apenas escolher dia e período (manhã/tarde), sem definir um horário fixo, o que gera imprecisão na agenda.
Considerando que são 4 horas minimo para fazer a sessão de fotos, vale colocar aberto e o cliente selecionar as 17:00 e não ser possível realizar ou as 11:00 e matar o dia do fotógrafo?

## Depois de selecionar dia/horário do fotógrafo, a tela mostra um card com o fotógrafo, porém sem foto de rosto e com poucos dados de identificação.
Necessário popular endpoint de confirmação da sessão de fotos com dados/foto do fotógrado.

## Ao cadastrar comodidades de um prédio, são exibidas comodidades típicas de imóvel residencial (suítes, dormitórios, cozinha, copa), em vez de comodidades comerciais adequadas (salas, recepção, vagas, auditório etc.).
Informe lista de comodidades, em ordem de aparecimento, por tipo de imóvel para que seja incluída
