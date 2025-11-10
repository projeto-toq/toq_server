./internal/adapter/right/mysql/user
├── activate_user_role.go
├── batch_update_user_last_activity.go
    ├── converters
    │   ├── agency_invite_domain_to_entity.go
    │   ├── agency_invite_entity_to_domain.go
    │   ├── user_domain_to_entity.go
    │   ├── user_entity_to_domain.go
    │   ├── user_role_domain_to_entity.go
    │   ├── user_role_entity_to_domain.go
    │   ├── user_validation_domain_to_entity.go
    │   ├── user_validation_entity_to_domain.go
    │   ├── user_with_role_entity_to_domain.go
    │   ├── wrong_signin_domain_to_entity.go
    │   └── wrong_signin_entity_to_domain.go
├── create_agency_invite.go
├── create_agency_relationship.go
├── create_user.go
├── create_user_role.go
├── deactivate_all_user_roles.go
├── delete_agency_realtor_relation.go
├── delete_expired_validations.go
├── delete_invite.go
├── delete_user_role.go
├── delete_user_roles_by_userid.go
├── delete_validation.go
├── delete_wrong_signin_by_userid.go
├── entities
│   ├── agency_invite_entity.go
│   ├── realtors_agency_entity.go
│   ├── user_entity.go
│   ├── user_role_entity.go
│   ├── user_validation_entity.go
│   ├── user_with_role_entity.go
│   └── wrong_signin_entity.go
├── exists_email_for_another_user.go
├── exists_phone_for_another_user.go
├── get_active_user_role_by_user_id.go
├── get_agency_of_realtor.go
├── get_invite_by_phone_number.go
├── get_realtors_by_agency.go
├── get_user_by_id.go
├── get_user_by_nationalid.go
├── get_user_by_phone_number.go
├── get_user_role_by_user_id_and_role_id.go
├── get_user_roles_by_user_id.go
├── get_user_validations.go
├── get_users_by_role_and_status.go
├── get_wrong_signin_by_userid.go
├── has_user_duplicate.go
├── list_all_users.go
├── list_users.go
├── reset_wrong_signin_attempts.go
├── rows_helpers.go
├── scan_helpers.go
├── update_agency_invite_by_id.go
├── update_user_by_id.go
├── update_user_last_activity.go
├── update_user_password_by_id.go
├── update_user_role.go
├── update_user_role_status.go
├── update_user_role_status_tx.go
├── update_user_validations.go
├── update_wrong_signin.go
├── user_adapter.go
└── user_blocking.go

3 directories, 63 files
