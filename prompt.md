A estrutura de bucket do GCS mudou agora teremos:
toq_server_users_media/
├── {user_id}/                    # ✅ Estrutura atual mantida
│   ├── photo.jpg                 # ✅ Foto atual (manter compatibilidade)
│   │── thumbnails/           # Diferentes resoluções
│   │   ├── small.jpg         # 150x150
│   │   ├── medium.jpg        # 300x300
│   │   └── large.jpg         # 600x600

portanto responder ao rpc GetProfile com a urlAssinada da photo no campo photo nã faz sentido, pois temos 4 tipos diferentes.

***Baseado que este projeto tem como princípios e que sempre devem ser respeitados:***
1) Utilização das melhroes práticas GO;(https://go.dev/talks/2013/bestpractices.slide#1, https://google.github.io/styleguide/go/)
2) Arquitetura hexagonal;
3) Não utilização de MOCK e implementação efetiva sempre;
4) Manter padrão de desenvolvimento entre as funções;
5) Preservar a injeção de dependencias atualmente implementada;

Assim, verifique as alterações necessárias para implementar o rpc GetProfileThumbnails com
message GetProfileThumbnailsRequest {
}

message GetProfileThumbnailsResponse {
    string originalUrl = 1;   // photo.jpg
    string smallUrl = 2;      // thumbnails/small.jpg  
    string mediumUrl = 3;     // thumbnails/medium.jpg
    string largeUrl = 4;      // thumbnails/large.jpg
}

***apresente as alterações e só implemente quando eu autorizar***
