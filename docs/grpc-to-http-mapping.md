# Mapeamento Completo gRPC → HTTP

## UserService (28 métodos)

### Autenticação (Public Endpoints)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| CreateOwner | /api/v1/auth/owner | POST | Criar conta de proprietário |
| CreateRealtor | /api/v1/auth/realtor | POST | Criar conta de corretor |
| CreateAgency | /api/v1/auth/agency | POST | Criar conta de imobiliária |
| SignIn | /api/v1/auth/signin | POST | Login do usuário |
| RefreshToken | /api/v1/auth/refresh | POST | Renovar token de acesso |
| RequestPasswordChange | /api/v1/auth/password/request | POST | Solicitar mudança de senha |
| ConfirmPasswordChange | /api/v1/auth/password/confirm | POST | Confirmar mudança de senha |
| SignOut | /api/v1/auth/signout | POST | Logout do usuário |

### Gestão de Perfil (Authenticated)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| GetProfile | /api/v1/user/profile | GET | Obter perfil do usuário |
| UpdateProfile | /api/v1/user/profile | PUT | Atualizar perfil |
| DeleteAccount | /api/v1/user/account | DELETE | Deletar conta |
| GetOnboardingStatus | /api/v1/user/onboarding | GET | Status do onboarding |
| GetUserRoles | /api/v1/user/roles | GET | Listar roles do usuário |
| GoHome | /api/v1/user/home | GET | Página inicial |
| UpdateOptStatus | /api/v1/user/opt-status | PUT | Atualizar status de opt-in |

### Gestão de Fotos
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| GetPhotoUploadURL | /api/v1/user/photo/upload-url | POST | URL para upload de foto |
| GetProfileThumbnails | /api/v1/user/profile/thumbnails | GET | Miniaturas do perfil |

### Mudança de Email
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| RequestEmailChange | /api/v1/user/email/request | POST | Solicitar mudança de email |
| ConfirmEmailChange | /api/v1/user/email/confirm | POST | Confirmar mudança de email |
| ResendEmailChangeCode | /api/v1/user/email/resend | POST | Reenviar código de email |

### Mudança de Telefone
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| RequestPhoneChange | /api/v1/user/phone/request | POST | Solicitar mudança de telefone |
| ConfirmPhoneChange | /api/v1/user/phone/confirm | POST | Confirmar mudança de telefone |
| ResendPhoneChangeCode | /api/v1/user/phone/resend | POST | Reenviar código de telefone |

### Gestão de Roles (Owner/Realtor only)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| AddAlternativeUserRole | /api/v1/user/role/alternative | POST | Adicionar role alternativo |
| SwitchUserRole | /api/v1/user/role/switch | POST | Trocar role ativo |

### Operações de Imobiliária (Agency only)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| GetDocumentsUploadURL | /api/v1/agency/documents/upload-url | POST | URL para upload de documentos |
| InviteRealtor | /api/v1/agency/invite-realtor | POST | Convidar corretor |
| GetRealtorsByAgency | /api/v1/agency/realtors | GET | Listar corretores da imobiliária |
| GetRealtorByID | /api/v1/agency/realtors/:id | GET | Obter corretor por ID |
| DeleteRealtorByID | /api/v1/agency/realtors/:id | DELETE | Remover corretor |

### Operações de Corretor (Realtor only)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| VerifyCreciImages | /api/v1/realtor/creci/verify | POST | Verificar imagens do CRECI |
| GetCreciUploadURL | /api/v1/realtor/creci/upload-url | POST | URL para upload do CRECI |
| AcceptInvitation | /api/v1/realtor/invitation/accept | POST | Aceitar convite |
| RejectInvitation | /api/v1/realtor/invitation/reject | POST | Rejeitar convite |
| GetAgencyOfRealtor | /api/v1/realtor/agency | GET | Obter imobiliária do corretor |
| DeleteAgencyOfRealtor | /api/v1/realtor/agency | DELETE | Sair da imobiliária |

---

## ListingService (24 métodos)

### Gestão Básica de Imóveis
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| GetAllListings | /api/v1/listings | GET | Listar todos os imóveis |
| StartListing | /api/v1/listings | POST | Iniciar novo anúncio |
| GetListing | /api/v1/listings/:id | GET | Obter imóvel específico |
| UpdateListing | /api/v1/listings/:id | PUT | Atualizar imóvel |
| DeleteListing | /api/v1/listings/:id | DELETE | Deletar imóvel |
| SearchListing | /api/v1/listings/search | GET | Buscar imóveis |

### Configuração e Opções
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| GetOptions | /api/v1/listings/options | GET | Obter opções de imóvel |
| GetBaseFeatures | /api/v1/listings/features/base | GET | Obter características base |

### Workflow de Imóveis (Owner side)
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| EndUpdateListing | /api/v1/listings/:id/end-update | POST | Finalizar atualização |
| GetListingStatus | /api/v1/listings/:id/status | GET | Status do imóvel |
| ApproveListing | /api/v1/listings/:id/approve | POST | Aprovar imóvel |
| RejectListing | /api/v1/listings/:id/reject | POST | Rejeitar imóvel |
| SuspendListing | /api/v1/listings/:id/suspend | POST | Suspender imóvel |
| ReleaseListing | /api/v1/listings/:id/release | POST | Liberar imóvel |
| CopyListing | /api/v1/listings/:id/copy | POST | Copiar imóvel |

### Operações de Corretor
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| ShareListing | /api/v1/listings/:id/share | POST | Compartilhar imóvel |
| GetFavoriteListings | /api/v1/listings/favorites | GET | Imóveis favoritos |
| AddFavoriteListing | /api/v1/listings/:id/favorite | POST | Adicionar aos favoritos |
| RemoveFavoriteListing | /api/v1/listings/:id/favorite | DELETE | Remover dos favoritos |

### Gestão de Visitas
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| RequestVisit | /api/v1/listings/:id/visit/request | POST | Solicitar visita |
| GetAllVisits | /api/v1/visits | GET | Listar todas as visitas |
| GetVisits | /api/v1/listings/:id/visits | GET | Visitas de um imóvel |
| CancelVisit | /api/v1/visits/:id | DELETE | Cancelar visita |
| ConfirmVisitDone | /api/v1/visits/:id/confirm | POST | Confirmar visita realizada |
| ApproveVisting | /api/v1/visits/:id/approve | POST | Aprovar visita |
| RejectVisting | /api/v1/visits/:id/reject | POST | Rejeitar visita |

### Gestão de Ofertas
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| CreateOffer | /api/v1/listings/:id/offers | POST | Criar oferta |
| GetAllOffers | /api/v1/offers | GET | Listar todas as ofertas |
| GetOffers | /api/v1/listings/:id/offers | GET | Ofertas de um imóvel |
| UpdateOffer | /api/v1/offers/:id | PUT | Atualizar oferta |
| SendOffer | /api/v1/offers/:id/send | POST | Enviar oferta |
| CancelOffer | /api/v1/offers/:id | DELETE | Cancelar oferta |
| ApproveOffer | /api/v1/offers/:id/approve | POST | Aprovar oferta |
| RejectOffer | /api/v1/offers/:id/reject | POST | Rejeitar oferta |

### Avaliações
| gRPC Method | HTTP Endpoint | HTTP Method | Descrição |
|-------------|---------------|-------------|-----------|
| EvaluateRealtor | /api/v1/realtors/:id/evaluate | POST | Avaliar corretor |
| EvaluateOwner | /api/v1/owners/:id/evaluate | POST | Avaliar proprietário |

---

## Resumo
- **Total UserService**: 28 métodos → 28 endpoints HTTP
- **Total ListingService**: 24 métodos → 24 endpoints HTTP  
- **Total Geral**: 52 métodos gRPC → 52 endpoints HTTP

Todos os métodos gRPC foram mapeados para endpoints HTTP RESTful seguindo as melhores práticas de design de APIs.
