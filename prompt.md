


Como a URL assinada para foto será sempre user-%d-bucket + photo, não há necessidade de salvar na tabela users. Assim retirei o campo da tabela.

***Baseado que este projeto tem como princípios e que sempre devem ser respeitados:***
1) Utilização das melhroes práticas GO;(https://go.dev/talks/2013/bestpractices.slide#1, https://google.github.io/styleguide/go/)
2) Arquitetura hexagonal;
3) Não utilização de MOCK e implementação efetiva sempre;
4) Manter padrão de desenvolvimento entre as funções;
5) Preservar a injeção de dependencias atualmente implementada;

faça a implementação sugerida e altere o getprofile/updateprofile para gerar sempre uma nova Url assinada ao invés de buscar no banco pelo campo photo
***apresente as alterações e só implemente quando eu autorizar***

[UserProfile] UserProfile.fromGrpc - photoUrl recebida: "https://storage.googleapis.com/user-2-bucket/photo.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Credential=gcs-reader%40generated-arena-468217-v3.iam.gserviceaccount.com%2F20250813%2Fauto%2Fstorage%2Fgoog4_request&X-Goog-Date=20250813T194742Z&X-Goog-Expires=3599&X-Goog-Signature=9350de36a7cee9126f953fd721df3cdec02fa2bbf8860f19ed5e27eb602d3a4081c640b2097afb2fd7bf07b5d6bf772482e906ddde28735bac09e95fd065ddb24ed9cf2925d47db3788dc9b1d198c3526e0f9c9d5e73c3923c0c527906d338be35492750a98bc84ba8147c2de6c298c8b41c92185e7703d2b3d8569108cd8055ba03e15672fd0255bcdeec2d237fe921dcd7e9dab26f8ace92867f37652f31742333bc201b0551966bdecb3ce1c4e24338d75a53695233bf6a3b708044766b2173e85ea7c325608d6ebcbdc1195e0737945444cdcd77564b42dd91971c8c3ed008a249b29316049e43c7adfe6ccdd5d776478b9a276477c363f2877456b880af&X-Goog-SignedHeaders=host"


[ProfilePhoto] _Exception (Exception: Failed to upload file: 403 <?xml version='1.0' encoding='UTF-8'?><Error><Code>SignatureDoesNotMatch</Code><Message>Access denied.</Message><Details>The request signature we calculated does not match the signature you provided. Check your Google secret key and signing method.</Details><StringToSign>GOOG4-RSA-SHA256